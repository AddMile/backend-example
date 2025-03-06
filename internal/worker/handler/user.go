package handlers

import (
	"context"
	"net/http"

	"github.com/AddMile/backend/internal"
	userJob "github.com/AddMile/backend/internal/job/user"
	"github.com/AddMile/backend/internal/shared/event"
	gen "github.com/AddMile/backend/internal/worker/codegen"

	pubsubkit "github.com/AddMile/backend/internal/kit/pubsub"
)

type UserHTTPHandler struct {
	processor *userJob.Processor
}

func NewUserHTTPHandler(processor *userJob.Processor) *UserHTTPHandler {
	return &UserHTTPHandler{processor: processor}
}

func (h *UserHTTPHandler) EmailUserCreated(
	ctx context.Context,
	request gen.EmailUserCreatedRequestObject,
) (gen.EmailUserCreatedResponseObject, error) {
	var event event.UserCreatedEvent
	if err := pubsubkit.DecodeMessageData(request.Body.Message.Data, &event); err != nil {
		return nil, err
	}
	err := h.processor.EmailUserCreatedJob(ctx, event)
	if err != nil {
		if internal.ClientErr(err) {
			return gen.EmailUserCreated400JSONResponse{N400BadRequestJSONResponse: gen.N400BadRequestJSONResponse{
				Code: http.StatusBadRequest, Error: err.Error(),
			}}, nil
		}

		return gen.EmailUserCreated500JSONResponse{N500InternalServerErrorJSONResponse: gen.N500InternalServerErrorJSONResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		}}, nil
	}

	return gen.EmailUserCreated204Response{}, nil
}
