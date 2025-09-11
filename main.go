package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var database *sql.DB

// creazione database
func initDatabase() {
	var err error
	database, err = sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatal(err)
	}
	//SQL per la creazione della tabella
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS users (
	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	name TEXT)
	);`
	_, err = database.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("%v: %v\n", err, sqlStmt)
	}
}

// Handler per la richiesta della root del server
func rootHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(writer, "Benvenuto nel server")
}

// Handler per la richiesta di ./users/create
func getUsersHandler(writer http.ResponseWriter, request *http.Request) {
	rows, err := database.Query("SELECT id, name FROM users")

	if err != nil {
		http.Error(writer, "Errore query Database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			http.Error(writer, "Errore lettura dati", http.StatusInternalServerError)
		}
		users = append(users, user)
	}
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(users)
}

// Handler per la richiesta ./users/create
func createUserHandler(writer http.ResponseWriter, request *http.Request) {
	//controllo che la richiesta sia POST
	if request.Method != http.MethodPost {
		http.Error(writer, "Metodo non permesso", http.StatusInternalServerError)
		return
	}

	var newUser User
	err := json.NewDecoder(request.Body).Decode(&newUser)
	if err != nil {
		http.Error(writer, "Errore nel parsing del json", http.StatusInternalServerError)
		return
	}

	result, err := database.Exec("INSERT INTO users(name) VALUES(?)", newUser.Name)
	if err != nil {
		http.Error(writer, "Errore nel parsing del json", http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	newUser.ID = int(id)

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(newUser)
}

func main() {
	initDatabase()
	defer database.Close()

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/users", getUsersHandler)
	http.HandleFunc("users/create", createUserHandler)

	fmt.Println("Server avviato su http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

	//commit1
}
