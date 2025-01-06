package task

import (
	"errors"
	"fmt"
	db "study-planner-api/internal/database"
	"study-planner-api/internal/model"
	"study-planner-api/internal/utils"
	"sync"
	"time"

	"gorm.io/gorm"
)

type Status string

const (
	StatusTodo       Status = "Todo"
	StatusInProgress Status = "In Progress"
	StatusCompleted  Status = "Completed"
	StatusExpired    Status = "Expired"
)

func (s Status) String() string {
	return string(s)
}

func StatusFromString(str string) (Status, error) {
	switch str {
	case string(StatusTodo):
		return StatusTodo, nil
	case string(StatusInProgress):
		return StatusInProgress, nil
	case string(StatusCompleted):
		return StatusCompleted, nil
	case string(StatusExpired):
		return StatusExpired, nil
	default:
		return *new(Status), fmt.Errorf("invalid status: %s", str)
	}
}

type Priority string

const (
	PriorityLow    Priority = "Low"
	PriorityMedium Priority = "Medium"
	PriorityHigh   Priority = "High"
)

func PriorityFromString(str string) (Priority, error) {
	switch str {
	case string(PriorityLow):
		return PriorityLow, nil
	case string(PriorityMedium):
		return PriorityMedium, nil
	case string(PriorityHigh):
		return PriorityHigh, nil
	default:
		return *new(Priority), fmt.Errorf("invalid priority: %s", str)
	}
}

type SortType struct {
	Field SortField
	Order SortOrder
}

type SortField string

const (
	SortFieldCreatedAt SortField = "created_at"
	SortFieldStartTime SortField = "end_time"
	SortFieldEndTime   SortField = "start_time"
	SortFieldPriority  SortField = "priority"
)

func SortFieldFromString(str string) (SortField, error) {
	switch str {
	case string(SortFieldCreatedAt):
		return SortFieldCreatedAt, nil
	case string(SortFieldStartTime):
		return SortFieldStartTime, nil
	case string(SortFieldEndTime):
		return SortFieldEndTime, nil
	case string(SortFieldPriority):
		return SortFieldPriority, nil
	default:
		return *new(SortField), fmt.Errorf("invalid sort field: %s", str)
	}
}

type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

func SortOrderFromString(str string) (SortOrder, error) {
	switch str {
	case string(SortOrderAsc):
		return SortOrderAsc, nil
	case string(SortOrderDesc):
		return SortOrderDesc, nil
	default:
		return *new(SortOrder), fmt.Errorf("invalid sort order: %s", str)
	}
}

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

type GetCriteria struct {
	UserID     int32
	Status     *Status
	Search     *string
	Priority   *Priority
	StartTime  *time.Time
	EndTime    *time.Time
	SortType   SortType
	Pagination utils.Pagination
}

func GetTasks(criteria *GetCriteria) ([]model.Task, error) {
	var tasks []model.Task

	constructQuery := func(db *gorm.DB) *gorm.DB {
		query := db.
			Model(&model.Task{}).
			Order(fmt.Sprintf("%s %s", criteria.SortType.Field, criteria.SortType.Order)).
			Where("user_id = ?", criteria.UserID)

		if criteria.Status != nil {
			query = query.Where("status = ?", criteria.Status)
		}
		if criteria.Search != nil {
			query = query.Where("name LIKE ?", fmt.Sprintf("%%%s%%", *criteria.Search))
		}
		if criteria.Priority != nil {
			query = query.Where("priority = ?", criteria.Priority)
		}
		if criteria.StartTime != nil {
			query = query.Where("start_time >= ?", criteria.StartTime)
		}
		if criteria.EndTime != nil {
			query = query.Where("end_time <= ?", criteria.EndTime)
		}
		return query
	}

	var paginationQueryErr, listQueryError error
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		result := db.Instance().
			Scopes(constructQuery).
			Scopes(utils.Paginate(criteria.Pagination)).
			Find(&tasks)

		listQueryError = result.Error
	}()

	go func() {
		defer wg.Done()
		paginationQueryErr = utils.GetPaginationInfo(
			&criteria.Pagination,
			db.Instance().Scopes(constructQuery),
		)
	}()

	wg.Wait()

	if paginationQueryErr != nil {
		return nil, paginationQueryErr
	}
	if listQueryError != nil {
		return nil, listQueryError
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
