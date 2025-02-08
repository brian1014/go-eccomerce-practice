package repository

import (
	"errors"
	"go-ecommerce-app/internal/domain"
	"log"

	"gorm.io/gorm"
)

type CatalogRepository interface {
	CreateCategory(c *domain.Category) error
	FindCategories() ([]*domain.Category, error)
	FindCategoryById(id int) (*domain.Category, error)
	EditCategory(c *domain.Category) (*domain.Category, error)
	DeleteCategory(id int) error

	CreateProduct(e *domain.Product) error
	FindProducts() ([]*domain.Product, error)
	FindProductById(id int) (*domain.Product, error)
	FindSellerProducts(id int) ([]*domain.Product, error)
	EditProduct(e *domain.Product) (*domain.Product, error)
	DeleteProduct(e *domain.Product) error
}

type catalogRepository struct {
	db *gorm.DB
}

func NewCatalogRepository(db *gorm.DB) CatalogRepository {
	return &catalogRepository{
		db: db,
	}
}

// CreateProduct implements CatalogRepository.
func (c *catalogRepository) CreateProduct(e *domain.Product) error {
	err := c.db.Model(&domain.Product{}).Create(e).Error
	if err != nil {
		log.Printf("err: %v", err)
		return errors.New("can not create product")
	}
	return nil
}

// FindProducts implements CatalogRepository.
func (c *catalogRepository) FindProducts() ([]*domain.Product, error) {
	var products []*domain.Product
	err := c.db.Find(&products).Error

	if err != nil {
		return nil, errors.New("products does not exists")
	}
	return products, nil
}

// FindProductById implements CatalogRepository.
func (c *catalogRepository) FindProductById(id int) (*domain.Product, error) {
	var product *domain.Product
	err := c.db.First(&product, id).Error

	if err != nil {
		return nil, errors.New("product does not exists")
	}
	return product, nil
}

// FindSellerProducts implements CatalogRepository.
func (c *catalogRepository) FindSellerProducts(id int) ([]*domain.Product, error) {
	var products []*domain.Product
	err := c.db.Where("user_id=?", id).Find(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}

// EditProduct implements CatalogRepository.
func (c *catalogRepository) EditProduct(e *domain.Product) (*domain.Product, error) {
	err := c.db.Save(&e).Error
	if err != nil {
		log.Printf("db_error: %v", err)
		return nil, errors.New("fail to update product")
	}
	return e, nil
}

// DeleteProduct implements CatalogRepository.
func (c *catalogRepository) DeleteProduct(e *domain.Product) error {
	err := c.db.Delete(&domain.Product{}, e.ID).Error
	if err != nil {
		return errors.New("product cannot delete")
	}
	return nil
}

// CreateCategory implements CatalogRepository.
func (c catalogRepository) CreateCategory(e *domain.Category) error {
	err := c.db.Create(&e).Error

	if err != nil {
		log.Printf("db_error: %v", err)
		return errors.New("create category failed")
	}
	return nil
}

// FindCategories implements CatalogRepository.
func (c catalogRepository) FindCategories() ([]*domain.Category, error) {
	var categories []*domain.Category

	err := c.db.Find(&categories).Error

	if err != nil {
		log.Printf("db_err: %v", err)
		return nil, err
	}

	return categories, nil
}

// FindCategoryById implements CatalogRepository.
func (c catalogRepository) FindCategoryById(id int) (*domain.Category, error) {
	var category *domain.Category
	err := c.db.First(&category, id).Error

	if err != nil {
		log.Printf("db_err: %v", err)
		return nil, errors.New("category does not exists")
	}
	return category, nil
}

// EditCategory implements CatalogRepository.
func (c catalogRepository) EditCategory(e *domain.Category) (*domain.Category, error) {
	err := c.db.Save(&e).Error

	if err != nil {
		log.Printf("db_err: %v", err)
		return nil, errors.New("failed to update category")
	}

	return e, nil
}

// DeleteCategory implements CatalogRepository.
func (c catalogRepository) DeleteCategory(id int) error {
	err := c.db.Delete(&domain.Category{}, id).Error

	if err != nil {
		log.Printf("unable to delete category")
		return err
	}

	return nil
}
