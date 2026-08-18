package service

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

var publicImageVendorName = regexp.MustCompile(`(?i)leonardo|leostudio|leo\s*studio|\bleo\b|\bkrea\b`)

func TestCollectLeoImageReferenceURLs(t *testing.T) {
	body := []byte(`{
		"image_url": "https://example.com/one.png",
		"image_urls": ["https://example.com/one.png", "https://example.com/two.png"],
		"images": [
			{"image_url": "https://example.com/three.png"},
			{"image_url": {"url": "https://example.com/four.png"}},
			{"url": "https://example.com/five.png"},
			"https://example.com/six.png"
		]
	}`)

	require.Equal(t, []string{
		"https://example.com/one.png",
		"https://example.com/two.png",
		"https://example.com/three.png",
		"https://example.com/four.png",
		"https://example.com/five.png",
		"https://example.com/six.png",
	}, collectLeoImageReferenceURLs(body))
}

func TestRewriteLeoImageUpstreamRequestMapsOpenAIEdits(t *testing.T) {
	body := []byte(`{
		"model": "Image Model A",
		"prompt": "keep the product",
		"n": 1,
		"aspect_ratio": "1:1",
		"images": [{"image_url": "https://example.com/source.png"}]
	}`)

	out, contentType, endpoint, err := rewriteLeoImageUpstreamRequest(body, "application/json", &OpenAIImagesRequest{Endpoint: openAIImagesEditsEndpoint})
	require.NoError(t, err)
	require.Equal(t, "application/json", contentType)
	require.Equal(t, openAIImagesGenerationsEndpoint, endpoint)
	require.Equal(t, "https://example.com/source.png", gjson.GetBytes(out, "image_urls.0").String())
	require.False(t, gjson.GetBytes(out, "images").Exists())
	require.False(t, gjson.GetBytes(out, "image_url").Exists())
	require.Equal(t, "1:1", gjson.GetBytes(out, "aspect_ratio").String())
	require.Equal(t, "Image Model A", gjson.GetBytes(out, "model").String())
}

func TestRewriteLeoImageUpstreamRequestKeepsNativeImageURLs(t *testing.T) {
	body := []byte(`{"model":"Image Model A","prompt":"draw","image_url":"https://example.com/a.png","image_urls":["https://example.com/b.png"]}`)

	out, _, endpoint, err := rewriteLeoImageUpstreamRequest(body, "application/json", &OpenAIImagesRequest{Endpoint: openAIImagesGenerationsEndpoint})
	require.NoError(t, err)
	require.Equal(t, openAIImagesGenerationsEndpoint, endpoint)
	require.Equal(t, []string{"https://example.com/a.png", "https://example.com/b.png"}, jsonStringArray(gjson.GetBytes(out, "image_urls")))
	require.False(t, gjson.GetBytes(out, "image_url").Exists())
}

func TestRewriteLeoImageUpstreamRequestRejectsMaskAndMultipart(t *testing.T) {
	_, _, _, err := rewriteLeoImageUpstreamRequest([]byte(`{"prompt":"x","images":[{"image_url":"https://example.com/a.png"}],"mask":{"image_url":"https://example.com/mask.png"}}`), "application/json", &OpenAIImagesRequest{Endpoint: openAIImagesEditsEndpoint})
	require.ErrorContains(t, err, "does not support mask")

	_, _, _, err = rewriteLeoImageUpstreamRequest([]byte(`{"prompt":"x"}`), "multipart/form-data; boundary=x", &OpenAIImagesRequest{Multipart: true, Endpoint: openAIImagesEditsEndpoint})
	require.ErrorContains(t, err, "multipart uploads are not supported")
}

func TestValidateLeoImageReferenceURL(t *testing.T) {
	require.NoError(t, validateLeoImageReferenceURL("https://cdn.example/ref.png"))
	require.ErrorContains(t, validateLeoImageReferenceURL("data:image/png;base64,AAAA"), "data URLs")
	require.ErrorContains(t, validateLeoImageReferenceURL("/local.png"), "absolute HTTP(S) URL")
}

func TestOpenAIImagesUpstreamErrorClientMessageHidesVendorNames(t *testing.T) {
	err := &OpenAIImagesUpstreamError{Message: "Leonardo.ai rejected the request through LeoStudio and Krea"}
	require.Equal(t, "Image service rejected the request through Image service and Image service", err.ClientMessage())
	require.NotRegexp(t, publicImageVendorName, err.ClientMessage())
}

func TestLeoImagePublicErrorMessagesHideVendorNames(t *testing.T) {
	for _, message := range []string{
		leoImageMultipartUnsupported,
		leoImageMaskUnsupported,
		leoImageDataURLUnsupported,
		leoImageURLRequired,
		leoImageJSONRequired,
		(&LeoImageRequestError{message: leoImageMaskUnsupported}).Error(),
	} {
		require.NotRegexp(t, publicImageVendorName, message)
	}
}

func TestValidateLeoImageParsedRequest(t *testing.T) {
	require.NoError(t, validateLeoImageParsedRequest(PlatformOpenAI, &OpenAIImagesRequest{Multipart: true}, nil))
	require.ErrorContains(t, validateLeoImageParsedRequest(PlatformLeo, &OpenAIImagesRequest{Multipart: true}, nil), "multipart")
	require.ErrorContains(t, validateLeoImageParsedRequest(PlatformLeo, &OpenAIImagesRequest{HasMask: true}, []byte(`{}`)), "mask")
}

func jsonStringArray(value gjson.Result) []string {
	items := value.Array()
	out := make([]string, 0, len(items))
	for _, item := range items {
		out = append(out, item.String())
	}
	return out
}
