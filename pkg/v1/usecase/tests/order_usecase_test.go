package use_case_tests

import (
	"context"
	"fmt"
	"github.com/nawafswe/orders-service/internal/models"
	ordersMock "github.com/nawafswe/orders-service/mocks/github.com/nawafswe/orders-service/pkg/v1"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
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
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()
				o, err := orderMock.PlaceOrder(ctx, *order)
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
		"FailPlaceOrderDueToInvalidItemQuantity": {
			Description: "Should fail place order due to invalid item quantities",
			Before: func(orderMock *ordersMock.MockOrderUseCase, order *models.Order) {
				orderMock.On("PlaceOrder", mock.Anything, *order).Return(models.Order{}, fmt.Errorf("supplied quantity for item with sku %v, should be greater than zero, received is %v", order.Items[0].Sku, order.Items[0].OrderedQuantity))
			},
			Data: &models.Order{
				CustomerId: 1,
				Status:     "New",
				GrandTotal: 10,
				Items: []models.OrderedItem{
					{
						OrderedItemId:   1,
						OrderedQuantity: -10,
						Price:           1,
						Sku:             "12su",
					},
				},
			},
			Assert: func(t *testing.T, orderMock *ordersMock.MockOrderUseCase, order *models.Order) {
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()
				_, err := orderMock.PlaceOrder(ctx, *order)
				t.Logf("received err: %v\n", err)
				if err != nil {
					t.Errorf("failed place order, expected error with %v but got nil\n", err)
				}
				orderMock.AssertCalled(t, "PlaceOrder", mock.Anything, *order)
				orderMock.AssertExpectations(t)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			orderUseCaseMock := ordersMock.NewMockOrderUseCase(t)
			t.Logf("running %s", name)
			test.Before(orderUseCaseMock, test.Data)
			test.Assert(t, orderUseCaseMock, test.Data)
		})
	}
}
