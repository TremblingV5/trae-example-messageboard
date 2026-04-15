package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

type Message struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

var db *sql.DB

func InitDB() error {
	var err error
	db, err = sql.Open("sqlite", "messages.db")
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Create messages table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}

func GetMessages() ([]Message, error) {
	rows, err := db.Query("SELECT id, content, created_at FROM messages ORDER BY created_at DESC")
	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %w", err)
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		err := rows.Scan(&msg.ID, &msg.Content, &msg.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return messages, nil
}

func CreateMessage(content string) (Message, error) {
	result, err := db.Exec("INSERT INTO messages (content) VALUES (?)", content)
	if err != nil {
		return Message{}, fmt.Errorf("failed to insert message: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Message{}, fmt.Errorf("failed to get last insert id: %w", err)
	}

	// Get the created message with timestamp
	var msg Message
	err = db.QueryRow("SELECT id, content, created_at FROM messages WHERE id = ?", id).Scan(&msg.ID, &msg.Content, &msg.CreatedAt)
	if err != nil {
		return Message{}, fmt.Errorf("failed to get created message: %w", err)
	}

	return msg, nil
}
