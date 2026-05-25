package httpx

import (
	"encoding/json"
	"mime"
	"net/http"
)

func ReceiveJSON[T Validate](w http.ResponseWriter, r *http.Request) *T {
	if r.ContentLength < 1 {
		problem := config.ProblemForLengthRequired()
		RespondWithJSON(w, http.StatusLengthRequired, problem)
		return nil
	}

	maximumRequestBodySize := config.MaximumRequestBodySize()

	if r.ContentLength > maximumRequestBodySize && maximumRequestBodySize > 0 {
		problem := config.ProblemForRequestEntityTooLarge()
		RespondWithJSON(w, http.StatusRequestEntityTooLarge, problem)
		return nil
	}

	contentType := r.Header.Get("Content-Type")
	mimeType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		problem := config.ProblemForUnsupportedMediaType("application/json", contentType)
		RespondWithJSON(w, http.StatusUnsupportedMediaType, problem)
		return nil
	}

	if mimeType != "application/json" {
		problem := config.ProblemForUnsupportedMediaType("application/json", contentType)
		RespondWithJSON(w, http.StatusUnsupportedMediaType, problem)
		return nil
	}

	var requestModel T
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestModel); err != nil {
		problem := config.ProblemForInvalidJSON(err)
		RespondWithJSON(w, http.StatusBadRequest, problem)
		return nil
	}

	if err := requestModel.Validate(); err != nil {
		problem := config.ProblemForUnprocessableEntity(err)
		RespondWithJSON(w, http.StatusUnprocessableEntity, problem)
		return nil
	}

	return &requestModel
}
