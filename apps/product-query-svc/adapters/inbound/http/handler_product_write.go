package httpadapter

import (
	"context"
	"net/http"

	"github.com/fightingBald/GoTuto/apps/product-query-svc/domain"
)

func (s *Server) CreateProduct(ctx context.Context, request CreateProductRequestObject) (CreateProductResponseObject, error) {
	if request.Body == nil {
		payload := newErrorPayload("INVALID_JSON", "invalid request body")
		return CreateProduct400JSONResponse{Code: payload.Code, Message: payload.Message, Details: payload.Details}, nil
	}

	body := request.Body
	cents := amountToCents(body.Price)
	product, err := domain.NewProduct(body.Name, cents, nil)
	if err != nil {
		status, payload := errorPayloadFromDomain(err)
		if status == http.StatusBadRequest {
			return CreateProduct400JSONResponse{Code: payload.Code, Message: payload.Message, Details: payload.Details}, nil
		}
		return nil, err
	}

	id, err := s.products.Create(ctx, product)
	if err != nil {
		status, payload := errorPayloadFromDomain(err)
		if status == http.StatusBadRequest {
			return CreateProduct400JSONResponse{Code: payload.Code, Message: payload.Message, Details: payload.Details}, nil
		}
		return nil, err
	}

	created := *product
	created.ID = id

	return CreateProduct201JSONResponse(presentProduct(&created)), nil
}

func (s *Server) UpdateProduct(ctx context.Context, request UpdateProductRequestObject) (UpdateProductResponseObject, error) {
	if request.Id <= 0 {
		payload := newErrorPayload("INVALID_ID", "id must be a positive integer")
		return UpdateProduct400JSONResponse{Code: payload.Code, Message: payload.Message, Details: payload.Details}, nil
	}
	if request.Body == nil {
		payload := newErrorPayload("INVALID_JSON", "invalid request body")
		return UpdateProduct400JSONResponse{Code: payload.Code, Message: payload.Message, Details: payload.Details}, nil
	}

	body := request.Body
	cents := amountToCents(body.Price)
	product, err := domain.NewProduct(body.Name, cents, nil)
	if err != nil {
		status, payload := errorPayloadFromDomain(err)
		if status == http.StatusBadRequest {
			return UpdateProduct400JSONResponse{Code: payload.Code, Message: payload.Message, Details: payload.Details}, nil
		}
		return nil, err
	}

	product.ID = request.Id
	updated, err := s.products.Update(ctx, product)
	if err != nil {
		status, payload := errorPayloadFromDomain(err)
		switch status {
		case http.StatusBadRequest:
			return UpdateProduct400JSONResponse{Code: payload.Code, Message: payload.Message, Details: payload.Details}, nil
		case http.StatusNotFound:
			return UpdateProduct404JSONResponse{Code: payload.Code, Message: payload.Message, Details: payload.Details}, nil
		default:
			return nil, err
		}
	}

	return UpdateProduct200JSONResponse(presentProduct(updated)), nil
}

func (s *Server) DeleteProductByID(ctx context.Context, request DeleteProductByIDRequestObject) (DeleteProductByIDResponseObject, error) {
	if request.Id <= 0 {
		payload := newErrorPayload("INVALID_ID", "id must be a positive integer")
		return DeleteProductByID400JSONResponse{Code: payload.Code, Message: payload.Message, Details: payload.Details}, nil
	}

	if err := s.products.Remove(ctx, request.Id); err != nil {
		status, payload := errorPayloadFromDomain(err)
		switch status {
		case http.StatusBadRequest:
			return DeleteProductByID400JSONResponse{Code: payload.Code, Message: payload.Message, Details: payload.Details}, nil
		case http.StatusNotFound:
			return DeleteProductByID404JSONResponse{Code: payload.Code, Message: payload.Message, Details: payload.Details}, nil
		default:
			return nil, err
		}
	}

	return DeleteProductByID204Response{}, nil
}
