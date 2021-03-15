package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type memcachedHandler struct {
	Client *memcache.Client
}

//Result none
type Result struct {
	Surl string `json:"surl"` //Short url
	Eurl string `json:"eurl"` //Expanded url
}

func (h *memcachedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	surl := r.URL.Query().Get("url")
	if surl == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	surl = strings.Trim(surl, " ")
	mkey := "m_" + surl
	h.Client.Timeout = 5 * time.Second
	_eurl, err := h.Client.Get(mkey)
	eurl := surl
	cachehit := "false"
	if err == nil {
		eurl = string(_eurl.Value)
		cachehit = "true"
	} else {
		req, err := http.NewRequest("GET", surl, nil)
		if err == nil {
			req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.3 Safari/605.1.15")
			req.Header.Set("Accept", "*/*")
			tr := &http.Transport{
				DisableKeepAlives:  true,
				MaxIdleConns:       1,
				IdleConnTimeout:    30 * time.Second,
				DisableCompression: false,
			}
			client := &http.Client{Transport: tr}
			resp, err := client.Do(req)
			if err == nil {
				eurl = resp.Request.URL.String()
				h.Client.Set(&memcache.Item{Key: mkey, Value: []byte(eurl), Expiration: 0})
			} else {
				log.Println("Failed to process request", err)
			}
		} else {
			log.Println("Failed to get request", err)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Cache-Hit", cachehit)
	result := Result{Surl: surl, Eurl: eurl}
	json.NewEncoder(w).Encode(result)
}

func serveProxy(w http.ResponseWriter, r *http.Request) {
	surl := r.URL.Query().Get("url")
	if surl == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	surl = strings.Trim(surl, " ")
	contentType := "text/html; charset=utf-8"

	resp, err := http.Get(surl)
	if err != nil {
		log.Println("Failed to process request", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	contentType = resp.Header.Get("Content-Type")

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sb := string(body)
	w.Header().Set("Content-Type", contentType)
	w.Write([]byte(sb))
}

func main() {
	memserver := func() string {
		_memserver, exists := os.LookupEnv("MEMCACHED_SERVER")
		if !exists {
			_memserver = "memcached.iavian.net:11211"
		}
		return _memserver
	}()

	os.Stdout.WriteString(memserver)
	os.Stdout.WriteString("\n")

	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = ":8080"
	} else {
		port = ":" + port
	}
	http.Handle("/", &memcachedHandler{Client: memcache.New(memserver)})
	http.HandleFunc("/proxy", serveProxy)
	http.ListenAndServe(port, nil)
}
