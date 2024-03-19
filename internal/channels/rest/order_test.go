package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"tech-challenge-order/internal/canonical"
	"tech-challenge-order/internal/service"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterGroup(t *testing.T) {
	endpoint := "/order"

	type Given struct {
		group        *echo.Group
		orderService service.OrderService
	}
	type Expected struct {
		err        assert.ErrorAssertionFunc
		statusCode int
	}
	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"given valid group, should register endpoints successfully": {
			given: Given{
				group:        echo.New().Group("/order"),
				orderService: &OrderServiceMock{},
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusNotFound,
			},
		},
	}

	for _, tc := range tests {
		p := order{
			service: tc.given.orderService,
		}
		p.RegisterGroup(tc.given.group)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, endpoint+"/123", nil)
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues("123")

		e.ServeHTTP(rec, req)
		statusCode := rec.Result().StatusCode

		assert.Equal(t, tc.expected.statusCode, statusCode)
	}
}

func TestCreate(t *testing.T) {
	endpoint := "/order"

	type Given struct {
		request      *http.Request
		orderService service.OrderService
	}
	type Expected struct {
		err        assert.ErrorAssertionFunc
		statusCode int
	}
	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"given normal json income must process normally": {
			given: Given{
				request: createJsonRequest(http.MethodPost, endpoint, OrderRequest{
					OrderItems: []OrderItem{
						{
							ProductId: "product_id",
						},
					},
				}),
				orderService: mockOrderServiceForCreate1("product_id", nil, 1),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusOK,
			},
		},
		"given error creating, must return error": {
			given: Given{
				request: createJsonRequest(http.MethodPost, endpoint, OrderRequest{
					OrderItems: []OrderItem{
						{
							ProductId: "product_id",
						},
					},
				}),
				orderService: mockOrderServiceForCreate1("product_id", errors.New("generic error"), 1),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusInternalServerError,
			},
		},
		"given status unrecognized, must return error": {
			given: Given{
				request: createJsonRequest(http.MethodPost, endpoint, OrderRequest{
					OrderItems: []OrderItem{},
				}),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusBadRequest,
			},
		},
		"given wrong format must return error": {
			given: Given{
				request: createRequest(http.MethodPost, endpoint),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusBadRequest,
			},
		},
		"given invalid data, must return bad request": {
			given: Given{
				request: createRequest(http.MethodPost, endpoint),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusBadRequest,
			},
		},
	}

	for _, tc := range tests {
		rec := httptest.NewRecorder()

		orderSvc := order{
			service: tc.given.orderService,
		}

		err := orderSvc.Create(echo.New().NewContext(tc.given.request, rec))

		statusCode := rec.Result().StatusCode

		assert.Equal(t, tc.expected.statusCode, statusCode)

		tc.expected.err(t, err)
	}
}

func TestUpdate(t *testing.T) {
	endpoint := "/order"

	type Given struct {
		request      *http.Request
		pathParamID  string
		orderService service.OrderService
	}
	type Expected struct {
		err        assert.ErrorAssertionFunc
		statusCode int
	}
	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"given normal json income must process normally": {
			given: Given{
				pathParamID: "valid_ID",
				request: createJsonRequest(http.MethodPost, endpoint, OrderRequest{
					OrderItems: []OrderItem{
						{
							ProductId: "valid_ID",
						},
					},
				}),
				orderService: mockOrderServiceForUpdate1("valid_ID", nil),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusOK,
			},
		},
		"given error updating, must return error": {
			given: Given{
				pathParamID: "invalid_ID",
				request: createJsonRequest(http.MethodPost, endpoint, OrderRequest{
					OrderItems: []OrderItem{
						{
							ProductId: "invalid_ID",
						},
					},
				}),
				orderService: mockOrderServiceForUpdate1("invalid_ID", errors.New("generic error")),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusInternalServerError,
			},
		},
		"given wrong format must return error": {
			given: Given{
				pathParamID: "valid_ID",
				request:     createRequest(http.MethodPost, endpoint),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusBadRequest,
			},
		},
		"given invalid data, must return bad request": {
			given: Given{
				pathParamID: "invalid_ID",
				request:     createRequest(http.MethodPost, endpoint),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusBadRequest,
			},
		},
		"given empty id data, must return bad request": {
			given: Given{
				pathParamID: "",
				request:     createRequest(http.MethodPost, endpoint),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusBadRequest,
			},
		},
		"given invalid status, must return bad request": {
			given: Given{
				pathParamID: "valid_ID",
				request: createJsonRequest(http.MethodPost, endpoint, OrderRequest{
					OrderItems: []OrderItem{},
				}),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusBadRequest,
			},
		},
	}

	for _, tc := range tests {
		rec := httptest.NewRecorder()
		e := echo.New().NewContext(tc.given.request, rec)
		e.SetPath("/:id")
		e.SetParamNames("id")
		e.SetParamValues(tc.given.pathParamID)

		orderSvc := order{
			service: tc.given.orderService,
		}

		err := orderSvc.Update(e)

		statusCode := rec.Result().StatusCode

		assert.Equal(t, tc.expected.statusCode, statusCode)

		tc.expected.err(t, err)
	}
}

