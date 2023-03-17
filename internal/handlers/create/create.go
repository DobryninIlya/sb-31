package create

import (
	"encoding/json"
	"fmt"
	"io"
	http "net/http"

	user "main/internal/user"
)

type AnswerCreate struct {
	Id int `json:"id"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name=UserCreator
type UserCreator interface {
	CreateUser(u *user.User) int
}

func UserUnmarshal(content []byte) (*user.User, error) {
	var u user.User
	if err := json.Unmarshal(content, &u); err != nil {
		return &u, err
	}
	return &u, nil
}

func New(creator UserCreator) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("create method")
		content, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		u, err := UserUnmarshal(content)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			fmt.Println(err.Error())
			return
		}
		newUserId := creator.CreateUser(u)

		ans := AnswerCreate{newUserId}
		answer, err := json.Marshal(ans)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			fmt.Println(err.Error())
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write(answer)
	}
}
