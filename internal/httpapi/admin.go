package httpapi

import (
	"encoding/json"
	"net/http"
	"warrantyservice/internal/warranty"
)

func DecodeUpsert(r *http.Request) (warranty.UpsertInput, error) {
	var input warranty.UpsertInput
	err := json.NewDecoder(r.Body).Decode(&input)
	return input, err
}

func AdminTokenValid(r *http.Request, expected string) bool {
	return r.Header.Get("Authorization") == "Bearer "+expected
}
