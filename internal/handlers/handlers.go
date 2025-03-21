package handlers

import (
	"context"
	"net/http"

	"github.com/SkorikovGeorge/jobWorker/internal/consts"
	jsonhelpers "github.com/SkorikovGeorge/jobWorker/internal/json_helpers"
	"github.com/SkorikovGeorge/jobWorker/internal/redis"
	"github.com/SkorikovGeorge/jobWorker/internal/structs"
	"github.com/SkorikovGeorge/jobWorker/internal/workerpool"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func GetJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID := vars["job_id"]

	status, err := redis.RDB.GetJobStatus(context.Background(), jobID)
	if err != nil {
		if err.Error() == consts.ErrKeyNotFound {
			log.Error().Err(errors.Wrap(err, consts.ErrKeyNotFound))
			jsonhelpers.SendError(w, 404, consts.ErrKeyNotFound)
		} else {
			log.Error().Err(errors.Wrap(err, consts.ErrGetRedis)).Msg(errors.Wrap(err, consts.ErrGetRedis).Error())
			jsonhelpers.SendError(w, 500, consts.ErrGetRedis)
		}
		return
	}

	log.Info().Msgf("GetJob: job %s status %s", jobID, status["status"])

	if err := jsonhelpers.SendJSON(w, status); err != nil {
		log.Error().Err(err).Msg("Error sending JSON")
	}
}

func CreateJob(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Priority int `json:"priority"`
	}

	if err := jsonhelpers.ReadJSON(r, &req); err != nil {
		log.Error().Err(errors.Wrap(err, consts.ErrParsingJSON)).Msg(errors.Wrap(err, consts.ErrParsingJSON).Error())
		jsonhelpers.SendError(w, 500, consts.ErrSomething)
		return
	}

	if req.Priority == 0 {
		req.Priority = 5
	} else if 1 > req.Priority || req.Priority > 10 {
		jsonhelpers.SendError(w, 404, consts.ErrPriority)
		return
	}

	job := structs.NewJob(req.Priority)
	if err := redis.RDB.Enqueue(context.Background(), job); err != nil {
		log.Error().Err(errors.Wrap(err, consts.ErrEnqueueRedis)).Msg(errors.Wrap(err, consts.ErrEnqueueRedis).Error())
		return
	}

	if err := jsonhelpers.SendJSON(w, job); err != nil {
		log.Error().Err(err).Msg("Error sending JSON")
		return
	}
	log.Info().Msgf("CreateJob: created new job %s", job.ID)
}

func Pause(w http.ResponseWriter, r *http.Request) {
	log.Info().Msgf("Pause: workerpool is paused")
	workerpool.Workers.Paused.Store(true)
	res := struct {
		Msg string `json:"msg"`
	}{
		Msg: "Job processing paused",
	}
	if err := jsonhelpers.SendJSON(w, res); err != nil {
		log.Error().Err(err).Msg("Error sending JSON")
	}
}

func Resume(w http.ResponseWriter, r *http.Request) {
	log.Info().Msgf("Resume: workerpool is resumed")
	workerpool.Workers.Paused.Store(false)
	res := struct {
		Msg string `json:"msg"`
	}{
		Msg: "Job processing resumed",
	}
	if err := jsonhelpers.SendJSON(w, res); err != nil {
		log.Error().Err(err).Msg("Error sending JSON")
	}
}
