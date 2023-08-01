package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"networkCommunicationMin/models"
	"time"
)

const proxyAddress string = "localhost:3000"

var (
	urlServer        string
	url              string
	servers          = []string{"http://localhost:8080", "http://localhost:8081", "http://localhost:8082", "http://localhost:8088"}
	availableServers []string
	ms               map[string]int
)

func main() {
	http.HandleFunc("/", HandleProxy)
	log.Println("listening", proxyAddress)

	go pingServers()
	ms = make(map[string]int)

	if err := http.ListenAndServe(proxyAddress, nil); err != nil {
		log.Fatalln(err)
	}
}

func HandleProxy(w http.ResponseWriter, r *http.Request) {
	// HandleProxy - функция обработки прокси
	var req *http.Request
	var res *http.Response
	var err error
	var bodyBytes []byte

	if r.URL.Path != "/favicon.ico" {
		url = r.URL.Path
	}

	bodyBytes, err = io.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}
	defer r.Body.Close()
	body := bytes.NewBuffer(bodyBytes)

	if len(ms) == 0 {
		for _, val := range availableServers {
			ms[val] = 0
		}
	} else if len(ms) < len(availableServers) {
		for _, val := range availableServers {
			_, ok := ms[val]
			if ok != true {
				ms[val] = 0
			}
		}
	}

	n, s := models.FindMinID(ms)
	urlServer = s
	ms[s] = n + 1

	if r.Method == "GET" {
		req = getReq()
	} else if r.Method == "POST" {
		req = postReq(body)
	} else if r.Method == "DELETE" {
		req = deleteReq(body)
	} else if r.Method == "PUT" {
		req = putReq(body)
	}

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	if req != nil {
		res, err = client.Do(req)
		if err != nil {
			log.Fatalln(err)
		}
		defer res.Body.Close()

		bodyBytes, err = io.ReadAll(res.Body)
		if err != nil {
			log.Fatalln(err)
		}

		w.Write(bodyBytes)
	} else {
		w.Write([]byte("Сервер не доступен"))
	}
}

func pingServers() {
	// pingServers - функция проверки соединения с доступными серверами
	var checkServer int
	var isContains bool
	var indexServer int
	for {
		for _, val := range servers {
			checkServer = models.CheckServer(val + "/ping")
			isContains = models.Contains(availableServers, val)
			if isContains {
				indexServer = models.Find(availableServers, val)
			}

			if checkServer == 200 && !isContains {
				availableServers = append(availableServers, val)
			} else if checkServer != 200 && isContains {
				availableServers = append((availableServers)[:indexServer], (availableServers)[indexServer+1:]...)
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func getReq() *http.Request {
	// getReq - функция get запроса на сервер
	req, err := http.NewRequest(
		http.MethodGet,
		urlServer+url,
		nil,
	)
	if err != nil {
		log.Fatalln(err)
	}
	return req
}

func postReq(body io.Reader) *http.Request {
	// postReq - функция post запроса на сервер
	req, err := http.NewRequest(
		http.MethodPost,
		urlServer+url,
		body,
	)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Content-Type", "application/json")
	return req
}

func deleteReq(body io.Reader) *http.Request {
	// deleteReq - функция delete запроса на сервер
	req, err := http.NewRequest(
		http.MethodDelete,
		urlServer+url,
		body,
	)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Content-Type", "application/json")
	return req
}

func putReq(body io.Reader) *http.Request {
	// putReq - функция put запроса на сервер
	req, err := http.NewRequest(
		http.MethodPut,
		urlServer+url,
		body,
	)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Content-Type", "application/json")
	return req
}
