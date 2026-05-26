package httpx

import "net/http"

func RespondWithProblem(w http.ResponseWriter, statusCode int, problemType string, detail string) {
	problem := &Problem{
		Type:   problemType,
		Detail: detail,
	}

	RespondWithJSON(w, statusCode, problem)
}
