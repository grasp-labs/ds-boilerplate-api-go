package models

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRoutes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock instance of the Controller interface
	mockController := NewMockController(ctrl)

	// Set up the expected behavior for GetRoutes
	mockRoutes := Routes{}
	mockController.EXPECT().
		GetRoutes().
		Return(mockRoutes)

	result := mockController.GetRoutes()

	assert.Equal(t, mockRoutes, result)
}
