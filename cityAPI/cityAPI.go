package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/justinas/alice"
)

type city struct {
	Name string
	Area uint64
}

func filterContentType(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Currently in the content type checker middleware")
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte("415 - Usupported Media Type. JSON required"))
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func setServerTimeCookie(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Setting Cookie to every response
		cookie := &http.Cookie{
			Name:    "Server-Time",
			Value:   strconv.FormatInt(time.Now().Unix(), 10),
			Path:    "/",
			Expires: time.Now().Add(time.Hour * 24),
			MaxAge:  86400,
		}
		log.Println("Currently in the set server time middleware")
		http.SetCookie(w, cookie)
		handler.ServeHTTP(w, r)
	})
}

func mainLogic(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var tempCity city
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&tempCity)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 - Bad Request"))
			return
		}
		defer r.Body.Close()

		fmt.Printf("Got %s city with area of %d sq miles!\n", tempCity.Name, tempCity.Area)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("201 - Created"))
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 - Method Not Allowed"))
	}
}

func main() {
	mainLogicHandler := http.HandlerFunc(mainLogic)
	chain := alice.New(setServerTimeCookie, filterContentType).Then(mainLogicHandler)
	http.Handle("/city", chain)
	http.ListenAndServe(":80", nil)
}
