package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type Handler struct {
	repo        *Repository
	eventStream *EventStream
}

func NewHandler(repo *Repository, stream *EventStream) *Handler {
	return &Handler{
		repo:        repo,
		eventStream: stream,
	}
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		userId := strings.Split(request.URL.Path, "/")[1]

		response, err := h.handleGetUser(userId)
		if err != nil {
			writer.WriteHeader(404)
			return
		}

		marshal, err := json.Marshal(response)
		if err != nil {
			return
		}

		writer.Write(marshal)
		return

	case http.MethodPost:
		var body CreateUserRequest
		err := json.NewDecoder(request.Body).Decode(&body)
		if err != nil {
			writer.WriteHeader(400)
			return
		}

		response, err := h.handleCreateUser(body)
		if err != nil {
			writer.WriteHeader(400)
			return
		}

		marshal, err := json.Marshal(response)
		if err != nil {
			writer.WriteHeader(500)
			return
		}

		writer.Write(marshal)
		return
	}
}

type GetUserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *Handler) handleGetUser(userID string) (*GetUserResponse, error) {
	user, err := h.repo.getById(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	return &GetUserResponse{
		ID:    user.id,
		Name:  user.name,
		Email: user.email,
	}, nil
}

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateUserResponse struct {
	UserID string `json:"user_id"`
}

func (h *Handler) handleCreateUser(request CreateUserRequest) (*CreateUserResponse, error) {
	ctx := context.Background()
	if res, err := h.repo.getByEmail(ctx, request.Email); res != nil {
		return nil, err
	}

	userID := uuid.NewString()
	err := h.repo.create(ctx, userID, request.Name, request.Email)
	if err != nil {
		return nil, err
	}

	_ = h.eventStream.publish(ctx, []byte(fmt.Sprintf("user_created:%s", userID)))

	return &CreateUserResponse{UserID: userID}, nil
}
