package create

import (
	"bytes"
	"encoding/json"
	"main/internal/user"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type MockUsecase struct {
	User   *user.User
	status int
}

func (s *MockUsecase) CreateUser(user *user.User) int {
	s.User = user
	return 1
}

func TestNew(t *testing.T) {
	tests := []struct {
		testUser *user.User
		wantUser *user.User
		status   int
	}{
		{
			testUser: &user.User{
				Name:   "name",
				Age:    10,
				Friend: []int{},
			},
			wantUser: &user.User{
				Name:   "name",
				Age:    10,
				Friend: []int{},
			},
			status: 201,
		},
		{
			testUser: &user.User{
				Name:   "ave maria",
				Age:    666,
				Friend: []int{0, 1},
			},
			wantUser: &user.User{
				Name:   "ave maria",
				Age:    666,
				Friend: []int{0, 1},
			},
			status: 201,
		},
		{
			testUser: &user.User{
				Name:   "deus velt",
				Age:    -10000,
				Friend: []int{0, 1},
			},
			wantUser: nil,
			status:   400,
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			stub := &MockUsecase{}
			handler := http.HandlerFunc(New(stub))

			data, _ := json.Marshal(tt.testUser)
			body := bytes.NewReader(data)

			rr := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/users", body)

			handler.ServeHTTP(rr, req)

			if !reflect.DeepEqual(rr.Code, tt.status) {
				t.Logf("handler returned wron status code")
				t.Logf("got: %v", rr.Code)
				t.Errorf("want: %v", tt.status)
				t.Fail()
			}

			if !reflect.DeepEqual(stub.User, tt.wantUser) {
				t.Logf("handler decoded wrong user structure")
				t.Logf("got: %v", stub.User)
				t.Errorf("want: %v", tt.wantUser)
			}
		})
	}
}
