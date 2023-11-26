package main

import (
	"bytes"
	"flag"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"networkCommunicationMin/secondary_function"
	"time"
)

var (
	urlServer        string
	url              string
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
	log.Info("listening localhost:" + *port)

	go pingServers()
	ms = make(map[string]int)

	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		log.Fatalln(err)
	}
}

func handleProxy(w http.ResponseWriter, r *http.Request) {
	// HandleProxy - функция обработки прокси
	var req *http.Request
	var res *http.Response
	var err error
	var bodyBytes []byte

	if r.URL.Path != "/favicon.ico" {
		url = r.URL.Path
		log.Info("resource request: ", url)
	}

	bodyBytes, err = io.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
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

	n, s := secondary_function.FindMinID(ms)
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
			log.Error(err)
		}

		if res != nil && res.Body != nil {
			defer res.Body.Close()

			bodyBytes, err = io.ReadAll(res.Body)
			if err != nil {
				log.Error(err)
			}

			w.Write(bodyBytes)
		}
	} else {
		w.Write([]byte("Сервер не доступен"))
	}
}

func makeReq(body io.Reader, method string) *http.Request {
	var req *http.Request
	var err error

	if method != http.MethodGet {
		req, err = http.NewRequest(
			method,
			urlServer+url,
			body,
		)
		if err != nil {
			log.Error(err)
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(
			method,
			urlServer+url,
			nil,
		)
		if err != nil {
			log.Error(err)
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
			checkServer = secondary_function.CheckServer(val + "/ping")
			isContains = secondary_function.Contains(availableServers, val)

			if isContains {
				indexServer = secondary_function.Find(availableServers, val)
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
