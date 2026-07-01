package controllers

import (
	"net/http"

	"github.com/BrunoPolaski/internal-registry-mcp-server/internal/infra/controllers/dto"
	"github.com/BrunoPolaski/internal-registry-mcp-server/internal/infra/thirdparty/logger"
	"github.com/BrunoPolaski/internal-registry-mcp-server/internal/infra/thirdparty/validator"
	httphelper "github.com/BrunoPolaski/internal-registry-mcp-server/internal/interfaces/http"
	"github.com/BrunoPolaski/internal-registry-mcp-server/internal/services"
)

type AddressController struct {
	AddressService *services.AddressService
}

func NewAddressController(as *services.AddressService) *AddressController {
	return &AddressController{
		AddressService: as,
	}
}

// @Summary Create a new address
// @Description Creates a new address with the provided details.
// @Tags /address
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param address body dto.AddressDTO true "Address"
// @Success 201 {object} dto.AddressDTO
// @Failure 400 {object} rest_err.RestErr
// @Failure 401 {object} rest_err.RestErr
// @Router /address [post]
func (ac *AddressController) Create(w http.ResponseWriter, r *http.Request) {
	logger.Info("Creating address")

	var addrReq dto.AddressDTO
	if err := validator.ShouldBindJSON(r, &addrReq); err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	createdAddress, err := ac.AddressService.Create(r.Context(), &addrReq)
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	httphelper.JSON(dto.NewAddressDTO(createdAddress), w)
}

// @Summary Get an address by ID
// @Description Get an address by its ID.
// @Tags /address
// @Produce json
// @Security CookieAuth
// @Param id path string true "Address ID"
// @Success 200 {object} dto.AddressDTO
// @Failure 400 {object} rest_err.RestErr
// @Failure 401 {object} rest_err.RestErr
// @Failure 404 {object} rest_err.RestErr
// @Router /address/{id} [get]
func (ac *AddressController) GetById(w http.ResponseWriter, r *http.Request) {
	id, err := httphelper.IntPathParam(r, "id")
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	addr, err := ac.AddressService.GetById(r.Context(), id)
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	httphelper.JSON(dto.NewAddressDTO(addr), w)
}

// @Summary List all addresses
// @Description Get a list of all addresses.
// @Tags /address
// @Produce json
// @Security CookieAuth
// @Param limit query int false "Limit the number of results"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} dto.PaginatedResponse[dto.AddressDTO]
// @Failure 401 {object} rest_err.RestErr
// @Failure 400 {object} rest_err.RestErr
// @Router /address [get]
func (ac *AddressController) GetAll(w http.ResponseWriter, r *http.Request) {
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

	paginatedAddresses, err := ac.AddressService.GetAll(r.Context(), limit, offset, map[string]any{})
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	httphelper.JSON(paginatedAddresses, w)
}

// @Summary Delete an address
// @Description Delete an address by its ID.
// @Tags /address
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param id path string true "Address ID"
// @Success 204
// @Failure 400 {object} rest_err.RestErr
// @Failure 401 {object} rest_err.RestErr
// @Failure 404 {object} rest_err.RestErr
// @Router /address [delete]
func (ac *AddressController) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := httphelper.IntPathParam(r, "id")
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	err = ac.AddressService.Delete(r.Context(), id)
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
