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
	"time"

	"github.com/Barsenick/calculator/pkg/calc"
	"github.com/google/uuid"
)

type Test struct {
	name               string
	expression         string
	expected           float64
	wantError          bool
	expectedStatusCode int
}

type CalcResponse struct {
	ID    int    `json:"id"`
	Error string `json:"error,omitempty"`
}

type ExpressionResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Result string `json:"result"`
}

type RegistrationResponse struct {
	Status string `json:"status"`
	Token  string `json:"token,omitempty"`
}

var tests = []Test{
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

var authToken = "Bearer "

func TestCalcFunc(t *testing.T) {
	for i, test := range tests {
		result, err := calc.NormalCalc(test.expression)
		if test.wantError && err == nil {
			t.Fatalf("Expected error; got %v in test #%v", result, i+1)
		}
		if !test.wantError && err != nil {
			t.Fatalf("Unexpected error in test #%v: %v", i+1, err)
		}
		if result != test.expected {
			t.Fatalf("Expected %v; got %v in test #%v", test.expected, result, i+1)
		}
	}
}

func TestWebCalc(t *testing.T) {
	name := uuid.NewString()
	request, err := http.NewRequest(http.MethodPost, "http://localhost:8070/api/v1/register", bytes.NewReader([]byte(`{"login":"`+name+`", "password": "test"}`)))
	if err != nil {
		panic(err)
	}

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var reg_resp RegistrationResponse

	err = json.Unmarshal(responseBytes, &reg_resp)
	if err != nil {
		panic(err)
	}

	if reg_resp.Status != "200 OK" {
		t.Fatalf("registration status is not ok: %v", reg_resp.Status)
	}

	request, err = http.NewRequest(http.MethodPost, "http://localhost:8070/api/v1/login", bytes.NewReader([]byte(`{"login":"`+name+`", "password": "test"}`)))
	if err != nil {
		panic(err)
	}

	client = &http.Client{}

	response, err = client.Do(request)
	if err != nil {
		panic(err)
	}

	responseBytes, err = io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var login_resp RegistrationResponse

	err = json.Unmarshal(responseBytes, &login_resp)
	if err != nil {
		panic(err)
	}

	if login_resp.Status != "200 OK" {
		t.Fatalf("login status is not ok: %v", login_resp.Status)
	}

	token := login_resp.Token
	authToken += token

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err1 := http.NewRequest(http.MethodPost, "http://localhost:8070/api/v1/calculate", bytes.NewReader([]byte(fmt.Sprintf(`{"expression":"%v"}`, tt.expression))))
			if err1 != nil {
				panic(err1)
			}

			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Authorization", authToken)

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

			var cr CalcResponse
			err4 := json.Unmarshal(responseBytes, &cr)
			if err4 != nil {
				panic(err4)
			}

			id := strconv.Itoa(cr.ID)

			request2, err5 := http.NewRequest(http.MethodGet, "http://localhost:8070/api/v1/expressions?id="+id, nil)
			if err5 != nil {
				panic(err5)
			}

			haventGotResult := true
			attempts := 0

			for haventGotResult {
				if attempts >= 10 {
					t.Fatal("still pending after 10 attempts")
				} else {
					request2.Header.Set("Authorization", authToken)
					time.Sleep(50 * time.Millisecond)
					response2, err6 := client.Do(request2)
					if err6 != nil {
						panic(err6)
					}
					defer response2.Body.Close()
					request2.Header.Set("Authorization", authToken)

					response2Bytes, err7 := io.ReadAll(response2.Body)
					if err7 != nil {
						panic(err7)
					}

					var er ExpressionResponse
					err8 := json.Unmarshal(response2Bytes, &er)
					t.Log(string(response2Bytes))
					t.Log(string(authToken))
					if err8 != nil {
						panic(err8)
					}

					if er.Result == "pending" {
						attempts++
						continue
					} else if er.Result == calc.Err422.Error() || er.Result == calc.Err500.Error() {
						if !tt.wantError {
							t.Fatalf("unexpected error")
						}
						haventGotResult = false
						break
					}

					resfloat, errfloat := strconv.ParseFloat(er.Result, 64)
					if errfloat != nil {
						panic(errfloat)
					}

					if math.Abs(resfloat-tt.expected) > 0.0001 {
						t.Fatalf("expected %v, got %v", tt.expected, resfloat)
					}
					if er.Status != fmt.Sprint(tt.expectedStatusCode) {
						t.Fatalf("expected %v, got %v", tt.expectedStatusCode, er.Status)
					}
					if tt.wantError && er.Status == "200" {
						t.Fatalf("expected error, got success")
					}
					haventGotResult = false
					break
				}
			}
		})
	}
}

//}
