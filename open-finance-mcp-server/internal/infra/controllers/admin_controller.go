package controllers

import (
	"net/http"

	"github.com/BrunoPolaski/open-finance-mcp-server/internal/infra/controllers/dto"
	"github.com/BrunoPolaski/open-finance-mcp-server/internal/infra/controllers/request"
	"github.com/BrunoPolaski/open-finance-mcp-server/internal/infra/thirdparty/logger"
	"github.com/BrunoPolaski/open-finance-mcp-server/internal/infra/thirdparty/validator"
	httphelper "github.com/BrunoPolaski/open-finance-mcp-server/internal/interfaces/http"
	"github.com/BrunoPolaski/open-finance-mcp-server/internal/services"
)

type AdminController struct {
	AdminService *services.AdminService
}

func NewAdminController(as *services.AdminService) *AdminController {
	return &AdminController{
		AdminService: as,
	}
}

// @Summary Create a new admin
// @Description Creates a new admin with the provided details.
// @Tags /admin
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param personalInformation body request.PersonalInformationRequest true "Personal Information"
// @Success 201 {object} dto.AdminDTO
// @Failure 400 {object} rest_err.RestErr
// @Failure 401 {object} rest_err.RestErr
// @Router /admin [post]
func (ac *AdminController) Create(w http.ResponseWriter, r *http.Request) {
	logger.Info("Creating admin")

	var adminReq request.PersonalInformationRequest
	if err := validator.ShouldBindJSON(r, &adminReq); err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	createdAdmin, err := ac.AdminService.Create(r.Context(), &adminReq)
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	httphelper.JSON(dto.NewAdminDTO(createdAdmin), w)
}

// @Summary Get an admin by ID
// @Description Get an admin by its ID.
// @Tags /admin
// @Produce json
// @Security CookieAuth
// @Param id path string true "Admin ID"
// @Success 200 {object} dto.AdminDTO
// @Failure 400 {object} rest_err.RestErr
// @Failure 401 {object} rest_err.RestErr
// @Failure 404 {object} rest_err.RestErr
// @Router /admin/{id} [get]
func (ac *AdminController) GetById(w http.ResponseWriter, r *http.Request) {
	id, err := httphelper.IntPathParam(r, "id")
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	addr, err := ac.AdminService.GetById(r.Context(), id)
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	httphelper.JSON(dto.NewAdminDTO(addr), w)
}

// @Summary List all admins
// @Description Get a list of all admins.
// @Tags /admin
// @Produce json
// @Security CookieAuth
// @Param limit query int false "Limit the number of results"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} dto.PaginatedResponse[dto.AdminDTO]
// @Failure 401 {object} rest_err.RestErr
// @Failure 400 {object} rest_err.RestErr
// @Router /admin [get]
func (ac *AdminController) GetAll(w http.ResponseWriter, r *http.Request) {
	limit, err := httphelper.IntQueryParam(
		r,
		httphelper.IntQueryParamOptions{
			Key:          "limit",
			DefaultValue: 10,
		},
	)
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	offset, err := httphelper.IntQueryParam(
		r,
		httphelper.IntQueryParamOptions{
			Key:          "offset",
			DefaultValue: 0,
		},
	)
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	paginatedAdmines, err := ac.AdminService.GetAll(r.Context(), limit, offset, map[string]any{})
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	httphelper.JSON(paginatedAdmines, w)
}

// @Summary Delete an admin
// @Description Delete an admin by its ID.
// @Tags /admin
// @Security CookieAuth
// @Param id path string true "Admin ID"
// @Success 204
// @Failure 400 {object} rest_err.RestErr
// @Failure 401 {object} rest_err.RestErr
// @Failure 404 {object} rest_err.RestErr
// @Router /admin [delete]
func (ac *AdminController) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := httphelper.IntPathParam(r, "id")
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	err = ac.AdminService.Delete(r.Context(), id)
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