func TestUpdateStatus(t *testing.T) {
	endpoint := "/order"

	type Given struct {
		request        *http.Request
		pathParamID    string
		pathParamKey   string
		pathParamValue string
		orderService   service.OrderService
	}
	type Expected struct {
		err        assert.ErrorAssertionFunc
		statusCode int
	}
	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"given normal json income must process normally": {
			given: Given{
				pathParamID:    "valid_ID",
				pathParamKey:   "status",
				pathParamValue: "RECEIVED",
				request:        createJsonRequest(http.MethodPost, endpoint, OrderRequest{}),
				orderService:   mockOrderServiceForUpdateStatus("valid_ID", nil),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusOK,
			},
		},
		"given invalid data, must return bad request": {
			given: Given{
				pathParamID:    "invalid_ID",
				pathParamKey:   "status",
				pathParamValue: "RECEIVED",
				request:        createRequest(http.MethodPost, endpoint),
				orderService:   mockOrderServiceForUpdateStatus("invalid_ID", errors.New("generic error")),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusInternalServerError,
			},
		},
		"given empty id, must return bad request": {
			given: Given{
				pathParamID:    "",
				pathParamKey:   "status",
				pathParamValue: "RECEIVED",
				request:        createRequest(http.MethodPost, endpoint),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusBadRequest,
			},
		},
		"given error updating must return internal server error": {
			given: Given{
				pathParamID:    "invalid_ID_updt",
				pathParamKey:   "status",
				pathParamValue: "RECEIVED",
				request:        createJsonRequest(http.MethodPost, endpoint, OrderRequest{}),
				orderService:   mockOrderServiceForUpdateStatus("invalid_ID_updt", errors.New("generic error")),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusInternalServerError,
			},
		},
	}

	for _, tc := range tests {
		rec := httptest.NewRecorder()
		e := echo.New().NewContext(tc.given.request, rec)
		if tc.given.pathParamKey != "" {
			e.QueryParams().Add("id", tc.given.pathParamID)
			e.QueryParams().Add(tc.given.pathParamKey, tc.given.pathParamValue)
		}

		orderSvc := order{
			service: tc.given.orderService,
		}

		err := orderSvc.UpdateStatus(e)
		statusCode := rec.Result().StatusCode

		assert.Equal(t, tc.expected.statusCode, statusCode)

		tc.expected.err(t, err)
	}
}

func TestCheckout(t *testing.T) {
	endpoint := "/order"

	type Given struct {
		request      *http.Request
		pathParamID  string
		orderService service.OrderService
	}
	type Expected struct {
		err        assert.ErrorAssertionFunc
		statusCode int
	}
	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"given normal json income must process normally": {
			given: Given{
				pathParamID:  "valid_ID",
				request:      createJsonRequest(http.MethodPost, endpoint, OrderRequest{}),
				orderService: mockOrderServiceForCheckout("valid_ID", canonical.Order{}),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusOK,
			},
		},
		"given given error updating, must return error": {
			given: Given{
				pathParamID:  "invalid_ID",
				request:      createJsonRequest(http.MethodPost, endpoint, OrderRequest{}),
				orderService: mockOrderServiceForCheckout("valid_ID", canonical.Order{}),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusInternalServerError,
			},
		},
		"given empty id, must return error": {
			given: Given{
				pathParamID:  "",
				request:      createJsonRequest(http.MethodPost, endpoint, OrderRequest{}),
				orderService: mockOrderServiceForCheckout("valid_ID", canonical.Order{}),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusBadRequest,
			},
		},
	}

	for _, tc := range tests {
		rec := httptest.NewRecorder()
		e := echo.New().NewContext(tc.given.request, rec)

		e.QueryParams().Add("id", tc.given.pathParamID)

		orderSvc := order{
			service: tc.given.orderService,
		}

		err := orderSvc.CheckoutOrder(e)
		statusCode := rec.Result().StatusCode

		assert.Equal(t, tc.expected.statusCode, statusCode)

		tc.expected.err(t, err)
	}
}

