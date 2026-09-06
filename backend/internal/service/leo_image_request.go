package service

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/url"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

const (
	leoImageMultipartUnsupported = "This model accepts HTTP(S) image URLs via image_urls or images[].image_url; multipart uploads are not supported"
	leoImageMaskUnsupported      = "This model does not support mask; send image_urls or images[].image_url"
	leoImageDataURLUnsupported   = "This model does not accept data URLs; use an absolute HTTP(S) URL"
	leoImageURLRequired          = "image reference must be an absolute HTTP(S) URL"
	leoImageJSONRequired         = "This model requires a JSON request body"
)

type LeoImageRequestError struct {
	message string
}

func (e *LeoImageRequestError) Error() string {
	if e == nil {
		return ""
	}
	return SanitizeImageProviderMessage(e.message)
}

func newLeoImageRequestError(message string) error {
	return &LeoImageRequestError{message: message}
}

func isLeoImagePlatform(platform string) bool {
	return strings.EqualFold(strings.TrimSpace(platform), PlatformLeo)
}

func validateLeoImageParsedRequest(platform string, req *OpenAIImagesRequest, body []byte) error {
	if !isLeoImagePlatform(platform) || req == nil {
		return nil
	}
	if req.HasMask || strings.TrimSpace(req.MaskImageURL) != "" || req.MaskUpload != nil || gjson.GetBytes(body, "mask").Exists() {
		return newLeoImageRequestError(leoImageMaskUnsupported)
	}
	if req.Multipart {
		return nil
	}
	for _, imageURL := range collectLeoImageReferenceURLs(body) {
		if err := validateLeoImageReferenceURL(imageURL); err != nil {
			return err
		}
	}
	return nil
}

func BuildLeoMultipartImageRequest(req *OpenAIImagesRequest, imageURLs []string) ([]byte, error) {
	if req == nil || len(imageURLs) == 0 {
		return nil, newLeoImageRequestError(leoImageURLRequired)
	}
	if req.HasMask || strings.TrimSpace(req.MaskImageURL) != "" || req.MaskUpload != nil {
		return nil, newLeoImageRequestError(leoImageMaskUnsupported)
	}
	for _, imageURL := range imageURLs {
		if err := validateLeoImageReferenceURL(imageURL); err != nil {
			return nil, err
		}
	}

	payload := map[string]any{
		"model":      req.Model,
		"prompt":     req.Prompt,
		"n":          req.N,
		"image_urls": imageURLs,
	}
	setString := func(name, value string) {
		if strings.TrimSpace(value) != "" {
			payload[name] = value
		}
	}
	setString("size", req.Size)
	setString("response_format", req.ResponseFormat)
	setString("quality", req.Quality)
	setString("background", req.Background)
	setString("output_format", req.OutputFormat)
	setString("moderation", req.Moderation)
	setString("input_fidelity", req.InputFidelity)
	setString("style", req.Style)
	if req.Stream {
		payload["stream"] = true
	}
	if req.OutputCompression != nil {
		payload["output_compression"] = *req.OutputCompression
	}
	if req.PartialImages != nil {
		payload["partial_images"] = *req.PartialImages
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("encode multipart image request: %w", err)
	}
	req.Multipart = false
	req.ContentType = "application/json"
	req.InputImageURLs = append([]string(nil), imageURLs...)
	req.Uploads = nil
	req.Body = body
	return body, nil
}

