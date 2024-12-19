package application

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	calc "github.com/Barsenick/calculator/pkg/calc"
)

type Application struct {
}

func New() *Application {
	return &Application{}
}

type Request struct {
	Expression string `json:"expression"`
}

type ResponseSuccess struct {
	Result string `json:"result"`
}

type ResponseError struct {
	Error string `json:"error"`
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	request := new(Request)

	defer r.Body.Close()
	err1 := json.NewDecoder(r.Body).Decode(&request)
	if err1 != nil {
		log.Println(request)
		http.Error(w, "Invalid request body: "+request.Expression, http.StatusBadRequest)
		return
	}

	result, err2 := calc.Calc(request.Expression)
	if err2 != nil {
		response, err3 := json.Marshal(ResponseError{Error: fmt.Sprintf("%v", err2)})
		if err3 != nil {
			http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
			return
		}
		http.Error(w, string(response), http.StatusUnprocessableEntity)
	} else {
		response, err3 := json.Marshal(ResponseSuccess{Result: fmt.Sprintf("%v", result)})
		if err3 != nil {
			http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, string(response))
	}
}

// Функция запуска приложения
// тут будем читать введенную строку и после нажатия ENTER писать результат работы программы на экране
// если пользователь ввел exit - то останаваливаем приложение
// func (a *Application) Run() error {
// 	for {
// 		// читаем выражение для вычисления из командной строки
// 		log.Println("input expression")
// 		reader := bufio.NewReader(os.Stdin)
// 		text, err := reader.ReadString('\n')
// 		if err != nil {
// 			log.Println("failed to read expression from console")
// 		}
// 		// убираем пробелы, чтобы оставить только вычислемое выражение
// 		text = strings.TrimSpace(text)
// 		// выходим, если ввели команду "exit"
// 		if text == "exit" {
// 			log.Println("aplication was successfully closed")
// 			return nil
// 		}
// 		//вычисляем выражение
// 		result, err := calc.Calc(text)
// 		if err != nil {
// 			log.Println(text, " calculation failed with error: ", err)
// 		} else {
// 			log.Println(text, "=", result)
// 		}
// 	}
// }

func (a *Application) RunServer() error {
	http.HandleFunc("/api/v1/calculate", CalcHandler)

	log.Println("Starting server on :8080")
	return http.ListenAndServe(":8080", nil)
}
