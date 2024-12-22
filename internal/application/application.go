package application

import (
	"LMC/pkg/calculation"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type CalculationRequest struct {
	Expression string `json:"expression"`
}

type CalculationResponse struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

// Обрабатывает вычисления и генерирует ответы
func HandleCalculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not accepted", http.StatusMethodNotAllowed)
		return
	}

	var req CalculationRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		respondWithError(w, "Invalid request format", http.StatusUnprocessableEntity)
		return
	}

	if strings.TrimSpace(req.Expression) == "" {
		respondWithError(w, "Expression is not valid", http.StatusUnprocessableEntity)
		return
	}

	// Вычисление выражения
	result, err := calculation.Calc(req.Expression)
	if err != nil {
		if errors.Is(err, calculation.ErrInvalidExpression) {
			respondWithError(w, "Expression is not valid", http.StatusUnprocessableEntity)
		} else {
			respondWithError(w, "Server Error", http.StatusInternalServerError)
		}
		return
	}

	response := CalculationResponse{Result: fmt.Sprintf("%f", result)}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func respondWithError(w http.ResponseWriter, message string, statusCode int) {
	response := CalculationResponse{Error: message}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
