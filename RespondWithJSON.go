package httpx

import (
	"encoding/json"
	"net/http"
)

func RespondWithJSON[T any](w http.ResponseWriter, statusCode int, responseModel T) {
	rawJSON, err := json.Marshal(responseModel)
	if err != nil {
		problem := config.ProblemForUnexpectedError(err)
		RespondWithJSON(w, http.StatusInternalServerError, problem)
		return
	}

	RespondWithBytes(w, statusCode, "application/json", rawJSON)
}
