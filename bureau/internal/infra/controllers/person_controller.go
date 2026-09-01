package controllers

import (
	"net/http"

	"github.com/BrunoPolaski/bureau/internal/infra/controllers/dto"
	"github.com/BrunoPolaski/bureau/internal/infra/controllers/request"
	"github.com/BrunoPolaski/bureau/internal/infra/thirdparty/logger"
	"github.com/BrunoPolaski/bureau/internal/infra/thirdparty/validator"
	httphelper "github.com/BrunoPolaski/bureau/internal/interfaces/http"
	"github.com/BrunoPolaski/bureau/internal/services"
)

type PersonController struct {
	PersonService *services.PersonService
}

func NewPersonController(as *services.PersonService) *PersonController {
	return &PersonController{
		PersonService: as,
	}
}

// @Summary Create a new person
// @Description Creates a new person with the provided details.
// @Tags /person
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param personalInformation body request.PersonalInformationRequest true "Personal Information"
// @Success 201 {object} dto.PersonDTO
// @Failure 400 {object} rest_err.RestErr
// @Failure 401 {object} rest_err.RestErr
// @Router /person [post]
func (ac *PersonController) Create(w http.ResponseWriter, r *http.Request) {
	logger.Info("Creating person")

	var addrReq request.PersonalInformationRequest
	if err := validator.ShouldBindJSON(r, &addrReq); err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	createdPerson, err := ac.PersonService.Create(r.Context(), &addrReq)
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	httphelper.JSON(dto.NewPersonDTO(createdPerson), w)
}

// @Summary Get an person by ID
// @Description Get an person by its ID.
// @Tags /person
// @Produce json
// @Security CookieAuth
// @Param id path string true "Person ID"
// @Success 200 {object} dto.PersonDTO
// @Failure 400 {object} rest_err.RestErr
// @Failure 401 {object} rest_err.RestErr
// @Failure 404 {object} rest_err.RestErr
// @Router /person/{id} [get]
func (ac *PersonController) GetById(w http.ResponseWriter, r *http.Request) {
	id, err := httphelper.IntPathParam(r, "id")
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	addr, err := ac.PersonService.GetById(r.Context(), id)
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	httphelper.JSON(dto.NewPersonDTO(addr), w)
}

// @Summary Get an person by document
// @Description Get an person by document.
// @Tags /person
// @Produce json
// @Security CookieAuth
// @Param document path string true "Document"
// @Success 200 {object} dto.PersonDTO
// @Failure 400 {object} rest_err.RestErr
// @Failure 401 {object} rest_err.RestErr
// @Failure 404 {object} rest_err.RestErr
// @Router /person/{document} [get]
func (ac *PersonController) GetByDocument(w http.ResponseWriter, r *http.Request) {
	document, err := httphelper.PathParam(r, "document")
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	person, err := ac.PersonService.GetByDocument(r.Context(), document)
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	httphelper.JSON(dto.NewPersonDTO(person), w)
}

// @Summary List all persons
// @Description Get a list of all persons.
// @Tags /person
// @Produce json
// @Security CookieAuth
// @Param limit query int false "Limit the number of results"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} dto.PaginatedResponse[dto.PersonDTO]
// @Failure 401 {object} rest_err.RestErr
// @Failure 400 {object} rest_err.RestErr
// @Router /person [get]
func (ac *PersonController) GetAll(w http.ResponseWriter, r *http.Request) {
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

	paginatedPersones, err := ac.PersonService.GetAll(r.Context(), limit, offset, map[string]any{})
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	httphelper.JSON(paginatedPersones, w)
}

// @Summary Delete an person
// @Description Delete an person by its ID.
// @Tags /person
// @Security CookieAuth
// @Param id path string true "Person ID"
// @Success 204
// @Failure 400 {object} rest_err.RestErr
// @Failure 401 {object} rest_err.RestErr
// @Failure 404 {object} rest_err.RestErr
// @Router /person/{id} [delete]
func (ac *PersonController) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := httphelper.IntPathParam(r, "id")
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	err = ac.PersonService.Delete(r.Context(), id)
	if err != nil {
		httphelper.ErrorResponse(err, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
