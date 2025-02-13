package models

import (
	"errors"
	"time"
)

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

var ErrNoRecord = errors.New("models: подходящей записи не найдено")

type ErrorResponse struct {
	Error string `json:"error"`
}
