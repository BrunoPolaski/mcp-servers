package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/BrunoPolaski/bureau-mcp-server/internal/infra/controllers/dto"
	"github.com/mark3labs/mcp-go/mcp"
)

func (s *Server) GetPersonByIDTool() mcp.Tool {
	return mcp.NewTool(
		"get_customer_by_id",
		mcp.WithDescription(
			`
			Get a customer by their ID
			This returns a customer's bureau information, including their name, date of birth, and associated addresses.
			Example usage:
			{
				"id": 123
			}
			`,
		),
		mcp.WithOutputSchema[dto.PersonDTO](),
		mcp.WithInteger(
			"id",
			mcp.Description("The ID of the customer to retrieve"),
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
		"get_customer_by_document",
		mcp.WithDescription(
			`
			Get a customer by their document number
			This returns a customer's bureau information, including their name, date of birth, and associated addresses.
			Example usage:
			{
				"document": "12345678900"
			}
			`,
		),
		mcp.WithOutputSchema[dto.PersonDTO](),
		mcp.WithString(
			"document",
			mcp.Description("The document number of the customer to retrieve"),
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
		"get_all_customers",
		mcp.WithDescription(
			`
			List customers with pagination.
			This returns only a summary of each customer: their id, name and document.
			It does NOT return their full bureau information.
			To retrieve the complete data for a customer, call get_customer_by_id (using the
			returned "id") or get_customer_by_document (using the returned "document").

			Example usage:
			{
				"limit": 10,
				"offset": 0
			}

			To fetch all customers, set "limit" to 0.
			`,
		),
		mcp.WithOutputSchema[dto.PaginatedResponse[dto.PersonSummaryDTO]](),
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
			mcp.Description("Additional filters for querying customers. e.g. {\"name\": \"John Doe\"}"),
		),
	)
}

func (s *Server) HandleGetAllPersons(ctx context.Context, request mcp.CallToolRequest, args mcp.CallToolParams) (*dto.PaginatedResponse[dto.PersonSummaryDTO], error) {
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

	persons, restErr := s.personService.GetAllSummary(ctx, limit, offset, params)
	if restErr != nil {
		return nil, restErr
	}

	return persons, nil
}
