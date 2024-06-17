package forum

import (
	"fmt"
	"net/http"
	"strings"
)

func MyPostsPageHandler(w http.ResponseWriter, r *http.Request) {

	// Получение идентификатора сессии из куки
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		data := ViewData{Username1: "", Id1: "", Title: "Не авторизованный пользователь"}
		RenderTemplate(w, "myPostsPosts.html", data)
		return
	}
	sessionID := cookie.Value
	session, ok := sessionsMap[sessionID]
	if !ok {
		// Если сессия не найдена в карте, отправляем шаблон без данных о сессии
		data := ViewData{Username1: "", Id1: "", Title: "Не авторизованный пользователь"}
		RenderTemplate(w, "myPostsPosts.html", data)
		return
	}

	// Проверка соответствия идентификатора сессии идентификатору в куке
	if session.ID != sessionID {
		// Если идентификаторы не совпадают, считаем сессию недействительной
		data := ViewData{Username1: "", Id1: "", Title: "Не авторизованный пользователь"}
		RenderTemplate(w, "myPostsPosts.html", data)
		return
	}

	var posts []Post1

	usersquery := `
		SELECT id
		FROM users
		WHERE username = ?;
	`
	rows, err := db.Query(usersquery, session.Username)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	var id_user int
	for rows.Next() {
		if err := rows.Scan(&id_user); err != nil {
			fmt.Println(err)
			return
		}
	}
	if err := rows.Err(); err != nil {
		fmt.Println(err)
		return
	}
	
	query := `
	WITH posts_data AS (
		SELECT id, title, description, anons, author_id, author, date, image_url, category_id, category
		FROM posts
	),
	likes_data AS (
		SELECT post_id, like, dislike,
			ROW_NUMBER() OVER (PARTITION BY post_id ORDER BY id DESC) AS rn
		FROM post_likes
	),
	comments_data AS (
		SELECT post_id, COUNT(*) AS comment_count
		FROM comments
		GROUP BY post_id
	)
	SELECT 
		pd.id, pd.title, pd.description, pd.anons, pd.author_id, pd.author, pd.date, pd.image_url, pd.category_id, pd.category,
		ld.like, ld.dislike, ld.post_id, 
		COALESCE(cd.comment_count, 0) AS comment_count
	FROM 
		posts_data pd
	LEFT JOIN 
		likes_data ld ON pd.id = ld.post_id AND ld.rn = 1
	LEFT JOIN 
		comments_data cd ON pd.id = cd.post_id
	WHERE 
        pd.author_id = ?
	ORDER BY 
		pd.id DESC;
    `

	rows2, err := db.Query(query, id_user)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows2.Close()

	for rows2.Next() {
		var post Post1
		var categorys string
		var CategoryIDs string

		if err := rows2.Scan(&post.ID, &post.Title, &post.Description, &post.Anons, &post.AuthorID, &post.Author, &post.Date, &post.ImageURL, &CategoryIDs, &categorys, &post.Like, &post.DisLike, &post.PostID, &post.CommentLen); err != nil {
			fmt.Println(err)
			return
		}

		// Преобразуем строку категорий в массив строк
		post.Categorys = strings.Split(strings.Trim(categorys, "[]"), ",")
		post.CategoryIDs = strings.Split(strings.Trim(CategoryIDs, "[]"), ",")
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		fmt.Println(err)
		return
	}

	data := ViewData{Username1: session.Username, Id1: session.ID, Title: "Мои посты", Posts: posts}

	RenderTemplate(w, "myPostsPosts.html", data)
}
