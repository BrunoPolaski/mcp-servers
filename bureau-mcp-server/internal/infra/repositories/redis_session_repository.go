package repositories

import (
	"context"
	"encoding/json"

	"github.com/BrunoPolaski/bureau-mcp-server/internal/core/entities"
	"github.com/BrunoPolaski/bureau-mcp-server/internal/helper/sessions"
	"github.com/BrunoPolaski/bureau-mcp-server/internal/infra/repositories/interfaces"
	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/redis/go-redis/v9"
)

type redisSessionRepository struct {
	rdb *redis.Client
}

func NewRedisSessionRepository(rdb *redis.Client) interfaces.SessionRepository {
	return &redisSessionRepository{
		rdb: rdb,
	}
}

func (r *redisSessionRepository) Create(ctx context.Context, session *entities.Session) *rest_err.RestErr {
	data, err := json.Marshal(session)
	if err != nil {
		return rest_err.NewInternalServerError("malformed session data")
	}

	err = r.rdb.Set(ctx, session.UUID, data, sessions.AbsoluteTimeout).Err()
	if err != nil {
		return rest_err.NewInternalServerError("error while creating session")
	}

	return nil
}

func (r *redisSessionRepository) GetById(ctx context.Context, id string) (*entities.Session, *rest_err.RestErr) {
	data, err := r.rdb.Get(ctx, id).Bytes()
	if err == redis.Nil {
		return nil, rest_err.NewNotFoundError("session not found")
	} else if err != nil {
		return nil, rest_err.NewInternalServerError("error while fetching session")
	}

	var sess entities.Session
	if err := json.Unmarshal(data, &sess); err != nil {
		return nil, rest_err.NewInternalServerError("error while unmarshalling session")
	}

	return &sess, nil
}

func (r *redisSessionRepository) Delete(ctx context.Context, id string) *rest_err.RestErr {
	err := r.rdb.Del(ctx, id).Err()
	if err != nil {
		return rest_err.NewInternalServerError("error while deleting session")
	}
	return nil
}
