package task

import (
	"errors"
	db "study-planner-api/internal/database"
	"study-planner-api/internal/model"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

func CreateTask(task model.Task) (*model.Task, error) {
	result := db.Instance().
		Select("UserID", "Name", "Description", "Priority", "EstimatedTime", "Status", "StartTime", "EndTime").
		Create(&task)
	if result.Error != nil {
		return new(model.Task), result.Error
	}

	return &task, nil
}

func UpdateTask(task model.Task) error {
	result := db.Instance().
		Model(&model.Task{}).
		Where("id = ?", task.ID).
		Updates(&task).
		Limit(1)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrTaskNotFound
	}

	return nil
}

func GetAllTasks(userID int32) ([]model.Task, error) {
	var tasks []model.Task
	result := db.Instance().
		Model(&model.Task{}).
		Where("user_id = ?", userID).
		Find(&tasks)
	if result.Error != nil {
		return nil, result.Error
	}

	return tasks, nil
}

func DeleteTaskOfUser(taskId int32, userId int32) error {
	result := db.Instance().
		Where("id = ? and user_id = ?", taskId, userId).
		Delete(&model.Task{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrTaskNotFound
	}

	return nil
}
