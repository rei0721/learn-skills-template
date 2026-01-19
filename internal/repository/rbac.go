package repository

import "context"

type RBACRepository interface {
	AssignRole(ctx context.Context, user, role, domain string) error
	RevokeRole(ctx context.Context, user, role, domain string) error
	GetRoles(ctx context.Context, user, domain string) ([]string, error)

	AddPolicy(ctx context.Context, sub, domain, obj, act string) error
	RemovePolicy(ctx context.Context, sub, domain, obj, act string) error
	ListPolicies(ctx context.Context, sub, domain, obj, act string) ([][]string, error)

	Enforce(ctx context.Context, sub, domain, obj, act string) (bool, error)
}
