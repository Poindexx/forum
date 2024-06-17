package forum

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func CreateCommentLikesTable() {
	createCommentLikesTableSQL := `
        CREATE TABLE IF NOT EXISTS comment_likes (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            comment_id INTEGER,
            user_id INTEGER DEFAULT 0,
            like INTEGER DEFAULT 0,
            dislike INTEGER DEFAULT 0,
			state INTEGER DEFAULT 0,
            FOREIGN KEY(comment_id) REFERENCES comments(id),
            FOREIGN KEY(user_id) REFERENCES users(id)
        )
    `
	_, err := db.Exec(createCommentLikesTableSQL)
	if err != nil {
		log.Fatalf("Error creating comment_likes table: %v\n", err)
	}
}

func UpdateLikesComHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var req LikeUpdateComRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var oldLikeCount, oldDislikeCount, oldUser_id, state int
	row1 := db.QueryRow("SELECT like, dislike, user_id, state FROM comment_likes WHERE comment_id = ? ORDER BY id DESC LIMIT 1", req.ComID)
	_ = row1.Scan(&oldLikeCount, &oldDislikeCount, &oldUser_id, &state)

	var query string
	if oldUser_id == 0 {
		if req.Type == "like" {
			oldLikeCount++
			state = 1
		} else if req.Type == "dislike" {
			oldDislikeCount++
			state = 2
		}
	} else {
		if oldUser_id == req.AuthorID {
			query = `
			DELETE FROM comment_likes
			WHERE comment_id = ? AND user_id = ?
			`
			_, err := db.Exec(query, req.ComID, req.AuthorID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			switch req.Type {
			case "like":
				switch state {
				case 0:
					oldLikeCount++
					state = 1
				case 1:
					oldLikeCount--
					state = 0
				case 2:
					oldLikeCount++
					oldDislikeCount--
					state = 1
				}
			case "dislike":
				switch state {
				case 0:
					oldDislikeCount++
					state = 2
				case 1:
					oldDislikeCount++
					oldLikeCount--
					state = 2
				case 2:
					oldDislikeCount--
					state = 0
				}
			}
		} else {
			var state1 int
			row2 := db.QueryRow("SELECT state FROM comment_likes WHERE comment_id = ? AND user_id = ? ORDER BY id DESC LIMIT 1", req.ComID, req.AuthorID)
			err1 := row2.Scan(&state1)
			tt := err1 == nil

			if tt {
				state = state1
				query = `
				DELETE FROM comment_likes
				WHERE comment_id = ? AND user_id = ?
				`
				_, err := db.Exec(query, req.ComID, req.AuthorID)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				state = 0
			}
			switch req.Type {
			case "like":
				switch state {
				case 0:
					oldLikeCount++
					state = 1
				case 1:
					oldLikeCount--
					state = 0
				case 2:
					oldLikeCount++
					oldDislikeCount--
					state = 1
				}
			case "dislike":
				switch state {
				case 0:
					oldDislikeCount++
					state = 2
				case 1:
					oldDislikeCount++
					oldLikeCount--
					state = 2
				case 2:
					oldDislikeCount--
					state = 0
				}
			}
		}
	}

	_, err := db.Exec("INSERT INTO comment_likes (comment_id, user_id, like, dislike, state) VALUES (?, ?, ?, ?, ?)", req.ComID, req.AuthorID, oldLikeCount, oldDislikeCount, state)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	idddi, _ := strconv.Atoi(req.ComID)
	resp := LikeComUpdateResponse{
		NewLikeCount:    oldLikeCount,
		NewDislikeCount: oldDislikeCount,
		NewComid:        idddi,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
