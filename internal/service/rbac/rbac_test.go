package rbac

import (
	"context"
	stdErrors "errors"
	"testing"

	"github.com/rei0721/go-scaffold/internal/repository"
	appErrors "github.com/rei0721/go-scaffold/types/errors"
)

type stubRepo struct {
	assignFn     func(ctx context.Context, user, role, domain string) error
	revokeFn     func(ctx context.Context, user, role, domain string) error
	getRolesFn   func(ctx context.Context, user, domain string) ([]string, error)
	addPolicyFn  func(ctx context.Context, sub, domain, obj, act string) error
	removePolFn  func(ctx context.Context, sub, domain, obj, act string) error
	listPolicyFn func(ctx context.Context, sub, domain, obj, act string) ([][]string, error)
	enforceFn    func(ctx context.Context, sub, domain, obj, act string) (bool, error)
}

func (s stubRepo) AssignRole(ctx context.Context, user, role, domain string) error {
	return s.assignFn(ctx, user, role, domain)
}
func (s stubRepo) RevokeRole(ctx context.Context, user, role, domain string) error {
	return s.revokeFn(ctx, user, role, domain)
}
func (s stubRepo) GetRoles(ctx context.Context, user, domain string) ([]string, error) {
	return s.getRolesFn(ctx, user, domain)
}
func (s stubRepo) AddPolicy(ctx context.Context, sub, domain, obj, act string) error {
	return s.addPolicyFn(ctx, sub, domain, obj, act)
}
func (s stubRepo) RemovePolicy(ctx context.Context, sub, domain, obj, act string) error {
	return s.removePolFn(ctx, sub, domain, obj, act)
}
func (s stubRepo) ListPolicies(ctx context.Context, sub, domain, obj, act string) ([][]string, error) {
	return s.listPolicyFn(ctx, sub, domain, obj, act)
}
func (s stubRepo) Enforce(ctx context.Context, sub, domain, obj, act string) (bool, error) {
	return s.enforceFn(ctx, sub, domain, obj, act)
}

func TestRBACService_AddPolicy_Validate(t *testing.T) {
	var _ repository.RBACRepository = stubRepo{}
	svc := NewRBACService(stubRepo{})

	if err := svc.AddPolicy(context.Background(), "", "", "users", "read"); err == nil {
		t.Fatalf("expected error")
	} else {
		var biz *appErrors.BizError
		if !stdErrors.As(err, &biz) || biz.Code != appErrors.ErrInvalidParams {
			t.Fatalf("expected BizError invalid params, got %T %v", err, err)
		}
	}
}

func TestRBACService_Enforce_Validate(t *testing.T) {
	svc := NewRBACService(stubRepo{})

	_, err := svc.Enforce(context.Background(), "", "", "users", "read")
	if err == nil {
		t.Fatalf("expected error")
	}
}
