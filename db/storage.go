package db

import (
	"database/sql"
	"fmt"
	"networkCommunicationMin/models"

	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
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
	// ConnectToBD - функция подключения к БД
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
	// DataReading - метод получения пользователей из БД
	rows, err := s.DB.Query("select t.id, t.name, t.age, t.friends from users t")
	if err != nil {
		return nil, fmt.Errorf("method 'GetUsers' Cause: %w", err)
	}
	defer rows.Close()

	store := make(map[int]models.User)
	for rows.Next() {
		var u models.User
		if err = rows.Scan(&u.ID, &u.Name, &u.Age, pq.Array(&u.Friends)); err != nil {
			return nil, fmt.Errorf("method 'GetUsers' Cause: %w", err)
		}
		store[u.ID] = u
	}
	return store, nil
}

func (s *Storage) GetUserById(id int) (models.User, error) {
	// DataReading - метод получения пользователя по конкретному id из БД
	u := models.User{}
	sql := "select t.id, t.name, t.age, t.friends from users t where id = $1"
	if err := s.DB.QueryRow(sql, id).Scan(&u.ID, &u.Name, &u.Age, pq.Array(&u.Friends)); err != nil {
		return models.User{}, fmt.Errorf("method 'GetUserById' Cause: %w", err)
	}
	return u, nil
}

func (s *Storage) SaveUser(name string, age int, friends []int64) (int, error) {
	// DataRecording - метод записи пользователя в БД
	var id int
	sql := "insert into users (name, age, friends) values ($1, $2, $3) returning id"
	if err := s.DB.QueryRow(sql, name, age, pq.Array(friends)).Scan(&id); err != nil {
		return -1, fmt.Errorf("method 'SaveUser' Cause: %w", err)
	}
	return id, nil
}

func (s *Storage) UpdateUser(id int, name string, age int, friends []int64) error {
	// DataUpdating - метод обновления пользователя в БД
	sql := "update users set name = $1, age = $2, friends = $3 where id = $4"
	if _, err := s.DB.Exec(sql, name, age, pq.Array(friends), id); err != nil {
		return fmt.Errorf("method 'UpdateUser' Cause: %w", err)
	}
	return nil
}

func (s *Storage) DeleteUser(id int) error {
	// DataDeleting - метод удаления пользователя в БД
	if _, err := s.DB.Exec("delete from users where id = $1", id); err != nil {
		return fmt.Errorf("method 'DeleteUser' Cause: %w", err)
	}
	return nil
}
