package rbac

import (
	"context"

	"github.com/rei0721/go-scaffold/internal/service"
)

type RBACService interface {
	service.Service

	AssignRoleToUser(ctx context.Context, userID int64, role, domain string) error
	RevokeRoleFromUser(ctx context.Context, userID int64, role, domain string) error
	GetUserRoles(ctx context.Context, userID int64, domain string) ([]string, error)

	AddPolicy(ctx context.Context, role, domain, object, action string) error
	RemovePolicy(ctx context.Context, role, domain, object, action string) error
	ListPolicies(ctx context.Context, role, domain, object, action string) ([][]string, error)

	Enforce(ctx context.Context, subject, domain, object, action string) (bool, error)
}
