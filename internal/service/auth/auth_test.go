package auth

import (
	"context"
	stdErrors "errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/rei0721/go-scaffold/internal/models"
	"github.com/rei0721/go-scaffold/pkg/crypto"
	"github.com/rei0721/go-scaffold/pkg/dbtx"
	jwtpkg "github.com/rei0721/go-scaffold/pkg/jwt"
	"github.com/rei0721/go-scaffold/pkg/rbac"
	"github.com/rei0721/go-scaffold/types"
	appErrors "github.com/rei0721/go-scaffold/types/errors"
	"gorm.io/gorm"
)

type stubAuthRepo struct {
	mu         sync.Mutex
	byUsername map[string]*models.DBUser
	byID       map[int64]*models.DBUser
	seq        int64
}

func newStubAuthRepo() *stubAuthRepo {
	return &stubAuthRepo{
		byUsername: map[string]*models.DBUser{},
		byID:       map[int64]*models.DBUser{},
		seq:        0,
	}
}

func (r *stubAuthRepo) FindUserByUsername(ctx context.Context, username string) (*models.DBUser, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if u, ok := r.byUsername[username]; ok {
		cp := *u
		return &cp, nil
	}
	return nil, nil
}

func (r *stubAuthRepo) FindUserByEmail(ctx context.Context, email string) (*models.DBUser, error) {
	return nil, nil
}

func (r *stubAuthRepo) FindUserByID(ctx context.Context, userID int64) (*models.DBUser, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if u, ok := r.byID[userID]; ok {
		cp := *u
		return &cp, nil
	}
	return nil, nil
}

func (r *stubAuthRepo) CreateUser(ctx context.Context, tx *gorm.DB, user *models.DBUser) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.seq++
	user.ID = r.seq
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	cp := *user
	r.byUsername[user.Username] = &cp
	r.byID[user.ID] = &cp
	return nil
}

func (r *stubAuthRepo) UpdateUserPassword(ctx context.Context, tx *gorm.DB, userID int64, hashedPassword string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	u, ok := r.byID[userID]
	if !ok {
		return gorm.ErrRecordNotFound
	}
	u.Password = hashedPassword
	u.UpdatedAt = time.Now()
	r.byUsername[u.Username] = u
	return nil
}

func (r *stubAuthRepo) UpdateUser(ctx context.Context, tx *gorm.DB, user *models.DBUser) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.byID[user.ID]; !ok {
		return gorm.ErrRecordNotFound
	}
	cp := *user
	r.byID[user.ID] = &cp
	r.byUsername[user.Username] = &cp
	return nil
}

type stubTxManager struct{}

func (m stubTxManager) WithTx(ctx context.Context, fn dbtx.TxFunc) error {
	return fn(nil)
}

func (m stubTxManager) WithTxOptions(ctx context.Context, opts *dbtx.TxOptions, fn dbtx.TxFunc) error {
	return fn(nil)
}

func (m stubTxManager) GetDB() *gorm.DB {
	return nil
}

type stubCrypto struct{}

func (c stubCrypto) HashPassword(password string) (string, error) {
	return "hash:" + password, nil
}

func (c stubCrypto) VerifyPassword(hashedPassword, password string) error {
	if hashedPassword != "hash:"+password {
		return stdErrors.New("invalid password")
	}
	return nil
}

func (c stubCrypto) UpdateConfig(opts ...crypto.Option) error {
	return nil
}

type stubJWT struct{}

func (j stubJWT) GenerateToken(userID int64, username string) (string, error) {
	return "token:" + username, nil
}

func (j stubJWT) ValidateToken(tokenString string) (*jwtpkg.Claims, error) {
	return nil, stdErrors.New("not implemented")
}

func (j stubJWT) RefreshToken(tokenString string) (string, error) {
	return "", stdErrors.New("not implemented")
}

type stubRBAC struct {
	mu       sync.Mutex
	assigned []string
}

var _ rbac.RBAC = (*stubRBAC)(nil)

func (s *stubRBAC) Enforce(sub, obj, act string) (bool, error)                { return false, nil }
func (s *stubRBAC) EnforceWithDomain(sub, dom, obj, act string) (bool, error) { return false, nil }

func (s *stubRBAC) AddRoleForUser(user, role string) error {
	return s.AddRoleForUserInDomain(user, role, "")
}
func (s *stubRBAC) AddRoleForUserInDomain(user, role, domain string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.assigned = append(s.assigned, fmt.Sprintf("%s|%s|%s", user, role, domain))
	return nil
}

func (s *stubRBAC) DeleteRoleForUser(user, role string) error                     { return nil }
func (s *stubRBAC) DeleteRoleForUserInDomain(user, role, domain string) error     { return nil }
func (s *stubRBAC) GetRolesForUser(user string) ([]string, error)                 { return nil, nil }
func (s *stubRBAC) GetRolesForUserInDomain(user, domain string) ([]string, error) { return nil, nil }
func (s *stubRBAC) GetUsersForRole(role string) ([]string, error)                 { return nil, nil }

