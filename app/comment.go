package forum

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func CreateCommentsTable() {
	createCommentsTableSQL := `
        CREATE TABLE IF NOT EXISTS comments (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            post_id INTEGER,
            user_id INTEGER,
            user TEXT,
            comment TEXT,
            created_at TEXT,
            FOREIGN KEY(post_id) REFERENCES posts(id),
            FOREIGN KEY(user_id) REFERENCES users(id)
            FOREIGN KEY(user) REFERENCES users(username)
        )
    `
	_, err := db.Exec(createCommentsTableSQL)
	if err != nil {
		log.Fatalf("Error creating comments table: %v\n", err)
	}
}

func UpdateCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	tx, err1 := db.Begin()
	if err1 != nil {
		fmt.Println(err1)
	}
	defer func() {
		if err1 != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	var comm Comments
	err := json.NewDecoder(r.Body).Decode(&comm)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	dateComm := time.Now().Format("02.01.2006 15:04:05")
	insertCommentSQL := `INSERT INTO comments (post_id, user_id, user, comment, created_at) VALUES (?, ?, ?, ?, ?)`
	_, err = tx.Exec(insertCommentSQL, comm.PostID, comm.AuthorID, comm.Author, comm.Comment, dateComm)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error inserting comment", http.StatusInternalServerError)
		return
	}

	var commentID int
	err = tx.QueryRow("SELECT last_insert_rowid()").Scan(&commentID)
	if err != nil {
		fmt.Println(err)
	}

	insertPostLikeSQL := `INSERT INTO comment_likes (comment_id) VALUES (?)`
	_, err = tx.Exec(insertPostLikeSQL, commentID)
	if err != nil {
		fmt.Println(err)
	}

	resp := CommResponse{
		Com: "ss",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
	}
}

func GetAllCommentsProcess(w http.ResponseWriter, r *http.Request) {
	var req CommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	data_c, err := GetAllComments(req.PostID)
	if err != nil {
		log.Printf("Error getting posts: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data_c); err != nil {
		log.Printf("Error encoding posts to JSON: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func GetAllComments(postID string) ([]Comments1, error) {
	var commentss []Comments1

	query := `
        WITH RankedCommentLikes AS (
            SELECT 
                comment_id, 
                like, 
                dislike, 
                ROW_NUMBER() OVER (PARTITION BY comment_id ORDER BY id DESC) AS rn
            FROM 
                comment_likes
        )
        SELECT 
            c.id AS comment_id, 
            c.post_id, 
            c.user_id, 
            c.user, 
            c.comment, 
            c.created_at, 
            rcl.like, 
            rcl.dislike
        FROM 
            comments c
        JOIN 
            RankedCommentLikes rcl ON c.id = rcl.comment_id
        WHERE 
            c.post_id = ? AND rcl.rn = 1
		ORDER BY 
			c.id DESC;
    `

	rows, err := db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment Comments1

		if err := rows.Scan(&comment.Id, &comment.PostID, &comment.AuthorID, &comment.Author, &comment.Comment, &comment.Date, &comment.Like, &comment.DisLike); err != nil {
			return nil, err
		}

		commentss = append(commentss, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return commentss, nil
}
