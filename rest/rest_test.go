package rest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"networkCommunicationMin/mocks"
	"networkCommunicationMin/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestService_GetAll(t *testing.T) {
	type fields struct {
		storage func(ctrl *gomock.Controller) *mocks.MockStorage
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}

	users := map[int]models.User{
		1: {ID: 1, Name: "Tom", Age: 30, Friends: nil},
		2: {ID: 2, Name: "Sara", Age: 28, Friends: nil},
	}
	usersEmpty := map[int]models.User{}
	url := "http://localhost:3000/get_all"
	resultSuccess := "{\"result\":[\"id 1: name is Tom, age is 30 and friends are []\",\"id 2: name is Sara, age is 28 and friends are []\"],\"errors\":null}\n"
	someErr := fmt.Errorf("some error")
	resultError := "{\"result\":null,\"errors\":[\"get users: " + someErr.Error() + "\"]}\n"

	tests := []struct {
		name       string
		fields     fields
		args       args
		wantCode   int
		wantResult string
	}{
		{
			name: "success: users received",
			fields: fields{
				storage: func(ctrl *gomock.Controller) *mocks.MockStorage {
					mock := mocks.NewMockStorage(ctrl)
					mock.EXPECT().GetUsers().Return(users, nil)
					return mock
				},
			},
			args: args{
				httptest.NewRecorder(),
				httptest.NewRequest(http.MethodGet, url, nil),
			},
			wantCode:   http.StatusOK,
			wantResult: resultSuccess,
		},
		{
			name: "failed: users not received",
			fields: fields{
				storage: func(ctrl *gomock.Controller) *mocks.MockStorage {
					mock := mocks.NewMockStorage(ctrl)
					mock.EXPECT().GetUsers().Return(usersEmpty, someErr)
					return mock
				},
			},
			args: args{
				httptest.NewRecorder(),
				httptest.NewRequest(http.MethodGet, url, nil),
			},
			wantCode:   http.StatusInternalServerError,
			wantResult: resultError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := &Service{
				Storage: tt.fields.storage(ctrl),
			}

			s.GetAll(tt.args.w, tt.args.r)

			assert.Equal(t, tt.wantCode, tt.args.w.Code, "status code should be setted")
			assert.Equal(t, tt.wantResult, tt.args.w.Body.String(), "result should be setted")
		})
	}
}
