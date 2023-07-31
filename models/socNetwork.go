package models

type User struct {
	ID      int     `json:"-"`
	Name    string  `json:"name"`
	Age     int     `json:"age"`
	Friends []int64 `json:"friends"`
}

type Friends struct {
	Source_id int `json:"source_id"`
	Target_id int `json:"target_id"`
}

type ChangeAge struct {
	New_age int `json:"new_age"`
}
