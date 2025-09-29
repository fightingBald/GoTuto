package httpadapter

import "context"

func (s *Server) CreateProduct(ctx context.Context, request CreateProductRequestObject) (CreateProductResponseObject, error) {
	product, err := newProductFromCreateBody(request.Body)
	if err != nil {
		if resp, handled := createProductError(err); handled {
			return resp, nil
		}
		return nil, err
	}

	id, err := s.products.Create(ctx, product)
	if err != nil {
		if resp, handled := createProductError(err); handled {
			return resp, nil
		}
		return nil, err
	}

	product.ID = id

	return okCreateProduct(product), nil
}

func (s *Server) UpdateProduct(ctx context.Context, request UpdateProductRequestObject) (UpdateProductResponseObject, error) {
	product, err := newProductFromUpdateBody(request.Id, request.Body)
	if err != nil {
		if resp, handled := updateProductError(err); handled {
			return resp, nil
		}
		return nil, err
	}

	updated, err := s.products.Update(ctx, product)
	if err != nil {
		if resp, handled := updateProductError(err); handled {
			return resp, nil
		}
		return nil, err
	}

	return okUpdateProduct(updated), nil
}

func (s *Server) DeleteProductByID(ctx context.Context, request DeleteProductByIDRequestObject) (DeleteProductByIDResponseObject, error) {
	if err := s.products.Remove(ctx, request.Id); err != nil {
		if resp, handled := deleteProductError(err); handled {
			return resp, nil
		}
		return nil, err
	}

	return okDeleteProduct(), nil
}
