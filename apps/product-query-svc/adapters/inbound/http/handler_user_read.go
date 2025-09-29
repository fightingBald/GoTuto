package httpadapter

import (
	"context"
	"net/http"
)

func (s *Server) GetUserByID(ctx context.Context, request GetUserByIDRequestObject) (GetUserByIDResponseObject, error) {
	if request.Id <= 0 {
		payload := newErrorPayload("INVALID_ID", "id must be a positive integer")
		return GetUserByID400JSONResponse{Code: payload.Code, Message: payload.Message, Details: payload.Details}, nil
	}

	user, err := s.users.FetchByID(ctx, request.Id)
	if err != nil {
		status, payload := errorPayloadFromDomain(err)
		switch status {
		case http.StatusBadRequest:
			return GetUserByID400JSONResponse{Code: payload.Code, Message: payload.Message, Details: payload.Details}, nil
		case http.StatusNotFound:
			return GetUserByID404JSONResponse{Code: payload.Code, Message: payload.Message, Details: payload.Details}, nil
		default:
			return nil, err
		}
	}

	presented := presentUser(user)
	return GetUserByID200JSONResponse(presented), nil
}
