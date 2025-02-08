package repository

import (
	"errors"
	"go-ecommerce-app/internal/domain"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository interface {
	CreateUser(u domain.User) (domain.User, error)
	FindUser(email string) (domain.User, error)
	FindUserById(id uint) (domain.User, error)
	UpdateUser(id uint, u domain.User) (domain.User, error)
	CreateBankAccount(e domain.BackAccount) error

	//cart
	FindCartItems(uId uint) ([]domain.Cart, error)
	FindCartItem(uId uint, pId uint) (domain.Cart, error)
	CreateCart(c domain.Cart) error
	UpdateCart(c domain.Cart) error
	DeleteCartById(id uint) error
	DeleteCartItems(uId uint) error

	// Order
	CreateOrder(o domain.Order) error
	FindOrders(uId uint) ([]domain.Order, error)
	FindOrderById(id uint, uId uint) (domain.Order, error)

	// profile
	CreateProfile(e domain.Address) error
	UpdateProfile(e domain.Address) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// CreateOrder implements UserRepository.
func (r *userRepository) CreateOrder(o domain.Order) error {
	err := r.db.Create(&o).Error
	if err != nil {
		log.Printf("error on creating order %v", err)
		return errors.New("failed to create order")
	}
	return nil
}

// FindOrders implements UserRepository.
func (r *userRepository) FindOrders(uId uint) ([]domain.Order, error) {
	var orders []domain.Order
	err := r.db.Where("user_id=?", uId).Find(&orders).Error
	if err != nil {
		log.Printf("error on fetching orders %v", err)
		return nil, errors.New("failed to fetching orders")
	}

	return orders, nil
}

// FindOrderById implements UserRepository.
func (r *userRepository) FindOrderById(id uint, uId uint) (domain.Order, error) {
	var order domain.Order
	err := r.db.Preload("Items").Where("id=? AND user_id=?", id, uId).First(&order).Error
	if err != nil {
		log.Printf("error on fetching order by id %v", err)
		return domain.Order{}, errors.New("failed to fetching order by id")
	}
	return order, nil
}

// CreateProfile implements UserRepository.
func (r *userRepository) CreateProfile(e domain.Address) error {
	err := r.db.Create(&e).Error
	if err != nil {
		log.Printf("error on creating profile with address %v", err)
		return errors.New("failed to create profile")
	}
	return nil
}

// UpdateProfile implements UserRepository.
func (r *userRepository) UpdateProfile(e domain.Address) error {
	err := r.db.Where("user_id=?", e.UserId).Updates(e).Error
	if err != nil {
		log.Printf("error on update profile with address %v", err)
		return errors.New("failed to update profile")
	}
	return err
}

// FindCartItems implements UserRepository.
func (r *userRepository) FindCartItems(uId uint) ([]domain.Cart, error) {
	var carts []domain.Cart
	err := r.db.Where("user_id=?", uId).Find(&carts).Error
	return carts, err
}

// FindCartItem implements UserRepository.
func (r *userRepository) FindCartItem(uId uint, pId uint) (domain.Cart, error) {
	cartItem := domain.Cart{}
	err := r.db.Where("user_id=? AND product_id=?", uId, pId).First(&cartItem).Error
	return cartItem, err
}

// CreateCart implements UserRepository.
func (r *userRepository) CreateCart(c domain.Cart) error {
	return r.db.Create(&c).Error
}

// UpdateCart implements UserRepository.
func (r *userRepository) UpdateCart(c domain.Cart) error {
	var cart domain.Cart
	err := r.db.Model(&cart).Clauses(clause.Returning{}).Where("id=?", c.ID).Updates(c).Error
	return err
}

// DeleteCartById implements UserRepository.
func (r *userRepository) DeleteCartById(id uint) error {
	err := r.db.Delete(&domain.Cart{}, id).Error
	return err
}

// DeleteCartItems implements UserRepository.
func (r *userRepository) DeleteCartItems(uId uint) error {
	err := r.db.Where("user_id=?", uId).Delete(&domain.Cart{}).Error
	return err
}

// CreateBankAccount implements UserRepository.
func (r userRepository) CreateBankAccount(e domain.BackAccount) error {
	return r.db.Create(&e).Error
}

func (r userRepository) CreateUser(usr domain.User) (domain.User, error) {
	err := r.db.Create(&usr).Error

	if err != nil {
		log.Printf("create user error %v", err)
		return domain.User{}, errors.New("failed to create user")
	}

	return usr, nil
}

func (r userRepository) FindUser(email string) (domain.User, error) {
	var user domain.User

	err := r.db.Preload("Address").First(&user, "email=?", email).Error

	if err != nil {
		log.Printf("find user error %v", err)
		return domain.User{}, errors.New("user does not exists")
	}

	return user, nil
}

func (r userRepository) FindUserById(id uint) (domain.User, error) {
	var user domain.User

	err := r.db.Preload("Address").Preload("Cart").Preload("Orders").First(&user, id).Error

	if err != nil {
		log.Printf("find user by id error %v", err)
		return domain.User{}, errors.New("user does not exists")
	}

	return user, nil
}

func (r userRepository) UpdateUser(id uint, u domain.User) (domain.User, error) {
	var user domain.User

	err := r.db.Model(&user).Clauses(clause.Returning{}).Where("id=?", id).Updates(u).Error

	if err != nil {
		log.Printf("error on update  %v", err)
		return domain.User{}, errors.New("failed update user")
	}
	return user, nil
}
