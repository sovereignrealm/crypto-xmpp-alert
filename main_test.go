package main

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockFileReader struct {
	FileContent []byte
	ReadError   error
}

func (m *MockFileReader) ReadFile(filename string) ([]byte, error) {
	if m.ReadError != nil {
		return nil, m.ReadError
	}
	return m.FileContent, nil
}

type MockFileWriter struct {
	WriteError error
}

func (m *MockFileWriter) WriteFile(filename string, data []byte, perm fs.FileMode) error {
	if m.WriteError != nil {
		return m.WriteError
	}
	return nil
}

func TestCalculateTotalValues(t *testing.T) {
	transactions := []CryptoTransaction{
		{PurchasedPrice: 100.0, CryptoAmount: 1.0},
		{PurchasedPrice: 200.0, CryptoAmount: 2.0},
		{PurchasedPrice: 300.0, CryptoAmount: 3.0},
	}

	totalAmountInvested, totalCryptoAsset := calculateTotalValues(transactions)

	assert.Equal(t, 600.0, totalAmountInvested)
	assert.Equal(t, 6.0, totalCryptoAsset)
}

func TestCalculatePercentageChange(t *testing.T) {
	moneyNow := 1000.0
	totalAmountInvested := 2000.0
	totalGained := calculatePercentageChange(totalAmountInvested, moneyNow)
	assert.Equal(t, -50.0, totalGained)
	moneyNow = 3000.0
	totalAmountInvested = 1000.0
	totalGained = calculatePercentageChange(totalAmountInvested, moneyNow)
	assert.Equal(t, 200.0, totalGained)
}

func TestIsFirstTime(t *testing.T) {
	data := []byte("false")
	result := isFirstTime(data)
	assert.False(t, result)

	data = []byte("true")
	result = isFirstTime(data)
	assert.True(t, result)

	data = []byte("invalid")
	result = isFirstTime(data)
	assert.True(t, result) // Default to true for invalid data
}

func TestReadFile(t *testing.T) {
	mockReader := &MockFileReader{
		FileContent: []byte("Mocked file content"),
	}
	data, err := readFile(mockReader, "Bitcoin")
	assert.NoError(t, err)
	assert.NotNil(t, data)
}

func TestWriteFile(t *testing.T) {
	mockFileWriter := &MockFileWriter{}
	assert.NotPanics(t, func() {
		writeFile(mockFileWriter, "Bitcoin")
	}, "The function should not panic")
}
