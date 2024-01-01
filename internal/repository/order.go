package repository

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	user         User `gorm:"foreignKey:UserID"`
	UserID       uint
	OrderedItems []OrderedItem `gorm:"foreignKey:OrderID"`
}

type OrderedItem struct {
	gorm.Model
	// order Order `gorm:"foreignKey:OrderID"`
	OrderID   uint    `gorm:"not null"`
	product   Product `gorm:"foreignKey:ProductID"`
	ProductID uint    `gorm:"not null"`
	Quantity  int
}

type OrderRequest struct {
	CartItems []CartItem `json:"cart_items"`
}

type CartItem struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required"`
}

type OrderRepository interface {
	CreateOrder(uint, OrderRequest) (Order, error)
	// CreateOrderedItem(item CartItem) error
}

type concreteOrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository() OrderRepository {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to Database")
		return nil
	}

	// update order model
	db.AutoMigrate(&Order{}, &OrderedItem{})

	return &concreteOrderRepository{
		db: db,
	}
}

func (repo *concreteOrderRepository) CreateOrder(userID uint, orderRequest OrderRequest) (Order, error) {
	// check product id exists
	var order Order
	var products []Product
	productIDs := make([]uint, len(orderRequest.CartItems))
	for i := 0; i < len(orderRequest.CartItems); i++ {
		productIDs[i] = orderRequest.CartItems[i].ProductID
	}
	if err := repo.db.Find(&products, productIDs).Error; err != nil {
		return order, err
	}

	// create each item
	// for i := 0; i < len(orderRequest.CartItems); i++ {
	// 	repo.CreateOrderedItem(orderRequest.CartItems[i])
	// }
	orderedItems := make([]OrderedItem, len(products))
	for i := 0; i < len(products); i++ {
		orderedItems[i].ProductID = products[i].ID
		orderedItems[i].Quantity = orderRequest.CartItems[i].Quantity
	}
	order.OrderedItems = orderedItems
	order.UserID = userID
	return order, repo.db.Create(&order).Error
}

// func (repo *concreteOrderRepository) CreateOrderedItem(item CartItem) error {
// 	repo.db
// }
