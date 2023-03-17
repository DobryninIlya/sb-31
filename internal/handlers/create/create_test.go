package create

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"main/internal/handlers/create/mocks"
	"main/internal/user"
	"net/http"
	"net/http/httptest"
	"testing"
)

type userParams struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Friends []int  `json:"friends"`
}

func TestNew(t *testing.T) {
	type args struct {
		creator UserCreator
		user    user.User
	}
	tests := []struct {
		name string
		args args
		want func(w http.ResponseWriter, r *http.Request)
	}{
		{
			name: "base test",
			args: args{
				user: user.User{
					Name:   "lalala",
					Age:    66,
					Friend: []int{},
				},
			},
		},
	}
	bytesRepresentation, err := json.Marshal(userParams{"test", 2, []int{}})
	if err != nil {
		log.Fatalln(err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userCreator := mocks.NewUserCreator(t)
			userCreator.On("CreateUser", tt.args.user).Return(20)
			got := New(tt.args.creator)
			srv := httptest.NewServer(http.HandlerFunc(got))

			resp, err := http.Post(srv.URL, "text/json", bytes.NewBuffer(bytesRepresentation))
			if err != nil {
				t.Log(err)
				t.Fail()
			}
			defer resp.Body.Close()
			textBytes, err := io.ReadAll(resp.Body)
			fmt.Println(textBytes)

			//if got := New(tt.args.creator); !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("New() = %v, want %v", got, tt.want)
			//}
		})
	}
}
