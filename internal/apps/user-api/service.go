package user_api

import (
	"OrderUserProject/internal/models"
	"OrderUserProject/internal/repository"
	"errors"
	"github.com/google/uuid"
	"time"
)

type UserService struct {
	Repository *repository.UserRepository
}

func NewUserService(repository *repository.UserRepository) *UserService {
	userService := &UserService{Repository: repository}
	return userService
}

type IUserService interface {
	GetAll() ([]models.User, error)
	GetUserById(id string) (models.User, error)
	Insert(user models.User) (models.User, error)
	Update(user models.User) (bool, error)
	Delete(id string) (bool, error)
	InvoiceRegularAddressCheck(user models.User) models.User
}

func (b *UserService) GetAll() ([]models.User, error) {
	result, err := b.Repository.GetAll()

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (b *UserService) GetUserById(id string) (models.User, error) {

	result, err := b.Repository.GetUserById(id)

	if err != nil {
		return result, err
	}

	return result, nil
}

func (b *UserService) Insert(user models.User) (models.User, error) {

	// Create id and created date value
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt

	result, err := b.Repository.Insert(user)

	if err != nil || result == false {
		return user, err
	}

	return user, nil
}

func (b *UserService) Update(user models.User) (bool, error) {
	// to create updated date value
	user.UpdatedAt = time.Now()

	result, err := b.Repository.Update(user)

	if err != nil || result == false {
		return false, err
	}

	return true, nil
}

func (b *UserService) Delete(id string) (bool, error) {
	result, err := b.Repository.Delete(id)

	if err != nil || result == false {
		return false, err
	}

	return true, nil
}

func (b *UserService) InvoiceRegularAddressCheck(user models.User) (models.User, error) {
	// Invoice and regular addresses check
	hasDefaultInvoice := false
	hasDefaultRegular := false
	for _, addressRequest := range user.Addresses {
		if addressRequest.Default.IsDefaultRegularAddress {
			hasDefaultRegular = true
		}
		if addressRequest.Default.IsDefaultInvoiceAddress {
			hasDefaultInvoice = true
		}
	}

	if !hasDefaultInvoice {
		return user, errors.New("At least one address chosen as a default invoice address!")
	}

	if !hasDefaultRegular {
		return user, errors.New("At least one address chosen as a default regular address!")
	}

	return user, nil
}
