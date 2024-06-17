package forum

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func GetSortedPostsProcess(w http.ResponseWriter, r *http.Request) {
	var req SortedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var posts []Post1
	var filters []string
	var params []interface{}

	if len(req.Category_ids) > 0 {
		categoryConditions := make([]string, len(req.Category_ids))
		for i, id := range req.Category_ids {
			categoryConditions[i] = "pd.category_id LIKE ?"
			params = append(params, "%"+id+"%")
		}
		filters = append(filters, fmt.Sprintf("(%s)", strings.Join(categoryConditions, " OR ")))
	}

	// Фильтрация по тексту
	if req.Text_dis != "" {
		filters = append(filters, "(pd.title LIKE ? OR pd.description LIKE ? OR pd.anons LIKE ?)")
		likeText := "%" + req.Text_dis + "%"
		params = append(params, likeText, likeText, likeText)
	}

	// Фильтрация по дате
	if req.Start_date != "" && req.End_date != "" {
		filters = append(filters, "pd.date BETWEEN ? AND ?")
		params = append(params, req.Start_date, req.End_date)
	} else if req.Start_date != "" {
		filters = append(filters, "pd.date >= ?")
		params = append(params, req.Start_date)
	} else if req.End_date != "" {
		filters = append(filters, "pd.date <= ?")
		params = append(params, req.End_date)
	}

	// Основной SQL-запрос
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
    `

	if len(filters) > 0 {
		query += "WHERE " + strings.Join(filters, " AND ")
	}

	// Сортировка результатов
	switch req.Sort_post {
	case "1":
		query += " ORDER BY pd.id DESC"
	case "2":
		query += " ORDER BY ld.like DESC"
	case "3":
		query += " ORDER BY cd.comment_count DESC"
	default:
		query += " ORDER BY pd.id DESC"
	}

	rows, err := db.Query(query, params...)
	if err != nil {
		fmt.Println("Error executing query:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var post Post1
		var categorys string
		var CategoryIDs string

		if err := rows.Scan(&post.ID, &post.Title, &post.Description, &post.Anons, &post.AuthorID, &post.Author, &post.Date, &post.ImageURL, &CategoryIDs, &categorys, &post.Like, &post.DisLike, &post.PostID, &post.CommentLen); err != nil {
			fmt.Println("Error scanning row:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		post.Categorys = strings.Split(strings.Trim(categorys, "[]"), ",")
		post.CategoryIDs = strings.Split(strings.Trim(CategoryIDs, "[]"), ",")
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Rows error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправка результата как JSON-ответ клиенту
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		fmt.Println("Error encoding response:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
