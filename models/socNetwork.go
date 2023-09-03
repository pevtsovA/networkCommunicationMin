package models

type User struct {
	ID      int     `json:"-"`
	Name    string  `json:"name"`
	Age     int     `json:"age"`
	Friends []int64 `json:"friends"`
}

type Friends struct {
	SourceId int `json:"source_id"`
	TargetId int `json:"target_id"`
}

type ChangeAge struct {
	NewAge int `json:"new_age"`
}
