package controllers

import (
	"net/http"

	"github.com/BrunoPolaski/internal-registry/internal/infra/controllers/dto"
	httphelper "github.com/BrunoPolaski/internal-registry/internal/interfaces/http"
	"github.com/BrunoPolaski/internal-registry/internal/services"
)

type UserController struct {
	UserService *services.UserService
}

func NewUserController(us *services.UserService) *UserController {
	return &UserController{
		UserService: us,
	}
}

// @Summary Get a user by ID
// @Description Get a user by its ID.
// @Tags /user
// @Produce json
// @Security CookieAuth
// @Param id path string true "User ID"
// @Success 200 {object} dto.UserDTO
// @Failure 400 {object} rest_err.RestErr
// @Failure 401 {object} rest_err.RestErr
// @Failure 404 {object} rest_err.RestErr
// @Router /user/{id} [get]
func (uc *UserController) GetById(w http.ResponseWriter, r *http.Request) {
	id, err := httphelper.IntPathParam(r, "id")
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	user, err := uc.UserService.GetById(r.Context(), id)
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	httphelper.JSON(dto.NewUserDTO(user), w)
}

// @Summary List all users
// @Description Get a list of all users.
// @Tags /user
// @Produce json
// @Security CookieAuth
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} dto.PaginatedResponse[dto.UserDTO]
// @Failure 400 {object} rest_err.RestErr
// @Failure 401 {object} rest_err.RestErr
// @Router /user [get]
func (uc *UserController) GetAll(w http.ResponseWriter, r *http.Request) {
	limit, err := httphelper.IntQueryParam(r,
		httphelper.IntQueryParamOptions{
			Key:          "limit",
			DefaultValue: 10,
		},
	)
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	offset, err := httphelper.IntQueryParam(r,
		httphelper.IntQueryParamOptions{
			Key:          "offset",
			DefaultValue: 0,
		},
	)
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	users, err := uc.UserService.GetAll(r.Context(), limit, offset, map[string]any{})
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	httphelper.JSON(users, w)
}

// @Summary Delete a user
// @Description Delete a user by its ID.
// @Tags /user
// @Param id path string true "User ID"
// @Success 204
// @Failure 400 {object} rest_err.RestErr
// @Failure 401 {object} rest_err.RestErr
// @Failure 404 {object} rest_err.RestErr
// @Router /user/{id} [delete]
func (uc *UserController) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := httphelper.IntPathParam(r, "id")
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	err = uc.UserService.Delete(r.Context(), id)
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
