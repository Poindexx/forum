package forum

import (
	"encoding/json"
	"log"
	"net/http"
)

func CreateUserTable() {
	// SQL запрос для создания таблицы пользователей
	createUserTableSQL := `
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT UNIQUE,
            password TEXT,
            phone TEXT UNIQUE,
			email TEXT UNIQUE
        )
    `
	// Выполняем SQL запрос
	_, err := db.Exec(createUserTableSQL)
	if err != nil {
		log.Fatalf("Error creating users table: %v\n", err)
	}
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := GetUsers()
	if err != nil {
		log.Printf("Error getting users: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Отправляем полученные категории в формате JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Printf("Error encoding users to JSON: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func GetUsers() ([]User, error) {
	var users []User

	query := "SELECT id, username FROM users"

	// Выполняем SQL запрос
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