func TestGet(t *testing.T) {
	endpoint := "/order/"

	type Given struct {
		request        *http.Request
		pathParamKey   string
		pathParamValue string
		orderService   service.OrderService
	}
	type Expected struct {
		err        assert.ErrorAssertionFunc
		statusCode int
	}
	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"given clean request returns valid order and status 200": {
			given: Given{
				request: createRequest(http.MethodGet, endpoint),
				orderService: mockOrderServiceForGetAll([]canonical.Order{{
					ID: "1234",
				}}),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusOK,
			},
		},
		"given valid id returns valid order and status 200": {
			given: Given{
				request:        createRequest(http.MethodGet, endpoint),
				pathParamKey:   "id",
				pathParamValue: "1234",
				orderService: mockOrderServiceForGetByID("1234", &canonical.Order{
					ID: "1234",
				}),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusOK,
			},
		},
		"given valid status returns valid order and status 200": {
			given: Given{
				request:        createRequest(http.MethodGet, endpoint),
				pathParamKey:   "status",
				pathParamValue: "RECEIVED",
				orderService: mockOrderServiceForGetByStatus(canonical.ORDER_RECEIVED, []canonical.Order{{
					ID: "1234",
				}}),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusOK,
			},
		},
		"given valid request returns valid order and status 200": {
			given: Given{
				request:        createRequest(http.MethodGet, endpoint),
				pathParamKey:   "status",
				pathParamValue: "RECEIVED",
				orderService: mockOrderServiceForGetByStatus(canonical.ORDER_RECEIVED, []canonical.Order{{
					ID: "1234",
				}, {
					ID: "1234",
				}, {
					ID: "1234",
				}}),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusOK,
			},
		},
		"given empty id returns no order and status 400": {
			given: Given{
				request:      createRequest(http.MethodGet, endpoint),
				orderService: mockOrderServiceForGetAll_error(nil),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusInternalServerError,
			},
		},
		"given invalic id returns no order and status 404": {
			given: Given{
				request:      createRequest(http.MethodGet, endpoint),
				orderService: mockOrderServiceForGetAll(nil),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusNotFound,
			},
		},
	}

	for name, tc := range tests {
		t.Log(name)
		rec := httptest.NewRecorder()
		e := echo.New().NewContext(tc.given.request, rec)

		if tc.given.pathParamKey != "" {
			e.QueryParams().Add(tc.given.pathParamKey, tc.given.pathParamValue)
		}
		orderSvc := order{
			service: tc.given.orderService,
		}

		err := orderSvc.Get(e)
		statusCode := rec.Result().StatusCode

		assert.Equal(t, tc.expected.statusCode, statusCode)

		tc.expected.err(t, err)
	}
}

func mockOrderServiceForUpdate1(id string, errToReturn error) *OrderServiceMock {
	mockOrderSvc := new(OrderServiceMock)

	mockOrderSvc.On("Update", id).Return(errToReturn)

	return mockOrderSvc
}

func mockOrderServiceForUpdateStatus(id string, errToReturn error) *OrderServiceMock {
	mockOrderSvc := new(OrderServiceMock)

	mockOrderSvc.On("UpdateStatus", id).Return(errToReturn)

	return mockOrderSvc
}

func mockOrderServiceForCheckout(id string, orderReturned canonical.Order) *OrderServiceMock {
	mockOrderSvc := new(OrderServiceMock)

	mockOrderSvc.
		On("CheckoutOrder", mock.Anything, id).
		Return(&orderReturned, nil)

	mockOrderSvc.
		On("CheckoutOrder", mock.Anything, "invalid_ID").
		Return(nil, errors.New(""))

	mockOrderSvc.
		On("CheckoutOrder", mock.Anything, "invalid_ID_updt").
		Return(nil, errors.New(""))

	return mockOrderSvc
}

func mockOrderServiceForGetByStatus(status canonical.OrderStatus, orderReturned []canonical.Order) *OrderServiceMock {
	mockOrderSvc := new(OrderServiceMock)

	mockOrderSvc.
		On("GetByStatus", mock.Anything, status).
		Return(orderReturned, nil)

	return mockOrderSvc
}

func mockOrderServiceForGetByID(orderID string, orderReturned *canonical.Order) *OrderServiceMock {
	mockOrderSvc := new(OrderServiceMock)

	mockOrderSvc.
		On("GetByID", mock.Anything, orderID).
		Return(orderReturned, nil)

	return mockOrderSvc
}

func mockOrderServiceForGetAll(orderReturned []canonical.Order) *OrderServiceMock {
	mockOrderSvc := new(OrderServiceMock)

	mockOrderSvc.
		On("GetAll", mock.Anything).
		Return(orderReturned, nil).Maybe()
	return mockOrderSvc
}

func mockOrderServiceForGetAll_error(orderReturned []canonical.Order) *OrderServiceMock {
	mockOrderSvc := new(OrderServiceMock)

	mockOrderSvc.
		On("GetAll", mock.Anything).
		Return(orderReturned, errors.New(""))

	return mockOrderSvc
}

func mockOrderServiceForCreate1(idInput string, errReturn error, times int) *OrderServiceMock {
	mockOrderSvc := new(OrderServiceMock)

	mockOrderSvc.On("Create", mock.Anything, mock.MatchedBy(func(id canonical.Order) bool {
		_, ok := id.OrderItems[idInput]

		return ok
	})).Return(errReturn).Times(times)

	return mockOrderSvc
}

func createRequest(method, endpoint string) *http.Request {
	req := createJsonRequest(method, endpoint, nil)
	req.Header.Del("Content-Type")
	return req
}

func createJsonRequest(method, endpoint string, request interface{}) *http.Request {
	json, _ := json.Marshal(request)
	req := httptest.NewRequest(method, endpoint, bytes.NewReader(json))
	req.Header.Set("Content-Type", "application/json")

	token, _ := generateToken("")
	req.Header.Set("authorization", "Berear "+token)
	return req
}

func generateToken(userId string) (string, error) {
	permissions := jwt.MapClaims{}
	permissions["userId"] = userId
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permissions)

	return token.SignedString([]byte(""))
}
