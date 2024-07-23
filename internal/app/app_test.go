package app

import (
	"context"
	"errors"
	"github.com/ARUMANDESU/uniclubs-comments-service/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thejerf/slogassert"
	"log/slog"
	"testing"
)

type MockStarter struct {
	mock.Mock
}

func (m *MockStarter) Start(ctx context.Context, callback func(error)) {
	args := m.Called(ctx)
	callback(args.Error(0))
}

type MockStopper struct {
	mock.Mock
}

func (m *MockStopper) Stop(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestApp_StartsAllStarters(t *testing.T) {
	mockStarter := new(MockStarter)
	mockStarter.On("Start", mock.Anything).Return(nil)

	appInstance := &App{
		log:      logger.Plug(),
		starters: []Starter{mockStarter},
	}

	appInstance.Start()

	mockStarter.AssertCalled(t, "Start", mock.Anything)
}

func TestApp_StopsAllStoppers(t *testing.T) {
	mockStopper := new(MockStopper)
	mockStopper.On("Stop", mock.Anything).Return(nil)

	appInstance := &App{
		log:      logger.Plug(),
		stoppers: []Stopper{mockStopper},
	}

	appInstance.Stop()

	mockStopper.AssertCalled(t, "Stop", mock.Anything)
}

func TestApp_StartFailsWhenStarterFails(t *testing.T) {
	mockStarter := new(MockStarter)
	mockStarter.On("Start", mock.Anything).Return(errors.New("start error"))

	handler := slogassert.New(t, slog.LevelWarn, nil)
	log := slog.New(handler)

	appInstance := &App{
		log:      log,
		starters: []Starter{mockStarter},
	}

	err := appInstance.Start()

	assert.Error(t, err)
	handler.AssertSomeMessage("failed to start service")
}

func TestApp_StopFailsWhenStopperFails(t *testing.T) {
	mockStopper := new(MockStopper)
	mockStopper.On("Stop", mock.Anything).Return(errors.New("stop error"))

	handler := slogassert.New(t, slog.LevelWarn, nil)
	log := slog.New(handler)

	appInstance := &App{
		log:      log,
		stoppers: []Stopper{mockStopper},
	}

	appInstance.Stop()

	handler.AssertSomeMessage("failed to stop service")
}
