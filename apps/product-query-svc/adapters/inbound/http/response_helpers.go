package httpadapter

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/fightingBald/GoTuto/apps/product-query-svc/domain"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, errorBody{Code: code, Message: message})
}

type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details *[]struct {
		Field  *string `json:"field,omitempty"`
		Reason *string `json:"reason,omitempty"`
	} `json:"details,omitempty"`
}

func classifyDomainError(err error) (int, string) {
	switch {
	case errors.Is(err, domain.ErrValidation):
		return http.StatusBadRequest, "VALIDATION"
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound, "NOT_FOUND"
	case errors.Is(err, domain.ErrForbidden):
		return http.StatusForbidden, "FORBIDDEN"
	default:
		return http.StatusInternalServerError, "INTERNAL"
	}
}

func domainErrorMessage(status int, err error) string {
	if status == http.StatusInternalServerError {
		return http.StatusText(status)
	}
	if errors.Is(err, domain.ErrValidation) {
		if split := strings.SplitN(err.Error(), "\n", 2); len(split) == 2 {
			return split[1]
		}
	}
	return err.Error()
}

func errorPayloadFromDomain(err error) (int, errorBody) {
	status, code := classifyDomainError(err)
	return status, errorBody{Code: code, Message: domainErrorMessage(status, err)}
}

func createProductError(err error) (CreateProductResponseObject, bool) {
	status, payload := errorPayloadFromDomain(err)
	switch status {
	case http.StatusBadRequest:
		return CreateProduct400JSONResponse{
			Code:    payload.Code,
			Message: payload.Message,
			Details: payload.Details,
		}, true
	default:
		return nil, false
	}
}

func updateProductError(err error) (UpdateProductResponseObject, bool) {
	status, payload := errorPayloadFromDomain(err)
	switch status {
	case http.StatusBadRequest:
		return UpdateProduct400JSONResponse{
			Code:    payload.Code,
			Message: payload.Message,
			Details: payload.Details,
		}, true
	case http.StatusNotFound:
		return UpdateProduct404JSONResponse{
			Code:    payload.Code,
			Message: payload.Message,
			Details: payload.Details,
		}, true
	default:
		return nil, false
	}
}

func deleteProductError(err error) (DeleteProductByIDResponseObject, bool) {
	status, payload := errorPayloadFromDomain(err)
	switch status {
	case http.StatusBadRequest:
		return DeleteProductByID400JSONResponse{
			Code:    payload.Code,
			Message: payload.Message,
			Details: payload.Details,
		}, true
	case http.StatusNotFound:
		return DeleteProductByID404JSONResponse{
			Code:    payload.Code,
			Message: payload.Message,
			Details: payload.Details,
		}, true
	default:
		return nil, false
	}
}

func getProductError(err error) (GetProductByIDResponseObject, bool) {
	status, payload := errorPayloadFromDomain(err)
	switch status {
	case http.StatusBadRequest:
		return GetProductByID400JSONResponse{
			Code:    payload.Code,
			Message: payload.Message,
			Details: payload.Details,
		}, true
	case http.StatusNotFound:
		return GetProductByID404JSONResponse{
			Code:    payload.Code,
			Message: payload.Message,
			Details: payload.Details,
		}, true
	default:
		return nil, false
	}
}

func searchProductsError(err error) (SearchProductsResponseObject, bool) {
	status, payload := errorPayloadFromDomain(err)
	if status == http.StatusBadRequest {
		return SearchProducts400JSONResponse{
			Code:    payload.Code,
			Message: payload.Message,
			Details: payload.Details,
		}, true
	}
	return nil, false
}

func getUserError(err error) (GetUserByIDResponseObject, bool) {
	status, payload := errorPayloadFromDomain(err)
	switch status {
	case http.StatusBadRequest:
		return GetUserByID400JSONResponse{
			Code:    payload.Code,
			Message: payload.Message,
			Details: payload.Details,
		}, true
	case http.StatusNotFound:
		return GetUserByID404JSONResponse{
			Code:    payload.Code,
			Message: payload.Message,
			Details: payload.Details,
		}, true
	default:
		return nil, false
	}
}

