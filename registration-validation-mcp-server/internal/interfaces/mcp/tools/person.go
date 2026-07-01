package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/BrunoPolaski/registration-validation-mcp-server/internal/infra/controllers/dto"
	"github.com/mark3labs/mcp-go/mcp"
)

func (s *Server) GetPersonByIDTool() mcp.Tool {
	return mcp.NewTool(
		"get_person_by_id",
		mcp.WithDescription(
			`
			Get a person by their ID
			This returns a person's registration validation data, including their name, date of birth, and associated addresses.
			Example usage:
			{
				"id": 123
			}
			`,
		),
		mcp.WithOutputSchema[dto.PersonDTO](),
		mcp.WithInteger(
			"id",
			mcp.Description("The ID of the person to retrieve"),
		),
	)
}

func (s *Server) HandleGetPersonByID(ctx context.Context, request mcp.CallToolRequest, args mcp.CallToolParams) (*dto.PersonDTO, error) {
	id, err := request.RequireInt("id")
	if err != nil {
		return nil, err
	} else if id <= 0 {
		return nil, fmt.Errorf("invalid ID: must be a positive integer")
	}

	person, restErr := s.personService.GetById(ctx, uint(id))
	if restErr != nil {
		return nil, restErr
	}

	return dto.NewPersonDTO(person), nil
}

func (s *Server) GetPersonByDocumentTool() mcp.Tool {
	return mcp.NewTool(
		"get_person_by_document",
		mcp.WithDescription(
			`
			Get a person by their document number
			This returns a person's registration validation data, including their name, date of birth, and associated addresses.
			Example usage:
			{
				"document": "12345678900"
			}
			`,
		),
		mcp.WithOutputSchema[dto.PersonDTO](),
		mcp.WithString(
			"document",
			mcp.Description("The document number of the person to retrieve"),
		),
	)
}

func (s *Server) HandleGetPersonByDocument(ctx context.Context, request mcp.CallToolRequest, args mcp.CallToolParams) (*dto.PersonDTO, error) {
	document, err := request.RequireString("document")
	if err != nil {
		return nil, err
	} else if document == "" {
		return nil, fmt.Errorf("invalid document: cannot be empty")
	}

	person, restErr := s.personService.GetByDocument(ctx, document)
	if restErr != nil {
		return nil, restErr
	}

	return dto.NewPersonDTO(person), nil
}

func (s *Server) GetAllPersonsTool() mcp.Tool {
	return mcp.NewTool(
		"get_all_persons",
		mcp.WithDescription(
			`
			List persons with pagination
			Example usage:
			{
				"limit": 10,
				"offset": 0
			}

			To fetch all persons, set "limit" to 0.
			`,
		),
		mcp.WithOutputSchema[dto.PaginatedResponse[dto.PersonDTO]](),
		mcp.WithInteger(
			"limit",
			mcp.Description("Limit the number of results"),
			mcp.DefaultNumber(10),
			mcp.Min(0),
		),
		mcp.WithInteger(
			"offset",
			mcp.Description("Offset for pagination"),
			mcp.DefaultNumber(0),
			mcp.Min(0),
		),
		mcp.WithObject(
			"params",
			mcp.Description("Additional filters for querying persons. e.g. {\"name\": \"John Doe\"}"),
		),
	)
}

func (s *Server) HandleGetAllPersons(ctx context.Context, request mcp.CallToolRequest, args mcp.CallToolParams) (*dto.PaginatedResponse[dto.PersonDTO], error) {
	arguments := map[string]any{}
	if request.Params.Arguments != nil {
		var ok bool
		arguments, ok = request.Params.Arguments.(map[string]any)
		if !ok {
			return nil, errors.New("invalid arguments")
		}
	}

	params, _ := arguments["params"].(map[string]any)
	limit := request.GetInt("limit", 10)
	offset := request.GetInt("offset", 0)
	if limit < 0 {
		return nil, fmt.Errorf("invalid limit: must be a non-negative integer")
	}
	if offset < 0 {
		return nil, fmt.Errorf("invalid offset: must be a non-negative integer")
	}

	persons, restErr := s.personService.GetAll(ctx, limit, offset, params)
	if restErr != nil {
		return nil, restErr
	}

	return persons, nil
}
