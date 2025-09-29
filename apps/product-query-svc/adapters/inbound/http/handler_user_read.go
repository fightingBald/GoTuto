package httpadapter

import "context"

func (s *Server) GetUserByID(ctx context.Context, request GetUserByIDRequestObject) (GetUserByIDResponseObject, error) {
	user, err := s.users.FetchByID(ctx, request.Id)
	if err != nil {
		if resp, handled := getUserError(err); handled {
			return resp, nil
		}
		return nil, err
	}

	return okGetUser(user), nil
}
