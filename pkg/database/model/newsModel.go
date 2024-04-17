package model

import (
	"newsaggr/pkg/database"
	"regexp"
)

type News struct {
	ID      int    // номер записи
	Title   string // заголовок публикации
	Content string `gorm:"unique"` // содержание публикации
	PubTime int64  // время публикации
	Link    string // ссылка на источник
}

// FindAll - все записи
func FindAll() ([]*News, error) {
	db := database.GetDB()
	query := "SELECT * FROM news"
	rows, err := db.Raw(query).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var news []*News

	for rows.Next() {
		row := &News{}
		err := rows.Scan(&row.ID, &row.Title, &row.Content, &row.PubTime, &row.Link)
		if err != nil {
			return nil, err
		}
		news = append(news, row)
	}

	return news, nil
}

// Create - создание новой категории
func (n *News) Create() error {
	db := database.GetDB()
	// Регулярное выражение для удаления HTML тегов
	htmlTagRegex := regexp.MustCompile("<[^>]*>")

	// Удаление HTML тегов из строки
	title := htmlTagRegex.ReplaceAllString(n.Title, "")
	content := htmlTagRegex.ReplaceAllString(n.Content, "")

	query := "INSERT INTO news (title, content, pub_time, link) VALUES (?, ?, ?, ?)"
	result := db.Exec(query, title, content, n.PubTime, n.Link)
	return result.Error
}

// FindOne - поиск записи по id
func (n *News) FindOne(id int) error {
	db := database.GetDB()
	query := "SELECT * FROM news WHERE (id = ?)"
	result := db.Raw(query, id).Scan(&n)
	return result.Error
}

// FindLimit - поиск последних N записей
func FindLimit(limit int) ([]*News, error) {
	db := database.GetDB()
	query := "SELECT * FROM news ORDER BY id DESC LIMIT ?;"
	var news []*News

	rows, err := db.Raw(query, limit).Rows()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		row := &News{}
		err := rows.Scan(&row.ID, &row.Title, &row.Content, &row.PubTime, &row.Link)
		if err != nil {
			return nil, err
		}
		news = append(news, row)
	}

	return news, nil
}

// Delete - удаление новости по ID
func Delete(id uint) error {
	db := database.GetDB()
	query := "DELETE FROM news WHERE id = ?"
	result := db.Exec(query, id)

	return result.Error
}
