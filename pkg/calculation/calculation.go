package calculation

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type CalculationRequest struct {
	Expression string `json:"expression"`
}

type CalculationResponse struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func main() {
	http.HandleFunc("/api/v1/calculate", HandleCalculate)
	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", nil)
}

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
	result, err := Calc(req.Expression)
	if err != nil {
		if errors.Is(err, ErrInvalidExpression) {
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

func Calc(expression string) (float64, error) {
	expression = strings.ReplaceAll(expression, " ", "")
	if expression == "" {
		return 0, ErrInvalidExpression
	}

	postfix, err := infixToPostfix(expression)
	if err != nil {
		return 0, ErrInvalidExpression
	}

	return evaluatePostfix(postfix)
}

func infixToPostfix(expression string) ([]string, error) {
	precedence := map[byte]int{'(': 1, '+': 2, '-': 2, '*': 3, '/': 3}
	var output []string
	var stack []byte

	for i := 0; i < len(expression); i++ {
		ch := expression[i]

		if isDigit(ch) {
			num := string(ch)
			for i+1 < len(expression) && isDigit(expression[i+1]) {
				i++
				num += string(expression[i])
			}
			output = append(output, num)
		} else if ch == '(' {
			stack = append(stack, ch)
		} else if ch == ')' {
			for len(stack) > 0 && stack[len(stack)-1] != '(' {
				output = append(output, string(stack[len(stack)-1]))
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				return nil, ErrInvalidExpression
			}
			stack = stack[:len(stack)-1]
		} else {
			for len(stack) > 0 && precedence[stack[len(stack)-1]] >= precedence[ch] {
				output = append(output, string(stack[len(stack)-1]))
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, ch)
		}
	}

	for len(stack) > 0 {
		output = append(output, string(stack[len(stack)-1]))
		stack = stack[:len(stack)-1]
	}

	return output, nil
}

func evaluatePostfix(postfix []string) (float64, error) {
	var stack []float64

	for _, token := range postfix {
		if isDigitString(token) {
			num, _ := strconv.ParseFloat(token, 64)
			stack = append(stack, num)
		} else {
			if len(stack) < 2 {
				return 0, ErrInvalidExpression
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			switch token {
			case "+":
				stack = append(stack, a+b)
			case "-":
				stack = append(stack, a-b)
			case "*":
				stack = append(stack, a*b)
			case "/":
				if b == 0 {
					return 0, ErrDivideByZero
				}
				stack = append(stack, a/b)
			default:
				return 0, ErrInvalidExpression
			}
		}
	}

	if len(stack) != 1 {
		return 0, ErrInvalidExpression
	}
	return stack[0], nil
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isDigitString(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
