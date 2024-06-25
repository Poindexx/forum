package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	forum "forum/app"

	_ "github.com/mattn/go-sqlite3"
)

// Route определяет структуру маршрута
type Route struct {
	Path    string
	Handler http.HandlerFunc
	Method  string
}

// CustomRouter определяет структуру кастомного маршрутизатора
type CustomRouter struct {
	routes []Route
}

// AddRoute добавляет маршрут в кастомный маршрутизатор
func (cr *CustomRouter) AddRoute(path string, handler http.HandlerFunc, method string) {
	cr.routes = append(cr.routes, Route{Path: path, Handler: handler, Method: method})
}

// ServeHTTP обрабатывает HTTP-запросы и вызывает соответствующие обработчики
func (cr *CustomRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range cr.routes {
		if matchRoute(route.Path, r.URL.Path) && r.Method == route.Method {
			route.Handler.ServeHTTP(w, r)
			return
		}
	}
	if matchRoute("/photo/", r.URL.Path) {
		http.StripPrefix("/photo/", http.FileServer(http.Dir("photo"))).ServeHTTP(w, r)
		return
	}
	http.Redirect(w, r, "/error", http.StatusSeeOther)
}

// matchRoute проверяет, соответствует ли путь шаблону маршрута
func matchRoute(routePath, requestPath string) bool {
	if routePath == requestPath {
		return true
	}

	// Обработка динамического маршрута /Id/{id}
	if strings.HasPrefix(routePath, "/Id/") && strings.HasPrefix(requestPath, "/Id/") {
		return true
	}

	// Обработка динамического маршрута /Id/{id}
	if strings.HasPrefix(routePath, "/photo/") && strings.HasPrefix(requestPath, "/photo/") {
		return true
	}

	// Обработка динамического маршрута /Categorys/{id}
	if strings.HasPrefix(routePath, "/Categorys/") && strings.HasPrefix(requestPath, "/Categorys/") {
		return true
	}

	// Обработка динамического маршрута /Categorys/{id}
	if strings.HasPrefix(routePath, "/Author/") && strings.HasPrefix(requestPath, "/Author/") {
		return true
	}

	return false
}

func main() {
	// Обработка статических файлов
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))

	// Создаем базу данных SQLite
	forum.SetupDB()

	// Создаем кастомный маршрутизатор
	router := &CustomRouter{}

	// Определение маршрутов
	router.AddRoute("/", forum.HomePage, "GET")
	router.AddRoute("/login", forum.LoginPage, "GET")
	router.AddRoute("/register-process", forum.RegisterProcess, "POST")
	router.AddRoute("/login-process", forum.LoginProcess, "POST")
	router.AddRoute("/exit", forum.ExitProcess, "GET")
	router.AddRoute("/create-category", forum.CreateCategoryHandler, "POST")
	router.AddRoute("/get-categories", forum.GetAllCategoryProcess, "GET")
	router.AddRoute("/get-comments", forum.GetAllCommentsProcess, "POST")
	router.AddRoute("/create-post", forum.CreatePostHandler, "POST")
	router.AddRoute("/getSortedPost", forum.GetSortedPostsProcess, "POST")
	router.AddRoute("/GetAllPostsProcess", forum.GetAllPostsProcess, "GET")
	router.AddRoute("/post-dell", forum.PostDelProcess, "POST")
	router.AddRoute("/get-user-data", forum.GetUsersHandler, "GET")
	router.AddRoute("/Id/", forum.HandleIdRequest, "GET")
	router.AddRoute("/Categorys/", forum.HandleCategorysRequest, "GET")
	router.AddRoute("/Author/", forum.HandleAuthorRequest, "GET")
	router.AddRoute("/My_likes", forum.LikePageHandler, "GET")
	router.AddRoute("/My_posts", forum.MyPostsPageHandler, "GET")
	router.AddRoute("/updateLikes", forum.UpdateLikesHandler, "POST")
	router.AddRoute("/updateLikesCom", forum.UpdateLikesComHandler, "POST")
	router.AddRoute("/updateComment", forum.UpdateCommentHandler, "POST")
	router.AddRoute("/login/google", forum.GoogleLoginHandler, "POST")
	router.AddRoute("/auth/google/callback", forum.GithubLoginHandler, "POST")
	router.AddRoute("/login/github", forum.GithubLoginHandler, "POST")
	router.AddRoute("/auth/github/callback", forum.GithubCallbackHandler, "POST")
	router.AddRoute("/error", forum.Error404Handler, "GET")

	// Используем кастомный маршрутизатор
	http.Handle("/", router)

	// Запускаем веб-сервер
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
