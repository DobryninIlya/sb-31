package getFriends

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	http "net/http"
	"strconv"
)

type AnswerFriends struct {
	Friends []int `json:"friends"`
}

type GetterFriends interface {
	GetFriends(userId int) (friends []int, err error)
}

func New(getter GetterFriends) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("get user friends method")
		userIdS := chi.URLParam(r, "userId")
		userId, err := strconv.Atoi(userIdS)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		var friends AnswerFriends
		friendsSlise, err := getter.GetFriends(userId)
		friends.Friends = friendsSlise
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Друзей у пользователя не найдено"))
			return
		}
		result, err := json.Marshal(friends)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(result)
		return
	}
}
