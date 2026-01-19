package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rei0721/go-scaffold/internal/middleware"
	"github.com/rei0721/go-scaffold/internal/service/rbac"
	"github.com/rei0721/go-scaffold/pkg/logger"
	appErrors "github.com/rei0721/go-scaffold/types/errors"
	"github.com/rei0721/go-scaffold/types/result"
)

type RBACHandler struct {
	svc    rbac.RBACService
	logger logger.Logger
}

func NewRBACHandler(svc rbac.RBACService, log logger.Logger) *RBACHandler {
	return &RBACHandler{svc: svc, logger: log}
}

type assignRoleRequest struct {
	Role   string `json:"role" binding:"required,min=1,max=100"`
	Domain string `json:"domain" binding:"omitempty,max=100"`
}

type policyRequest struct {
	Object string `json:"object" binding:"required,min=1,max=200"`
	Action string `json:"action" binding:"required,min=1,max=100"`
	Domain string `json:"domain" binding:"omitempty,max=100"`
}

type enforceRequest struct {
	Subject string `json:"subject" binding:"required,min=1,max=200"`
	Object  string `json:"object" binding:"required,min=1,max=200"`
	Action  string `json:"action" binding:"required,min=1,max=100"`
	Domain  string `json:"domain" binding:"omitempty,max=100"`
}

func (h *RBACHandler) AssignRoleToUser(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil || userID <= 0 {
		result.BadRequest(c, "invalid userId")
		return
	}

	var req assignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.BadRequest(c, "invalid request")
		return
	}

	if err := h.svc.AssignRoleToUser(c.Request.Context(), userID, req.Role, req.Domain); err != nil {
		h.writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, result.Success(gin.H{"assigned": true}))
}

func (h *RBACHandler) RevokeRoleFromUser(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil || userID <= 0 {
		result.BadRequest(c, "invalid userId")
		return
	}

	var req assignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.BadRequest(c, "invalid request")
		return
	}

	if err := h.svc.RevokeRoleFromUser(c.Request.Context(), userID, req.Role, req.Domain); err != nil {
		h.writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, result.Success(gin.H{"revoked": true}))
}

func (h *RBACHandler) GetUserRoles(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil || userID <= 0 {
		result.BadRequest(c, "invalid userId")
		return
	}

	domain := c.Query("domain")
	roles, err := h.svc.GetUserRoles(c.Request.Context(), userID, domain)
	if err != nil {
		h.writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, result.Success(gin.H{"roles": roles}))
}

func (h *RBACHandler) AddPolicyToRole(c *gin.Context) {
	role := c.Param("role")
	if role == "" {
		result.BadRequest(c, "invalid role")
		return
	}

	var req policyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.BadRequest(c, "invalid request")
		return
	}

	if err := h.svc.AddPolicy(c.Request.Context(), role, req.Domain, req.Object, req.Action); err != nil {
		h.writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, result.Success(gin.H{"added": true}))
}

func (h *RBACHandler) RemovePolicyFromRole(c *gin.Context) {
	role := c.Param("role")
	if role == "" {
		result.BadRequest(c, "invalid role")
		return
	}

	domain := c.Query("domain")
	object := c.Query("object")
	action := c.Query("action")
	if object == "" || action == "" {
		result.BadRequest(c, "object and action are required")
		return
	}

	if err := h.svc.RemovePolicy(c.Request.Context(), role, domain, object, action); err != nil {
		h.writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, result.Success(gin.H{"removed": true}))
}

func (h *RBACHandler) ListPolicies(c *gin.Context) {
	role := c.Query("role")
	domain := c.Query("domain")
	object := c.Query("object")
	action := c.Query("action")

	policies, err := h.svc.ListPolicies(c.Request.Context(), role, domain, object, action)
	if err != nil {
		h.writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, result.Success(gin.H{"policies": policies}))
}

func (h *RBACHandler) Enforce(c *gin.Context) {
	var req enforceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.BadRequest(c, "invalid request")
		return
	}

	ok, err := h.svc.Enforce(c.Request.Context(), req.Subject, req.Domain, req.Object, req.Action)
	if err != nil {
		h.writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, result.Success(gin.H{"allowed": ok}))
}

func (h *RBACHandler) writeError(c *gin.Context, err error) {
	traceID := middleware.GetTraceID(c)

	var biz *appErrors.BizError
	if errors.As(err, &biz) {
		httpStatus := http.StatusInternalServerError
		if biz.Code >= 1000 && biz.Code < 2000 {
			httpStatus = http.StatusBadRequest
		} else if biz.Code >= 2000 && biz.Code < 3000 {
			httpStatus = http.StatusUnprocessableEntity
		} else if biz.Code >= 3000 && biz.Code < 4000 {
			if biz.Code == appErrors.ErrUnauthorized || biz.Code == appErrors.ErrInvalidToken || biz.Code == appErrors.ErrTokenExpired {
				httpStatus = http.StatusUnauthorized
			} else {
				httpStatus = http.StatusForbidden
			}
		} else if biz.Code >= 4000 && biz.Code < 5000 {
			httpStatus = http.StatusNotFound
		}

		c.JSON(httpStatus, result.ErrorWithTrace(biz.Code, biz.Message, traceID))
		return
	}

	if h.logger != nil {
		h.logger.Error("rbac request failed", "error", err, "traceId", traceID)
	}
	c.JSON(http.StatusInternalServerError, result.ErrorWithTrace(appErrors.ErrInternalServer, "internal server error", traceID))
}
