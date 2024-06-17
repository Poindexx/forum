package forum

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

func Error404Handler(w http.ResponseWriter, r *http.Request) {
	data := ViewData{Title: "ERROR 404 page"}
	RenderTemplate(w, "error404.html", data)
}

func SetupDB() {
	// Открываем соединение с базой данных SQLite
	database, err := sql.Open("sqlite3", "messenger.db")
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}
	db = database

	// Создаем таблицу пользователей, если она не существует
	CreateUserTable()
	CreateCategoryTable()
	CreateCommentLikesTable()
	CreateCommentsTable()
	CreatePostLikesTable()
	CreatePostTable()
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	posts, err := GetAllPosts()
	if err != nil {
		log.Printf("Error getting posts: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// Получение идентификатора сессии из куки
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		data := ViewData{Username1: "", Id1: "", Title: "FORUM", Posts: posts}
		RenderTemplate(w, "index.html", data)
		return
	}
	sessionID := cookie.Value
	session, ok := sessionsMap[sessionID]
	if !ok {
		// Если сессия не найдена в карте, отправляем шаблон без данных о сессии
		data := ViewData{Username1: "", Id1: "", Title: "FORUM", Posts: posts}
		RenderTemplate(w, "index.html", data)
		return
	}

	// Проверка соответствия идентификатора сессии идентификатору в куке
	if session.ID != sessionID {
		// Если идентификаторы не совпадают, считаем сессию недействительной
		data := ViewData{Username1: "", Id1: "", Title: "FORUM", Posts: posts}
		RenderTemplate(w, "index.html", data)
		return
	}

	data := ViewData{Username1: session.Username, Id1: session.ID, Title: "FORUM", Posts: posts}
	RenderTemplate(w, "index.html", data)
}

func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles("web/templates/base.html", "web/templates/"+tmpl)
	if err != nil {
		log.Printf("Error parsing template %s: %v\n", tmpl, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = t.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		log.Printf("Error executing template %s: %v\n", tmpl, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func ExitProcess(w http.ResponseWriter, r *http.Request) {
	// Получение идентификатора сессии из куки
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		// Если кука не найдена, отправляем шаблон без данных о сессии
		RenderTemplate(w, "index.html", nil)
		return
	}

	sessionID := cookie.Value
	for key, session := range sessionsMap {
		if session.ID == sessionID {
			delete(sessionsMap, key)
			break
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
