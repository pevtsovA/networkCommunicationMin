package secondary_function

import (
	"math"
	"net/http"
	"net/url"
	"networkCommunicationMin/models"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func FindUser(a []int64, x models.User) int {
	// FindUser - функция нахождения индекса указанного пользователя
	for i, n := range a {
		if int64(x.ID) == n {
			return i
		}
	}
	return -1
}

func Logger(handler http.Handler) http.Handler {
	// Logger - функция логирования
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/ping" {
			log.Info("resource request: ", request.URL.Path)
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

func GetURL(path, endpoint string) string {
	base, err := url.Parse(path)
	if err != nil {
		log.Error("parse url: ", err)
	}
	base.Path, err = url.JoinPath(base.Path, endpoint)
	if err != nil {
		log.Error("join path url: ", err)
	}

	return base.String()
}
