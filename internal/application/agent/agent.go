package application

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/Barsenick/calculator/pkg/calc"
)

func SolveOperation(op rune, arg1, arg2 float64) (float64, error) {
	switch op {
	case '+':
		return arg1 + arg2, nil
	case '-':
		return arg1 - arg2, nil
	case '*':
		return arg1 * arg2, nil
	case '/':
		if arg2 == 0 {
			return 0, calc.Err422
		}
		return arg1 / arg2, nil
	case '^':
		return math.Pow(arg1, arg2), nil
	}
	return 0, calc.Err422
}

func StartAgent() {
	log.Println("agent started")
	url := "http://localhost:8080/internal/task"

	for {
		time.Sleep(10 * time.Millisecond)
		request, err1 := http.NewRequest(http.MethodGet, url, nil)
		if err1 != nil {
			panic(err1)
		}

		client := &http.Client{}
		response, err2 := client.Do(request)
		if err2 != nil {
			panic(err2)
		}

		if response.StatusCode == 200 {
			task := calc.Task{}
			bodyBytes, err := io.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}

			err3 := json.Unmarshal(bodyBytes, &task)
			if err3 != nil {
				response.Body.Close()
				panic(err3)
			}
			if task.Operation != 0 {
				calc.Tasks.M.Lock()
				id := len(calc.Tasks.Tasks) - 1
				calc.Tasks.M.Unlock()

				c := make(chan string, 1)
				var res float64
				var errop error
				go func() {
					res, errop = SolveOperation(task.Operation, task.Arg1, task.Arg2)
					if errop != nil {
						c <- errop.Error()
						return
					}
					c <- "success"
				}()

				select {
				case opRes := <-c:
					if opRes != "success" {
						tr := calc.TaskResult{TaskID: id, Result: 0, Error: opRes}
						js, err5 := json.Marshal(tr)
						if err5 != nil {
							response.Body.Close()
							panic(err5)
						}
						request, err6 := http.NewRequest(http.MethodPost, url, bytes.NewReader(js))
						if err6 != nil {
							response.Body.Close()
							panic(err6)
						}
						client := &http.Client{}
						client.Do(request)
					} else {
						tr := calc.TaskResult{TaskID: task.TaskID, Result: res}
						js, err5 := json.Marshal(tr)
						if err5 != nil {
							response.Body.Close()
							panic(err5)
						}
						request, err6 := http.NewRequest(http.MethodPost, url, bytes.NewReader(js))
						if err6 != nil {
							response.Body.Close()
							panic(err6)
						}
						client := &http.Client{}
						client.Do(request)
						response.Body.Close()
					}
				case <-time.After(time.Duration(task.OperationTime) * time.Millisecond):
					tr := calc.TaskResult{TaskID: id, Result: 0, Error: calc.ErrTimeout.Error()}
					js, err5 := json.Marshal(tr)
					if err5 != nil {
						response.Body.Close()
						continue
					}
					request, err6 := http.NewRequest(http.MethodPost, url, bytes.NewReader(js))
					if err6 != nil {
						response.Body.Close()
						continue
					}
					client := &http.Client{}
					client.Do(request)
				}
			}
		}
	}
}