func okCreateProduct(product *domain.Product) CreateProductResponseObject {
	return CreateProduct201JSONResponse(presentProduct(product))
}

func okUpdateProduct(product *domain.Product) UpdateProductResponseObject {
	return UpdateProduct200JSONResponse(presentProduct(product))
}

func okDeleteProduct() DeleteProductByIDResponseObject {
	return DeleteProductByID204Response{}
}

func okGetProduct(product *domain.Product) GetProductByIDResponseObject {
	return GetProductByID200JSONResponse(presentProduct(product))
}

func okSearchProducts(items []domain.Product, page, pageSize, total int) SearchProductsResponseObject {
	return SearchProducts200JSONResponse(ProductList{
		Items:    presentProducts(items),
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	})
}

func okGetUser(user *domain.User) GetUserByIDResponseObject {
	return GetUserByID200JSONResponse(presentUser(user))
}

func listCommentsError(err error) (ListProductCommentsResponseObject, bool) {
	status, payload := errorPayloadFromDomain(err)
	switch status {
	case http.StatusBadRequest:
		return ListProductComments400JSONResponse{
			Code:    payload.Code,
			Message: payload.Message,
			Details: payload.Details,
		}, true
	case http.StatusNotFound:
		return ListProductComments404JSONResponse{
			Code:    payload.Code,
			Message: payload.Message,
			Details: payload.Details,
		}, true
	default:
		return nil, false
	}
}

func createCommentError(err error) (CreateProductCommentResponseObject, bool) {
	status, payload := errorPayloadFromDomain(err)
	switch status {
	case http.StatusBadRequest:
		return CreateProductComment400JSONResponse{
			Code:    payload.Code,
			Message: payload.Message,
			Details: payload.Details,
		}, true
	case http.StatusNotFound:
		return CreateProductComment404JSONResponse{
			Code:    payload.Code,
			Message: payload.Message,
			Details: payload.Details,
		}, true
	default:
		return nil, false
	}
}

func updateCommentError(err error) (UpdateProductCommentResponseObject, bool) {
	status, payload := errorPayloadFromDomain(err)
	switch status {
	case http.StatusBadRequest:
		return UpdateProductComment400JSONResponse{
			Code:    payload.Code,
			Message: payload.Message,
			Details: payload.Details,
		}, true
	case http.StatusForbidden:
		return UpdateProductComment403JSONResponse{
			Code:    payload.Code,
			Message: payload.Message,
			Details: payload.Details,
		}, true
	case http.StatusNotFound:
		return UpdateProductComment404JSONResponse{
			Code:    payload.Code,
			Message: payload.Message,
			Details: payload.Details,
		}, true
	default:
		return nil, false
	}
}

func deleteCommentError(err error) (DeleteProductCommentResponseObject, bool) {
	status, payload := errorPayloadFromDomain(err)
	switch status {
	case http.StatusBadRequest:
		return DeleteProductComment400JSONResponse{
			Code:    payload.Code,
			Message: payload.Message,
			Details: payload.Details,
		}, true
	case http.StatusForbidden:
		return DeleteProductComment403JSONResponse{
			Code:    payload.Code,
			Message: payload.Message,
			Details: payload.Details,
		}, true
	case http.StatusNotFound:
		return DeleteProductComment404JSONResponse{
			Code:    payload.Code,
			Message: payload.Message,
			Details: payload.Details,
		}, true
	default:
		return nil, false
	}
}

func okListComments(items []domain.Comment) ListProductCommentsResponseObject {
	return ListProductComments200JSONResponse(CommentList{Items: presentComments(items)})
}

func okCreateComment(comment *domain.Comment) CreateProductCommentResponseObject {
	return CreateProductComment201JSONResponse(presentComment(comment))
}

func okUpdateComment(comment *domain.Comment) UpdateProductCommentResponseObject {
	return UpdateProductComment200JSONResponse(presentComment(comment))
}

func okDeleteComment() DeleteProductCommentResponseObject {
	return DeleteProductComment204Response{}
}
