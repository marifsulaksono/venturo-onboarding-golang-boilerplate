package models

import (
	"simple-crud-rnd/structs"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var productID uuid.UUID
var productDetailID uuid.UUID

// Setup mock DB
func setupMockDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: mockDB,
	}), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}
	return db, mock, nil
}

// Test Create Product
func TestCreateProduct_Success(t *testing.T) {
	// Setup
	db, mock, err := setupMockDB()
	assert.NoError(t, err)

	// Initialize the model
	pm := ProductModel{db: db}

	// Generate UUIDs for product and details
	productID = uuid.New()
	productDetailID = uuid.New()

	// Sample payload
	payload := structs.ProductWithDetailCreateOrUpdate{
		Product: structs.Product{
			ID:          productID,
			Name:        "Test Product",
			Price:       100.00,
			Description: "Test description",
		},
		Details: []structs.ProductDetailRequest{
			{
				ID:          productDetailID,
				ProductID:   productID,
				Type:        "Topping",
				Description: "Test Detail",
				Price:       10.00,
				IsAdded:     true,
			},
		},
	}

	// Expected SQL queries and their mock responses
	mock.ExpectBegin()

	// Mock the product create query
	mock.ExpectExec(`INSERT INTO "m_product"`).WithArgs(productID, payload.Name, payload.Price, payload.Description).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Mock the product detail create query
	mock.ExpectExec(`INSERT INTO "m_product_detail"`).WithArgs(productDetailID, productID, "Topping", "Test Detail", 10.00).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Mock commit
	mock.ExpectCommit()

	// Call the method
	createdProduct, err := pm.Create(&payload)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, createdProduct)
	assert.Equal(t, "Test Product", createdProduct.Name)

	// Ensure all expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

// Test Update Product
func TestUpdateProduct_Success(t *testing.T) {
	// Setup
	db, mock, err := setupMockDB()
	assert.NoError(t, err)

	// Initialize the model
	pm := ProductModel{db: db}

	// Ensure productID and productDetailID are set from the Create test
	assert.NotZero(t, productID)
	assert.NotZero(t, productDetailID)

	// Sample payload for update
	payload := structs.ProductWithDetailCreateOrUpdate{
		Product: structs.Product{
			ID:          productID,
			Name:        "Updated Product",
			Price:       150.00,
			Description: "Updated description",
		},
		Details: []structs.ProductDetailRequest{
			{
				ID:          productDetailID,
				ProductID:   productID,
				Type:        "Topping",
				Description: "Updated Detail",
				Price:       15.00,
				IsUpdated:   true,
			},
		},
	}

	// Expected SQL queries and their mock responses
	mock.ExpectBegin()

	// Mock the product update query
	mock.ExpectExec(`UPDATE "m_product" SET`).WithArgs(payload.Name, payload.Price, payload.Description, productID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Mock the product detail update query
	mock.ExpectExec(`UPDATE "m_product_detail" SET`).WithArgs(payload.Details[0].Type, payload.Details[0].Description, payload.Details[0].Price, productDetailID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Mock commit
	mock.ExpectCommit()

	// Call the update method
	updatedProduct, err := pm.Update(&payload)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, updatedProduct)
	assert.Equal(t, "Updated Product", updatedProduct.Name)

	// Ensure all expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}
