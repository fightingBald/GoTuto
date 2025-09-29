package httpadapter

import "context"

func (s *Server) GetProductByID(ctx context.Context, request GetProductByIDRequestObject) (GetProductByIDResponseObject, error) {
	product, err := s.products.FetchByID(ctx, request.Id)
	if err != nil {
		if resp, handled := getProductError(err); handled {
			return resp, nil
		}
		return nil, err
	}

	return okGetProduct(product), nil
}

func (s *Server) SearchProducts(ctx context.Context, request SearchProductsRequestObject) (SearchProductsResponseObject, error) {
	filters := newSearchFilters(request.Params)

	items, total, err := s.products.Search(ctx, filters.query, filters.page, filters.pageSize)
	if err != nil {
		if resp, handled := searchProductsError(err); handled {
			return resp, nil
		}
		return nil, err
	}

	return okSearchProducts(items, filters.page, filters.pageSize, total), nil
}
