package rbac

import (
	"context"
	"strconv"

	"github.com/rei0721/go-scaffold/internal/repository"
	"github.com/rei0721/go-scaffold/internal/service"
	"github.com/rei0721/go-scaffold/types/errors"
)

type rbacService struct {
	service.BaseService[repository.RBACRepository]
}

func NewRBACService(repo repository.RBACRepository) RBACService {
	s := &rbacService{}
	s.BaseService.SetRepository(repo)
	return s
}

func (s *rbacService) AssignRoleToUser(ctx context.Context, userID int64, role, domain string) error {
	if role == "" {
		return errors.NewBizError(errors.ErrInvalidParams, "role is required")
	}
	return s.Repo.AssignRole(ctx, strconv.FormatInt(userID, 10), role, domain)
}

func (s *rbacService) RevokeRoleFromUser(ctx context.Context, userID int64, role, domain string) error {
	if role == "" {
		return errors.NewBizError(errors.ErrInvalidParams, "role is required")
	}
	return s.Repo.RevokeRole(ctx, strconv.FormatInt(userID, 10), role, domain)
}

func (s *rbacService) GetUserRoles(ctx context.Context, userID int64, domain string) ([]string, error) {
	return s.Repo.GetRoles(ctx, strconv.FormatInt(userID, 10), domain)
}

func (s *rbacService) AddPolicy(ctx context.Context, role, domain, object, action string) error {
	if role == "" || object == "" || action == "" {
		return errors.NewBizError(errors.ErrInvalidParams, "role, object and action are required")
	}
	return s.Repo.AddPolicy(ctx, role, domain, object, action)
}

func (s *rbacService) RemovePolicy(ctx context.Context, role, domain, object, action string) error {
	if role == "" || object == "" || action == "" {
		return errors.NewBizError(errors.ErrInvalidParams, "role, object and action are required")
	}
	return s.Repo.RemovePolicy(ctx, role, domain, object, action)
}

func (s *rbacService) ListPolicies(ctx context.Context, role, domain, object, action string) ([][]string, error) {
	return s.Repo.ListPolicies(ctx, role, domain, object, action)
}

func (s *rbacService) Enforce(ctx context.Context, subject, domain, object, action string) (bool, error) {
	if subject == "" || object == "" || action == "" {
		return false, errors.NewBizError(errors.ErrInvalidParams, "subject, object and action are required")
	}
	return s.Repo.Enforce(ctx, subject, domain, object, action)
}
