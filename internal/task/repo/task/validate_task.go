package task

import (
	"github.com/barantoraman/microgate/internal/task/repo/entity"
	"github.com/barantoraman/microgate/pkg/validator"
)

func ValidateTask(v *validator.Validator, task entity.Task) {
	v.Check(task.Title != "", "title", "must be provided")
	v.Check(len(task.Title) <= 80, "title", "must bot be more 80")

	v.Check(task.Description != "", "description", "must be provided")
	v.Check(len(task.Description) <= 150, "description", "must bot be more 80")

	v.Check(task.Status != "", "status", "must be provided")
	v.Check(len(task.Status) <= 80, "status", "must bot be more 80")
}
