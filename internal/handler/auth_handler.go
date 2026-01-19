package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rei0721/go-scaffold/internal/middleware"
	"github.com/rei0721/go-scaffold/internal/service/auth"
	"github.com/rei0721/go-scaffold/pkg/logger"
	"github.com/rei0721/go-scaffold/types"
	appErrors "github.com/rei0721/go-scaffold/types/errors"
	"github.com/rei0721/go-scaffold/types/result"
)

type AuthHandler struct {
	svc    auth.AuthService
	logger logger.Logger
}

func NewAuthHandler(svc auth.AuthService, log logger.Logger) *AuthHandler {
	return &AuthHandler{svc: svc, logger: log}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req types.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.BadRequest(c, "invalid request")
		return
	}

	resp, err := h.svc.Register(c.Request.Context(), &req)
	if err != nil {
		h.writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, result.Success(resp))
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req types.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.BadRequest(c, "invalid request")
		return
	}

	resp, err := h.svc.Login(c.Request.Context(), &req)
	if err != nil {
		h.writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, result.Success(resp))
}

func (h *AuthHandler) writeError(c *gin.Context, err error) {
	traceID := middleware.GetTraceID(c)

	var biz *appErrors.BizError
	if errors.As(err, &biz) {
		httpStatus := http.StatusInternalServerError
		if biz.Code >= 1000 && biz.Code < 2000 {
			httpStatus = http.StatusBadRequest
		} else if biz.Code >= 2000 && biz.Code < 3000 {
			httpStatus = http.StatusUnprocessableEntity
		} else if biz.Code >= 3000 && biz.Code < 4000 {
			httpStatus = http.StatusUnauthorized
		} else if biz.Code >= 4000 && biz.Code < 5000 {
			httpStatus = http.StatusNotFound
		}

		c.JSON(httpStatus, result.ErrorWithTrace(biz.Code, biz.Message, traceID))
		return
	}

	if h.logger != nil {
		h.logger.Error("auth request failed", "error", err, "traceId", traceID)
	}
	c.JSON(http.StatusInternalServerError, result.ErrorWithTrace(appErrors.ErrInternalServer, "internal server error", traceID))
}
