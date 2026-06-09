package services

import (
	"context"
	"crypto/sha256"
	"time"

	"github.com/BrunoPolaski/bureau-mcp-server/internal/core/entities"
	valueobjects "github.com/BrunoPolaski/bureau-mcp-server/internal/core/entities/value_objects"
	"github.com/BrunoPolaski/bureau-mcp-server/internal/infra/controllers/dto"
	"github.com/BrunoPolaski/bureau-mcp-server/internal/infra/repositories"
	"github.com/BrunoPolaski/bureau-mcp-server/internal/infra/repositories/interfaces"
	"github.com/BrunoPolaski/bureau-mcp-server/internal/infra/thirdparty"
	"github.com/BrunoPolaski/bureau-mcp-server/internal/infra/thirdparty/jwt"
	"github.com/BrunoPolaski/go-rest-err/rest_err"
)

type AuthService struct {
	userRepository    interfaces.UserRepository
	sessionRepository interfaces.SessionRepository
	jwtAdapter        jwt.JWT
}

func NewAuthService(rf *repositories.RepositoryFactory, tpf *thirdparty.ThirdPartyFactory) *AuthService {
	return &AuthService{
		userRepository:    rf.UserRepository(),
		sessionRepository: rf.SessionRepository(),
		jwtAdapter:        tpf.JWT(),
	}
}

func (as *AuthService) Register(ctx context.Context, user *dto.UserDTO) (*entities.User, *rest_err.RestErr) {
	domain, err := user.ToEntity()
	if err != nil {
		return nil, err
	}

	if exists, _ := as.userRepository.GetByEmail(ctx, domain.Email.String()); exists != nil {
		return nil, rest_err.NewBadRequestError("email already exists")
	}

	domain, err = as.userRepository.Register(ctx, domain)
	if err != nil {
		return nil, err
	}

	return domain, nil
}

func (as *AuthService) SignIn(ctx context.Context, email, pwd string) (*entities.User, *string, *rest_err.RestErr) {
	password := valueobjects.NewPassword(pwd)

	domainEmail := valueobjects.NewEmail(email)

	user, err := as.userRepository.GetByEmail(ctx, domainEmail.String())
	if err != nil {
		return nil, nil, err
	}

	if !user.Password.CompareHashWithPassword(password.String()) {
		return nil, nil, rest_err.NewUnauthorizedError("invalid credentials")
	}

	token := as.jwtAdapter.GenerateToken()
	hash := sha256.Sum256([]byte(token))
	expiresAt := time.Now().Add(24 * time.Hour)
	createdAt := time.Now()

	tokenStruct := &entities.Token{
		TokenHash:  string(hash[:]),
		CreatedAt:  createdAt,
		ExpiresAt:  &expiresAt,
		LastUsedAt: &createdAt,
	}

	if err := as.sessionRepository.Create(ctx, tokenStruct); err != nil {
		return nil, nil, err
	}

	return user, &token, nil
}
