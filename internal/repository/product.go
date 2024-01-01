package repository

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// model
type Product struct {
	gorm.Model
	Name     string `json:"name" binding:"required"`
	Quantity int    `json:"quantity" binding:"required"`
}

type ProductRepository interface {
	GetAllProduct() ([]Product, error)
	GetById(id int) (Product, error)
	AddProduct(product *Product) error
	UpdateProduct(updateproduct Product) (Product, error)
	Delete(id int) error
}

type concreteProductRepository struct {
	db *gorm.DB
}

// interface can store the pointer of struct or struct
// so we cannot return the pointer of interface
func NewProductRepository() ProductRepository {
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

	// update product model
	db.AutoMigrate(&Product{})

	return &concreteProductRepository{
		db: db,
	}
}

func (repo *concreteProductRepository) AddProduct(product *Product) error {
	return repo.db.Create(&product).Error
}

func (repo *concreteProductRepository) GetAllProduct() ([]Product, error) {
	var products []Product
	if err := repo.db.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (repo *concreteProductRepository) GetById(id int) (Product, error) {
	var product Product
	if err := repo.db.First(&product, id).Error; err != nil {
		return product, err
	}
	return product, nil
}

func (repo *concreteProductRepository) UpdateProduct(updateproduct Product) (Product, error) {
	var product Product
	if err := repo.db.First(&product, updateproduct.ID).Error; err != nil {
		return product, err
	}

	if err := repo.db.Model(&product).Updates(updateproduct).Error; err != nil {
		return product, err
	}
	return product, nil
}

func (repo *concreteProductRepository) Delete(id int) error {
	return repo.db.Delete(&Product{}, id).Error
}
