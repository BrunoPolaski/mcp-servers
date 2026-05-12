package tools

// func HandleGetAddressResourceTemplate(
// 	ctx context.Context,
// 	request mcp.ReadResourceRequest,
// ) (*mcp.ReadResourceResult, error) {
// 	id, err := parseIDFromURI(request.Params.URI)
// 	if err != nil {
// 		return nil, err
// 	}

// 	ctx, err = requireSessionFromArgs(ctx, tpf.Redis(), request.Params.Arguments)
// 	if err != nil {
// 		return nil, err
// 	}

// 	address, restErr := addressService.GetById(ctx, id)
// 	if restErr != nil {
// 		return nil, restErr
// 	}

// 	return resourceJSON(request.Params.URI, dto.NewAddressDTO(address))
// }
