package jsonhelpers

import (
	"encoding/json"
	"net/http"

	internalErr "github.com/SkorikovGeorge/jobWorker/internal/errors"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func SendError(w http.ResponseWriter, errCode int, errResp, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errCode)

	errResponse := ErrorResponse{
		Error: errResp,
	}
	if err := json.NewEncoder(w).Encode(errResponse); err != nil {
		log.Error().Err(errors.Wrap(err, internalErr.EncodingJSON)).Msg(errors.Wrap(err, internalErr.EncodingJSON).Error())
	}
}

func ReadJSON(r *http.Request, data interface{}) error {
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Error().Err(errors.Wrap(err, internalErr.ClosingBodyJSON)).Msg(errors.Wrap(err, internalErr.ClosingBodyJSON).Error())
		}
	}()
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return errors.Wrap(err, internalErr.ParsingJSON)
	}
	return nil
}

func SendJSON(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")

	code := http.StatusOK
	w.WriteHeader(code)

	if data == nil {
		return nil
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Error().Err(errors.Wrap(err, internalErr.EncodingJSON)).Msg(errors.Wrap(err, internalErr.EncodingJSON).Error())
		SendError(w, http.StatusInternalServerError, errors.Wrap(err, internalErr.EncodingJSON).Error(), errors.Wrap(err, internalErr.EncodingJSON).Error())
		return errors.Wrap(err, internalErr.ParsingJSON)
	}
	return nil
}
