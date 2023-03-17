package delete

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	http "net/http"
	"strconv"
)

type Answer struct {
	Username string `json:"username,omitempty"`
}

type UserDeleter interface {
	DeleteUser(uId int) (username string, err error)
}

func New(deleter UserDeleter) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("delete user method")
		userId := chi.URLParam(r, "userId")
		uId, err := strconv.Atoi(userId)
		if err != nil || uId < 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		username, err := deleter.DeleteUser(uId)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		var a Answer
		a.Username = username
		answer, err := json.Marshal(a)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(answer)
	}
}
