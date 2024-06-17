package forum

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

func LoginPage(w http.ResponseWriter, r *http.Request) {
	data := ViewData{Username1: "", Id1: "", Title: "Login & Reg"}
	RenderTemplate(w, "login.html", data)
}

func getUserByUsername(username string) ([]byte, error) {
	getUserSQL := `SELECT password FROM users WHERE username = ?`
	var password []byte
	err := db.QueryRow(getUserSQL, username).Scan(&password)
	if err != nil {
		return nil, err
	}
	return password, nil
}

func LoginProcess(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("login_name")
	password := r.FormValue("login_password")

	if len(username) <= 0 || len(password) <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка, введите логин и пароль"})
		return
	}

	// Получение хэша пароля пользователя из базы данных
	storedPassword, err := getUserByUsername(username)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка, неверный логин или пароль"})
		return
	}

	// Сравнение хэша пароля из базы данных с введенным паролем
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка, неверный логин или пароль"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка при сравнении паролей"})
		return
	}

	existingSession, exists := GetSessionByUsername(username)
	if exists {
		// Удаление предыдущей сессии для этого пользователя
		delete(sessionsMap, existingSession.ID)
	}

	// Создание идентификатора сессии
	sessionID, err := uuid.NewV4()
	if err != nil {
		// Ошибка при создании UUID
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

func GetSessionByUsername(username string) (Session, bool) {
	for _, session := range sessionsMap {
		if session.Username == username {
			return session, true
		}
	}
	return Session{}, false
}
