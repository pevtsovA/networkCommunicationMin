package dataBase

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"log"
	"networkCommunicationMin/models"
)

type Storage struct {
	DB *sql.DB
}

func (s *Storage) GetUsers() (map[int]models.User, error) {
	// DataReading - функция получения пользователей из БД
	rows, err := s.DB.Query("select * from users")
	if err != nil {
		log.Println("Error! failed to get users from database:", err.Error())
		return nil, fmt.Errorf("failed to get users from database: %w", err)
	}
	defer rows.Close()

	store := make(map[int]models.User)
	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.ID, &u.Name, &u.Age, pq.Array(&u.Friends))
		if err != nil {
			log.Println("Error! failed to get users:", err.Error())
			return nil, fmt.Errorf("failed to get users: %w", err)
		}
		store[u.ID] = u
	}
	return store, nil
}

func (s *Storage) SaveUser(name string, age int, friends []int64) (int, error) {
	// DataRecording - функция записи пользователя в БД
	var id int
	err := s.DB.QueryRow("insert into users (name, age, friends) values ($1, $2, $3) returning id", name, age, pq.Array(friends)).Scan(&id)
	if err != nil {
		log.Println("Error! failed to save users:", err.Error())
		return -1, fmt.Errorf("failed to save user: %w", err)
	}
	return id, nil
}

func (s *Storage) UpdateUser(id int, name string, age int, friends []int64) error {
	// DataUpdating - функция обновления пользователя в БД
	_, err := s.DB.Exec("update users set name = $1, age = $2, friends = $3 where id = $4", name, age, pq.Array(friends), id)
	if err != nil {
		log.Println("Error! failed to update users:", err.Error())
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (s *Storage) DeleteUser(id int) error {
	// DataDeleting - функция удаления пользователя в БД
	_, err := s.DB.Exec("delete from users where id = $1", id)
	if err != nil {
		log.Println("Error! failed to delete users:", err.Error())
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
