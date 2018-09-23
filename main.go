package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/satori/go.uuid"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"-"`
	Age      int    `json:"age"`
	Score    int32  `json:"score"`
}

type UserRequest struct {
	User
	Password string `json:"password"`
}

var sessions = make(map[string]string)
var users = map[string]User{
	"a.ostapenko@corp.mail.ru": User{
		Email:    "a.ostapenko@corp.mail.ru",
		Password: "111",
		Age:      21,
		Score:    72,
	},
	"d.dorofeev@corp.mail.ru": User{
		Email:    "d.dorofeev@corp.mail.ru",
		Password: "232323",
		Age:      24,
		Score:    100500,
	},
	"s.volodin@corp.mail.ru": User{
		Email:    "s.volodin@corp.mail.ru",
		Password: "22",
		Age:      21,
		Score:    0,
	},
	"a.tyuldyukov@corp.mail.ru": User{
		Email:    "a.tyuldyukov@corp.mail.ru",
		Password: "kek",
		Age:      6,
		Score:    7,
	},
}

func main() {
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {

		slice := make([]User, 0, 4)
		for _, user := range users {
			slice = append(slice, user)
		}

		resp, err := json.Marshal(&slice)
		if err != nil {
			log.Printf("cannot marshal: %s\n", err)
		}

		log.Println("request", r.URL)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	})

	http.HandleFunc("/me", func(w http.ResponseWriter, r *http.Request) {
		id, err := r.Cookie("session_id")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		email, ok := sessions[id.Value]
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		user := users[email]
		resp, err := json.Marshal(&user)
		if err != nil {
			log.Printf("cannot marshal: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Println("request", r.URL)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)

	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		userReq := &UserRequest{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(userReq)
		if err != nil {
			log.Printf("cannot marshal: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user, ok := users[userReq.Email]
		if !ok {
			log.Printf("User %s not found", user.Email)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if user.Password != userReq.Password {
			log.Printf("User %s: wrong password", user.Email)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sessionID := uuid.NewV4().String()
		sessions[sessionID] = userReq.Email

		cookie := &http.Cookie{
			Name:  "session_id",
			Value: sessionID,
		}
		http.SetCookie(w, cookie)
		w.WriteHeader(http.StatusOK)

	})

	http.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		request := &UserRequest{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(request)
		if err != nil {
			log.Printf("cannot marshal: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, exists := users[request.Email]; request.Email == "" || request.Password == "" || exists {
			log.Printf("cannot create user: exists %v, email: '%s'", exists, request.Email)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		users[request.Email] = User{
			Email:    request.Email,
			Password: request.Password,
			Age:      request.Age,
		}
		sessionID := uuid.NewV4().String()
		sessions[sessionID] = request.Email

		cookie := &http.Cookie{
			Name:  "session_id",
			Value: sessionID,
		}
		http.SetCookie(w, cookie)
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := fmt.Sprintf("can not get %s", r.URL)
		log.Printf(err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err))
	})

	log.Println("try to listen on http://127.0.0.1:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("cannot listen: %s", err)
	}
}
