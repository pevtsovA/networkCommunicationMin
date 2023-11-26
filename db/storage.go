package db

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"networkCommunicationMin/models"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

type Storage struct {
	DB *sql.DB
}

func NewStorage(DB *sql.DB) *Storage {
	return &Storage{DB: DB}
}

func ConnectToBD(password string, isDocker bool) (*sql.DB, func()) {
	var connStr string
	if isDocker {
		connStr = fmt.Sprintf("host=db user=postgres password=%s dbname=socnetworkdb port=5432 sslmode=disable TimeZone=Europe/Moscow", password)
	} else {
		connStr = fmt.Sprintf("user=postgres password=%s dbname=socnetworkdb sslmode=disable", password)
	}

	dbConnect, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return dbConnect, func() { dbConnect.Close() }
}

func (s *Storage) GetUsers() (map[int]models.User, error) {
	// DataReading - функция получения пользователей из БД
	rows, err := s.DB.Query("select * from users")
	if err != nil {
		return nil, fmt.Errorf("method 'GetUsers' Cause: %v", err)
	}
	defer rows.Close()

	store := make(map[int]models.User)
	for rows.Next() {
		var u models.User
		err = rows.Scan(&u.ID, &u.Name, &u.Age, pq.Array(&u.Friends))
		if err != nil {
			return nil, fmt.Errorf("method 'GetUsers' Cause: %v", err)
		}
		store[u.ID] = u
	}
	return store, nil
}

func (s *Storage) GetUserById(id int) (models.User, error) {
	// DataReading - функция получения пользователя по конкретному id из БД
	row := s.DB.QueryRow("select * from users t where id = $1", id)
	u := models.User{}
	err := row.Scan(&u.ID, &u.Name, &u.Age, pq.Array(&u.Friends))
	if err != nil {
		return models.User{}, fmt.Errorf("method 'GetUserById' Cause: %v", err)
	}
	return u, nil
}

func (s *Storage) SaveUser(name string, age int, friends []int64) (int, error) {
	// DataRecording - функция записи пользователя в БД
	var id int
	err := s.DB.QueryRow("insert into users (name, age, friends) values ($1, $2, $3) returning id", name, age, pq.Array(friends)).Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("method 'SaveUser' Cause: %v", err)
	}
	return id, nil
}

func (s *Storage) UpdateUser(id int, name string, age int, friends []int64) error {
	// DataUpdating - функция обновления пользователя в БД
	_, err := s.DB.Exec("update users set name = $1, age = $2, friends = $3 where id = $4", name, age, pq.Array(friends), id)
	if err != nil {
		return fmt.Errorf("method 'UpdateUser' Cause: %v", err)
	}
	return nil
}

func (s *Storage) DeleteUser(id int) error {
	// DataDeleting - функция удаления пользователя в БД
	_, err := s.DB.Exec("delete from users where id = $1", id)
	if err != nil {
		return fmt.Errorf("method 'DeleteUser' Cause: %v", err)
	}
	return nil
}
