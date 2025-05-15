# Golang Calculator Server

This server is a simple calculator that can evaluate mathematical expressions. It provides multiple endpoints for handling calculations, retrieving expressions, and managing tasks. The server interacts with an agent that sends requests to ask for new tasks, solves them, and sends the results back to the server.  User authentication is now required to access the calculator's features. All expressions and user data are stored in a SQLite3 database.

![How does this work?](https://github.com/user-attachments/assets/26364cd8-28cd-490a-81aa-87b52968bbb4)

## Features

- **Rich support**: Supports basic arithmetic operations (+, -, *, /, ^) and decimal points
- **Practical**: Handles real numbers and supports parentheses
- **Easy to handle errors**: Provides comprehensive error handling
- **Multiple endpoints**: Manages tasks and expressions through various endpoints
- **User Authentication**: Requires registration and login to access endpoints
- **Secure Password Storage**: Stores hashed passwords in the database
- **Token-based Authentication**:  Uses tokens for authenticating requests
- **Data Persistence**: Stores expressions and user data in a SQLite3 database

## Endpoints

### API Endpoints

- **`/api/v1/register`**: Accepts POST requests with user credentials in JSON format (`{"login": "user", "password":"password"}`) to register a new user.
- **`/api/v1/login`**: Accepts POST requests with user credentials in JSON format (`{"login": "user", "password":"password"}`) to authenticate a user and retrieve a JWT token.
- **`/api/v1/calculate`**: Accepts POST requests containing an expression in JSON and returns the result or error in JSON. Requires a valid JWT token in the `Authorization` header.
- **`/api/v1/expressions`**: Retrieves a list of all expressions evaluated by the server.  Requires a valid JWT token in the `Authorization` header.

- **`/internal/task`**: for agents and server communication.

### Web Page Endpoints

- **`/register`**: Displays the user registration web page.
- **`/login`**: Displays the user login web page.
- **`/calculate`**: Displays the calculator web page where users can input expressions and see results. Requires a valid JWT token to be stored in a cookie.
- **`/expressions`**: Displays a list of all expressions evaluated by the server. Requires a valid JWT token to be stored in a cookie.
- **`/expression`**: Displays details of a specific expression by ID. Requires a valid JWT token to be stored in a cookie.

# Setup
You must have Golang installed.
1. **Clone the Repository:**
   Open your terminal and run the following command to clone the repository:
   
```
git clone https://github.com/Barsenick/calculator.git
```

2. **Run the orchestrator:**
Start the calculator server by executing:
```
go run cmd/orchestrator/main.go
```

3. **Run the agents:**
Start the agents by executing:
```
go run cmd/agents/main.go
```

# Usage

## Registration

First, register a new user:

**Windows (PowerShell):**
```powershell
Invoke-RestMethod -Method Post -Uri http://localhost:8080/api/v1/register -ContentType 'application/json' -Body '{"login": "user", "password":"password"}'
```

**Linux (Bash):**
```bash
curl -X POST -H "Content-Type: application/json" -d '{"login": "user", "password":"password"}' http://localhost:8080/api/v1/register
```

This will create a new user account.

## Login

Next, log in to get a JWT token:

**Windows (PowerShell):**
```powershell
Invoke-RestMethod -Method Post -Uri http://localhost:8080/api/v1/login -ContentType 'application/json' -Body '{"login": "user", "password":"password"}'
```

**Linux (Bash):**
```bash
curl -X POST -H "Content-Type: application/json" -d '{"login": "user", "password":"password"}' http://localhost:8080/api/v1/login
```

This will return a JSON response containing your JWT token:

```json
{"token": "YOUR_JWT_TOKEN"}
```

## Making Authenticated Requests

You must include the JWT token in the `Authorization` header of each subsequent request. Replace `YOUR_JWT_TOKEN` with the actual token you received from the login request.

**Windows (PowerShell):**
```powershell
Invoke-RestMethod -Method Post -Uri http://localhost:8080/api/v1/calculate -ContentType 'application/json' -Headers @{"Authorization" = "Bearer YOUR_JWT_TOKEN"} -Body '{"expression": "5*(22.5+2.5)-2^3"}'
```

**Linux (Bash):**
```bash
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer YOUR_JWT_TOKEN" -d '{"expression": "5*(22.5+2.5)-2^3"}' http://localhost:8080/api/v1/calculate
```

Returns: ```{"result":"117"}```.

## Error 401

If the JWT is invalid or missing, the server will return error 401. Make sure to include "Bearer " before your token.

## Error 422

If the expression is not valid, the server will return error 422.

**Windows (PowerShell):**
```powershell
Invoke-RestMethod -Method Post -Uri http://localhost:8080/api/v1/calculate -ContentType 'application/json' -Headers @{"Authorization" = "Bearer YOUR_JWT_TOKEN"} -Body '{"expression": "2+++2"}'
```

**Linux (Bash):**
```bash
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer YOUR_JWT_TOKEN" -d '{"expression": "2+++2"}' http://localhost:8080/api/v1/calculate
```
Returns: ```{"error":"invalid expression"}```.

## Error 500

If an internal error occurs during the calculation, the server will return error 500.

**Windows (PowerShell):**
```powershell
Invoke-RestMethod -Method Post -Uri http://localhost:8080/api/v1/calculate -ContentType 'application/json' -Headers @{"Authorization" = "Bearer YOUR_JWT_TOKEN"} -Body '{"expression": "1+99999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999"}'
```

**Linux (Bash):**
```bash
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer YOUR_JWT_TOKEN" -d '{"expression": "1+99999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999"}' http://localhost:8080/api/v1/calculate
```
Returns: ```{"error":"internal server error"}``` (the number is too big for float64).
