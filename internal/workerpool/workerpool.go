package workerpool

import (
	"context"
	"math/rand/v2"
	"time"

	"github.com/SkorikovGeorge/jobWorker/internal/consts"
	"github.com/SkorikovGeorge/jobWorker/internal/redis"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func (wp *WorkerPool) SetupWorkers() {
	for i := 1; i <= wp.Config.WorkersCount; i++ {
		wp.wg.Add(1)
		go wp.startWorker(i)
	}
}

func (wp *WorkerPool) startWorker(workerNum int) {
	log.Info().Msgf("WorkerPool: worker #%d started", workerNum)
	defer wp.wg.Done()
	for {
		if wp.Paused.Load() {
			time.Sleep(wp.Config.PauseDelay)
			continue
		}
		job, err := redis.RDB.Dequeue(wp.ShutdownCtx)
		if err != nil {
			if err.Error() == consts.ErrEmptyQueue {
				time.Sleep(wp.Config.EmptyDelay)
				continue
			} else if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				log.Info().Msgf("WorkerPool: worker #%d: shutdown signal received, stopping...", workerNum)
				return
			} else {
				log.Error().Err(errors.Wrap(err, consts.ErrRunningJob))
				return
			}
		}

		jobTimeoutCtx, cancel := context.WithTimeout(wp.ShutdownCtx, wp.Config.JobTimeout)
		log.Info().Msgf("WorkerPool: worker #%d recieved job %s", workerNum, job.ID)
		if err := PerformJob(jobTimeoutCtx, job.ID, nil); err != nil {
			log.Warn().Msgf("WorkerPool: unable to perform job %s", job.ID)
			retriesCount, err := redis.RDB.IncrementRetries(context.Background(), job.ID)
			if err != nil {
				log.Error().Err(errors.Wrap(err, consts.ErrRunningJob))
			}
			if retriesCount <= wp.Config.MaxRetries {
				log.Warn().Msgf("WorkerPool: job %s will be retried, attempt %d/%d", job.ID, retriesCount, wp.Config.MaxRetries)
				time.Sleep(wp.Config.RetryDelay)
				if err := redis.RDB.Enqueue(context.Background(), job); err != nil {
					log.Error().Err(errors.Wrap(err, consts.ErrRunningJob))
				} else {
					log.Info().Msgf("WorkerPool: job %s re-enqueued for retry", job.ID)
				}
			} else {
				log.Error().Msgf("Max retries reached for job %s. Job failed", job.ID)
				if err := redis.RDB.SetJobStatus(context.Background(), job.ID, consts.StatusFailed); err != nil {
					log.Error().Err(errors.Wrap(err, consts.ErrRunningJob)).Msg(errors.Wrap(err, consts.ErrRunningJob).Error())
				}
			}
		} else {
			if err := redis.RDB.SetJobStatus(context.Background(), job.ID, consts.StatusDone); err != nil {
				log.Error().Err(errors.Wrap(err, consts.ErrRunningJob)).Msg(errors.Wrap(err, consts.ErrRunningJob).Error())
			}
			log.Info().Msgf("WorkerPool: worker #%d completed job id %s", workerNum, job.ID)
		}
		cancel()
	}
}

func (wp *WorkerPool) Shutdown(ctx context.Context) error {
	wp.CancelCtx()

	done := make(chan int)
	go func() {
		wp.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Info().Msg("WorkerPool: worker pool shutdown gracefully")
		return nil
	case <-ctx.Done():
		log.Info().Msg("WorkerPool: worker pool shutdown with timeout")
		return ctx.Err()
	}
}

func PerformJob(ctx context.Context, name string, jobData []byte) error {
	log.Info().Msgf("Started job %s at %d", name, time.Now().UnixMilli())

	// для тестированя retry с некоторой вероятностью появляется ошибка
	if rand.IntN(100) >= 90 {
		return errors.New(consts.ErrSomething)
	}

	// do some work. Random sleep time; max = 3s
	select {
	case <-time.After(time.Duration(rand.Int64N(int64(3 * time.Second)))):
		log.Info().Msgf("Job: finished job %s at %d", name, time.Now().UnixMilli())
		return nil
	case <-ctx.Done():
		log.Warn().Msgf("Job: job %s cancelled due to timeout", name)
		return context.DeadlineExceeded
	}
}
