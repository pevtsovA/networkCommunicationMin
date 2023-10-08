package contract

import (
	"networkCommunicationMin/models"
)

type Storage interface {
	GetUsers() (map[int]models.User, error)
	GetUserById(id int) (models.User, error)
	SaveUser(name string, age int, friends []int64) (int, error)
	UpdateUser(id int, name string, age int, friends []int64) error
	DeleteUser(id int) error
}
