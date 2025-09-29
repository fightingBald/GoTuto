package httpadapter

import (
	"context"
	"net/http"
)

func (s *Server) GetProductByID(ctx context.Context, request GetProductByIDRequestObject) (GetProductByIDResponseObject, error) {
	if request.Id <= 0 {
		payload := newErrorPayload("INVALID_ID", "id must be a positive integer")
		return GetProductByID400JSONResponse{Code: payload.Code, Message: payload.Message, Details: payload.Details}, nil
	}

	product, err := s.products.FetchByID(ctx, request.Id)
	if err != nil {
		status, payload := errorPayloadFromDomain(err)
		switch status {
		case http.StatusBadRequest:
			return GetProductByID400JSONResponse{Code: payload.Code, Message: payload.Message, Details: payload.Details}, nil
		case http.StatusNotFound:
			return GetProductByID404JSONResponse{Code: payload.Code, Message: payload.Message, Details: payload.Details}, nil
		default:
			return nil, err
		}
	}

	return GetProductByID200JSONResponse(presentProduct(product)), nil
}

func (s *Server) SearchProducts(ctx context.Context, request SearchProductsRequestObject) (SearchProductsResponseObject, error) {
	params := request.Params

	query := ""
	if params.Q != nil {
		query = *params.Q
	}
	if query != "" && len(query) < 3 {
		payload := newErrorPayload("INVALID_QUERY", "q must be at least 3 characters if provided")
		return SearchProducts400JSONResponse{Code: payload.Code, Message: payload.Message, Details: payload.Details}, nil
	}

	page := 1
	if params.Page != nil {
		page = *params.Page
	}

	pageSize := 20
	if params.PageSize != nil {
		pageSize = *params.PageSize
	}

	items, total, err := s.products.Search(ctx, query, page, pageSize)
	if err != nil {
		status, payload := errorPayloadFromDomain(err)
		if status == http.StatusBadRequest {
			return SearchProducts400JSONResponse{Code: payload.Code, Message: payload.Message, Details: payload.Details}, nil
		}
		return nil, err
	}

	resp := ProductList{
		Items:    presentProducts(items),
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}
	return SearchProducts200JSONResponse(resp), nil
}
