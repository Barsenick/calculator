package application

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"math"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/Barsenick/calculator/pkg/calc"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/mattn/go-sqlite3"
)

var (
	ErrNoToken                = errors.New("no token found")
	ErrUniqueConstraintFailed = errors.New("UNIQUE constraint failed: users.login")
	ErrInvalidToken           = errors.New("token is invalid")
)

const hmacSampleSecret = "calculator_service_signature2"

type Request struct {
	Expression string `json:"expression"`
}

type Expressions struct {
	Expressions []Expression `json:"expressions"`
}

type RegistrationRequest struct {
	Name     string `json:"login"`
	Password string `json:"password"`
}

type RegistrationResponse struct {
	Status string `json:"status"`
	Token  string `json:"token,omitempty"`
}

type Expression struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Result  string `json:"result"`
	OwnerID int64  `json:"-"`
}

type Response struct {
	ID string `json:"id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Application struct {
}

func New() *Application {
	return &Application{}
}

type User struct {
	ID             int64
	Name           string
	Password       string
	OriginPassword string
}

var DB *sql.DB

func generateErrorResponse(w http.ResponseWriter, errMsg string, code int) {
	json, err := json.Marshal(ErrorResponse{Error: errMsg})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	fmt.Fprint(w, string(json))
}

func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", 404)
}

func getToken(r *http.Request, claims *jwt.MapClaims) (*jwt.Token, error) {
	if claims == nil {
		tokenCookie, err := r.Cookie("token")
		if err != nil || tokenCookie.Value == "" {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				return nil, ErrNoToken
			} else {
				tokenString := strings.TrimPrefix(authHeader, "Bearer ")
				token, err := parseToken(tokenString)
				if err != nil || !token.Valid {
					return nil, ErrInvalidToken
				}
				return token, nil
			}
		} else {
			token, err := parseToken(tokenCookie.Value)
			if err != nil || !token.Valid {
				return nil, ErrInvalidToken
			}
			return token, nil
		}
	} else {
		tokenCookie, err := r.Cookie("token")
		if err != nil || tokenCookie.Value == "" {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				return nil, ErrNoToken
			} else {
				tokenString := strings.TrimPrefix(authHeader, "Bearer ")
				token, err := parseTokenWithClaims(tokenString, claims)
				if err != nil || !token.Valid {
					return nil, ErrInvalidToken
				}
				return token, nil
			}
		} else {
			token, err := parseTokenWithClaims(tokenCookie.Value, claims)
			if err != nil || !token.Valid {
				return nil, ErrInvalidToken
			}
			return token, nil
		}
	}

}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := getToken(r, nil)
		if err != nil || !token.Valid {
			http.Redirect(w, r, "/login", http.StatusFound)
		}

		next.ServeHTTP(w, r)
	})
}

func auth(tokenstr string) *jwt.Token {
	if tokenstr == "" {
		return nil
	}

	token, err := parseToken(tokenstr)
	if err != nil || !token.Valid {
		return nil
	}

	return token
}

func authWithClaims(header string, claims *jwt.MapClaims) *jwt.Token {
	if header == "" {
		return nil
	}

	tokenString := strings.TrimPrefix(header, "Bearer ")
	token, err := parseTokenWithClaims(tokenString, claims)
	if err != nil || !token.Valid {
		return nil
	}

	return token
}

func parseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(hmacSampleSecret), nil
	})
}

func parseTokenWithClaims(tokenString string, claims *jwt.MapClaims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(hmacSampleSecret), nil
	})
}

func (u User) ComparePassword(u2 User) error {
	err := compare(u2.Password, u.OriginPassword)
	if err != nil {
		return err
	}

	return nil
}

func createUsersTable(ctx context.Context, DB *sql.DB) error {
	const usersTable = `
	CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		login TEXT UNIQUE,
		password TEXT
	);`

	if _, err := DB.ExecContext(ctx, usersTable); err != nil {
		return err
	}

	return nil
}

func getUserExpressions(ctx context.Context, DB *sql.DB, userID int64) ([]Expression, error) {
	selectExpressions := `
	SELECT id, status, result FROM expressions
	WHERE ownerID = ?;`

	rows, err := DB.QueryContext(ctx, selectExpressions, userID)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var expressions []Expression
	for rows.Next() {
		var expression Expression
		if err := rows.Scan(&expression.ID, &expression.Status, &expression.Result); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}
		expressions = append(expressions, expression)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		return nil, err
	}

	return expressions, nil
}

func createExpressionsTable(ctx context.Context, DB *sql.DB) error {
	const expressionsTable = `
	CREATE TABLE IF NOT EXISTS expressions(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		ownerID TEXT,
		status TEXT,
		result TEXT,
		FOREIGN KEY(ownerID) REFERENCES users(id)
	);`

	if _, err := DB.ExecContext(ctx, expressionsTable); err != nil {
		return err
	}

	return nil
}

func insertUser(ctx context.Context, DB *sql.DB, user *User) (int64, error) {
	var q = `
	INSERT INTO users (login, password) values ($1, $2)
	`
	result, err := DB.ExecContext(ctx, q, user.Name, user.Password)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func insertExpression(ctx context.Context, DB *sql.DB, expr *Expression) (int64, error) {
	var q = `
	INSERT INTO expressions (status, result, ownerID) values ($1, $2, $3)
	`
	result, err := DB.ExecContext(ctx, q, &expr.Status, &expr.Result, &expr.OwnerID)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func modifyExpression(ctx context.Context, DB *sql.DB, expr *Expression) (int64, error) {
	var q = `
	UPDATE expressions SET status = $1, result = $2, ownerID = $3 WHERE id = $f;
	`
	result, err := DB.ExecContext(ctx, q, &expr.Status, &expr.Result, &expr.OwnerID, &expr.ID)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func selectUser(ctx context.Context, DB *sql.DB, login string) (User, error) {
	var (
		user User
		err  error
	)

	var q = "SELECT id, login, password FROM users WHERE login=$1"
	err = DB.QueryRowContext(ctx, q, login).Scan(&user.ID, &user.Name, &user.Password)
	return user, err
}

func selectUserByID(ctx context.Context, DB *sql.DB, id int64) (User, error) {
	var (
		user User
		err  error
	)

	var q = "SELECT id, login, password FROM users WHERE id=$1"
	err = DB.QueryRowContext(ctx, q, id).Scan(&user.ID, &user.Name, &user.Password)
	return user, err
}

func selectExpression(ctx context.Context, DB *sql.DB, id string) (Expression, error) {
	var (
		expr Expression
		err  error
	)

	var q = "SELECT id, status, result, ownerID FROM expressions WHERE id=$1"
	err = DB.QueryRowContext(ctx, q, id).Scan(&expr.ID, &expr.Status, &expr.Result, &expr.OwnerID)
	return expr, err
}

func generate(s string) (string, error) {
	saltedBytes := []byte(s)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	hash := string(hashedBytes[:])
	return hash, nil
}

func generateToken(id int64) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"nbf": now.Unix(),
		"exp": now.Add(168 * time.Hour).Unix(),
		"iat": now.Unix(),
	})

	tokenString, err := token.SignedString([]byte(hmacSampleSecret))
	if err != nil {
		return "", nil
	}
	return tokenString, nil
}

func compare(hash string, s string) error {
	incoming := []byte(s)
	existing := []byte(hash)
	return bcrypt.CompareHashAndPassword(existing, incoming)
}

func OpenDB() (*sql.DB, error) {
	ctx := context.TODO()

	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		return nil, err
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	if err = createUsersTable(ctx, db); err != nil {
		return nil, err
	}

	if err = createExpressionsTable(ctx, db); err != nil {
		return nil, err
	}

	return db, nil
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	tr := calc.TaskResult{}
	err1 := json.NewDecoder(r.Body).Decode(&tr)
	calc.Tasks.M.Lock()
	defer calc.Tasks.M.Unlock()
	if err1 != nil {
		if len(calc.Tasks.Tasks) == len(calc.TaskResults.TaskResults) {
			fmt.Fprint(w, "{}")
		} else {
			task := calc.Tasks.Tasks[len(calc.Tasks.Tasks)-1]
			js, err := json.Marshal(task)
			if err != nil {
				http.Error(w, err.Error(), 500)
			}
			calc.TaskResults.TaskResults = append(calc.TaskResults.TaskResults, calc.TaskResult{TaskID: task.TaskID, Result: 0, Error: ""})
			fmt.Fprintf(w, "%v", string(js))
		}
		return
	}

	calc.TaskResults.TaskResults[tr.TaskID].Result = tr.Result
	calc.TaskResults.TaskResults[tr.TaskID].Error = tr.Error
	calc.Wg.Done()
}

func ApiCalcHandler(w http.ResponseWriter, r *http.Request) {
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

	ctx := context.TODO()

	claims := jwt.MapClaims{}
	token, err := getToken(r, &claims)
	if err != nil || !token.Valid {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	ownerID := int64(math.Floor(claims["id"].(float64)))

	expr := Expression{ID: "-1", Status: "201", Result: "pending", OwnerID: ownerID}

	id, err := insertExpression(ctx, DB, &expr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	expr.ID = fmt.Sprint(id)

	jsonid, err3 := json.Marshal(calc.ID{ID: id})
	if err3 != nil {
		http.Error(w, err3.Error(), http.StatusInternalServerError)
		return
	}

	_, errwr := fmt.Fprint(w, string(jsonid))
	if errwr != nil {
		http.Error(w, errwr.Error(), http.StatusInternalServerError)
	}

	go func(DB *sql.DB) {
		ctx := context.TODO()

		res, errCalc := calc.Calc(ClientRequest.Expression)
		if errCalc != nil {
			if errCalc == calc.Err422 {
				expr.Status = "422"
			} else {
				expr.Status = "500"
			}
			expr.Result = errCalc.Error()
		} else {
			expr.Status = "200"
			expr.Result = fmt.Sprintf("%f", res)
		}
		_, err := modifyExpression(ctx, DB, &expr)
		if err != nil {
			log.Println(err.Error())
		}
	}(DB)
}

func ApiExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}

	ctx := context.TODO()

	claims := jwt.MapClaims{}
	token, err := getToken(r, &claims)
	if err != nil || !token.Valid {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	uid := int64(math.Floor(claims["id"].(float64)))

	idstr := r.URL.Query().Get("id")
	if idstr != "" {
		expr, err := selectExpression(ctx, DB, idstr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
		}

		if expr.OwnerID == uid {
			json, err2 := json.Marshal(expr)
			if err2 != nil {
				http.Error(w, err2.Error(), 500)
				return
			}
			fmt.Fprintf(w, "%v", string(json))
		} else {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	} else {
		actualexprs, err := getUserExpressions(ctx, DB, uid)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		exprs := Expressions{Expressions: actualexprs}
		json, err := json.Marshal(exprs)
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

func ApiRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}

	ClientRequest := new(RegistrationRequest)

	defer r.Body.Close()
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(bodyBytes, &ClientRequest)
	if err != nil {
		http.Error(w, "Invalid ClientRequest", http.StatusBadRequest)
		return
	}

	if len(ClientRequest.Name) < 3 && len(ClientRequest.Password) < 3 {
		generateErrorResponse(w, "Username and password must be at least 3 letters long", 401)
		return
	} else if len(ClientRequest.Name) < 3 {
		generateErrorResponse(w, "Username must be at least 3 letters long", 401)
		return
	} else if len(ClientRequest.Password) < 3 {
		generateErrorResponse(w, "Password must be at least 3 letters long", 401)
		return
	}

	ctx := context.TODO()

	hash, err := generate(ClientRequest.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := User{-1, ClientRequest.Name, hash, ClientRequest.Password}

	uid, err := insertUser(ctx, DB, &user)
	if err != nil {
		if err.Error() == ErrUniqueConstraintFailed.Error() {
			generateErrorResponse(w, "Username had already been taken", 401)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	token, err := generateToken(uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(RegistrationResponse{Status: "200 OK", Token: token})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(json)+"")
}

func RegistrationPageHandler(w http.ResponseWriter, r *http.Request) {
	_, err := getToken(r, &jwt.MapClaims{})
	if err == ErrNoToken || err == ErrInvalidToken {
		tmpl, err := template.ParseFiles("../html_templates/html/register.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		http.Redirect(w, r, "/calculate", http.StatusFound)
	}
}

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	_, err := getToken(r, &jwt.MapClaims{})
	if err == ErrNoToken || err == ErrInvalidToken {
		tmpl, err := template.ParseFiles("../html_templates/html/login.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		http.Redirect(w, r, "/calculate", http.StatusFound)
	}
}

func ApiLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}

	ClientRequest := new(RegistrationRequest)

	defer r.Body.Close()
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(bodyBytes, &ClientRequest)
	if err != nil {
		http.Error(w, "Invalid ClientRequest", http.StatusBadRequest)
		return
	}

	user := User{-1, ClientRequest.Name, "", ClientRequest.Password}

	ctx := context.TODO()

	userFromDB, err := selectUser(ctx, DB, ClientRequest.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			generateErrorResponse(w, "This username doesnt exist", 401)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	user.ID = userFromDB.ID

	if err := user.ComparePassword(userFromDB); err == nil {
		token, err := generateToken(user.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json, err := json.Marshal(RegistrationResponse{Status: "200 OK", Token: token})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, string(json)+"")

	} else if err == bcrypt.ErrMismatchedHashAndPassword {
		generateErrorResponse(w, "Invalid password", 401)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// func headerMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
// 		w.Header().Add("Content-Type", "application/json")
// 		next.ServeHTTP(w, req)
// 	})
// }

func panicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				log.Println(string(debug.Stack()))
			}
		}()
		next.ServeHTTP(w, req)
	})
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		w.Header().Set("Access-Control-Allow-Origin", origin)
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			return
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func withMiddlewareFunc(handler http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) http.HandlerFunc {
	h := http.Handler(handler)
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h.ServeHTTP
}

func pathHandler(w http.ResponseWriter, r *http.Request) {
	middlewares := []func(http.Handler) http.Handler{panicRecovery /*headerMiddleware,*/, CORSMiddleware, authMiddleware}
	switch r.URL.Path {
	case "/api/v1/register":
		withMiddlewareFunc(ApiRegistrationHandler, middlewares[:2]...)(w, r)
	case "/api/v1/login":
		withMiddlewareFunc(ApiLoginHandler, middlewares[:2]...)(w, r)
	case "/internal/task":
		withMiddlewareFunc(TasksHandler, middlewares[:1]...)(w, r)
	case "/api/v1/calculate":
		withMiddlewareFunc(ApiCalcHandler, middlewares...)(w, r)
	case "/api/v1/expressions":
		withMiddlewareFunc(ApiExpressionsHandler, middlewares...)(w, r)
	case "/register":
		withMiddlewareFunc(RegistrationPageHandler, middlewares[:2]...)(w, r)
	case "/login":
		withMiddlewareFunc(LoginPageHandler, middlewares[:2]...)(w, r)
	case "/calculate":
		withMiddlewareFunc(CalcPageHandler, middlewares...)(w, r)
	case "/expressions":
		withMiddlewareFunc(ExpressionsPageHandler, middlewares...)(w, r)
	case "/expression":
		withMiddlewareFunc(ExpressionPageHandler, middlewares...)(w, r)
	case "/everything":
		withMiddlewareFunc(EverythingPageHandler, middlewares...)(w, r)
	default:
		r.Header.Add("Content-Type", "")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func (a *Application) RunServer() error {
	//middlewares := []func(http.Handler) http.Handler{panicRecovery /*headerMiddleware,*/, CORSMiddleware, authMiddleware}
	mux := http.NewServeMux()

	// mux.HandleFunc("/api/v1/register", withMiddlewareFunc(ApiRegistrationHandler, middlewares[:2]...))
	// mux.HandleFunc("/api/v1/login", withMiddlewareFunc(ApiLoginHandler, middlewares[:2]...))

	// mux.HandleFunc("/internal/task", withMiddlewareFunc(TasksHandler, middlewares[:1]...))

	// mux.HandleFunc("/api/v1/calculate", withMiddlewareFunc(ApiCalcHandler, middlewares...))
	// mux.HandleFunc("/api/v1/expressions", withMiddlewareFunc(ApiExpressionsHandler, middlewares...))

	// mux.HandleFunc("/register", withMiddlewareFunc(RegistrationPageHandler, middlewares[:2]...))
	// mux.HandleFunc("/login", withMiddlewareFunc(LoginPageHandler, middlewares[:2]...))

	// mux.HandleFunc("/calculate", withMiddlewareFunc(CalcPageHandler, middlewares...))
	// mux.HandleFunc("/expressions", withMiddlewareFunc(ExpressionsPageHandler, middlewares...))
	// mux.HandleFunc("/expression", withMiddlewareFunc(ExpressionPageHandler, middlewares...))

	// mux.HandleFunc("/everything", withMiddlewareFunc(EverythingPageHandler, middlewares...))

	mux.HandleFunc("/", pathHandler)

	mux.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("../html_templates/css"))))
	mux.Handle("/js/", http.StripPrefix("/js", http.FileServer(http.Dir("../html_templates/js"))))
	mux.Handle("/icons/", http.StripPrefix("/icons", http.FileServer(http.Dir("../html_templates/icons"))))

	log.Println("Starting server on", calc.Port)
	err := http.ListenAndServe(calc.Port, mux)

	return err
}
