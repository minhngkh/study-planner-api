package handler

import (
	"context"
	"errors"
	"study-planner-api/internal/api"
	"study-planner-api/internal/model"
	"study-planner-api/internal/task"
)

const (
	TaskSortByDefault        = task.SortFieldCreatedAt
	TaskSortOrderDefault     = task.SortOrderDesc
	TaskPageDefault      int = 1
	TaskLimitDefault     int = 10
)

func (s *Handler) GetTasks(ctx context.Context, request api.GetTasksRequestObject) (api.GetTasksResponseObject, error) {
	authInfo := api.AuthInfoOfRequest(ctx)

	criteria := task.GetCriteria{
		UserID: authInfo.ID,
		Search: request.Params.Search,
	}

	if request.Params.Page != nil {
		criteria.Pagination.Page = *request.Params.Page
	} else {
		criteria.Pagination.Page = TaskPageDefault
	}

	if request.Params.Limit != nil {
		criteria.Pagination.Limit = *request.Params.Limit
	} else {
		criteria.Pagination.Limit = TaskLimitDefault
	}

	if request.Params.Status != nil {
		status, err := task.StatusFromString(*request.Params.Status)
		if err != nil {
			return api.GetTasks400Response{}, nil
		}

		criteria.Status = &status
	}

	if request.Params.Priority != nil {
		priority, err := task.PriorityFromString(*request.Params.Priority)
		if err != nil {
			return api.GetTasks400Response{}, nil
		}

		criteria.Priority = &priority
	}

	if request.Params.StartDate != nil {
		criteria.StartTime = &request.Params.StartDate.Time
	}

	if request.Params.EndDate != nil {
		criteria.EndTime = &request.Params.EndDate.Time
	}

	if request.Params.SortBy != nil {
		sortBy, err := task.SortFieldFromString(string(*request.Params.SortBy))
		if err != nil {
			return api.GetTasks400Response{}, nil
		}

		criteria.SortType.Field = sortBy
	} else {
		criteria.SortType.Field = TaskSortByDefault
	}

	if request.Params.SortOrder != nil {
		sortOrder, err := task.SortOrderFromString(string(*request.Params.SortOrder))
		if err != nil {
			return api.GetTasks400Response{}, nil
		}

		criteria.SortType.Order = sortOrder
	} else {
		criteria.SortType.Order = TaskSortOrderDefault
	}

	tasks, err := task.GetTasks(&criteria)
	if err != nil {
		return nil, err
	}

	apiTasks := make([]api.Task, len(tasks))
	for i, t := range tasks {
		apiTasks[i] = api.Task{
			Id:            &t.ID,
			Name:          &t.Name,
			Description:   t.Description,
			StartTime:     t.StartTime,
			EndTime:       t.EndTime,
			Status:        &t.Status,
			UserId:        t.UserID,
			CreatedAt:     t.CreatedAt,
			UpdatedAt:     t.UpdatedAt,
			EstimatedTime: t.EstimatedTime,
			Priority:      &t.Priority,
		}
	}

	return api.GetTasks200JSONResponse{
		Data: &apiTasks,
		Pagination: &api.PaginationResponse{
			Total:      &criteria.Pagination.Total,
			Page:       &criteria.Pagination.Page,
			Limit:      &criteria.Pagination.Limit,
			TotalPages: &criteria.Pagination.TotalPages,
		},
	}, nil
}

// PostTasks implements api.StrictServerInterface.
func (s *Handler) PostTasks(ctx context.Context, request api.PostTasksRequestObject) (api.PostTasksResponseObject, error) {
	authInfo := api.AuthInfoOfRequest(ctx)

	newTask := model.Task{
		UserID:        &authInfo.ID,
		Name:          request.Body.Name,
		Description:   request.Body.Description,
		StartTime:     request.Body.StartTime,
		EndTime:       request.Body.EndTime,
		Status:        request.Body.Status,
		Priority:      request.Body.Priority,
		EstimatedTime: request.Body.EstimatedTime,
	}

	resTask, err := task.CreateTask(newTask)
	if err != nil {
		return nil, err
	}

	return api.PostTasks201JSONResponse{
		Id:            &resTask.ID,
		CreatedAt:     resTask.CreatedAt,
		Description:   resTask.Description,
		EndTime:       resTask.EndTime,
		EstimatedTime: resTask.EstimatedTime,
		Name:          &resTask.Name,
		Priority:      &resTask.Priority,
		StartTime:     resTask.StartTime,
		Status:        &resTask.Status,
		UpdatedAt:     resTask.UpdatedAt,
		UserId:        resTask.UserID,
	}, nil
}

// PutTasksId implements api.StrictServerInterface.
func (s *Handler) PutTasksId(ctx context.Context, request api.PutTasksIdRequestObject) (api.PutTasksIdResponseObject, error) {
	authInfo := api.AuthInfoOfRequest(ctx)

	taskToUpdate := model.Task{
		ID:     request.Id,
		UserID: &authInfo.ID,
	}

	if request.Body.Name != nil {
		taskToUpdate.Name = *request.Body.Name
	}
	if request.Body.Description != nil {
		taskToUpdate.Description = request.Body.Description
	}
	if request.Body.Priority != nil {
		taskToUpdate.Priority = *request.Body.Priority
	}
	if request.Body.EstimatedTime != nil {
		taskToUpdate.EstimatedTime = request.Body.EstimatedTime
	}
	if request.Body.Status != nil {
		taskToUpdate.Status = *request.Body.Status
	}
	if request.Body.StartTime != nil {
		taskToUpdate.StartTime = request.Body.StartTime
	}
	if request.Body.EndTime != nil {
		taskToUpdate.EndTime = request.Body.EndTime
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
