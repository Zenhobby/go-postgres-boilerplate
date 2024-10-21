package app

import (
	"encoding/json"
	"fmt"
	"go-postgres-boilerplate/dao"
	"log"
	"net/http"
	"time"

	"database/sql"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	PersonDAO dao.PersonDAO
}

func NewApp(personDAO dao.PersonDAO) *App {
	return &App{
		PersonDAO: personDAO,
	}
}

func (a *App) SetupRouter() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/person", a.GetAllPersons).Methods("GET")
	router.HandleFunc("/person", a.CreatePerson).Methods("POST")
	router.HandleFunc("/person/{id}", a.GetPerson).Methods("GET")
	router.HandleFunc("/person/{id}", a.UpdatePerson).Methods("PUT")
	router.HandleFunc("/person/{id}", a.DeletePerson).Methods("DELETE")

	log.Print("App Router Initialized!")
	return router
}

func (a *App) GetAllPersons(w http.ResponseWriter, r *http.Request) {
	people, err := a.PersonDAO.GetAllPersons()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(people)
}

func (a *App) CreatePerson(w http.ResponseWriter, r *http.Request) {
	var person dao.Person
	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = a.PersonDAO.GetPersonByName(person.Name)
	if err == nil {
		http.Error(w, fmt.Sprintf("%s already exists in the database!", person.Name), http.StatusConflict)
		return
	} else if err != sql.ErrNoRows {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	person.Timestamp = time.Now()
	person.UID = uuid.New().String()
	err = a.PersonDAO.Save(&person)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(person)
}

func (a *App) GetPerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	person, err := a.PersonDAO.GetPersonById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Person not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(person)
}

func (a *App) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var person dao.Person
	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	person.UID = id
	err = a.PersonDAO.Save(&person)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(person)
}

func (a *App) DeletePerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	err := a.PersonDAO.DeletePerson(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
