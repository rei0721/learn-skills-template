package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/rei0721/go-scaffold/pkg/rbac"
)

type rbacRepository struct {
	r rbac.RBAC
}

func NewRBACRepository(r rbac.RBAC) RBACRepository {
	return &rbacRepository{r: r}
}

var errRBACNotInitialized = errors.New("rbac not initialized")

func (repo *rbacRepository) AssignRole(ctx context.Context, user, role, domain string) error {
	if repo.r == nil {
		return errRBACNotInitialized
	}
	if domain == "" {
		return repo.r.AddRoleForUser(user, role)
	}
	return repo.r.AddRoleForUserInDomain(user, role, domain)
}

func (repo *rbacRepository) RevokeRole(ctx context.Context, user, role, domain string) error {
	if repo.r == nil {
		return errRBACNotInitialized
	}
	if domain == "" {
		return repo.r.DeleteRoleForUser(user, role)
	}
	return repo.r.DeleteRoleForUserInDomain(user, role, domain)
}

func (repo *rbacRepository) GetRoles(ctx context.Context, user, domain string) ([]string, error) {
	if repo.r == nil {
		return nil, errRBACNotInitialized
	}
	if domain == "" {
		return repo.r.GetRolesForUser(user)
	}
	return repo.r.GetRolesForUserInDomain(user, domain)
}

func (repo *rbacRepository) AddPolicy(ctx context.Context, sub, domain, obj, act string) error {
	if repo.r == nil {
		return errRBACNotInitialized
	}
	if domain == "" {
		return repo.r.AddPolicy(sub, obj, act)
	}
	return repo.r.AddPolicyWithDomain(sub, domain, obj, act)
}

func (repo *rbacRepository) RemovePolicy(ctx context.Context, sub, domain, obj, act string) error {
	if repo.r == nil {
		return errRBACNotInitialized
	}
	if domain == "" {
		return repo.r.RemovePolicy(sub, obj, act)
	}
	return repo.r.RemovePolicyWithDomain(sub, domain, obj, act)
}

func (repo *rbacRepository) ListPolicies(ctx context.Context, sub, domain, obj, act string) ([][]string, error) {
	if repo.r == nil {
		return nil, errRBACNotInitialized
	}

	if sub == "" && domain == "" && obj == "" && act == "" {
		return repo.r.GetPolicy(), nil
	}

	if domain != "" {
		policies := repo.r.GetFilteredPolicy(0, sub)
		var out [][]string
		for _, p := range policies {
			if len(p) < 4 {
				continue
			}
			if sub != "" && p[0] != sub {
				continue
			}
			if p[1] != domain {
				continue
			}
			if obj != "" && p[2] != obj {
				continue
			}
			if act != "" && p[3] != act {
				continue
			}
			out = append(out, p)
		}
		return out, nil
	}

	filters := make([]string, 0, 3)
	filters = append(filters, sub)
	filters = append(filters, obj)
	filters = append(filters, act)
	for len(filters) > 0 && strings.TrimSpace(filters[len(filters)-1]) == "" {
		filters = filters[:len(filters)-1]
	}
	return repo.r.GetFilteredPolicy(0, filters...), nil
}

func (repo *rbacRepository) Enforce(ctx context.Context, sub, domain, obj, act string) (bool, error) {
	if repo.r == nil {
		return false, errRBACNotInitialized
	}
	if domain == "" {
		return repo.r.Enforce(sub, obj, act)
	}
	return repo.r.EnforceWithDomain(sub, domain, obj, act)
}
