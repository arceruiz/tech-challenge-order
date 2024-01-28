package rest_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"tech-challenge-order/internal/canonical"
	"tech-challenge-order/internal/channels/rest"
	"tech-challenge-order/internal/mocks"
	"tech-challenge-order/internal/service"
	"testing"
	"time"

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
				orderService: &mocks.OrderServiceMock{},
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusNotFound,
			},
		},
	}

	for _, tc := range tests {
		p := rest.NewOrderChannel(tc.given.orderService)
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
				request: createJsonRequest(http.MethodPost, endpoint, rest.OrderRequest{
					ID:         "",
					CustomerID: "",
					Status:     "",
					CreatedAt:  time.Time{},
					UpdatedAt:  time.Time{},
					Total:      0,
					OrderItems: []rest.OrderItem{},
				}),
				orderService: mockOrderServiceForCreate(canonical.Order{}, canonical.Order{}),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusOK,
			},
		},
		"given error creating, must return error": {
			given: Given{
				request:      createJsonRequest(http.MethodPost, endpoint, rest.OrderRequest{}),
				orderService: mockOrderServiceForCreateError(canonical.Order{}, canonical.Order{}),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusInternalServerError,
			},
		},
		"given status unrecognized, must return error": {
			given: Given{
				request: createJsonRequest(http.MethodPost, endpoint, rest.OrderRequest{
					ID:         "",
					CustomerID: "",
					Status:     "asdasd",
					CreatedAt:  time.Time{},
					UpdatedAt:  time.Time{},
					Total:      0,
					OrderItems: []rest.OrderItem{},
				}),
				orderService: mockOrderServiceForCreate(canonical.Order{}, canonical.Order{}),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusBadRequest,
			},
		},
		"given wrong format must return error": {
			given: Given{
				request:      createRequest(http.MethodPost, endpoint),
				orderService: mockOrderServiceForCreate(canonical.Order{}, canonical.Order{}),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusBadRequest,
			},
		},
		"given invalid data, must return bad request": {
			given: Given{
				request:      createRequest(http.MethodPost, endpoint),
				orderService: mockOrderServiceForCreate(canonical.Order{}, canonical.Order{}),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusBadRequest,
			},
		},
	}

	for _, tc := range tests {
		rec := httptest.NewRecorder()
		err := rest.NewOrderChannel(tc.given.orderService).Create(echo.New().NewContext(tc.given.request, rec))
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
				pathParamID:  "valid_ID",
				request:      createJsonRequest(http.MethodPost, endpoint, rest.OrderRequest{}),
				orderService: mockOrderServiceForUpdate("valid_ID", canonical.Order{}),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusOK,
			},
		},
		"given given error updating, must return error": {
			given: Given{
				pathParamID:  "invalid_ID",
				request:      createJsonRequest(http.MethodPost, endpoint, rest.OrderRequest{}),
				orderService: mockOrderServiceForUpdate("valid_ID", canonical.Order{}),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusInternalServerError,
			},
		},
		"given wrong format must return error": {
			given: Given{
				pathParamID:  "valid_ID",
				request:      createRequest(http.MethodPost, endpoint),
				orderService: mockOrderServiceForUpdate("valid_ID", canonical.Order{}),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusBadRequest,
			},
		},
		"given invalid data, must return bad request": {
			given: Given{
				pathParamID:  "invalid_ID",
				request:      createRequest(http.MethodPost, endpoint),
				orderService: mockOrderServiceForUpdate("valid_ID", canonical.Order{}),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusBadRequest,
			},
		},
		"given empty id data, must return bad request": {
			given: Given{
				pathParamID:  "",
				request:      createRequest(http.MethodPost, endpoint),
				orderService: mockOrderServiceForUpdate("valid_ID", canonical.Order{}),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusBadRequest,
			},
		},
		"given invalid status, must return bad request": {
			given: Given{
				pathParamID: "valid_ID",
				request: createJsonRequest(http.MethodPost, endpoint, rest.OrderRequest{
					ID:         "",
					CustomerID: "",
					Status:     "asdasd",
					CreatedAt:  time.Time{},
					UpdatedAt:  time.Time{},
					Total:      0,
					OrderItems: []rest.OrderItem{},
				}),
				orderService: mockOrderServiceForUpdate("valid_ID", canonical.Order{}),
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
		err := rest.NewOrderChannel(tc.given.orderService).Update(e)
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
				request:        createJsonRequest(http.MethodPost, endpoint, rest.OrderRequest{}),
				orderService:   mockOrderServiceForUpdate("valid_ID", canonical.Order{}),
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
				orderService:   mockOrderServiceForUpdate("valid_ID", canonical.Order{}),
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
				orderService:   mockOrderServiceForUpdate("valid_ID", canonical.Order{}),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusBadRequest,
			},
		},
		"given wring status, must return bad request": {
			given: Given{
				pathParamID:    "valid_ID",
				pathParamKey:   "status",
				pathParamValue: "assdasdasdasd",
				request:        createRequest(http.MethodPost, endpoint),
				orderService:   mockOrderServiceForUpdate("valid_ID", canonical.Order{}),
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
				request:        createJsonRequest(http.MethodPost, endpoint, rest.OrderRequest{}),
				orderService:   mockOrderServiceForUpdate("valid_ID", canonical.Order{}),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusInternalServerError,
			},
		},
		"given error retrieving order, must return not found": {
			given: Given{
				pathParamID:    "invalid_ID_nil",
				pathParamKey:   "status",
				pathParamValue: "RECEIVED",
				request:        createJsonRequest(http.MethodPost, endpoint, rest.OrderRequest{}),
				orderService:   mockOrderServiceForUpdate("valid_ID", canonical.Order{}),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusNotFound,
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
		err := rest.NewOrderChannel(tc.given.orderService).UpdateStatus(e)
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
				request:      createJsonRequest(http.MethodPost, endpoint, rest.OrderRequest{}),
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
				request:      createJsonRequest(http.MethodPost, endpoint, rest.OrderRequest{}),
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
				request:      createJsonRequest(http.MethodPost, endpoint, rest.OrderRequest{}),
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
		err := rest.NewOrderChannel(tc.given.orderService).CheckoutOrder(e)
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
				orderService: mockOrderServiceForGetAll("1234", []canonical.Order{{
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
				orderService: mockOrderServiceForGetAll_error("1234", nil),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusInternalServerError,
			},
		},
		"given invalic id returns no order and status 404": {
			given: Given{
				request:      createRequest(http.MethodGet, endpoint),
				orderService: mockOrderServiceForGetAll("1234", nil),
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
		err := rest.NewOrderChannel(tc.given.orderService).Get(e)
		statusCode := rec.Result().StatusCode

		assert.Equal(t, tc.expected.statusCode, statusCode)

		tc.expected.err(t, err)
	}
}

func mockOrderServiceForRemove(id string, orderReturned canonical.Order) *mocks.OrderServiceMock {
	mockOrderSvc := new(mocks.OrderServiceMock)

	mockOrderSvc.
		On("Remove", mock.Anything, id).
		Return(nil)

	mockOrderSvc.
		On("Remove", mock.Anything, "invalid_ID").
		Return(errors.New(""))

	return mockOrderSvc
}

func mockOrderServiceForUpdate(id string, orderReturned canonical.Order) *mocks.OrderServiceMock {
	mockOrderSvc := new(mocks.OrderServiceMock)

	mockOrderSvc.
		On("GetByID", mock.Anything, id).
		Return(&orderReturned, nil)

	mockOrderSvc.
		On("GetByID", mock.Anything, "invalid_ID_updt").
		Return(&orderReturned, nil)

	mockOrderSvc.
		On("GetByID", mock.Anything, "invalid_ID").
		Return(&orderReturned, errors.New(""))

	mockOrderSvc.
		On("GetByID", mock.Anything, "invalid_ID_nil").
		Return(nil, nil)

	mockOrderSvc.
		On("Update", mock.Anything, id, orderReturned).
		Return(nil)

	mockOrderSvc.
		On("Update", mock.Anything, "invalid_ID", orderReturned).
		Return(errors.New(""))

	mockOrderSvc.
		On("Update", mock.Anything, "invalid_ID_updt", orderReturned).
		Return(errors.New(""))

	return mockOrderSvc
}

func mockOrderServiceForCheckout(id string, orderReturned canonical.Order) *mocks.OrderServiceMock {
	mockOrderSvc := new(mocks.OrderServiceMock)

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

func mockOrderServiceForGetByStatus(status canonical.OrderStatus, orderReturned []canonical.Order) *mocks.OrderServiceMock {
	mockOrderSvc := new(mocks.OrderServiceMock)

	mockOrderSvc.
		On("GetByStatus", mock.Anything, status).
		Return(orderReturned, nil)

	return mockOrderSvc
}

func mockOrderServiceForGetByID(orderID string, orderReturned *canonical.Order) *mocks.OrderServiceMock {
	mockOrderSvc := new(mocks.OrderServiceMock)

	mockOrderSvc.
		On("GetByID", mock.Anything, orderID).
		Return(orderReturned, nil)

	return mockOrderSvc
}

func mockOrderServiceForGetAll(orderID string, orderReturned []canonical.Order) *mocks.OrderServiceMock {
	mockOrderSvc := new(mocks.OrderServiceMock)

	mockOrderSvc.
		On("GetAll", mock.Anything).
		Return(orderReturned, nil).Maybe()
	return mockOrderSvc
}

func mockOrderServiceForGetAll_error(orderID string, orderReturned []canonical.Order) *mocks.OrderServiceMock {
	mockOrderSvc := new(mocks.OrderServiceMock)

	mockOrderSvc.
		On("GetAll", mock.Anything).
		Return(orderReturned, errors.New(""))

	return mockOrderSvc
}

func mockOrderServiceForCreate(orderReceived, orderReturned canonical.Order) *mocks.OrderServiceMock {
	mockOrderSvc := new(mocks.OrderServiceMock)
	mockOrderSvc.On("Create", mock.Anything, orderReceived).Return(nil)
	mockOrderSvc.On("Create", mock.Anything, canonical.Order{
		ID:         "invalid_ID",
		CustomerID: "",
		Status:     0,
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
		Total:      0,
		OrderItems: []canonical.OrderItem{},
	}).Return(errors.New(""))
	return mockOrderSvc
}
func mockOrderServiceForCreateError(orderReceived, orderReturned canonical.Order) *mocks.OrderServiceMock {
	mockOrderSvc := new(mocks.OrderServiceMock)
	mockOrderSvc.On("Create", mock.Anything, orderReceived).Return(errors.New(""))
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
	return req
}
