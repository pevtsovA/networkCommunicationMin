package dataBase

import (
	"database/sql"
	"github.com/lib/pq"
	"log"
	"networkCommunicationMin/models"
)

type Service struct {
	DB *sql.DB
}

func (s *Service) DataReading() map[int]models.User {
	// DataReading - функция чтения данных из БД
	rows, err := s.DB.Query("select * from users")
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	store := make(map[int]models.User)
	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.ID, &u.Name, &u.Age, pq.Array(&u.Friends))
		if err != nil {
			log.Fatalln(err)
		}
		store[u.ID] = u
	}
	return store
}

func (s *Service) DataRecording(name string, age int, friends []int64) (id int) {
	// DataRecording - функция записи данных в БД
	s.DB.QueryRow("insert into users (name, age, friends) values ($1, $2, $3) returning id", name, age, pq.Array(friends)).Scan(&id)
	return
}

func (s *Service) DataUpdating(id int, name string, age int, friends []int64) {
	// DataUpdating - функция обновления данных в БД
	_, err := s.DB.Exec("update users set name = $1, age = $2, friends = $3 where id = $4", name, age, pq.Array(friends), id)
	if err != nil {
		panic(err)
	}
	return
}

func (s *Service) DataDeleting(id int) {
	// DataDeleting - функция удаления данных в БД
	_, err := s.DB.Exec("delete from users where id = $1", id)
	if err != nil {
		panic(err)
	}
	return
}
