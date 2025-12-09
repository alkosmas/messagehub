package storage

import (
	"database/sql"

	"github.com/alkosmas/messagehub/internal/domain"
	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func NewPostgres(connectionString string) (*Repository, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Create table
	query := `
    CREATE TABLE IF NOT EXISTS messages (
        id VARCHAR(36) PRIMARY KEY,
        type VARCHAR(20) NOT NULL,
        recipient VARCHAR(255) NOT NULL,
        body TEXT,
        status VARCHAR(20) NOT NULL,
        provider VARCHAR(50),
        created_at TIMESTAMP,
        error_message TEXT
    );`

	_, err = db.Exec(query)
	if err != nil {
		return nil, err
	}

	return &Repository{db: db}, nil
}

func (r *Repository) Save(msg *domain.Message) error {
	query := `
        INSERT INTO messages (id, type, recipient, body, status, created_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
	_, err := r.db.Exec(query, msg.ID, msg.Type, msg.To, msg.Body, msg.Status, msg.CreatedAt)
	return err
}

func (r *Repository) UpdateStatus(id string, status string, provider string, errorMsg string) error {
	query := `
        UPDATE messages 
        SET status = $1, provider = $2, error_message = $3 
        WHERE id = $4
    `
	_, err := r.db.Exec(query, status, provider, errorMsg, id)
	return err
}
