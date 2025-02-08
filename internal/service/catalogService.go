package service

import (
	"errors"
	"go-ecommerce-app/config"
	"go-ecommerce-app/internal/domain"
	"go-ecommerce-app/internal/dto"
	"go-ecommerce-app/internal/helper"
	"go-ecommerce-app/internal/repository"
)

type CatalogService struct {
	Repo   repository.CatalogRepository
	Auth   helper.Auth
	Config config.AppConfig
}

func (s CatalogService) CreateCategory(input dto.CreateCategoryRequest) error {
	err := s.Repo.CreateCategory(&domain.Category{
		Name:         input.Name,
		ImageUrl:     input.ImageUrl,
		DisplayOrder: input.DisplayOrder,
	})
	return err
}

func (s CatalogService) EditCategory(id int, input dto.CreateCategoryRequest) (*domain.Category, error) {
	existingCategory, err := s.Repo.FindCategoryById(id)
	if err != nil {
		return nil, errors.New("category not exists")
	}

	if len(input.Name) > 0 {
		existingCategory.Name = input.Name
	}

	if input.ParendId > 0 {
		existingCategory.ParentId = input.ParendId
	}

	if len(input.ImageUrl) > 0 {
		existingCategory.ImageUrl = input.ImageUrl
	}

	if input.DisplayOrder > 0 {
		existingCategory.DisplayOrder = input.DisplayOrder
	}

	categoryUpdated, err := s.Repo.EditCategory(existingCategory)
	if err != nil {
		return nil, errors.New("unable to update category")
	}

	return categoryUpdated, nil
}

func (s CatalogService) DeleteCategory(id int) error {
	err := s.Repo.DeleteCategory(id)
	if err != nil {
		return errors.New("unable to delete category")
	}
	return nil
}

func (s CatalogService) GetCategories() ([]*domain.Category, error) {
	categories, err := s.Repo.FindCategories()
	if err != nil {
		return nil, errors.New("unable to get categories")
	}

	return categories, nil
}

func (s CatalogService) GetCategory(id int) (*domain.Category, error) {
	cat, err := s.Repo.FindCategoryById(id)

	if err != nil {
		return nil, errors.New("category not exists")
	}

	return cat, nil
}

func (s CatalogService) CreateProduct(input dto.CreateProductRequest, user domain.User) error {
	err := s.Repo.CreateProduct(&domain.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		CategoryId:  input.CategoryId,
		ImageUrl:    input.ImageUrl,
		UserId:      int(user.ID),
		Stock:       uint(input.Stock),
	})
	return err
}

func (s CatalogService) EditProduct(id int, input dto.CreateProductRequest, user domain.User) (*domain.Product, error) {
	existsProduct, err := s.Repo.FindProductById(id)
	if err != nil {
		return nil, errors.New("product does not exists")
	}

	// verify product owner
	if existsProduct.UserId != int(user.ID) {
		return nil, errors.New("you dont have manage rights of this product")
	}

	if len(input.Name) > 0 {
		existsProduct.Name = input.Name
	}

	if len(input.Description) > 0 {
		existsProduct.Description = input.Description
	}

	if input.Price > 0 {
		existsProduct.Price = input.Price
	}

	if input.CategoryId > 0 {
		existsProduct.CategoryId = input.CategoryId
	}

	updateProduct, err := s.Repo.EditProduct(existsProduct)
	if err != nil {
		return nil, errors.New("unable to update product")
	}

	return updateProduct, nil
}

func (s CatalogService) DeleteProduct(id int, user domain.User) error {
	existsProduct, err := s.Repo.FindProductById(id)
	if err != nil {
		return errors.New("product does not exists")
	}

	// verify product owner
	if existsProduct.UserId != int(user.ID) {
		return errors.New("you dont have manage rights of this product")
	}

	err = s.Repo.DeleteProduct(existsProduct)
	if err != nil {
		return errors.New("product cannot delete")
	}

	return nil
}

func (s CatalogService) GetProducts() ([]*domain.Product, error) {
	products, err := s.Repo.FindProducts()
	if err != nil {
		return nil, errors.New("products not exists")
	}
	return products, nil
}

func (s CatalogService) GetProductById(id int) (*domain.Product, error) {
	product, err := s.Repo.FindProductById(id)
	if err != nil {
		return nil, errors.New("product not exists")
	}
	return product, nil
}

func (s CatalogService) GetSellerProducts(id int) ([]*domain.Product, error) {
	products, err := s.Repo.FindSellerProducts(id)
	if err != nil {
		return nil, errors.New("products does not exists")
	}

	return products, err
}

func (s CatalogService) UpdateProductStock(e domain.Product) (*domain.Product, error) {
	product, err := s.Repo.FindProductById(int(e.ID))
	if err != nil {
		return nil, errors.New("product not exists")
	}

	// verify product owner
	if product.UserId != e.UserId {
		return nil, errors.New("you dont have manage rights of this product")
	}

	product.Stock = e.Stock
	editProduct, err := s.Repo.EditProduct(product)
	if err != nil {
		return nil, errors.New("unable to update stock of product")
	}

	return editProduct, nil
}
