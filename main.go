package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Todo struct {
	Name        string `json:"name"`
	Description string `json:"desc"`
	Done        bool   `json:"done"`
}

type Todos []Todo

type todoHandler struct {
	sync.Mutex
	todos Todos
}

func (th *todoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		th.get(w, r)
	case "POST":
		th.post(w, r)
	case "PUT", "PATCH":
		th.put(w, r)
	case "DELETE":
		th.delete(w, r)
	default:
		respondWithError(w, http.StatusMethodNotAllowed, "invalid method")
	}
}

func (th *todoHandler) get(w http.ResponseWriter, r *http.Request) {
	defer th.Unlock()
	th.Lock()
	id, err := idFromUrl(r)
	if err != nil {
		respondWithJSON(w, http.StatusOK, th.todos)
		return
	}
	if id >= len(th.todos) || id < 0 {
		respondWithError(w, http.StatusNotFound, "not found")
		return
	}
	respondWithJSON(w, http.StatusOK, th.todos[id])
}
func (th *todoHandler) post(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		respondWithError(w, http.StatusUnsupportedMediaType, "content type 'application/json' required")
		return
	}
	var todo Todo
	err = json.Unmarshal(body, &todo)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer th.Unlock()
	th.Lock()
	th.todos = append(th.todos, todo)
	respondWithJSON(w, http.StatusCreated, todo)
}
func (th *todoHandler) put(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	id, err := idFromUrl(r)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		respondWithError(w, http.StatusUnsupportedMediaType, "content type 'application/json' required")
		return
	}
	var todo Todo
	err = json.Unmarshal(body, &todo)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer th.Unlock()
	th.Lock()
	if id >= len(th.todos) || id < 0 {
		respondWithError(w, http.StatusNotFound, "not found")
		return
	}
	if todo.Name != "" {
		th.todos[id].Name = todo.Name
	}
	if todo.Description != "" {
		th.todos[id].Description = todo.Description
	}
	th.todos[id].Done = todo.Done

	respondWithJSON(w, http.StatusOK, th.todos[id])
}
func (th *todoHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := idFromUrl(r)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "not found")
		return
	}
	defer th.Unlock()
	th.Lock()
	if id >= len(th.todos) || id < 0 {
		respondWithError(w, http.StatusNotFound, "not found")
		return
	}
	if id < len(th.todos)-1 {
		th.todos[len(th.todos)-1], th.todos[id] = th.todos[id], th.todos[len(th.todos)-1]
	}
	th.todos = th.todos[:len(th.todos)-1]
	respondWithJSON(w, http.StatusNoContent, "")
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"error": msg})
}

func respondWithJSON(w http.ResponseWriter, code int, data interface{}) {
	response, _ := json.Marshal(data)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func idFromUrl(r *http.Request) (int, error) {
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 3 {
		return 0, errors.New("not found")
	}
	id, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		return 0, errors.New("not found")
	}
	return id, nil
}

func newTodoHandler() *todoHandler {
	return &todoHandler{
		todos: Todos{
			Todo{"Task 1", "The first task", false},
			Todo{"Task 2", "The second task", false},
			Todo{"Task 3", "The third task", false},
		},
	}
}

func main() {
	port := ":8080"
	th := newTodoHandler()
	http.Handle("/todos", th)
	http.Handle("/todos/", th)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})
	fmt.Println("Starting server on port", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
