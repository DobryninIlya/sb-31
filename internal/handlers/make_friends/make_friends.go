package makeFriends

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"io"
	http "net/http"
	"strconv"
)

type FriendTarget struct {
	Target    int    `json:"target_id,omitempty"`
	Username1 string `json:"username1,omitempty"`
	Username2 string `json:"username2,omitempty"`
}

type FriendMaker interface {
	MakeFriends(user1 int, user2 int) (friends [2]string, err error)
}

func New(friendMaker FriendMaker) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("make friend method")
		userId := chi.URLParam(r, "userId")
		uId, err := strconv.Atoi(userId)
		if err != nil || uId < 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		content, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		var tar FriendTarget
		if err := json.Unmarshal(content, &tar); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			fmt.Println(err.Error())
			return
		}
		friends, err := friendMaker.MakeFriends(uId, tar.Target)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Пользователи либо уже дружат, либо не найдены."))
			return
		}
		tar.Username1 = friends[0]
		tar.Username2 = friends[1]
		answer, err := json.Marshal(tar)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(answer)
		return
	}
}
