package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
)

type Request struct {
	Expression string `json:"expression"`
}

type Response struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func isValidExpression(expression string) bool {
	re := regexp.MustCompile(`^[\d\s\+\-\*/\(\)]+$`)
	return re.MatchString(expression)
}

func calculateExpression(expression string) (float64, error) {
	result, err := eval(expression)
	if err != nil {
		return 0, err
	}
	return result, nil
}

// eval - простая реализация вычисления выражения (можно заменить на более надежную)
func eval(expr string) (float64, error) {
	// Здесь можно использовать библиотеку для более безопасного вычисления выражений
	return strconv.ParseFloat(expr, 64) // Пример, замените на вашу логику
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Expression == "" {
		http.Error(w, `{"error": "Expression is not valid"}`, http.StatusUnprocessableEntity)
		return
	}

	if !isValidExpression(req.Expression) {
		http.Error(w, `{"error": "Expression is not valid"}`, http.StatusUnprocessableEntity)
		return
	}

	result, err := calculateExpression(req.Expression)
	if err != nil {
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	resp := Response{Result: fmt.Sprintf("%f", result)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/api/v1/calculate", calculateHandler)
	fmt.Println("Сервер запущен на порту :8080")
	http.ListenAndServe(":8080", nil)
}
