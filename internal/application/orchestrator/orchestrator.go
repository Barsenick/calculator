package application

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/Barsenick/calculator/pkg/calc"
)

type Request struct {
	Expression string `json:"expression"`
}

type Response struct {
	ID string `json:"id"`
}

type Application struct {
}

func New() *Application {
	return &Application{}
}

func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	tr := calc.TaskResult{}
	err1 := json.NewDecoder(r.Body).Decode(&tr)
	if err1 != nil {
		calc.Tasks.M.Lock()
		if len(calc.Tasks.Tasks) == len(calc.TaskResults.TaskResults) {
			fmt.Fprint(w, "{}")
		} else {
			task := calc.Tasks.Tasks[len(calc.Tasks.Tasks)-1]
			js, err := json.Marshal(task)
			if err != nil {
				http.Error(w, err.Error(), 500)
			}
			calc.TaskResults.M.Lock()
			calc.TaskResults.TaskResults = append(calc.TaskResults.TaskResults, calc.TaskResult{TaskID: task.TaskID, Result: 0, Error: ""})
			calc.TaskResults.M.Unlock()
			fmt.Fprintf(w, "%v", string(js))
		}
		calc.Tasks.M.Unlock()
		return
	}

	calc.TaskResults.M.Lock()
	calc.TaskResults.TaskResults[tr.TaskID].Result = tr.Result
	calc.TaskResults.TaskResults[tr.TaskID].Error = tr.Error
	calc.TaskResults.M.Unlock()
	calc.Wg.Done()

}

func ApiCalcHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)

	if r.Method == http.MethodOptions {
		return
	}

	ClientRequest := new(Request)

	clientIP := r.RemoteAddr
	xForwardedFor := r.Header.Get("X-Forwarded-For")

	if xForwardedFor != "" {
		clientIP = xForwardedFor
	}

	defer r.Body.Close()
	bodyBytes, err1 := io.ReadAll(r.Body)
	if err1 != nil {
		log.Printf("Error reading request body from %s: %s\n", clientIP, err1.Error())
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	err2 := json.Unmarshal(bodyBytes, &ClientRequest)
	if err2 != nil {
		log.Printf("Invalid ClientRequest body from %s: %s\n", clientIP, string(bodyBytes))
		http.Error(w, "Invalid ClientRequest body: "+ClientRequest.Expression, http.StatusBadRequest)
		return
	}

	log.Printf("Request from %s: %s\n", clientIP, ClientRequest.Expression)

	calc.Exprs.M.Lock()
	id := len(calc.Exprs.Expressions)
	jsonid, err3 := json.Marshal(calc.ID{ID: id})
	if err3 != nil {
		http.Error(w, err3.Error(), http.StatusInternalServerError)
		return
	}
	calc.Exprs.M.Unlock()

	_, errwr := fmt.Fprint(w, string(jsonid))
	if errwr != nil {
		http.Error(w, errwr.Error(), http.StatusInternalServerError)
	}

	go func() {
		calc.Exprs.M.Lock()
		calc.Exprs.Expressions = append(calc.Exprs.Expressions, calc.Expression{ID: id, Status: 201, Result: "pending"})
		calc.Exprs.M.Unlock()

		res, errCalc := calc.Calc(ClientRequest.Expression)
		calc.Exprs.M.Lock()
		if errCalc != nil {
			if errCalc == calc.Err422 {
				calc.Exprs.Expressions[id].Status = 422
			} else {
				calc.Exprs.Expressions[id].Status = 500
			}
			calc.Exprs.Expressions[id].Result = errCalc.Error()
		} else {
			calc.Exprs.Expressions[id].Status = 200
			calc.Exprs.Expressions[id].Result = fmt.Sprintf("%f", res)
		}
		calc.Exprs.M.Unlock()
	}()
}

func ApiExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)

	if r.Method == http.MethodOptions {
		return
	}

	idstr := r.URL.Query().Get("id")
	if idstr != "" {
		id, err1 := strconv.Atoi(idstr)
		if err1 != nil {
			http.Error(w, err1.Error(), 500)
		} else {
			calc.Exprs.M.Lock()
			if id >= len(calc.Exprs.Expressions) || id < 0 {
				http.Error(w, http.StatusText(404), 404)
			} else {
				expr := calc.Exprs.Expressions[id]
				json, err2 := json.Marshal(expr)
				if err2 != nil {
					http.Error(w, err2.Error(), 500)
				} else {
					fmt.Fprintf(w, "%v", string(json))
				}
			}
			calc.Exprs.M.Unlock()
		}
	} else {
		calc.Exprs.M.Lock()
		json, err := json.Marshal(calc.Exprs)
		calc.Exprs.M.Unlock()
		if err != nil {
			http.Error(w, err.Error(), 500)
		} else {
			fmt.Fprintf(w, "%v", string(json))
		}
	}
}

func CalcPageHandler(w http.ResponseWriter, r *http.Request) {
	// Render the calculate.html template
	tmpl, err := template.ParseFiles("../html_templates/html/calculate.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ExpressionsPageHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)

	if r.Method == http.MethodOptions {
		return
	}

	tmpl, err := template.ParseFiles("../html_templates/html/expressions.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ExpressionPageHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)

	tmpl, err := template.ParseFiles("../html_templates/html/expression.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func EverythingPageHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)

	tmpl, err := template.ParseFiles("../html_templates/html/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *Application) RunServer() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/internal/task", TasksHandler)

	mux.HandleFunc("/api/v1/calculate", ApiCalcHandler)
	mux.HandleFunc("/api/v1/expressions", ApiExpressionsHandler)

	mux.HandleFunc("/calculate", CalcPageHandler)
	mux.HandleFunc("/expressions", ExpressionsPageHandler)
	mux.HandleFunc("/expression", ExpressionPageHandler)
	mux.HandleFunc("/everything", EverythingPageHandler)

	mux.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("../html_templates/css"))))
	mux.Handle("/js/", http.StripPrefix("/js", http.FileServer(http.Dir("../html_templates/js"))))
	//mux.Handle("/icons/", http.StripPrefix("/icons", http.FileServer(http.Dir("../html_templates/icons"))))

	log.Println("Starting server on", calc.Port)
	err2 := http.ListenAndServe(calc.Port, mux)

	return err2
}
