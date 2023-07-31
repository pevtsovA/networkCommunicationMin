package models

import (
	"fmt"
	"log"
	"math"
	"net/http"
)

func (u *User) ToSting() string {
	// ToSting - функция преобразования вывода данных на экран
	return fmt.Sprintf("name is %s, age is %d and friends are %v \n", u.Name, u.Age, u.Friends)
}

func FindUser(a []int64, x User) int {
	// FindUser - функция нахождения индекса указанного пользователя
	for i, n := range a {
		if int64(x.ID) == n {
			return i
		}
	}
	return -1
}

func MaxUserID(u map[int]User) int {
	// MaxUserID - функция нахождения максимального индекса пользователя
	index := 0
	for _, i := range u {
		if index < i.ID {
			index = i.ID
		}
	}
	return index
}

func MinUserID(u map[int]User) int {
	// MinUserID - функция нахождения минимального индекса пользователя
	index := math.MaxInt
	for _, i := range u {
		if index > i.ID {
			index = i.ID
		}
	}
	return index
}

func Logger(handler http.Handler) http.Handler {
	// Logger - функция логирования
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/ping" {
			log.Println(request.URL.Path)
		}
		handler.ServeHTTP(writer, request)
	})
}

func CheckServer(url string) (result int) {
	// CheckServer - функция проверки работоспособности сервера
	resp, err := http.Get(url)
	if err != nil {
		return -1
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return resp.StatusCode
	}
	return resp.StatusCode
}

func Contains(a []string, x string) bool {
	// Contains - функция указывает, содержится ли x в a.
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func Find(a []string, x string) int {
	// Find - функция возвращает индекс элемента, если x содержится в a.
	for i, n := range a {
		if x == n {
			return i
		}
	}
	return -1
}

func FindMinID(m map[string]int) (int, string) {
	// FindMinID - функция нахождения минимального значения ключа
	value := math.MaxInt
	var key string
	for s, i := range m {
		if value > i {
			value = i
			key = s
		}
	}
	return value, key
}
