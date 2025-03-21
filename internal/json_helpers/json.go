package jsonhelpers

import (
	"encoding/json"
	"net/http"
	"strings"

	internalErr "github.com/SkorikovGeorge/jobWorker/internal/consts"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func SendError(w http.ResponseWriter, errCode int, errResp string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errCode)

	errResponse := ErrorResponse{
		Error: errResp,
	}
	if err := json.NewEncoder(w).Encode(errResponse); err != nil {
		log.Error().Err(errors.Wrap(err, internalErr.ErrEncodingJSON)).Msg(errors.Wrap(err, internalErr.ErrEncodingJSON).Error())
	}
}

func ReadJSON(r *http.Request, data interface{}) error {
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Error().Err(errors.Wrap(err, internalErr.ErrClosingBodyJSON)).Msg(errors.Wrap(err, internalErr.ErrClosingBodyJSON).Error())
		}
	}()
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		if strings.Contains(err.Error(), "EOF") {
			return nil
		}
		return errors.Wrap(err, internalErr.ErrParsingJSON)
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
		log.Error().Err(errors.Wrap(err, internalErr.ErrEncodingJSON)).Msg(errors.Wrap(err, internalErr.ErrEncodingJSON).Error())
		SendError(w, http.StatusInternalServerError, errors.Wrap(err, internalErr.ErrEncodingJSON).Error())
	}
	return nil
}

func ToString(data interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Error().Err(errors.Wrap(err, internalErr.ErrEncodingJSON)).Msg(errors.Wrap(err, internalErr.ErrEncodingJSON).Error())
		return "", errors.Wrap(err, internalErr.ErrEncodingJSON)
	}
	return string(jsonData), nil
}

func FromString(jsonString string, data interface{}) error {
	if err := json.Unmarshal([]byte(jsonString), data); err != nil {
		log.Error().Err(errors.Wrap(err, internalErr.ErrParsingJSON)).Msg(errors.Wrap(err, internalErr.ErrParsingJSON).Error())
		return errors.Wrap(err, internalErr.ErrParsingJSON)
	}
	return nil
}
