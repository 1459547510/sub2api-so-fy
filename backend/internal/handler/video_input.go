package handler

import (
	"errors"
	"net"
	"net/http"

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
	file, _, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": "image file is required"}})
		return
	}
	defer func() { _ = file.Close() }()
	input, err := h.store.Save(file)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrVideoInputTooLarge) {
			status = http.StatusRequestEntityTooLarge
		} else if errors.Is(err, service.ErrVideoInputUnsupportedType) {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"upload_id": input.Token, "image_url": input.URL, "content_type": input.ContentType, "size": input.Size})
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
