package calc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"testing"
)

type Test struct {
	name               string
	expression         string
	expected           float64
	wantError          bool
	expectedStatusCode int
}

type Response struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func TestCalc(t *testing.T) {
	tests := []Test{
		{"Simple One Number 1",
			"5",
			5,
			false,
			200},
		{"Simple One Number 2",
			"97525673.85739572",
			97525673.85739572,
			false,
			200},
		{"Simple Addition",
			"5+3",
			8,
			false,
			200},

		{"Division By Zero",
			"8/0",
			0,
			true,
			422},

		{"Nested Calculations",
			"35 + (10 - 2 * 5) + (6 / 2 * 5 - 10 + 2) * (2 * 3)",
			77,
			false,
			200},

		{"Order of Operations",
			"10 + 15 - (2 + 3) * 2",
			15,
			false,
			200},

		{"Multiplication and Addition",
			"5 * 8 + 4 * 6 + 15 - 14",
			65,
			false,
			200},

		{"Mixed Operations",
			"5-2+62-4+8/2-1",
			5 - 2 + 62 - 4 + 8/2 - 1,
			false,
			200},

		{"Division in Nested Expressions",
			"35 + (10 - 2 5) + (6 / 0 * 5 - 10 + 2) * (2 * 3)",
			0,
			true,
			422},

		{"Complex Expression",
			"(11437 + 128 * 31) / 237 - 37",
			28,
			false,
			200},

		{"Expression with Zero Division",
			"93478+23657-(52253/0)",
			0,
			true,
			422},

		{"Comparative Division Subtraction",
			"(37296 / 37 - 17780 / 35) / 250",
			(37296/37 - 17780/35) / 250,
			false,
			200},

		{"Unbalanced Parentheses",
			"5+(5-4",
			0,
			true,
			422},
		{"Nested Unbalanced Parentheses",
			"35 + (10 - 2 * 5) + (6 / 3 * 5 - 10 + 2 * (2 * 3)",
			0,
			true,
			422},

		{"Easy",
			"5+7",
			5 + 7,
			false,
			200},

		{"Operator at End",
			"5+7/",
			0,
			true,
			422},

		{"Invalid Start Operator",
			"*5+7",
			0,
			true,
			422},

		{"Simply Invalid Input",
			"valid input",
			0,
			true,
			422},

		{"Empty Input",
			"",
			0,
			true,
			422},

		{"Floating Point Addition",
			"5.1 + 5.2",
			10.3,
			false,
			200},

		{"Exponentiation",
			"5^2",
			25,
			false,
			200},

		{"Simple Cube",
			"2^3",
			8,
			false,
			200},

		{"Zero Exponentiation",
			"0^2",
			0,
			false,
			200},

		{"Negative Base Exponent",
			"-2^3",
			0,
			true,
			422},

		{"Exponent with Empty Base",
			"^2",
			0,
			true,
			422},

		{"Exponentiation With Parentheses",
			"5^(2+1)",
			125,
			false,
			200},

		{"Compound Exponentiation",
			"(2^3)^2",
			64,
			false,
			200},

		{"Identity Exponentiation",
			"(2^3)^1",
			8,
			false,
			200},

		{"Fractional Exponent",
			"(2^3)^(1/3)",
			math.Pow(math.Pow(2, 3), 1.0/3.0),
			false,
			200},

		{"Complex Multiplication",
			"5*(22+3)-2",
			5*(22+3) - 2,
			false,
			200},

		{"Big Number",
			"1+999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999",
			0,
			true,
			500},
	}

	url := "http://localhost:8080/api/v1/calculate"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err1 := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(fmt.Sprintf(`{"expression":"%v"}`, tt.expression))))
			if err1 != nil {
				panic(err1)
			}

			request.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			response, err2 := client.Do(request)
			if err2 != nil {
				panic(err2)
			}
			defer response.Body.Close()

			responseBytes, err3 := io.ReadAll(response.Body)
			if err3 != nil {
				panic(err3)
			}

			var result Response
			err4 := json.Unmarshal(responseBytes, &result)
			if err4 != nil {
				panic(err4)
			}

			if response.StatusCode != tt.expectedStatusCode {
				t.Errorf("Handler returned wrong status code: expected %v, got %v", tt.expectedStatusCode, response.StatusCode)
			}
			if tt.wantError {
				if result.Result != "" {
					t.Errorf("Handler returned wrong result: got %v want empty string", result.Result)
				}
			} else {
				result, err := strconv.ParseFloat(result.Result, 64)
				if err != nil {
					t.Fatalf("Error parsing result: %v", err)
				}
				if math.Abs(result-tt.expected) > 1e-6 {
					t.Errorf("Handler returned wrong result: got %v want %v", result, tt.expected)
				}
			}
		})
	}
}
