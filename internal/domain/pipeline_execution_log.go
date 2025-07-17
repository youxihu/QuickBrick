package domain

import "time"

type PipelineExecutionLog struct {
	ID            int64     `json:"id"`
	Env           string    `json:"env"`
	Type          string    `json:"type"` // frontend/backend
	EventType     string    `json:"event_type"`
	PipelineName  string    `json:"pipeline_name"`
	UsernameEmail string    `json:"username_email"`
	CommitID      string    `json:"commit_id"`
	ProjectURL    string    `json:"project_url"`
	Status        string    `json:"status"` // success/failure/running
	CreatedAt     time.Time `json:"created_at"`
}
