package handlers

import (
	"net/http"

	jsonhelpers "github.com/SkorikovGeorge/jobWorker/internal/json_helpers"
	"github.com/rs/zerolog/log"
)

var resp = struct {
	Msg string `json:"msg"`
}{
	Msg: "hello world",
}

// !!! delete this !!!

func GetJob(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("get job") // ! id

	body := resp

	if err := jsonhelpers.SendJSON(w, body); err != nil {
		log.Error().Err(err).Msg("Error sending JSON")
		return
	}
}

func CreateJob(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("create new job")

	body := resp

	if err := jsonhelpers.SendJSON(w, body); err != nil {
		log.Error().Err(err).Msg("Error sending JSON")
		return
	}
}
