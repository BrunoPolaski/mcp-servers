// mcpServer.AddResourceTemplate(mcp.ResourceTemplate{
// 	URITemplate: mustURITemplate("internalregistry://analyst/{id}"),
// 	Name:        "Analyst",
// 	Description: "Analyst profile by id.",
// 	MIMEType:    "application/json",
// }, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
// 	id, err := parseIDFromURI(request.Params.URI)
// 	if err != nil {
// 		return nil, err
// 	}

// 	analyst, restErr := analystService.GetById(ctx, id)
// 	if restErr != nil {
// 		return nil, restErr
// 	}

//		return resourceJSON(request.Params.URI, dto.NewAnalystDTO(analyst))
//	})
package tools
