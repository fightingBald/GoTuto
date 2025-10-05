package httpadapter

import "github.com/fightingBald/GoTuto/apps/product-query-svc/domain"

const (
	defaultPage     = 1
	defaultPageSize = 20
)

type searchFilters struct {
	query    string
	page     int
	pageSize int
}

func newSearchFilters(params SearchProductsParams) searchFilters {
	filters := searchFilters{
		page:     defaultPage,
		pageSize: defaultPageSize,
	}

	if params.Q != nil {
		filters.query = *params.Q
	}
	if params.Page != nil {
		filters.page = *params.Page
	}
	if params.PageSize != nil {
		filters.pageSize = *params.PageSize
	}

	return filters
}

func newProductFromCreateBody(body *CreateProductJSONRequestBody) (*domain.Product, error) {
	if body == nil {
		return nil, domain.ValidationError("invalid request body")
	}
	return domain.NewProduct(body.Name, amountToCents(body.Price), nil)
}

func newProductFromUpdateBody(id int64, body *UpdateProductJSONRequestBody) (*domain.Product, error) {
	if body == nil {
		return nil, domain.ValidationError("invalid request body")
	}
	product, err := domain.NewProduct(body.Name, amountToCents(body.Price), nil)
	if err != nil {
		return nil, err
	}
	product.ID = id
	return product, nil
}

func commentCreateInput(body *CreateProductCommentJSONRequestBody) (int64, string, error) {
	if body == nil {
		return 0, "", domain.ValidationError("invalid request body")
	}
	return body.UserId, body.Content, nil
}

func commentUpdateInput(body *UpdateProductCommentJSONRequestBody) (int64, string, error) {
	if body == nil {
		return 0, "", domain.ValidationError("invalid request body")
	}
	return body.UserId, body.Content, nil
}
