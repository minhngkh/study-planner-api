package handler

import (
	"context"
	"errors"
	"study-planner-api/internal/api"
	"study-planner-api/internal/model"
	"study-planner-api/internal/task"
)

// GetTasks implements api.StrictServerInterface.
func (s *Handler) GetTasks(ctx context.Context, request api.GetTasksRequestObject) (api.GetTasksResponseObject, error) {
	authInfo := api.AuthInfoOfRequest(ctx)

	tasks, err := task.GetAllTasks(authInfo.ID)
	if err != nil {
		return nil, err
	}

	apiTasks := make([]api.Task, len(tasks))
	for i, t := range tasks {
		apiTasks[i] = api.Task{
			Id:          &t.ID,
			Name:        &t.Name,
			Description: &t.Description,
			StartTime:   &t.StartTime,
			EndTime:     &t.EndTime,
			Status:      &t.Status,
			UserId:      &t.UserID,
		}
	}

	return api.GetTasks200JSONResponse(apiTasks), nil
}

// PostTasks implements api.StrictServerInterface.
func (s *Handler) PostTasks(ctx context.Context, request api.PostTasksRequestObject) (api.PostTasksResponseObject, error) {
	authInfo := api.AuthInfoOfRequest(ctx)

	newTask := model.Task{
		UserID:        authInfo.ID,
		Name:          request.Body.Name,
		Description:   *request.Body.Description,
		StartTime:     *request.Body.StartTime,
		EndTime:       *request.Body.EndTime,
		Status:        request.Body.Status,
		Priority:      request.Body.Priority,
		EstimatedTime: *request.Body.EstimatedTime,
	}

	resTask, err := task.CreateTask(newTask)
	if err != nil {
		return nil, err
	}

	return api.PostTasks201JSONResponse{
		Id:            &resTask.ID,
		CreatedAt:     &resTask.CreatedAt,
		Description:   &resTask.Description,
		EndTime:       &resTask.EndTime,
		EstimatedTime: &resTask.EstimatedTime,
		Name:          &resTask.Name,
		Priority:      &resTask.Priority,
		StartTime:     &resTask.StartTime,
		Status:        &resTask.Status,
		UpdatedAt:     &resTask.UpdatedAt,
		UserId:        &resTask.UserID,
	}, nil
}

// PutTasksId implements api.StrictServerInterface.
func (s *Handler) PutTasksId(ctx context.Context, request api.PutTasksIdRequestObject) (api.PutTasksIdResponseObject, error) {
	authInfo := api.AuthInfoOfRequest(ctx)

	taskToUpdate := model.Task{
		ID:            request.Id,
		UserID:        authInfo.ID,
		Name:          *request.Body.Name,
		Description:   *request.Body.Description,
		StartTime:     *request.Body.StartTime,
		EndTime:       *request.Body.EndTime,
		Status:        *request.Body.Status,
		Priority:      *request.Body.Priority,
		EstimatedTime: *request.Body.EstimatedTime,
	}

	err := task.UpdateTask(taskToUpdate)
	if err != nil {
		if errors.Is(err, task.ErrTaskNotFound) {
			return api.PutTasksId404JSONResponse{}, nil
		} else {
			return nil, err
		}
	}

	return api.PutTasksId200Response{}, nil
}

// DeleteTasksId implements api.StrictServerInterface.
func (s *Handler) DeleteTasksId(ctx context.Context, request api.DeleteTasksIdRequestObject) (api.DeleteTasksIdResponseObject, error) {
	authInfo := api.AuthInfoOfRequest(ctx)

	err := task.DeleteTaskOfUser(request.Id, authInfo.ID)
	if err != nil {
		if errors.Is(err, task.ErrTaskNotFound) {
			return api.DeleteTasksId404JSONResponse{}, nil
		} else {
			return nil, err
		}
	}

	return api.DeleteTasksId204Response{}, nil
}
