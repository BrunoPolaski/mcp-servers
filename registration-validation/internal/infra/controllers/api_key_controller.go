package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/BrunoPolaski/registration-validation/internal/infra/controllers/dto"
	httphelper "github.com/BrunoPolaski/registration-validation/internal/interfaces/http"
	"github.com/BrunoPolaski/registration-validation/internal/services"
)

type ApiKeyController struct {
	ApiKeyService *services.ApiKeyService
}

func NewApiKeyController(us *services.ApiKeyService) *ApiKeyController {
	return &ApiKeyController{
		ApiKeyService: us,
	}
}

// @Summary Create a new apiKey
// @Description Creates a new apiKey with the provided details.
// @Tags /api-key
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param apiKey body dto.ApiKeyDTO true "ApiKey"
// @Success 201 {object} dto.ApiKeyDTO
// @Failure 400 {object} rest_err.RestErr
// @Failure 401 {object} rest_err.RestErr
// @Router /api-key [post]
func (uc *ApiKeyController) Create(w http.ResponseWriter, r *http.Request) {
	usr := dto.ApiKeyDTO{}
	json.NewDecoder(r.Body).Decode(&usr)

	createdApiKey, err := uc.ApiKeyService.Create(r.Context(), &usr)
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}
	httphelper.JSONWithStatus(
		dto.NewApiKeyDTO(createdApiKey),
		http.StatusCreated,
		w,
	)
}

// @Summary Get an apiKey by UUID
// @Description Get an apiKey by its UUID.
// @Tags /api-key
// @Produce json
// @Security CookieAuth
// @Param uuid path string true "ApiKey UUID"
// @Success 200 {object} dto.ApiKeyDTO
// @Failure 400 {object} rest_err.RestErr
// @Failure 401 {object} rest_err.RestErr
// @Failure 404 {object} rest_err.RestErr
// @Router /api-key/{uuid} [get]
func (uc *ApiKeyController) GetById(w http.ResponseWriter, r *http.Request) {
	uuid, err := httphelper.PathParam(r, "uuid")
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	apiKey, err := uc.ApiKeyService.GetById(r.Context(), uuid)
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	httphelper.JSON(dto.NewApiKeyDTO(apiKey), w)
}

// @Summary List all apiKeys
// @Description Get a list of all apiKeys.
// @Tags /api-key
// @Produce json
// @Security CookieAuth
// @Success 200 {array} dto.PaginatedResponse[dto.ApiKeyDTO]
// @Failure 400 {object} rest_err.RestErr
// @Failure 401 {object} rest_err.RestErr
// @Failure 404 {object} rest_err.RestErr
// @Router /api-key [get]
func (ac *ApiKeyController) GetAll(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	limitStr := params.Get("limit")
	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	offsetStr := params.Get("offset")
	offset := 0
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	paginatedApiKeys, err := ac.ApiKeyService.GetAll(r.Context(), limit, offset, map[string]any{})
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	httphelper.JSON(paginatedApiKeys, w)
}

// @Summary Delete a apiKey
// @Description Delete a apiKey by its UUID.
// @Tags /api-key
// @Produce json
// @Security CookieAuth
// @Param uuid path string true "ApiKey UUID"
// @Success 204
// @Failure 400 {object} rest_err.RestErr
// @Failure 401 {object} rest_err.RestErr
// @Failure 404 {object} rest_err.RestErr
// @Router /api-key/{uuid} [delete]
func (uc *ApiKeyController) Delete(w http.ResponseWriter, r *http.Request) {
	uuid, err := httphelper.PathParam(r, "uuid")
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	err = uc.ApiKeyService.Delete(r.Context(), uuid)
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
