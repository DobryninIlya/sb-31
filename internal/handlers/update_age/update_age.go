package updateAge

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"io"
	http "net/http"
	"strconv"
)

type NewAge struct {
	Age int `json:"new_age"`
}

type AgeUpdater interface {
	UpdateAge(userId int, newAge int) error
}

func New(updater AgeUpdater) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("update user age method")
		userIdS := chi.URLParam(r, "userId")
		userId, err := strconv.Atoi(userIdS)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		}
		content, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		var age NewAge
		if err := json.Unmarshal(content, &age); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		err = updater.UpdateAge(userId, age.Age)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Пользователя с таким id не существует."))
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
