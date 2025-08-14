package task

import (
	"github.com/barantoraman/microgate/internal/task/repo/entity"
	"github.com/barantoraman/microgate/pkg/validator"
)

func ValidateTask(v *validator.Validator, task entity.Task) {
	v.Check(task.Title != "", "title", "must be provided")
	v.Check(len(task.Title) >= 3, "title", "must be at least 3 characters long")
	v.Check(len(task.Title) <= 80, "title", "must not be more than 80 characters")
	if task.Description != "" {
		v.Check(len(task.Description) <= 150, "description", "must not be more than 150 characters")
	}
	v.Check(task.Status != "", "status", "must be provided")
	v.Check(task.UserID > 0, "user_id", "must be valid")
}
