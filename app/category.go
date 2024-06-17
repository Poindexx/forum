package forum

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func CreateCategoryTable() {
	createCategoryTableSQL := `
        CREATE TABLE IF NOT EXISTS categories (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT UNIQUE
        )
    `
	_, err := db.Exec(createCategoryTableSQL)
	if err != nil {
		log.Fatalf("Error creating categories table: %v\n", err)
	}
}

func CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var category Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	// Проверка наличия названия категории
	if category.Name == "" {
		http.Error(w, "Category name is required", http.StatusBadRequest)
		return
	}

	// Добавление категории в базу данных
	err = AddCategory(category.Name)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error adding category: %v", err), http.StatusInternalServerError)
		return
	}

	// Отправка успешного ответа
	w.WriteHeader(http.StatusCreated)
}

func AddCategory(name string) error {
	insertCategorySQL := `INSERT INTO categories (name) VALUES (?)`
	_, err := db.Exec(insertCategorySQL, name)
	if err != nil {
		// Проверка, является ли ошибка нарушением уникального ограничения
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return fmt.Errorf("Category with the same name already exists")
		}
		return err
	}
	return nil
}

func GetAllCategoryProcess(w http.ResponseWriter, r *http.Request) {
	categories, err := GetAllCategories()
	if err != nil {
		log.Printf("Error getting categories: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Отправляем полученные категории в формате JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(categories); err != nil {
		log.Printf("Error encoding categories to JSON: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func GetAllCategories() ([]Category, error) {
	var categories []Category

	// SQL запрос для выборки всех категорий
	query := "SELECT id, name FROM categories"

	// Выполняем SQL запрос
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Итерируем по результатам запроса
	for rows.Next() {
		var category Category
		// Сканируем данные из строки результата в структуру Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, err
		}
		// Добавляем категорию в слайс
		categories = append(categories, category)
	}

	// Проверяем ошибки, которые могли возникнуть в процессе итерации по результатам
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}
