package repository

import (
	"bitcoin-exchange-rate/internal/model"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestCreateSubscriber(t *testing.T) {
	err := godotenv.Load("../../.env.test")
	require.NoError(t, err, "Failed to load .env.test file")

	testFilePath := os.Getenv("REPO_TEST_FILE_PATH")
	subscriberRepository := NewSubscriberFileRepository(testFilePath)

	testEmail := "test@gmail.com"
	testSubscriber := model.NewSubscriber(testEmail)

	require.NoError(t, subscriberRepository.Create(testSubscriber), "failed to add subscriber '%s' to file '%s'", testEmail, testFilePath)

	testSubscribers, err := subscriberRepository.GetAll()
	require.NoError(t, err, "failed to get all subscribers '%s' from file '%s': %v", testEmail, testFilePath)

	assert.Equal(t, testEmail, testSubscribers[1].GetEmail(), "subscriber '%s' doesn't exist in file '%s'", testEmail, testFilePath)

	err = ClearFile(testFilePath)
	require.NoError(t, err, "failed to clear file '%s': %v", testFilePath, err)
}
