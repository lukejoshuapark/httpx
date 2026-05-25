package httpx

import "net/http"

func UnexpectedError(w http.ResponseWriter, err error) {
	problem := config.ProblemForUnexpectedError(err)
	RespondWithJSON(w, http.StatusInternalServerError, problem)
}
