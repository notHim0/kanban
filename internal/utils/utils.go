package utils

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/notHim0/kanban/internal/types"
)

//loadSchema loads a JSON schema from a JSON file
func LoadSchema(filePath string) (string, error){
	data, err := os.ReadFile(filePath)

	if err != nil{
		return "", err
	}

	return string(data), nil
}

//general error response handler
func RespondWithError(w http.ResponseWriter, code int, message string){
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(types.ErrorResponse{Message: message})
}