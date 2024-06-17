package forum

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func CreatePostTable() {
	// SQL запрос для создания таблицы постов
	createPostTableSQL := `
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT,
			description TEXT,
			anons TEXT,
			author_id INTEGER,
			author TEXT,
			date TEXT,
			image_url TEXT,
			category_id TEXT,
			category TEXT
		)
    `
	// Выполняем SQL запрос
	_, err := db.Exec(createPostTableSQL)
	if err != nil {
		log.Fatalf("Error creating posts table: %v\n", err)
	}
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Парсинг JSON тела запроса
	var post Post
	err := json.NewDecoder(r.Body).Decode(&post)

	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// Проверка наличия заголовка поста
	if post.Title == "" {
		http.Error(w, "Post title is required", http.StatusBadRequest)
		fmt.Println("Post title is required")
		return
	}

	// Проверка количества категорий
	if len(post.Categorys) == 0 {
		http.Error(w, "Выберите хотя бы одну категорию.", http.StatusBadRequest)
		fmt.Println("No categories selected")
		return
	}
	fmt.Println(post)

	// Декодирование Base64 изображения
	imageData, err := base64.StdEncoding.DecodeString(post.ImageBase64)
	if err != nil {
		http.Error(w, "Error decoding Base64 image", http.StatusBadRequest)
		return
	}

	// Создание уникального имени файла и сохранение в папку photo
	fileName := fmt.Sprintf("%d-%s", time.Now().UnixNano(), post.ImageName)
	filePath := filepath.Join("photo", fileName)

	// Логирование пути файла
	log.Printf("Saving file to: %s\n", filePath)

	// Убедиться, что папка существует
	if _, err := os.Stat("photo"); os.IsNotExist(err) {
		log.Println("Creating directory: photo")
		err = os.Mkdir("photo", os.ModePerm)
		if err != nil {
			log.Printf("Error creating directory: %v\n", err)
			http.Error(w, "Unable to create directory", http.StatusInternalServerError)
			return
		}
	}

	// Сохранение файла
	err = ioutil.WriteFile(filePath, imageData, 0644)
	if err != nil {
		log.Printf("Error saving file: %v\n", err)
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}
	post.ImageName = "/" + filePath
	// Сохранение поста в базе данных
	err = SavePost(post)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error saving post: %v", err), http.StatusInternalServerError)
		fmt.Println("Error saving post:", err)
		return
	}

	// Отправка успешного ответа
	w.WriteHeader(http.StatusCreated)
}

func SavePost(post Post) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	categoryIDsStr := strings.Join(post.CategoryIDs, ",")
	categorysStr := strings.Join(post.Categorys, ",")

	insertPostSQL := `
        INSERT INTO posts (title, description, anons, author_id, author, date, image_url, category_id, category)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
    `

	post.Date = time.Now().Format("02.01.2006 15:04:05")
	_, err = tx.Exec(insertPostSQL, post.Title, post.Description, post.Anons, post.AuthorID, post.Author, post.Date, post.ImageName, categoryIDsStr, categorysStr)
	if err != nil {
		return err
	}

	var postID int
	err = tx.QueryRow("SELECT last_insert_rowid()").Scan(&postID)
	if err != nil {
		return err
	}

	insertPostLikeSQL := `INSERT INTO post_likes (post_id) VALUES (?)`
	_, err = tx.Exec(insertPostLikeSQL, postID)
	if err != nil {
		return err
	}

	return nil
}

func GetAllPostsProcess(w http.ResponseWriter, r *http.Request) {
	data_c, err := GetAllPosts()
	if err != nil {
		log.Printf("Error getting posts: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Отправляем полученные категории в формате JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data_c); err != nil {
		log.Printf("Error encoding posts to JSON: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func GetAllPosts() ([]Post1, error) {
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
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post Post1
		var categorys string
		var CategoryIDs string

		if err := rows.Scan(&post.ID, &post.Title, &post.Description, &post.Anons, &post.AuthorID, &post.Author, &post.Date, &post.ImageURL, &CategoryIDs, &categorys, &post.Like, &post.DisLike, &post.PostID, &post.CommentLen); err != nil {
			fmt.Println(err)
			return nil, err
		}

		// Преобразуем строку категорий в массив строк
		post.Categorys = strings.Split(strings.Trim(categorys, "[]"), ",")
		post.CategoryIDs = strings.Split(strings.Trim(CategoryIDs, "[]"), ",")

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	data_b = append(data_b, posts...)
	return posts, nil
}
