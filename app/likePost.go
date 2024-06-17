package forum

import (
	"encoding/json"
	"log"
	"net/http"
)

func CreatePostLikesTable() {
	createPostLikesTableSQL := `
        CREATE TABLE IF NOT EXISTS post_likes (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            post_id INTEGER,
            user_id INTEGER DEFAULT 0,
            like INTEGER DEFAULT 0,
            dislike INTEGER DEFAULT 0,
            state INTEGER DEFAULT 0,
            FOREIGN KEY(post_id) REFERENCES posts(id),
            FOREIGN KEY(user_id) REFERENCES users(id)
        )
    `
	_, err := db.Exec(createPostLikesTableSQL)
	if err != nil {
		log.Fatalf("Error creating post_likes table: %v\n", err)
	}
}

func UpdateLikesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req LikeUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var oldLikeCount, oldDislikeCount, oldUser_id, state int
	row1 := db.QueryRow("SELECT like, dislike, user_id, state FROM post_likes WHERE post_id = ? ORDER BY id DESC LIMIT 1", req.PostID)
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
			DELETE FROM post_likes
			WHERE post_id = ? AND user_id = ?
			`
			_, err := db.Exec(query, req.PostID, req.AuthorID)
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
			row2 := db.QueryRow("SELECT state FROM post_likes WHERE post_id = ? AND user_id = ? ORDER BY id DESC LIMIT 1", req.PostID, req.AuthorID)
			err1 := row2.Scan(&state1)
			tt := err1 == nil

			if tt {
				state = state1
				query = `
				DELETE FROM post_likes
				WHERE post_id = ? AND user_id = ?
				`
				_, err := db.Exec(query, req.PostID, req.AuthorID)
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

	_, err := db.Exec("INSERT INTO post_likes (post_id, user_id, like, dislike, state) VALUES (?, ?, ?, ?, ?)", req.PostID, req.AuthorID, oldLikeCount, oldDislikeCount, state)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := LikeUpdateResponse{
		NewLikeCount:    oldLikeCount,
		NewDislikeCount: oldDislikeCount,
		NewPostid:       req.PostID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
