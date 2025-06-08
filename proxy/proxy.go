package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"net/http"
	"networkCommunicationMin/rest"
	secondary "networkCommunicationMin/secondary_function"
	"slices"
	"time"

	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
)

var (
	urlServer        string
	endpoint         string
	servers          []string
	availableServers []string
	ms               map[string]int
	isDocker         *bool
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	port := flag.String("port", "3000", "Listen server port")
	isDocker = flag.Bool("runviadocker", false, "Run the application")
	flag.Parse()

	http.HandleFunc("/", handleProxy)
	log.Info("listening localhost: ", *port)

	go pingServers()
	ms = make(map[string]int)

	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		log.Fatalln("listenAndServe proxy: ", err)
	}
}

func handleProxy(w http.ResponseWriter, r *http.Request) {
	// HandleProxy - функция обработки прокси
	var req *http.Request
	var res *http.Response
	var err error
	var bodyBytes []byte

	if r.URL.Path != "/favicon.ico" {
		endpoint = r.URL.Path
		log.Info("resource request: ", endpoint)
	}

	bodyBytes, err = io.ReadAll(r.Body)
	if err != nil {
		log.Error("handleProxy: read body: ", err)
	}
	defer r.Body.Close()
	body := bytes.NewBuffer(bodyBytes)

	if len(ms) == 0 {
		for _, val := range availableServers {
			ms[val] = 0
		}
	} else if len(ms) < len(availableServers) {
		for _, val := range availableServers {
			if _, ok := ms[val]; !ok {
				ms[val] = 0
			}
		}
	}

	n, s := secondary.FindMinID(ms)
	urlServer = s
	ms[s] = n + 1

	if r.Method == http.MethodGet {
		body = nil
	}

	req = makeReq(body, r.Method)

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	if req != nil {
		res, err = client.Do(req)
		if err != nil {
			log.Error("handleProxy: client do req: ", err)
		}

		if res != nil && res.Body != nil {
			defer res.Body.Close()

			var response rest.ResponsePayload
			if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
				log.Error("handleProxy: response decode: ", err)
			}

			if response.Errors != nil {
				render.Status(r, http.StatusBadRequest)
			}
			_ = render.Render(w, r, &rest.ResponsePayload{Result: response.Result, Errors: response.Errors})
		}
	} else {
		render.Status(r, http.StatusServiceUnavailable)
		_ = render.Render(w, r, &rest.ResponsePayload{Errors: []string{"Сервер не доступен"}})
	}
}

func makeReq(body io.Reader, method string) *http.Request {
	// makeReq - функция создания запроса
	var req *http.Request
	var err error

	urlSrv := secondary.GetURL(urlServer, endpoint)
	if method != http.MethodGet {
		req, err = http.NewRequest(
			method,
			urlSrv,
			body,
		)
		if err != nil {
			log.Error("makeReq: new request: ", err)
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(
			method,
			urlSrv,
			nil,
		)
		if err != nil {
			log.Error("makeReq: new request: ", err)
		}
	}
	return req
}

func pingServers() {
	// pingServers - функция проверки соединения с доступными серверами
	var checkServer int
	var isContains bool
	var indexServer int

	if *isDocker {
		servers = []string{"http://172.19.0.1:8080", "http://172.19.0.1:8081", "http://172.19.0.1:8082", "http://172.19.0.1:8088"}
	} else {
		servers = []string{"http://localhost:8080", "http://localhost:8081", "http://localhost:8082", "http://localhost:8088"}
	}

	for {
		for _, val := range servers {
			urlSrv := secondary.GetURL(val, "ping")
			checkServer = secondary.CheckServer(urlSrv)
			isContains = secondary.Contains(availableServers, val)

			if isContains {
				indexServer = secondary.Find(availableServers, val)
			}

			if checkServer == http.StatusOK && !isContains {
				availableServers = append(availableServers, val)
			} else if checkServer != http.StatusOK && isContains {
				availableServers = slices.Delete(availableServers, indexServer, indexServer+1)
			}
		}
		time.Sleep(5 * time.Second)
	}
}
