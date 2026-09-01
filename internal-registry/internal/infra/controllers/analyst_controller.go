package controllers

import (
	"net/http"

	"github.com/BrunoPolaski/internal-registry/internal/infra/controllers/dto"
	"github.com/BrunoPolaski/internal-registry/internal/infra/controllers/request"
	"github.com/BrunoPolaski/internal-registry/internal/infra/thirdparty/logger"
	"github.com/BrunoPolaski/internal-registry/internal/infra/thirdparty/validator"
	httphelper "github.com/BrunoPolaski/internal-registry/internal/interfaces/http"
	"github.com/BrunoPolaski/internal-registry/internal/services"
)

type AnalystController struct {
	AnalystService *services.AnalystService
}

func NewAnalystController(as *services.AnalystService) *AnalystController {
	return &AnalystController{
		AnalystService: as,
	}
}

// @Summary Create a new analyst
// @Description Creates a new analyst with the provided details.
// @Tags /analyst
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param personalInformation body request.PersonalInformationRequest true "Personal Information"
// @Success 201 {object} dto.AnalystDTO
// @Failure 400 {object} rest_err.RestErr
// @Failure 401 {object} rest_err.RestErr
// @Router /analyst [post]
func (ac *AnalystController) Create(w http.ResponseWriter, r *http.Request) {
	logger.Info("Creating analyst")

	var analystReq request.PersonalInformationRequest
	if err := validator.ShouldBindJSON(r, &analystReq); err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	createdAnalyst, err := ac.AnalystService.Create(r.Context(), &analystReq)
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	httphelper.JSON(dto.NewAnalystDTO(createdAnalyst), w)
}

// @Summary Get an analyst by ID
// @Description Get an analyst by its ID.
// @Tags /analyst
// @Produce json
// @Security CookieAuth
// @Param id path string true "Analyst ID"
// @Success 200 {object} dto.AnalystDTO
// @Failure 400 {object} rest_err.RestErr
// @Failure 401 {object} rest_err.RestErr
// @Failure 404 {object} rest_err.RestErr
// @Router /analyst/{id} [get]
func (ac *AnalystController) GetById(w http.ResponseWriter, r *http.Request) {
	id, err := httphelper.IntPathParam(r, "id")
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	analyst, err := ac.AnalystService.GetById(r.Context(), id)
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	httphelper.JSON(dto.NewAnalystDTO(analyst), w)
}

// @Summary List all analysts
// @Description Get a list of all analysts.
// @Tags /analyst
// @Produce json
// @Security CookieAuth
// @Param limit query int false "Limit the number of results"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} dto.PaginatedResponse[dto.AnalystDTO]
// @Failure 401 {object} rest_err.RestErr
// @Failure 400 {object} rest_err.RestErr
// @Router /analyst [get]
func (ac *AnalystController) GetAll(w http.ResponseWriter, r *http.Request) {
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

	paginatedAnalystes, err := ac.AnalystService.GetAll(r.Context(), limit, offset, map[string]any{})
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	httphelper.JSON(paginatedAnalystes, w)
}

// @Summary Delete an analyst
// @Description Delete an analyst by its ID.
// @Tags /analyst
// @Security CookieAuth
// @Param id path string true "Analyst ID"
// @Success 204
// @Failure 400 {object} rest_err.RestErr
// @Failure 401 {object} rest_err.RestErr
// @Failure 404 {object} rest_err.RestErr
// @Router /analyst/{id} [delete]
func (ac *AnalystController) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := httphelper.IntPathParam(r, "id")
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	err = ac.AnalystService.Delete(r.Context(), id)
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
