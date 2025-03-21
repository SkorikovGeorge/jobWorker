package redis

import (
	"context"

	"github.com/SkorikovGeorge/jobWorker/internal/consts"
	jsonhelpers "github.com/SkorikovGeorge/jobWorker/internal/json_helpers"
	"github.com/SkorikovGeorge/jobWorker/internal/structs"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

var RDB = Setup()

type Redis struct {
	Config *RedisConfig
	db     *redis.Client
}

func Setup() *Redis {
	db := redis.NewClient(&redis.Options{
		Addr: cfg.Address,
		DB:   cfg.DB,
	})

	if _, err := db.Ping(context.Background()).Result(); err != nil {
		log.Fatal().Msg(errors.Wrap(err, consts.ErrConnectingRedis).Error())
	}

	rdb := Redis{
		Config: &cfg,
		db:     db,
	}

	return &rdb
}

func (rdb *Redis) Enqueue(ctx context.Context, job *structs.Job) error {
	jobString, err := jsonhelpers.ToString(job)
	if err != nil {
		return errors.Wrap(err, consts.ErrEnqueueRedis)
	}
	if err = rdb.SetJobStatus(context.Background(), job.ID, job.Status); err != nil {
		return errors.Wrap(err, consts.ErrEnqueueRedis)
	}

	if err := rdb.db.ZAdd(ctx, rdb.Config.QueueName, &redis.Z{
		Score:  float64(job.Priority),
		Member: jobString,
	}).Err(); err != nil {
		log.Error().Err(errors.Wrap(err, consts.ErrEnqueueRedis)).Msg(errors.Wrap(err, consts.ErrEnqueueRedis).Error())
		return errors.Wrap(err, consts.ErrEnqueueRedis)
	}
	log.Info().Msgf("Redis: new job %s added to queue", job.ID)
	return nil
}

func (rdb *Redis) Dequeue(ctx context.Context) (*structs.Job, error) {
	result, err := rdb.db.ZPopMin(ctx, rdb.Config.QueueName, 1).Result()
	if err == redis.Nil {
		return nil, errors.New(consts.ErrEmptyQueue)
	} else if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return nil, err
	} else if err != nil {
		log.Error().Err(errors.Wrap(err, consts.ErrDequeueRedis)).Msg(errors.Wrap(err, consts.ErrDequeueRedis).Error())
		return nil, errors.Wrap(err, consts.ErrDequeueRedis)
	}

	if len(result) == 0 {
		return nil, errors.New(consts.ErrEmptyQueue)
	}

	jobString := result[0].Member.(string)

	var job structs.Job
	if err = jsonhelpers.FromString(jobString, &job); err != nil {
		return nil, errors.Wrap(err, consts.ErrDequeueRedis)
	}

	job.Status = consts.StatusRunning
	if err = rdb.SetJobStatus(context.Background(), job.ID, job.Status); err != nil {
		return nil, errors.Wrap(err, consts.ErrDequeueRedis)
	}

	log.Info().Msgf("Redis: job %s removed from queue", job.ID)
	return &job, nil
}

func (rdb *Redis) SetJobStatus(ctx context.Context, jobID string, status string) error {
	if err := rdb.db.HSet(ctx, jobID, "status", status).Err(); err != nil {
		log.Error().Err(errors.Wrap(err, consts.ErrSetRedis)).Msg(errors.Wrap(err, consts.ErrSetRedis).Error())
		return errors.Wrap(err, consts.ErrSetRedis)
	}

	if err := rdb.db.Expire(ctx, jobID, rdb.Config.JobStatusTTL).Err(); err != nil {
		log.Error().Err(errors.Wrap(err, consts.ErrTTLRedis))
	}

	return nil
}

func (rdb *Redis) IncrementRetries(ctx context.Context, jobID string) (int, error) {
	res := rdb.db.HIncrBy(ctx, jobID, "retries", 1)
	if err := res.Err(); err != nil {
		log.Error().Err(errors.Wrap(err, consts.ErrIncrementRedis))
		return 0, errors.Wrap(err, consts.ErrIncrementRedis)
	}

	retriesCount, err := res.Result()
	if err != nil {
		log.Error().Err(errors.Wrap(err, consts.ErrIncrementRedis))
		return 0, errors.Wrap(err, consts.ErrIncrementRedis)
	}

	return int(retriesCount), nil
}

func (rdb *Redis) GetJobStatus(ctx context.Context, jobID string) (map[string]string, error) {
	statusMap := make(map[string]string)
	res := rdb.db.HGet(ctx, jobID, "status")
	status, err := res.Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, errors.New(consts.ErrKeyNotFound)
		}
		log.Error().Err(errors.Wrap(err, consts.ErrGetRedis))
		return nil, errors.Wrap(err, consts.ErrGetRedis)
	}
	statusMap["status"] = status

	return statusMap, nil
}
