/*
	 	s.Handle("GET /person/{id}", middlewares.HandlerChain(
			personController.GetById,
			middlewares.SessionAuthMiddleware(tpf.Redis()),
		))

		s.Handle("GET /person", middlewares.HandlerChain(
			personController.GetAll,
			middlewares.SessionAuthMiddleware(tpf.Redis()),
			middlewares.UserTypeMiddleware(valueobjects.UserTypeAnalyst),
		))
*/
package tools

import (
	"context"
	"fmt"

	"github.com/BrunoPolaski/bureau-mcp-server/internal/infra/controllers/dto"
	"github.com/mark3labs/mcp-go/mcp"
)

func (s *Server) GetPersonByIDTool() mcp.Tool {
	return mcp.NewTool(
		"get_person_by_id",
		mcp.WithDescription(
			`
			Get a person by their ID
			This returns a persons's bureau information, including their name, date of birth, and associated addresses.
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
