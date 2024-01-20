package use_case_tests

import (
	"context"
	"errors"
	"github.com/nawafswe/orders-service/internal/models"
	messagesMock "github.com/nawafswe/orders-service/mocks/github.com/nawafswe/orders-service/pkg/messaging"
	ordersMock "github.com/nawafswe/orders-service/mocks/github.com/nawafswe/orders-service/pkg/v1"
	useCases "github.com/nawafswe/orders-service/pkg/v1/usecase"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"slices"
	"testing"
)

func TestOrderUseCase(t *testing.T) {
	tests := map[string]struct {
		Description    string
		Input          models.Order
		ExpectedErr    error
		ExpectedResult models.Order
	}{
		"PlaceOrderSuccessfully": {
			Description: "Should Place order successfully",
			ExpectedErr: nil,
			ExpectedResult: models.Order{
				Model:      gorm.Model{ID: 1},
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
			Input: models.Order{
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
		},
		"FailPlaceOrderDueToInvalidItemQuantity": {
			Description: "Should fail place order due to invalid item quantities",
			Input: models.Order{
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
			ExpectedResult: models.Order{},
			ExpectedErr:    errors.New("supplied quantity for item with sku 12su, should be greater than zero, received is -10"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Logf("running %s", name)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			pubSubMock := messagesMock.NewMockMessageService(t)
			ordersRepoMock := ordersMock.NewMockOrderRepo(t)
			ordersUseCase := useCases.NewOrderUseCase(ordersRepoMock, pubSubMock)
			// setting up mocks
			if test.ExpectedErr == nil {
				newOrder := test.Input
				newOrder.ID = 1
				ordersRepoMock.On("Create", mock.Anything, test.Input).Return(newOrder, nil)
				pubSubMock.On("PublishAsync", mock.Anything, "orderCreated", mock.Anything).Return(nil)
				pubSubMock.On("PublishAsync", mock.Anything, "orderStatusChanged", mock.Anything).Return(nil)
			}
			result, err := ordersUseCase.PlaceOrder(ctx, test.Input)
			if err == nil && test.ExpectedErr != nil {
				t.Errorf("expected error from %s is %v but got %v", name, test.ExpectedErr, err)
			}
			if result.ID != test.ExpectedResult.ID {
				t.Errorf("expected order id is %v, but got %v", test.ExpectedResult.ID, result.ID)
			}
			if !slices.Equal(result.Items, test.ExpectedResult.Items) {
				t.Errorf("created items not matched, expected is %v, but got %v", test.ExpectedResult.Items, result.Items)
			}
			if test.ExpectedErr == nil {
				ordersRepoMock.AssertExpectations(t)
				pubSubMock.AssertExpectations(t)
			}
		})
	}
}
