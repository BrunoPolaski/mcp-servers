package services

import (
	"context"
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/registration-validation-mcp-server/internal/core/entities"
	valueobjects "github.com/BrunoPolaski/registration-validation-mcp-server/internal/core/entities/value_objects"
	"github.com/BrunoPolaski/registration-validation-mcp-server/internal/helper/sessions"
	"github.com/BrunoPolaski/registration-validation-mcp-server/internal/infra/controllers/dto"
	"github.com/BrunoPolaski/registration-validation-mcp-server/internal/infra/repositories"
	"github.com/BrunoPolaski/registration-validation-mcp-server/internal/infra/repositories/interfaces"
	"github.com/BrunoPolaski/registration-validation-mcp-server/internal/infra/thirdparty"
	"github.com/BrunoPolaski/registration-validation-mcp-server/internal/infra/thirdparty/jwt"
	"github.com/google/uuid"
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

func (as *AuthService) SignIn(ctx context.Context, email, pwd string) (*string, *string, *rest_err.RestErr) {
	password := valueobjects.NewPassword(pwd)

	domainEmail := valueobjects.NewEmail(email)

	user, err := as.userRepository.GetByEmail(ctx, domainEmail.String())
	if err != nil {
		return nil, nil, err
	}

	if !user.Password.CompareHashWithPassword(password.String()) {
		return nil, nil, rest_err.NewUnauthorizedError("invalid credentials")
	}

	tokenID := uuid.New().String()
	expiresAt := time.Now().Add(sessions.AbsoluteTimeout)
	createdAt := time.Now()

	refreshTokenJWT, jwtErr := as.jwtAdapter.GenerateToken(tokenID, strconv.Itoa(int(user.ID)), expiresAt)
	if jwtErr != nil {
		return nil, nil, jwtErr
	}

	hash := sha256.Sum256([]byte(refreshTokenJWT))

	refreshTokenStruct := &entities.Token{
		ID:         tokenID,
		TokenHash:  fmt.Sprintf("%x", hash[:]),
		CreatedAt:  createdAt,
		ExpiresAt:  &expiresAt,
		LastUsedAt: &createdAt,
	}

	if createErr := as.sessionRepository.Create(ctx, refreshTokenStruct); createErr != nil {
		return nil, nil, createErr
	}

	accessToken, accessTokenErr := as.jwtAdapter.GenerateToken(uuid.NewString(), strconv.Itoa(int(user.ID)), time.Now().Add(time.Minute*15))
	if accessTokenErr != nil {
		return nil, nil, accessTokenErr
	}

	return &accessToken, &refreshTokenJWT, nil
}
