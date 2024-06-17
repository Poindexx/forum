package forum

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func HandleCategorysRequest(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/Categorys/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Categorys id", http.StatusBadRequest)
		return
	}

	CategorysIdHandler(w, r, id)
}

func CategorysIdHandler(w http.ResponseWriter, r *http.Request, id_Cat int) {
	var posts []Post1

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
	ORDER BY 
		pd.id DESC;
    `

	rows, err := db.Query(query)
	title_page := ""
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var post Post1
		var categorys string
		var CategoryIDs string

		if err := rows.Scan(&post.ID, &post.Title, &post.Description, &post.Anons, &post.AuthorID, &post.Author, &post.Date, &post.ImageURL, &CategoryIDs, &categorys, &post.Like, &post.DisLike, &post.PostID, &post.CommentLen); err != nil {
			fmt.Println(err)
			return
		}

		// Преобразуем строку категорий в массив строк
		post.Categorys = strings.Split(strings.Trim(categorys, "[]"), ",")
		post.CategoryIDs = strings.Split(strings.Trim(CategoryIDs, "[]"), ",")
		for i := range post.CategoryIDs {
			id_Cat1, _ := strconv.Atoi(post.CategoryIDs[i])
			if id_Cat1 == id_Cat {
				title_page = post.Categorys[i]
				posts = append(posts, post)
			}
		}

	}

	if err := rows.Err(); err != nil {
		fmt.Println(err)
		return
	}
	// Получение идентификатора сессии из куки
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		data := ViewData{Username1: "", Id1: "", Title: title_page, Posts: posts}
		RenderTemplate(w, "postsCatId.html", data)
		return
	}
	sessionID := cookie.Value
	session, ok := sessionsMap[sessionID]
	if !ok {
		// Если сессия не найдена в карте, отправляем шаблон без данных о сессии
		data := ViewData{Username1: "", Id1: "", Title: title_page, Posts: posts}
		RenderTemplate(w, "postsCatId.html", data)
		return
	}

	// Проверка соответствия идентификатора сессии идентификатору в куке
	if session.ID != sessionID {
		// Если идентификаторы не совпадают, считаем сессию недействительной
		data := ViewData{Username1: "", Id1: "", Title: title_page, Posts: posts}
		RenderTemplate(w, "postsCatId.html", data)
		return
	}

	data := ViewData{Username1: session.Username, Id1: session.ID, Title: title_page, Posts: posts}

	RenderTemplate(w, "postsCatId.html", data)

}
