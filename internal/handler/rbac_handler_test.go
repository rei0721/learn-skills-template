package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rei0721/go-scaffold/internal/service"
	"github.com/rei0721/go-scaffold/internal/service/rbac"
)

// mockRBACService 模拟 RBACService
type mockRBACService struct {
	service.Service
	assignRoleFunc   func(ctx context.Context, userID int64, role, domain string) error
	revokeRoleFunc   func(ctx context.Context, userID int64, role, domain string) error
	getUserRolesFunc func(ctx context.Context, userID int64, domain string) ([]string, error)
	addPolicyFunc    func(ctx context.Context, role, domain, object, action string) error
	removePolicyFunc func(ctx context.Context, role, domain, object, action string) error
	listPoliciesFunc func(ctx context.Context, role, domain, object, action string) ([][]string, error)
	enforceFunc      func(ctx context.Context, subject, domain, object, action string) (bool, error)
}

func (m *mockRBACService) AssignRoleToUser(ctx context.Context, userID int64, role, domain string) error {
	if m.assignRoleFunc != nil {
		return m.assignRoleFunc(ctx, userID, role, domain)
	}
	return nil
}

func (m *mockRBACService) RevokeRoleFromUser(ctx context.Context, userID int64, role, domain string) error {
	if m.revokeRoleFunc != nil {
		return m.revokeRoleFunc(ctx, userID, role, domain)
	}
	return nil
}

func (m *mockRBACService) GetUserRoles(ctx context.Context, userID int64, domain string) ([]string, error) {
	if m.getUserRolesFunc != nil {
		return m.getUserRolesFunc(ctx, userID, domain)
	}
	return nil, nil
}

func (m *mockRBACService) AddPolicy(ctx context.Context, role, domain, object, action string) error {
	if m.addPolicyFunc != nil {
		return m.addPolicyFunc(ctx, role, domain, object, action)
	}
	return nil
}

func (m *mockRBACService) RemovePolicy(ctx context.Context, role, domain, object, action string) error {
	if m.removePolicyFunc != nil {
		return m.removePolicyFunc(ctx, role, domain, object, action)
	}
	return nil
}

func (m *mockRBACService) ListPolicies(ctx context.Context, role, domain, object, action string) ([][]string, error) {
	if m.listPoliciesFunc != nil {
		return m.listPoliciesFunc(ctx, role, domain, object, action)
	}
	return nil, nil
}

func (m *mockRBACService) Enforce(ctx context.Context, subject, domain, object, action string) (bool, error) {
	if m.enforceFunc != nil {
		return m.enforceFunc(ctx, subject, domain, object, action)
	}
	return true, nil
}

// setupRouter 设置测试路由
func setupRouter(svc rbac.RBACService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewRBACHandler(svc, nil)

	// 注册路由
	rbacGroup := r.Group("/rbac")
	{
		rbacGroup.POST("/users/:userId/roles", h.AssignRoleToUser)
		rbacGroup.DELETE("/users/:userId/roles", h.RevokeRoleFromUser)
		rbacGroup.GET("/users/:userId/roles", h.GetUserRoles)
		rbacGroup.POST("/roles/:role/policies", h.AddPolicyToRole)
		rbacGroup.DELETE("/roles/:role/policies", h.RemovePolicyFromRole)
		rbacGroup.GET("/policies", h.ListPolicies)
		rbacGroup.POST("/enforce", h.Enforce)
	}

	return r
}

func TestRBACHandler_AssignRoleToUser(t *testing.T) {
	tests := []struct {
		name       string
		userID     string
		body       map[string]interface{}
		mockFunc   func(ctx context.Context, userID int64, role, domain string) error
		wantStatus int
	}{
		{
			name:   "success",
			userID: "1",
			body: map[string]interface{}{
				"role":   "admin",
				"domain": "domain1",
			},
			mockFunc: func(ctx context.Context, userID int64, role, domain string) error {
				if userID != 1 || role != "admin" || domain != "domain1" {
					return errors.New("invalid params")
				}
				return nil
			},
			wantStatus: http.StatusOK,
		},
		{
			name:   "invalid user id",
			userID: "abc",
			body: map[string]interface{}{
				"role": "admin",
			},
			mockFunc:   nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:   "invalid request body",
			userID: "1",
			body: map[string]interface{}{
				"role": "", // required
			},
			mockFunc:   nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:   "service error",
			userID: "1",
			body: map[string]interface{}{
				"role": "admin",
			},
			mockFunc: func(ctx context.Context, userID int64, role, domain string) error {
				return errors.New("internal error")
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockRBACService{assignRoleFunc: tt.mockFunc}
			r := setupRouter(svc)

			bodyBytes, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest(http.MethodPost, "/rbac/users/"+tt.userID+"/roles", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("AssignRoleToUser() status = %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestRBACHandler_GetUserRoles(t *testing.T) {
	tests := []struct {
		name       string
		userID     string
		mockFunc   func(ctx context.Context, userID int64, domain string) ([]string, error)
		wantStatus int
		wantRoles  []string
	}{
		{
			name:   "success",
			userID: "1",
			mockFunc: func(ctx context.Context, userID int64, domain string) ([]string, error) {
				return []string{"admin", "editor"}, nil
			},
			wantStatus: http.StatusOK,
			wantRoles:  []string{"admin", "editor"},
		},
		{
			name:       "invalid user id",
			userID:     "abc",
			mockFunc:   nil,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockRBACService{getUserRolesFunc: tt.mockFunc}
			r := setupRouter(svc)

			req, _ := http.NewRequest(http.MethodGet, "/rbac/users/"+tt.userID+"/roles", nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("GetUserRoles() status = %v, want %v", w.Code, tt.wantStatus)
			}

			if tt.wantStatus == http.StatusOK {
				var resp map[string]interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				data := resp["data"].(map[string]interface{})
				roles := data["roles"].([]interface{})
				if len(roles) != len(tt.wantRoles) {
					t.Errorf("GetUserRoles() roles count = %v, want %v", len(roles), len(tt.wantRoles))
				}
			}
		})
	}
}

func TestRBACHandler_Enforce(t *testing.T) {
	tests := []struct {
		name       string
		body       map[string]interface{}
		mockFunc   func(ctx context.Context, subject, domain, object, action string) (bool, error)
		wantStatus int
		wantAllow  bool
	}{
		{
			name: "allow",
			body: map[string]interface{}{
				"subject": "alice",
				"object":  "data",
				"action":  "read",
			},
			mockFunc: func(ctx context.Context, subject, domain, object, action string) (bool, error) {
				return true, nil
			},
			wantStatus: http.StatusOK,
			wantAllow:  true,
		},
		{
			name: "deny",
			body: map[string]interface{}{
				"subject": "bob",
				"object":  "data",
				"action":  "write",
			},
			mockFunc: func(ctx context.Context, subject, domain, object, action string) (bool, error) {
				return false, nil
			},
			wantStatus: http.StatusOK,
			wantAllow:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockRBACService{enforceFunc: tt.mockFunc}
			r := setupRouter(svc)

			bodyBytes, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest(http.MethodPost, "/rbac/enforce", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Enforce() status = %v, want %v", w.Code, tt.wantStatus)
			}

			if tt.wantStatus == http.StatusOK {
				var resp map[string]interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				data := resp["data"].(map[string]interface{})
				allowed := data["allowed"].(bool)
				if allowed != tt.wantAllow {
					t.Errorf("Enforce() allowed = %v, want %v", allowed, tt.wantAllow)
				}
			}
		})
	}
}
