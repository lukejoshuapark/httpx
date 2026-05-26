package httpx

import (
	"encoding/json"
	"mime"
	"net/http"
)

type Ptr[T any] interface {
	*T
	Validate
}

func ReceiveJSON[T any, P Ptr[T]](w http.ResponseWriter, r *http.Request) *T {
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
	var requestModelPointer P = &requestModel

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(requestModelPointer); err != nil {
		problem := config.ProblemForInvalidJSON(err)
		RespondWithJSON(w, http.StatusBadRequest, problem)
		return nil
	}

	if err := requestModelPointer.Validate(); err != nil {
		problem := config.ProblemForUnprocessableEntity(err)
		RespondWithJSON(w, http.StatusUnprocessableEntity, problem)
		return nil
	}

	return requestModelPointer
}
