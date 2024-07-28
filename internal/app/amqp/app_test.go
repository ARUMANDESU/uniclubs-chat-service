package amqpapp

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/ARUMANDESU/uniclubs-comments-service/internal/app/amqp/mocks"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/domain"
	"github.com/ARUMANDESU/uniclubs-comments-service/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type Suite struct {
	App             *App
	mockAmqp        *mocks.Amqp
	mockUserService *mocks.UserService
}

func NewSuite(t *testing.T) *Suite {
	s := &Suite{
		mockAmqp:        mocks.NewAmqp(t),
		mockUserService: mocks.NewUserService(t),
	}
	s.App = New(logger.Plug(), s.mockUserService, s.mockAmqp)
	return s
}

func TestApp_HandleUpdateUser(t *testing.T) {
	t.Run("successful user update", func(t *testing.T) {
		suite := NewSuite(t)
		user := domain.User{ID: 1, FirstName: "John Doe"}
		msg := amqp.Delivery{Body: []byte(`{"id":1,"first_name":"John Doe"}`)}

		suite.mockUserService.On("Update", mock.Anything, user).Return(nil)

		err := suite.App.HandleUpdateUser(msg)
		assert.NoError(t, err)
		suite.mockUserService.AssertExpectations(t)
	})

	t.Run("failed JSON unmarshal", func(t *testing.T) {
		suite := NewSuite(t)
		msg := amqp.Delivery{Body: []byte(`invalid json`)}

		err := suite.App.HandleUpdateUser(msg)
		assert.Error(t, err)
		suite.mockUserService.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("failed user update", func(t *testing.T) {
		suite := NewSuite(t)

		user := domain.User{ID: 1, FirstName: "John Doe"}
		msg := amqp.Delivery{Body: []byte(`{"id":1,"first_name":"John Doe"}`)}

		suite.mockUserService.On("Update", mock.Anything, user).Return(fmt.Errorf("update error"))

		err := suite.App.HandleUpdateUser(msg)
		assert.Error(t, err)
		suite.mockUserService.AssertExpectations(t)
	})
}

func TestApp_ConsumeMessages(t *testing.T) {

	t.Run("successful message consumption", func(t *testing.T) {
		suite := NewSuite(t)
		var wg sync.WaitGroup
		wg.Add(1)

		queue := "testQueue"
		routingKey := "testRoutingKey"
		handler := func(msg amqp.Delivery) error {
			wg.Done()
			return nil
		}

		suite.mockAmqp.On("Consume", queue, routingKey, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			handler := args.Get(2).(func(amqp.Delivery) error)
			go handler(amqp.Delivery{})
		})

		suite.App.consumeMessages(queue, routingKey, handler)
		wg.Wait()
		suite.mockAmqp.AssertExpectations(t)
	})

	t.Run("failed message consumption", func(t *testing.T) {
		suite := NewSuite(t)
		var wg sync.WaitGroup
		wg.Add(1)

		queue := "testQueue"
		routingKey := "testRoutingKey"
		handler := func(msg amqp.Delivery) error {
			wg.Done()
			return nil
		}

		suite.mockAmqp.On("Consume", queue, routingKey, mock.Anything).Return(fmt.Errorf("consume error")).Run(func(args mock.Arguments) {
			handler := args.Get(2).(func(amqp.Delivery) error)
			go handler(amqp.Delivery{})
		})

		suite.App.consumeMessages(queue, routingKey, handler)
		wg.Wait()
		suite.mockAmqp.AssertExpectations(t)
	})
}

func TestApp_Stop(t *testing.T) {
	t.Run("successful stop", func(t *testing.T) {
		suite := NewSuite(t)

		suite.mockAmqp.On("Close").Return(nil)

		err := suite.App.Stop(context.Background())
		assert.NoError(t, err)
		suite.mockAmqp.AssertExpectations(t)
	})

	t.Run("failed stop", func(t *testing.T) {
		suite := NewSuite(t)

		suite.mockAmqp.On("Close").Return(fmt.Errorf("close error"))

		err := suite.App.Stop(context.Background())
		assert.Error(t, err)
		suite.mockAmqp.AssertExpectations(t)
	})
}
