package controllers

import (
	"net/http"

	"github.com/BrunoPolaski/bureau-mcp-server/internal/helper/sessions"
	"github.com/BrunoPolaski/bureau-mcp-server/internal/infra/controllers/dto"
	"github.com/BrunoPolaski/bureau-mcp-server/internal/infra/thirdparty/logger"
	"github.com/BrunoPolaski/bureau-mcp-server/internal/infra/thirdparty/validator"
	httphelper "github.com/BrunoPolaski/bureau-mcp-server/internal/interfaces/http"
	"github.com/BrunoPolaski/bureau-mcp-server/internal/services"
	"github.com/BrunoPolaski/go-rest-err/rest_err"
)

type AuthController struct {
	AuthService *services.AuthService
}

func NewAuthController(as *services.AuthService) *AuthController {
	return &AuthController{
		AuthService: as,
	}
}

// @Summary Create a new user
// @Description Creates a new user with the provided details.
// @Tags /auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param user body dto.UserDTO true "User to create"
// @Success 201 {object} dto.UserDTO
// @Failure 400 {object} rest_err.RestErr
// @Failure 401 {object} rest_err.RestErr
// @Router /auth/register [post]
func (ac *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	logger.Info("Creating user")

	var userReq dto.UserDTO
	if err := validator.ShouldBindJSON(r, &userReq); err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	createdUser, err := ac.AuthService.Register(r.Context(), &userReq)
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	httphelper.JSONWithStatus(dto.NewUserDTO(createdUser), http.StatusCreated, w)
}

// @Summary Sign in
// @Description Signs in with the provided credentials
// @Tags /auth
// @Security BasicAuth && ApiKeyAuth
// @Header 200 {string} Set-Cookie "sid=UUID; Path=/; HttpOnly Secure"
// @Success 200
// @Failure 400 {object} rest_err.RestErr
// @Failure 401 {object} rest_err.RestErr
// @Failure 404 {object} rest_err.RestErr
// @Router /auth/signin [post]
func (ac *AuthController) SignIn(w http.ResponseWriter, r *http.Request) {
	logger.Info("Signing in")

	email, pwd, ok := r.BasicAuth()
	if !ok {
		httphelper.ErrorResponse(rest_err.NewBadRequestError("invalid auth header"), w)
		return
	}

	usr, session, err := ac.AuthService.SignIn(r.Context(), email, pwd)
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	// Set secure session cookie
	cookie := &http.Cookie{
		Name:     sessions.CookieName,
		Value:    session.UUID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,                    // Only send over HTTPS
		SameSite: http.SameSiteStrictMode, // CSRF protection
	}
	http.SetCookie(w, cookie)
	w.Header().Set("sid", session.UUID)

	httphelper.JSON(dto.NewUserDTO(usr), w)
}