func rewriteLeoImageUpstreamRequest(body []byte, contentType string, parsed *OpenAIImagesRequest) ([]byte, string, string, error) {
	if parsed != nil && parsed.Multipart {
		return nil, "", "", newLeoImageRequestError(leoImageMultipartUnsupported)
	}
	mediaType, _, err := mime.ParseMediaType(strings.TrimSpace(contentType))
	if err == nil && strings.EqualFold(mediaType, "multipart/form-data") {
		return nil, "", "", newLeoImageRequestError(leoImageMultipartUnsupported)
	}
	if !gjson.ValidBytes(body) {
		return nil, "", "", newLeoImageRequestError(leoImageJSONRequired)
	}
	if gjson.GetBytes(body, "mask").Exists() {
		return nil, "", "", newLeoImageRequestError(leoImageMaskUnsupported)
	}

	refs := collectLeoImageReferenceURLs(body)
	for _, imageURL := range refs {
		if err := validateLeoImageReferenceURL(imageURL); err != nil {
			return nil, "", "", err
		}
	}

	out := body
	if len(refs) > 0 {
		out, err = sjson.SetBytes(out, "image_urls", refs)
		if err != nil {
			return nil, "", "", fmt.Errorf("rewrite image_urls: %w", err)
		}
	}
	deletePaths := []string{"images", "image_url", "mask"}
	if extractLeoImageReferenceURL(gjson.GetBytes(body, "image")) != "" {
		deletePaths = append(deletePaths, "image")
	}
	for _, path := range deletePaths {
		if !gjson.GetBytes(out, path).Exists() {
			continue
		}
		out, err = sjson.DeleteBytes(out, path)
		if err != nil {
			return nil, "", "", fmt.Errorf("remove incompatible image field %s: %w", path, err)
		}
	}
	return out, "application/json", openAIImagesGenerationsEndpoint, nil
}

func collectLeoImageReferenceURLs(body []byte) []string {
	seen := map[string]struct{}{}
	var refs []string
	appendRef := func(raw string) {
		imageURL := strings.TrimSpace(raw)
		if imageURL == "" {
			return
		}
		if _, exists := seen[imageURL]; exists {
			return
		}
		seen[imageURL] = struct{}{}
		refs = append(refs, imageURL)
	}

	appendRef(extractLeoImageReferenceURL(gjson.GetBytes(body, "image_url")))
	for _, item := range gjson.GetBytes(body, "image_urls").Array() {
		appendRef(extractLeoImageReferenceURL(item))
	}
	appendRef(extractLeoImageReferenceURL(gjson.GetBytes(body, "image")))
	for _, item := range gjson.GetBytes(body, "images").Array() {
		appendRef(extractLeoImageReferenceURL(item))
	}
	return refs
}

func extractLeoImageReferenceURL(value gjson.Result) string {
	if !value.Exists() || value.Type == gjson.Null {
		return ""
	}
	if value.Type == gjson.String {
		return strings.TrimSpace(value.String())
	}
	if !value.IsObject() {
		return ""
	}
	for _, path := range []string{"image_url.url", "image_url", "url"} {
		if imageURL := strings.TrimSpace(value.Get(path).String()); imageURL != "" {
			return imageURL
		}
	}
	return ""
}

func validateLeoImageReferenceURL(raw string) error {
	imageURL := strings.TrimSpace(raw)
	if imageURL == "" {
		return newLeoImageRequestError(leoImageURLRequired)
	}
	if strings.HasPrefix(strings.ToLower(imageURL), "data:") {
		return newLeoImageRequestError(leoImageDataURLUnsupported)
	}
	parsed, err := url.Parse(imageURL)
	if err != nil || parsed == nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return newLeoImageRequestError(leoImageURLRequired)
	}
	return nil
}

func appendLeoNativeImageURLs(body []byte, req *OpenAIImagesRequest) {
	if req == nil {
		return
	}
	seen := map[string]struct{}{}
	for _, imageURL := range req.InputImageURLs {
		imageURL = strings.TrimSpace(imageURL)
		if imageURL != "" {
			seen[imageURL] = struct{}{}
		}
	}
	for _, imageURL := range collectLeoImageReferenceURLs(body) {
		if _, exists := seen[imageURL]; exists {
			continue
		}
		seen[imageURL] = struct{}{}
		req.InputImageURLs = append(req.InputImageURLs, imageURL)
	}
}
