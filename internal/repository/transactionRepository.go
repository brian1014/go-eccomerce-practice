package repository

import (
	"go-ecommerce-app/internal/domain"
	"go-ecommerce-app/internal/dto"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreatePayment(payment *domain.Payment) error
	FindInitialPayment(uId int) (*domain.Payment, error)
	UpdatePayment(payment *domain.Payment) error
	FindOrders(uId uint) ([]domain.OrderItem, error)
	FindOrderById(uId uint, id uint) (dto.SellerOrderDetails, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}

// UpdatePayment implements TransactionRepository.
func (t *transactionRepository) UpdatePayment(payment *domain.Payment) error {
	return t.db.Save(payment).Error
}

// FindPayment implements TransactionRepository.
func (t *transactionRepository) FindInitialPayment(uId int) (*domain.Payment, error) {
	var payment *domain.Payment
	err := t.db.First(&payment, "user_id=? AND status=?", uId, "initial").Order("created_at desc").Error
	return payment, err
}

func (t *transactionRepository) CreatePayment(payment *domain.Payment) error {
	return t.db.Create(payment).Error
}

func (t *transactionRepository) FindOrders(uId uint) ([]domain.OrderItem, error) {
	return []domain.OrderItem{}, nil
}

func (t *transactionRepository) FindOrderById(uId uint, id uint) (dto.SellerOrderDetails, error) {
	return dto.SellerOrderDetails{}, nil
}
