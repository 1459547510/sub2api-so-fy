package handler

import (
	"bytes"
	"errors"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

type VideoInputHandler struct {
	store *service.VideoInputStore
}

func NewVideoInputHandler(store *service.VideoInputStore) *VideoInputHandler {
	return &VideoInputHandler{store: store}
}

func (h *VideoInputHandler) Upload(c *gin.Context) {
	if _, ok := middleware2.GetAPIKeyFromContext(c); !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"message": "Invalid API key"}})
		return
	}
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": "multipart form is invalid"}})
		return
	}
	fieldName := ""
	for _, candidate := range []string{"file", "image", "video", "audio"} {
		if headers := c.Request.MultipartForm.File[candidate]; len(headers) > 0 {
			fieldName = candidate
			break
		}
	}
	if fieldName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": "image, video, or audio file is required"}})
		return
	}
	file, header, err := c.Request.FormFile(fieldName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": "uploaded media file could not be opened"}})
		return
	}
	defer func() { _ = file.Close() }()
	kind := service.VideoInputKindImage
	switch fieldName {
	case "video":
		kind = service.VideoInputKindVideo
	case "audio":
		kind = service.VideoInputKindAudio
	default:
		requestedKind := strings.ToLower(strings.TrimSpace(c.PostForm("media_type")))
		switch requestedKind {
		case string(service.VideoInputKindVideo):
			kind = service.VideoInputKindVideo
		case string(service.VideoInputKindAudio):
			kind = service.VideoInputKindAudio
		case "", string(service.VideoInputKindImage):
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": "media_type must be image, video, or audio"}})
			return
		}
	}
	input, err := h.store.SaveMedia(file, kind, header.Filename)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrVideoInputTooLarge) {
			status = http.StatusRequestEntityTooLarge
		} else if errors.Is(err, service.ErrVideoInputUnsupportedType) || errors.Is(err, service.ErrVideoInputUnsupportedDuration) {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": gin.H{"message": err.Error()}})
		return
	}
	response := gin.H{
		"upload_id":    input.Token,
		"media_url":    input.URL,
		"media_type":   input.Kind,
		"content_type": input.ContentType,
		"size":         input.Size,
	}
	switch input.Kind {
	case service.VideoInputKindImage:
		response["image_url"] = input.URL
	case service.VideoInputKindVideo:
		response["video_url"] = input.URL
	case service.VideoInputKindAudio:
		response["audio_url"] = input.URL
	}
	c.JSON(http.StatusOK, response)
}

func (h *VideoInputHandler) GetInternal(c *gin.Context) {
	remoteHost, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil || net.ParseIP(remoteHost) == nil || !net.ParseIP(remoteHost).IsLoopback() {
		c.Status(http.StatusNotFound)
		return
	}
	input, err := h.store.Open(c.Param("token"))
	if errors.Is(err, service.ErrVideoInputNotFound) {
		c.Status(http.StatusNotFound)
		return
	}
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Header("Cache-Control", "private, max-age=3600")
	c.Data(http.StatusOK, input.ContentType, input.Data)
}

func (h *VideoInputHandler) PrepareLeoMultipartImageRequest(parsed *service.OpenAIImagesRequest) ([]byte, []string, error) {
	if h == nil || h.store == nil {
		return nil, nil, service.ErrVideoInputPublicURLMissing
	}
	imageURLs := make([]string, 0, len(parsed.Uploads))
	tokens := make([]string, 0, len(parsed.Uploads))
	for _, upload := range parsed.Uploads {
		input, err := h.store.SaveMedia(bytes.NewReader(upload.Data), service.VideoInputKindImage, upload.FileName)
		if err != nil {
			return nil, tokens, err
		}
		publicURL, err := h.store.PublicImageURL(input.Token)
		if err != nil {
			return nil, append(tokens, input.Token), err
		}
		tokens = append(tokens, input.Token)
		imageURLs = append(imageURLs, publicURL)
	}
	body, err := service.BuildLeoMultipartImageRequest(parsed, imageURLs)
	return body, tokens, err
}

func (h *VideoInputHandler) MarkTerminal(tokens []string) {
	if h == nil || h.store == nil {
		return
	}
	for _, token := range tokens {
		_ = h.store.MarkTerminal(token, time.Now())
	}
}

func (h *VideoInputHandler) GetPublicImage(c *gin.Context) {
	input, err := h.store.Open(c.Param("token"))
	if errors.Is(err, service.ErrVideoInputNotFound) || err == nil && !strings.HasPrefix(input.ContentType, "image/") {
		c.Status(http.StatusNotFound)
		return
	}
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Header("Cache-Control", "no-store")
	c.Header("X-Content-Type-Options", "nosniff")
	if c.Request.Method == http.MethodHead {
		c.Header("Content-Type", input.ContentType)
		c.Header("Content-Length", strconv.FormatInt(input.Size, 10))
		c.Status(http.StatusOK)
		return
	}
	c.Data(http.StatusOK, input.ContentType, input.Data)
}
