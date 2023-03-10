package user_api

import (
	"OrderUserProject/internal/models"
	"OrderUserProject/internal/repository"
	"github.com/google/uuid"
	"time"
)

type UserService struct {
	Repository repository.IUserRepository
}

func NewUserService(repository repository.IUserRepository) *UserService {
	userService := &UserService{Repository: repository}
	return userService
}

type IUserService interface {
	GetAll() ([]models.User, error)
	GetUserById(id string) (models.User, error)
	Insert(user models.User) (models.User, error)
	Update(user models.User) (bool, error)
	Delete(id string) (bool, error)
}

func (b UserService) GetAll() ([]models.User, error) {
	result, err := b.Repository.GetAll()

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (b UserService) GetUserById(id string) (models.User, error) {

	result, err := b.Repository.GetUserById(id)

	if err != nil {
		return result, err
	}

	return result, nil
}

func (b UserService) Insert(user models.User) (models.User, error) {

	// to create id and created date value
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()

	result, err := b.Repository.Insert(user)

	if err != nil || result == false {
		return user, err
	}

	return user, nil
}

func (b UserService) Update(user models.User) (bool, error) {
	// to create updated date value
	user.UpdatedAt = time.Now()

	result, err := b.Repository.Update(user)

	if err != nil || result == false {
		return false, err
	}

	return true, nil
}

func (b UserService) Delete(id string) (bool, error) {
	result, err := b.Repository.Delete(id)

	if err != nil || result == false {
		return false, err
	}

	return true, nil
}
