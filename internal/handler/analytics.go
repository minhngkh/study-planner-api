package handler

import (
	"context"
	"study-planner-api/internal/api"
	"study-planner-api/internal/database"
	"study-planner-api/internal/focussession"
	"study-planner-api/internal/model"
	"study-planner-api/internal/task"
	"time"
)

// TODO: move services code to analytics package
func (s *Handler) GetAnalyticsFocus(ctx context.Context, request api.GetAnalyticsFocusRequestObject) (api.GetAnalyticsFocusResponseObject, error) {
	authInfo := api.AuthInfoOfRequest(ctx)
	userID := authInfo.ID

	// Apply date filters if present
	startDate := request.Params.StartDate
	endDate := request.Params.EndDate

	query := database.Instance().
		Model(&model.FocusSession{}).
		Joins("JOIN task ON focus_session.task_id = task.id").
		Select("COALESCE(SUM(focus_duration), 0) as total, COALESCE(SUM(estimated_time) * 60, 0) as estimated").
		Where("task.user_id = ?", userID).
		Where("focus_session.status <> ?", focussession.StatusActive.String())
	if startDate != nil {
		query = query.Where("focus_session.created_at >= ?", startDate.Time)
	}
	if endDate != nil {
		query = query.Where("focus_session.created_at <= ?", endDate.Time.Add(time.Hour*24))
	}

	var timeStats struct {
		Total     int32
		Estimated int32
	}
	err := query.First(&timeStats).Error
	if err != nil {
		return nil, err
	}

	// Daily time spent
	type DailyTime struct {
		Date  string
		Total int32
	}
	dailyQuery := database.Instance().
		Model(&model.FocusSession{}).
		Select("DATE(created_at) as date, COALESCE(SUM(focus_duration), 0) as total").
		Where("focus_session.status <> ?", focussession.StatusActive.String()).
		Group("DATE(created_at)")

	var dailyTimes []DailyTime
	err = dailyQuery.Scan(&dailyTimes).Error
	if err != nil {
		return nil, err
	}

	dailyTimeSpent := make(map[string]int)
	for _, dt := range dailyTimes {
		dailyTimeSpent[dt.Date] = int(dt.Total)
	}

	// Task status counts
	type TaskStatusCount struct {
		Status string
		Count  int
	}
	var taskStatusCounts []TaskStatusCount
	err = database.Instance().
		Model(&model.Task{}).
		Where("user_id = ?", userID).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&taskStatusCounts).Error
	if err != nil {
		return nil, err
	}

	statusCounts := make(map[string]int)
	for _, tsc := range []string{
		string(task.StatusTodo),
		string(task.StatusInProgress),
		string(task.StatusCompleted),
		string(task.StatusExpired),
	} {
		statusCounts[tsc] = 0
	}

	for _, tsc := range taskStatusCounts {
		statusCounts[tsc.Status] = tsc.Count
	}

	return api.GetAnalyticsFocus200JSONResponse{
		TotalTimeSpent:     &timeStats.Total,
		TotalEstimatedTime: &timeStats.Estimated,
		DailyTimeSpent:     &dailyTimeSpent,
		TaskStatusCounts:   &statusCounts,
	}, nil
}
