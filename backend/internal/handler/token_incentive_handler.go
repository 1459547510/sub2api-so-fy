package handler

import (
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// TokenIncentiveHandler handles user self-claim weekly token incentive rewards.
type TokenIncentiveHandler struct {
	service *service.TokenIncentiveService
}

func NewTokenIncentiveHandler(service *service.TokenIncentiveService) *TokenIncentiveHandler {
	return &TokenIncentiveHandler{service: service}
}

// GetStatus returns current-week token incentive progress.
func (h *TokenIncentiveHandler) GetStatus(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	status, err := h.service.GetStatus(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, status)
}

// Claim credits the current user once for the current week if eligible.
func (h *TokenIncentiveHandler) Claim(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	status, err := h.service.Claim(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, status)
}
