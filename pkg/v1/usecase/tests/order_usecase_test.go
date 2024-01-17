package use_case_tests

import (
	"context"
	"github.com/nawafswe/orders-service/internal/models"
	ordersMock "github.com/nawafswe/orders-service/mocks/github.com/nawafswe/orders-service/pkg/v1"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestOrderUseCase(t *testing.T) {
	tests := map[string]struct {
		Description string
		Before      func(orderMock *ordersMock.MockOrderUseCase, order *models.Order)
		Assert      func(t *testing.T, mock *ordersMock.MockOrderUseCase, order *models.Order)
		Data        *models.Order
	}{
		"PlaceOrderSuccessfully": {
			Description: "Should Place order successfully",
			Before: func(orderMock *ordersMock.MockOrderUseCase, order *models.Order) {
				orderMock.On("PlaceOrder", mock.Anything, *order).Return(*order, nil)
			},
			Data: &models.Order{
				CustomerId: 1,
				Status:     "New",
				GrandTotal: 10,
				Items: []models.OrderedItem{
					{
						OrderedItemId:   1,
						OrderedQuantity: 10,
						Price:           1,
						Sku:             "12su",
					},
				},
			},
			Assert: func(t *testing.T, orderMock *ordersMock.MockOrderUseCase, order *models.Order) {
				o, err := orderMock.PlaceOrder(context.Background(), *order)
				if err != nil {
					t.Errorf("failed place order, expected error with nil but got %v\n", err)
				}
				orderMock.AssertCalled(t, "PlaceOrder", mock.Anything, *order)
				orderMock.AssertExpectations(t)

				if len(o.Items) != 1 {
					t.Errorf("expected orded items 1 but got %v\n", len(o.Items))
				}
				if o.GrandTotal != order.GrandTotal {
					t.Errorf("expected grand total of %v, but got %v", order.GrandTotal, o.GrandTotal)
				}
			},
		},
	}

	orderUseCaseMock := ordersMock.NewMockOrderUseCase(t)
	for name, test := range tests {

		t.Run(name, func(t *testing.T) {
			t.Logf("running %s", name)
			test.Before(orderUseCaseMock, test.Data)
			test.Assert(t, orderUseCaseMock, test.Data)
		})
	}
}