func (s *stubRBAC) AddPolicy(sub, obj, act string) error                               { return nil }
func (s *stubRBAC) AddPolicyWithDomain(sub, domain, obj, act string) error             { return nil }
func (s *stubRBAC) RemovePolicy(sub, obj, act string) error                            { return nil }
func (s *stubRBAC) RemovePolicyWithDomain(sub, domain, obj, act string) error          { return nil }
func (s *stubRBAC) GetPolicy() [][]string                                              { return nil }
func (s *stubRBAC) GetFilteredPolicy(fieldIndex int, fieldValues ...string) [][]string { return nil }
func (s *stubRBAC) AddPolicies(rules [][]string) error                                 { return nil }
func (s *stubRBAC) RemovePolicies(rules [][]string) error                              { return nil }
func (s *stubRBAC) LoadPolicy() error                                                  { return nil }
func (s *stubRBAC) SavePolicy() error                                                  { return nil }
func (s *stubRBAC) ClearCache() error                                                  { return nil }
func (s *stubRBAC) Close() error                                                       { return nil }

func TestAuthService_Register_UsernameDuplicate(t *testing.T) {
	repo := newStubAuthRepo()
	_ = repo.CreateUser(context.Background(), nil, &models.DBUser{Username: "alice", Password: "hash:pw", Status: 1})

	svc := NewAuthService(repo)
	impl := svc.(*authService)
	impl.SetTxManager(stubTxManager{})
	impl.SetCrypto(stubCrypto{})

	_, err := svc.Register(context.Background(), &types.RegisterRequest{Username: "alice", Password: "pw"})
	if err == nil {
		t.Fatalf("expected error")
	}
	var biz *appErrors.BizError
	if !stdErrors.As(err, &biz) || biz.Code != appErrors.ErrDuplicateUsername {
		t.Fatalf("expected duplicate username BizError, got %T %v", err, err)
	}
}

func TestAuthService_Register_Success(t *testing.T) {
	repo := newStubAuthRepo()
	svc := NewAuthService(repo)
	impl := svc.(*authService)
	impl.SetTxManager(stubTxManager{})
	impl.SetCrypto(stubCrypto{})

	resp, err := svc.Register(context.Background(), &types.RegisterRequest{Username: "alice", Password: "pw"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil || resp.UserID == 0 || resp.Username != "alice" {
		t.Fatalf("unexpected response: %#v", resp)
	}
	if resp.Email != nil {
		t.Fatalf("expected nil email, got %v", *resp.Email)
	}
}

func TestAuthService_Register_AssignDefaultRole_WhenRBACEnabled(t *testing.T) {
	repo := newStubAuthRepo()
	svc := NewAuthService(repo)
	impl := svc.(*authService)
	impl.SetTxManager(stubTxManager{})
	impl.SetCrypto(stubCrypto{})

	r := &stubRBAC{}
	impl.SetRBAC(r)

	resp, err := svc.Register(context.Background(), &types.RegisterRequest{Username: "alice", Password: "pw"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.assigned) != 1 {
		t.Fatalf("expected 1 role assignment, got %d", len(r.assigned))
	}
	if r.assigned[0] != fmt.Sprintf("%d|user|", resp.UserID) {
		t.Fatalf("unexpected role assignment: %q", r.assigned[0])
	}
}

func TestAuthService_Login_Success(t *testing.T) {
	repo := newStubAuthRepo()
	_ = repo.CreateUser(context.Background(), nil, &models.DBUser{Username: "alice", Password: "hash:pw", Status: 1})

	svc := NewAuthService(repo)
	impl := svc.(*authService)
	impl.SetCrypto(stubCrypto{})
	impl.SetJWT(stubJWT{})

	resp, err := svc.Login(context.Background(), &types.LoginRequest{Username: "alice", Password: "pw"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil || resp.Token == "" || resp.User.Username != "alice" {
		t.Fatalf("unexpected response: %#v", resp)
	}
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	repo := newStubAuthRepo()
	_ = repo.CreateUser(context.Background(), nil, &models.DBUser{Username: "alice", Password: "hash:pw", Status: 1})

	svc := NewAuthService(repo)
	impl := svc.(*authService)
	impl.SetCrypto(stubCrypto{})
	impl.SetJWT(stubJWT{})

	_, err := svc.Login(context.Background(), &types.LoginRequest{Username: "alice", Password: "wrong"})
	if err == nil {
		t.Fatalf("expected error")
	}
	var biz *appErrors.BizError
	if !stdErrors.As(err, &biz) || biz.Code != appErrors.ErrUnauthorized {
		t.Fatalf("expected unauthorized BizError, got %T %v", err, err)
	}
}
