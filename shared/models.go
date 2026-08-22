package shared

import "time"

type Task struct {
	ID string `json:"id"`
	FileName string `json:"file_name"`
	OutputFormat string `json:"output_format"`
	Status string `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type TaskRequest struct {
	FileName string `json:"file_name"`
	OutputFormat string `json:"output_format"`
}