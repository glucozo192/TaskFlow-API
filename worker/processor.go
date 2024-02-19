package worker

import (
	"context"

	"github.com/glu-project/internal/user/golibs/database"
	"github.com/hibiken/asynq"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	DB     database.Ext
}
