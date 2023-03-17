package storage

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"main/internal/user"
	"sync"
)

type Storage interface {
	CreateDB()
	GetFriends(userId int) (friends []int, err error)
	MakeFriends(user1 int, user2 int) (friends [2]string, err error)
	CreateUser(u *user.User) int
	UpdateAge(userId int, newAge int) error
	DeleteUser(uId int) (username string, err error)
}

type MainSotrage struct {
	sync.Mutex
	db *sql.DB
}

func (ms *MainSotrage) CreateDB() {
	var err error
	ms.db, err = sql.Open("sqlite3", "C:\\Users\\mrwoo\\GolandProjects\\awesomeProject\\31\\main.db")
	if err != nil {
		log.Printf("%v", err)
		return
	}
}

func (ms *MainSotrage) GetFriends(userId int) (friends []int, err error) {
	sqlQuery := fmt.Sprintf("SELECT friend FROM Friend WHERE user_id = %v", userId)
	rows, err := ms.db.Query(sqlQuery)
	if err != nil {
		log.Fatal(err)
		err = errors.New("Пользователь не найден")
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var friend int
		err = rows.Scan(&friend)
		if err != nil {
			log.Fatal(err)
		}
		friends = append(friends, friend)
	}
	if friends == nil {
		return nil, errors.New("Друзей у пользователя не найдено")
	}
	return friends, nil

}

func (ms *MainSotrage) MakeFriends(user1 int, user2 int) (friends [2]string, err error) {
	ms.Lock()
	defer ms.Unlock()
	sqlQuery := fmt.Sprintf(`SELECT name FROM 'Users' WHERE id IN (%v, %v)`, user1, user2)
	rows, err := ms.db.Query(sqlQuery)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	i := 0
	for rows.Next() {
		var friend string
		fmt.Println(&friend)
		err = rows.Scan(&friends[i])
		if err != nil {
			log.Fatal(err)
		}
		i++
	}
	if i < 1 {
		err = errors.New("неверный идентификатор")
		return
	}
	sqlQuery = fmt.Sprintf(`INSERT INTO 'Friend' VALUES (%v, %v), (%v, %v)`, user1, user2, user2, user1)
	_, err = ms.db.Exec(sqlQuery)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlQuery)
		err = errors.New("Ошибка. Пользователи уже дружат")
		return
	}

	return
}

func (ms *MainSotrage) CreateUser(u *user.User) int {
	ms.Lock()
	sqlQuery := fmt.Sprintf(`INSERT INTO 'Users' (name, age) VALUES ('%v', %v)`, u.Name, u.Age)
	data, err := ms.db.Exec(sqlQuery)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlQuery)
	}
	newUserId, _ := data.LastInsertId()
	ms.Unlock()
	return int(newUserId)
}

func (ms *MainSotrage) UpdateAge(userId int, newAge int) error {
	ms.Lock()
	defer ms.Unlock()
	sqlQuery := fmt.Sprintf(`UPDATE Users SET age = %v WHERE id = %v`, newAge, userId)
	data, err := ms.db.Exec(sqlQuery)
	if err != nil {
		log.Fatal(err)
	}
	if i, err := data.RowsAffected(); err == nil {
		if i == 1 {
			return nil
		}
	}
	return errors.New("пользователь не найден")
}

func (ms *MainSotrage) DeleteUser(uId int) (username string, err error) {
	ms.Lock()
	defer ms.Unlock()
	sqlQuery := fmt.Sprintf(`SELECT name FROM Users WHERE id = %v`, uId)
	rows, err := ms.db.Query(sqlQuery)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&username)
		if err != nil {
			log.Fatal(err)
		}
	}
	sqlQuery = fmt.Sprintf(`DELETE FROM Users WHERE id = %v`, uId)
	data, err := ms.db.Exec(sqlQuery)
	if err != nil {
		log.Fatal(err)
	}
	if i, err_parse := data.RowsAffected(); err_parse == nil {
		if i == 1 {
			return
		}
	}
	return "", errors.New("пользователь не найден")
}

func NewStore() *MainSotrage {

	ms := MainSotrage{}
	ms.CreateDB()
	return &ms
}

var Repository = NewStore()
