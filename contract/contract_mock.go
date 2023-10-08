package contract

import (
	"fmt"
	"networkCommunicationMin/models"
)

type MockStorage struct{}

func (m *MockStorage) GetUsers() (map[int]models.User, error) {
	return map[int]models.User{
		1: {
			ID:      1,
			Name:    "Tom",
			Age:     30,
			Friends: nil,
		},
		2: {
			ID:      2,
			Name:    "Sara",
			Age:     28,
			Friends: nil,
		},
	}, nil
}

func (m *MockStorage) GetUserById(id int) (models.User, error) {
	if id == 3 {
		return models.User{
			ID:      3,
			Name:    "Sam",
			Age:     19,
			Friends: nil,
		}, nil
	}
	return models.User{}, fmt.Errorf("user does not exist")
}

func (m *MockStorage) SaveUser(name string, age int, friends []int64) (int, error) {
	u := models.User{
		ID:      1,
		Name:    "Oleg",
		Age:     22,
		Friends: nil,
	}
	if u.Name == name && u.Age == age && friends == nil {
		return u.ID, nil
	}
	return -1, nil
}

func (m *MockStorage) UpdateUser(id int, name string, age int, friends []int64) error {
	panic("implement me")
}

func (m *MockStorage) DeleteUser(id int) error {
	panic("implement me")
}
