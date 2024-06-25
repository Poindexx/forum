package forum

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func generateState() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	state, err := generateState()
	if err != nil {
		http.Error(w, "Failed to generate state", http.StatusInternalServerError)
		return
	}
	url := fmt.Sprintf("https://accounts.google.com/o/oauth2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=openid email&state=%s",
		googleClientID, "http://localhost:8080/auth/google/callback", state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")
	if state == "" || code == "" {
		http.Error(w, "State or code missing", http.StatusBadRequest)
		return
	}

	tokenURL := "https://oauth2.googleapis.com/token"
	data := url.Values{}
	data.Set("client_id", googleClientID)
	data.Set("client_secret", googleClientSecret)
	data.Set("redirect_uri", "http://localhost:8080/auth/google/callback")
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		http.Error(w, "Failed to get token", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var tokenResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&tokenResponse)
	token := tokenResponse["access_token"].(string)

	userInfoURL := fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s", token)
	resp, err = http.Get(userInfoURL)
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&userInfo)

	email := userInfo["email"].(string)
	HandleOAuthLogin(w, email)
}

func GithubLoginHandler(w http.ResponseWriter, r *http.Request) {
	state, err := generateState()
	if err != nil {
		http.Error(w, "Failed to generate state", http.StatusInternalServerError)
		return
	}
	url := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=read:user&state=%s",
		githubClientID, "http://localhost:8080/auth/github/callback", state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func GithubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")
	if state == "" || code == "" {
		http.Error(w, "State or code missing", http.StatusBadRequest)
		return
	}

	tokenURL := "https://github.com/login/oauth/access_token"
	data := url.Values{}
	data.Set("client_id", githubClientID)
	data.Set("client_secret", githubClientSecret)
	data.Set("redirect_uri", "http://localhost:8080/auth/github/callback")
	data.Set("code", code)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to get token", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var tokenResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&tokenResponse)
	token := tokenResponse["access_token"].(string)

	userInfoURL := "https://api.github.com/user"
	req, err = http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&userInfo)

	email := userInfo["email"].(string)
	HandleOAuthLogin(w, email)
}

func HandleOAuthLogin(w http.ResponseWriter, email string) {
	// Создание нового пользователя в базе данных, если его еще нет
	// Идентификатором будет являться email, полученный от провайдера OAuth
	username := email // Здесь можно использовать email как username, либо придумать другую логику

	existingSession, exists := GetSessionByUsername(username)
	if exists {
		// Удаление предыдущей сессии для этого пользователя
		delete(sessionsMap, existingSession.ID)
	}

	sessionID, err := generateState() // Используем state как идентификатор сессии
	if err != nil {
		log.Printf("Error generating session ID: %v\n", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Установка времени для сессии (например, 1 час)
	sessionDuration := time.Hour

	// Установка куки с идентификатором сессии
	http.SetCookie(w, &http.Cookie{
		Name:     "sessionID",
		Value:    sessionID,
		Path:     "/",
		MaxAge:   int(sessionDuration.Seconds()),  // Время в секундах
		Expires:  time.Now().Add(sessionDuration), // Время истечения
		HttpOnly: true,
	})

	// Сохранение идентификатора сессии в карте сессий
	sessionsMap[sessionID] = Session{ID: sessionID, Username: username, Authenticated: true}

	// Отправка ответа об успешном входе в формате JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"success": "Login successful", "sessionID": sessionID})
}
