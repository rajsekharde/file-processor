package shared

import (
	"encoding/json"
	"time"
)

// Redis hash
type Task struct {
	ID string `json:"id"`
	Type string `json:"type"`
	Payload json.RawMessage `json:"payload"`
	Status string `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// POST /task response
type TaskRequest struct {
	Type string `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// Payload for convert format
type ConvertFormatPayload struct {
	FileName string `json:"file_name"`
	OutputFormat string `json:"output_format"`
}

// Payload for resize image
type ResizeImagePayload struct {
	FileName string `json:"file_name"`
	TargetWidth int `json:"target_width"`
	TargetHeight int `json:"target_height"`
}