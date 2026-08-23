package shared

import (
	"encoding/json"
	"time"
)

type Task struct {
	ID string `json:"id"`
	Type string `json:"type"`
	Payload json.RawMessage `json:"payload"`
	Status string `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type ConvertFormatPayload struct {
	FileName string `json:"file_name"`
	OutputFormat string `json:"output_format"`
}

type TaskRequest struct {
	Type string `json:"type"`
	Payload json.RawMessage `json:"payload"`
}