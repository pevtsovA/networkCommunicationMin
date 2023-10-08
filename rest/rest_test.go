package rest

import (
	"io"
	"net/http"
	"net/http/httptest"
	"networkCommunicationMin/contract"
	"testing"
)

func TestService_GetAll(t *testing.T) {
	s := Service{
		Storage: &contract.MockStorage{},
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/get_all", nil)

	s.GetAll(w, r)

	responseString, _ := io.ReadAll(w.Body)
	expectedString := "name is Tom, age is 30 and friends are [] \nname is Sara, age is 28 and friends are [] \n"

	if string(responseString) != expectedString {
		t.Fatalf("expected: %s, got: %s", expectedString, string(responseString))
	}
}

func TestMockStorage_GetUserById(t *testing.T) {
	s := Service{
		Storage: &contract.MockStorage{},
	}

	var responseString string
	response, err := s.Storage.GetUserById(3)
	if err != nil {
		responseString = err.Error()
	} else {
		responseString = response.ToSting()
	}
	expectedString := "name is Sam, age is 19 and friends are [] \n"

	if responseString != expectedString {
		t.Fatalf("expected: %s, got: %s", expectedString, responseString)
	}
}
