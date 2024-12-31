package handler

import (
	"context"
	"errors"
	"study-planner-api/internal/api"
	"study-planner-api/internal/focussession"
)

// PostFocusSessions implements api.StrictServerInterface.
func (s *Handler) PostFocusSessions(ctx context.Context, request api.PostFocusSessionsRequestObject) (api.PostFocusSessionsResponseObject, error) {
	authInfo := api.AuthInfoOfRequest(ctx)

	session, err := focussession.CreateSession(focussession.NewSession{
		UserID:        authInfo.ID,
		TaskID:        request.Body.TaskId,
		TimerDuration: request.Body.TimerDuration,
		BreakDuration: request.Body.BreakDuration,
	})
	if err != nil {
		if errors.Is(err, focussession.ErrTaskNotBelongToUser) ||
			errors.Is(err, focussession.ErrTaskNotFound) {
			return api.PostFocusSessions404Response{}, nil
		}
		if errors.Is(err, focussession.ErrTaskNotInProgress) {
			return api.PostFocusSessions400Response{}, nil
		}

		return nil, err
	}

	return api.PostFocusSessions201JSONResponse{
		Id:            &session.ID,
		UserId:        &authInfo.ID,
		TaskId:        &session.TaskID,
		TimerDuration: &session.TimerDuration,
		BreakDuration: &session.BreakDuration,
		Status:        &session.Status,
		FocusDuration: &session.FocusDuration,
		CreatedAt:     &session.CreatedAt,
		UpdatedAt:     &session.UpdatedAt,
	}, nil
}

// PostFocusSessionsIdEnd implements api.StrictServerInterface.
func (s *Handler) PostFocusSessionsIdEnd(ctx context.Context, request api.PostFocusSessionsIdEndRequestObject) (api.PostFocusSessionsIdEndResponseObject, error) {
	authInfo := api.AuthInfoOfRequest(ctx)

	session := focussession.SessionToEnd{
		UserID:    authInfo.ID,
		SessionID: request.Id,
	}
	if request.Body != nil {
		session.EndedEarly = &focussession.EndEarly{
			FocusDuration: request.Body.FocusDuration,
		}
	}

	endedSession, err := focussession.EndSession(session)
	if err != nil {
		if errors.Is(err, focussession.ErrSessionNotFound) {
			return api.PostFocusSessionsIdEnd404Response{}, nil
		}
		if errors.Is(err, focussession.ErrSessionNotActive) {
			return api.PostFocusSessionsIdEnd400Response{}, nil
		}

		return nil, err
	}

	return api.PostFocusSessionsIdEnd200JSONResponse{
		Id:            &endedSession.ID,
		UserId:        &authInfo.ID,
		TaskId:        &endedSession.TaskID,
		TimerDuration: &endedSession.TimerDuration,
		BreakDuration: &endedSession.BreakDuration,
		Status:        &endedSession.Status,
		FocusDuration: &endedSession.FocusDuration,
		CreatedAt:     &endedSession.CreatedAt,
		UpdatedAt:     &endedSession.UpdatedAt,
	}, nil
}
