package db

import (
	"database/sql"
	"fmt"
	"networkCommunicationMin/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestStorage_GetUserById(t *testing.T) {
	db, mockDB, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}
	defer db.Close()

	type fields struct {
		db   *sql.DB
		mock func()
	}
	type args struct {
		id int
	}

	userID := 1
	userIDFailed := -1
	user := models.User{ID: 1, Name: "John", Age: 30, Friends: []int64{2, 3}}
	userEmpty := models.User{}
	errUserNotFound := fmt.Errorf("user not found")
	errMsg := fmt.Errorf("method 'GetUserById' Cause: %s", errUserNotFound.Error())

	rows := sqlmock.NewRows([]string{"id", "name", "age", "friends"}).
		AddRow(user.ID, user.Name, user.Age, pq.Array(user.Friends))
	sql := "select t.id, t.name, t.age, t.friends from users t where id = \\$1"

	tests := []struct {
		name           string
		fields         fields
		args           args
		want           models.User
		wantErrMessage error
		wantErr        bool
	}{
		{
			name: "success: user received",
			fields: fields{
				db: db,
				mock: func() {
					mockDB.ExpectQuery(sql).
						WithArgs(userID).
						WillReturnRows(rows)
				},
			},
			args: args{
				id: userID,
			},
			want: user,
		},
		{
			name: "failed: user not received",
			fields: fields{
				db: db,
				mock: func() {
					mockDB.ExpectQuery(sql).
						WithArgs(userIDFailed).
						WillReturnError(errUserNotFound)
				},
			},
			args: args{
				id: userIDFailed,
			},
			want:           userEmpty,
			wantErr:        true,
			wantErrMessage: errMsg,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				DB: tt.fields.db,
			}
			tt.fields.mock()

			got, err := s.GetUserById(tt.args.id)
			assert.Equal(t, tt.want, got, "result should be setted")
			assert.Equal(t, tt.wantErrMessage, err, "err should be setted")

			// Проверка на выполнение всех ожидаемых запросов
			if err := mockDB.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %v", err)
			}
		})
	}
}
