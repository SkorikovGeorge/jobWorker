package structs

import (
	"github.com/SkorikovGeorge/jobWorker/internal/consts"
	"github.com/google/uuid"
)

type Job struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Priority int    `json:"priority"`
}

func NewJob(priority int) *Job {
	id := uuid.New().String()
	return &Job{
		ID:       id,
		Status:   consts.StatusPending,
		Priority: priority,
	}
}
