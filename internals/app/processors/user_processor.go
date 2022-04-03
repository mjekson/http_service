package processors

import (
	"errors"

	"github.com/mjekson/http_service/internals/app/db"
	"github.com/mjekson/http_service/internals/app/models"
)

type UsersProcessor struct {
	storage *db.UsersStorage
}

func NewUsersProcessor(storage *db.UsersStorage) *UsersProcessor {
	processor := new(UsersProcessor)
	processor.storage = storage
	return processor
}

func (processor *UsersProcessor) CreateUser(user models.User) error {
	if user.Name == "" {
		return errors.New("Name should not be empty")
	}

	return processor.storage.CreateUser(user)

}

func (processor *UsersProcessor) FindUser(id int64) (models.User, error) {
	user := processor.storage.GetUserById(id)
	if user.Id != id {
		return user, errors.New("user not found")
	}
	return user, nil
}

func (processor *UsersProcessor) ListUsers(nameFilter string) ([]models.User, error) {
	return processor.storage.GetUsersList(nameFilter), nil
}
