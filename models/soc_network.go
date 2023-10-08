package models

import "fmt"

type User struct {
	ID      int     `json:"-"`
	Name    string  `json:"name"`
	Age     int     `json:"age"`
	Friends []int64 `json:"friends"`
}

func (u *User) ToSting() string {
	// ToSting - функция преобразования вывода данных на экран
	return fmt.Sprintf("name is %s, age is %d and friends are %v \n", u.Name, u.Age, u.Friends)
}

type Friends struct {
	SourceId int `json:"source_id"`
	TargetId int `json:"target_id"`
}

type ChangeAge struct {
	NewAge int `json:"new_age"`
}
