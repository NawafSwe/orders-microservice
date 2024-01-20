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
	"reflect"
	"slices"
	"testing"
)

func TestPlaceOrderUseCase(t *testing.T) {
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

func TestUpdateOrderStatusUseCase(t *testing.T) {
	tests := map[string]struct {
		Description string
		Input       struct {
			OrderId int64
			Status  string
		}
		ExpectedResult models.Order
		ExpectedErr    error
	}{
		"SuccessfullyUpdateOrderStatusFromNewToApproved": {
			Description: "Should successfully update order status from new to approved",
			Input: struct {
				OrderId int64
				Status  string
			}{
				OrderId: 1,
				Status:  "Approved",
			},
			ExpectedResult: models.Order{Model: gorm.Model{ID: 1}, Status: "Approved"},
			ExpectedErr:    nil,
		},

		"FailedToUpdateOrderStatusDueInvalidIdPassed": {
			Description: "Should fail update order status due invalid order id passed",
			Input: struct {
				OrderId int64
				Status  string
			}{
				OrderId: -100,
				Status:  "Approved",
			},
			ExpectedResult: models.Order{},
			ExpectedErr:    errors.New("invalid order id"),
		},

		"FailedToUpdateOrderStatusDueGivenInvalidStatus": {
			Description: "Should fail update order status, due invalid order status passed",
			Input: struct {
				OrderId int64
				Status  string
			}{
				OrderId: 1,
			},
			ExpectedResult: models.Order{},
			ExpectedErr:    models.InvalidStatusChangeErr{Message: "given status '' is invalid"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Logf("running %v", name)
			pubSubMock := messagesMock.NewMockMessageService(t)
			ordersRepoMock := ordersMock.NewMockOrderRepo(t)
			ordersUseCase := useCases.NewOrderUseCase(ordersRepoMock, pubSubMock)
			if test.ExpectedErr == nil {

				ordersRepoMock.On("UpdateOrderStatus", mock.Anything, test.Input.OrderId, test.Input.Status).Return(
					models.Order{
						Model:  gorm.Model{ID: uint(test.Input.OrderId)},
						Status: "Approved",
					}, nil)
				pubSubMock.On("PublishAsync", mock.Anything, "orderStatusChanged", mock.Anything).Return(nil)
			} else {
				var expectedErr models.InvalidStatusChangeErr
				if !errors.As(test.ExpectedErr, &expectedErr) {
					ordersRepoMock.On("UpdateOrderStatus", mock.Anything, test.Input.OrderId, test.Input.Status).Return(
						models.Order{
							Model:  gorm.Model{ID: uint(test.Input.OrderId)},
							Status: test.Input.Status,
						}, test.ExpectedErr)
				}
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			result, err := ordersUseCase.UpdateOrderStatus(ctx, test.Input.OrderId, test.Input.Status)
			if test.ExpectedErr == nil {
				ordersRepoMock.AssertExpectations(t)
				pubSubMock.AssertExpectations(t)

			} else {
				pubSubMock.AssertNumberOfCalls(t, "PublishAsync", 0)
			}
			if err == nil && test.ExpectedErr != nil {
				t.Errorf("expected error to be %v, but got %v", test.ExpectedErr, err)
			}

			if !reflect.DeepEqual(result, test.ExpectedResult) {
				t.Errorf("expected result to be %v, but got %v", test.ExpectedResult, result)
			}
		})
	}
}

//Handlers already been tested because I am reusing the functionality
//func TestHandleOrderApprovalUseCase(t *testing.T) {
//	tests := map[string]struct {
//		Description    string
//		ExpectedResult models.Order
//		ExpectedErr    error
//	}{
//		"SuccessfullyHandleOrderApproval": {
//			Description: "Should successfully handle order approval from a published message on topic orderApproved",
//			ExpectedResult: models.Order{
//				Model:      gorm.Model{ID: 1},
//				CustomerId: 1,
//				GrandTotal: 10,
//				Status:     "Approved",
//				Items: []models.OrderedItem{
//					{
//						Model:           gorm.Model{ID: 1},
//						OrderedItemId:   1,
//						Sku:             "pepsi12s",
//						Price:           10,
//						OrderedQuantity: 1,
//					},
//				},
//			},
//		},
//	}
//
//	for name, test := range tests {
//		t.Run(name, func(t *testing.T) {
//			t.Logf("Running %v ", name)
//			pubSubMocks := messagesMock.NewMockMessageService(t)
//			ordersRepoMocks := ordersMock.NewMockOrderRepo(t)
//			orderUseCase := useCases.NewOrderUseCase(ordersRepoMocks, pubSubMocks)
//
//			ctx, cancel:= context.WithCancel(context.Background())
//			defer cancel()
//			if test.ExpectedErr == nil {
//				ordersRepoMocks.On("Update", mock.Anything)
//			}
//			orderUseCase.HandleOrderApproval(ctx)
//
//
//		})
//	}
//}
