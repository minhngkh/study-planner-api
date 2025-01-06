package focussession

import (
	"errors"
	"study-planner-api/internal/database"
	"study-planner-api/internal/model"
	"study-planner-api/internal/task"

	"gorm.io/gorm/clause"
)

type Status string

const (
	StatusActive      Status = "active"
	StatusCompleted   Status = "completed"
	StatusEndedEearly Status = "ended_early"
)

func (s Status) String() string {
	return string(s)
}

var (
	ErrInvalidTimerDuration = errors.New("timer duration must be greater than 0")
	ErrTaskNotFound         = errors.New("task not found")
	ErrTaskNotBelongToUser  = errors.New("task does not belong to user")
	ErrTaskNotInProgress    = errors.New("task is not in progress")

	ErrSessionNotFound        = errors.New("session not found")
	ErrSessionNotActive       = errors.New("session is not active")
	ErrSessionNotBelongToUser = errors.New("session does not belong to user")
)

type NewSession struct {
	UserID        int32
	TaskID        int32
	TimerDuration int32
	BreakDuration *int32
}

func CreateSession(session NewSession) (model.FocusSession, error) {
	if session.TimerDuration <= 0 {
		return model.FocusSession{}, ErrInvalidTimerDuration
	}

	var taskInfo model.Task
	result := database.Instance().
		Model(&model.Task{}).
		Where("id = ?", session.TaskID).
		First(&taskInfo)
	if result.Error != nil {
		return model.FocusSession{}, result.Error
	}
	// if result.RowsAffected != 0 {
	// 	return model.FocusSession{}, ErrTaskNotFound
	// }

	if taskInfo.UserID == nil {
		return model.FocusSession{}, ErrTaskNotBelongToUser
	}

	if *taskInfo.UserID != session.UserID {
		return model.FocusSession{}, ErrTaskNotBelongToUser
	}
	if task.Status(taskInfo.Status) != task.StatusInProgress {
		return model.FocusSession{}, ErrTaskNotInProgress
	}

	newSession := model.FocusSession{
		TaskID:        &session.TaskID,
		TimerDuration: session.TimerDuration,
		Status:        string(StatusActive),
	}
	if session.BreakDuration != nil {
		newSession.BreakDuration = session.BreakDuration
	}

	result = database.Instance().
		Model(&model.FocusSession{}).
		Create(&newSession)

	if result.Error != nil {
		return model.FocusSession{}, result.Error
	}

	return newSession, nil
}

type EndEarly struct {
	FocusDuration int32
}

type SessionToEnd struct {
	UserID     int32
	SessionID  int32
	EndedEarly *EndEarly
}

func EndSession(session SessionToEnd) (model.FocusSession, error) {
	if session.EndedEarly != nil && session.EndedEarly.FocusDuration <= 0 {
		return model.FocusSession{}, errors.New("invalid focus duration")
	}

	var sessionInfo struct {
		UserID int32
		model.FocusSession
	}
	result := database.Instance().
		Model(&model.FocusSession{}).
		Select("task.user_id, focus_session.*").
		Joins("INNER JOIN task ON focus_session.task_id = task.id").
		Where("focus_session.id = ?", session.SessionID).
		First(&sessionInfo)
	if result.Error != nil {
		return model.FocusSession{}, result.Error
	}
	// if result.RowsAffected != 0 {
	// 	return model.FocusSession{}, ErrSessionNotFound
	// }
	if sessionInfo.UserID != session.UserID {
		return model.FocusSession{}, ErrSessionNotBelongToUser
	}
	if Status(sessionInfo.Status) != StatusActive {
		return model.FocusSession{}, ErrSessionNotActive
	}
	if session.EndedEarly != nil && session.EndedEarly.FocusDuration >= sessionInfo.TimerDuration {
		return model.FocusSession{}, errors.New("invalid focus duration")

	}

	endedSession := model.FocusSession{}
	if session.EndedEarly != nil {
		endedSession.Status = string(StatusEndedEearly)
		endedSession.FocusDuration = &session.EndedEarly.FocusDuration
	} else {
		endedSession.Status = string(StatusCompleted)
		endedSession.FocusDuration = &sessionInfo.TimerDuration
	}

	result = database.Instance().
		Model(&model.FocusSession{}).
		Clauses(clause.Returning{}).
		Where("id = ?", session.SessionID).
		Updates(&endedSession)
	if result.Error != nil {
		return model.FocusSession{}, result.Error
	}
	if result.RowsAffected == 0 {
		return model.FocusSession{}, errors.New("should not be reachable")
	}

	return endedSession, nil
}

// type Analytics struct {
// 	TotalTimeSpent     int32
// 	TotalEstimatedTime int32
// 	DailyTimeSpent     map[string]int32
// 	TaskStatusCounts   map[string]int32
// 	AIFeedback         struct {
// 		Strengths        []string
// 		ImprovementAreas []string
// 		Motivation       string
// 	}
// }

// func GetAnalytics(userID int32, startDate, endDate *time.Time) (*Analytics, error) {
// 	// TODO: Implement
// 	// 1. Get all focus sessions in date range
// 	// 2. Calculate statistics
// 	// 3. Generate AI feedback
// 	return nil, errors.New("not implemented")
// }
