package handler

import (
	"context"
	"net/http"

	"github.com/AddMile/backend/internal"
	gen "github.com/AddMile/backend/internal/api/codegen"
	"github.com/AddMile/backend/internal/app/user"
	"github.com/AddMile/backend/internal/shared"

	contextkit "github.com/AddMile/backend/internal/kit/context"
)

type UserHTTPHandler struct {
	service *user.Service
}

func NewUserHTTPHandler(service *user.Service) *UserHTTPHandler {
	return &UserHTTPHandler{service: service}
}

func toUserJSON(u user.User) gen.User {
	return gen.User{
		Id:        u.ID,
		Language:  gen.Language(u.Language),
		CreatedAt: u.CreatedAt,
	}
}

func (h *UserHTTPHandler) UpsertUser(
	ctx context.Context,
	request gen.UpsertUserRequestObject,
) (gen.UpsertUserResponseObject, error) {
	params := user.UpsertUserParams{
		Email:    request.Body.Email,
		Language: shared.Language(request.Body.Language),
	}

	userID, err := h.service.UpsertUser(ctx, params)
	if err != nil {
		if internal.ClientErr(err) {
			return gen.UpsertUser400JSONResponse{N400BadRequestJSONResponse: gen.N400BadRequestJSONResponse{
				Code: http.StatusBadRequest, Error: err.Error(),
			}}, nil
		}

		return gen.UpsertUser500JSONResponse{N500InternalServerErrorJSONResponse: gen.N500InternalServerErrorJSONResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		}}, nil
	}

	return gen.UpsertUser201JSONResponse{Id: userID}, nil
}

func (h *UserHTTPHandler) GetUser(
	ctx context.Context,
	request gen.GetUserRequestObject,
) (gen.GetUserResponseObject, error) {
	userID := contextkit.MustUserID(ctx)

	u, err := h.service.User(ctx, user.UserParams{
		UserID: userID,
	})
	if err != nil {
		if internal.ClientErr(err) {
			return gen.GetUser400JSONResponse{N400BadRequestJSONResponse: gen.N400BadRequestJSONResponse{
				Code: http.StatusBadRequest, Error: err.Error(),
			}}, nil
		}

		return gen.GetUser500JSONResponse{N500InternalServerErrorJSONResponse: gen.N500InternalServerErrorJSONResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		}}, nil
	}

	return gen.GetUser200JSONResponse(toUserJSON(u)), nil
}
