package forum

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(username, password, phone, email string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("ошибка хеширования пароля: %v", err)
	}

	insertUserSQL := `INSERT INTO users (username, password, phone, email) VALUES (?, ?, ?, ?)`
	_, err = db.Exec(insertUserSQL, username, hashedPassword, phone, email)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			switch {
			case strings.Contains(err.Error(), "username"):
				return fmt.Errorf("такой логин уже существует")
			case strings.Contains(err.Error(), "phone"):
				return fmt.Errorf("такой номер уже существует")
			case strings.Contains(err.Error(), "email"):
				return fmt.Errorf("такой email уже существует")
			default:
				return fmt.Errorf("серверная ошибка")
			}
		}
		return err
	}
	return nil
}

func RegisterProcess(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("reg_name")
	password := r.FormValue("reg_password")
	phone := r.FormValue("reg_phone")
	email := r.FormValue("reg_email")
	if len(username) <= 0 || len(password) <= 0 || len(phone) <= 0 || len(email) <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка заполните все поля"})
		return
	}

	// Создание нового пользователя в базе данных
	err := CreateUser(username, password, phone, email)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("%v", err)})
		log.Printf("Error creating user: %v\n", err)
		return
	}

	existingSession, exists := GetSessionByUsername(username)
	if exists {
		// Удаление предыдущей сессии для этого пользователя
		delete(sessionsMap, existingSession.ID)
	}

	sessionID, err := uuid.NewV4()
	if err != nil {
		log.Printf("Error generating UUID: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Серверная ошибка"})
		return
	}

	// Установка времени для сессии (например, 1 час)
	sessionDuration := time.Hour

	// Установка куки с идентификатором сессии
	http.SetCookie(w, &http.Cookie{
		Name:     "sessionID",
		Value:    sessionID.String(),
		Path:     "/",
		MaxAge:   int(sessionDuration.Seconds()),  // Время в секундах
		Expires:  time.Now().Add(sessionDuration), // Время истечения
		HttpOnly: true,
	})

	// Сохранение идентификатора сессии в карте сессий
	sessionsMap[sessionID.String()] = Session{ID: sessionID.String(), Username: username, Authenticated: true}

	// Отправка ответа об успешном входе в формате JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"success": "Login successful", "sessionID": sessionID.String()})
}
