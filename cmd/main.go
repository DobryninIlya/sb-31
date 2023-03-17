package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"log"
	createHandler "main/internal/handlers/create"
	deleteHandler "main/internal/handlers/delete"
	getFriendsHandler "main/internal/handlers/get_friends"
	makeFriendsHandler "main/internal/handlers/make_friends"
	updateAgeHandler "main/internal/handlers/update_age"
	"main/internal/storage"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
)

type App struct {
	router *chi.Mux
	done   chan os.Signal
}

func NewApp() *App {
	ret := &App{
		router: chi.NewRouter(),
		done:   make(chan os.Signal, 1),
	}
	signal.Notify(ret.done, os.Interrupt)
	return ret
}

func (a *App) run(port string) {
	//a.router.Post("/user", createUser.Create)
	a.router.Route("/users", func(r chi.Router) {
		r.Post("/", createHandler.New(storage.Repository)) // POST /users
		r.Route("/{userId}", func(r chi.Router) {
			r.Post("/make_friends", makeFriendsHandler.New(storage.Repository))
			r.Delete("/", deleteHandler.New(storage.Repository))
			r.Get("/friends", getFriendsHandler.New(storage.Repository))
			r.Patch("/age", updateAgeHandler.New(storage.Repository))
		})
	})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Fatal(http.ListenAndServe(":"+port, a.router))
		for {
			select {
			case <-a.done:
				fmt.Println("Exiting")
				return
			}
		}
	}()
	wg.Wait()

}

func main() {
	if len(os.Args[1:]) == 0 {
		log.Fatal("Must input port number")
		return
	}
	port := os.Args[1]
	if portNum, err := strconv.Atoi(port); err == nil {
		if portNum > 10000 || portNum < 1000 {
			log.Fatal("Incorrect port number")
			return
		}
		var app = NewApp()
		app.run(port)
	} else {
		log.Fatal("Must input port number")
		return
	}

}
