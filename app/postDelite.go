package forum

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func PostDelProcess(w http.ResponseWriter, r *http.Request) {
	var req PostDelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println(err)
		http.Error(w, "Ошибка декодирования запроса", http.StatusBadRequest)
		return
	}

	query := `
		DELETE FROM posts WHERE id = ?;
	`
	result, err := db.Exec(query, req.PostId)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Ошибка выполнения запроса", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Ошибка получения результата выполнения запроса", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Пост не найден", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"success": "Пост успешно удален"})
}
