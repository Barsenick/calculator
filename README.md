# Golang Calculator Server

This server is a simple calculator that can evaluate mathematical expressions. It provides multiple endpoints for handling calculations, retrieving expressions, and managing tasks. The server interacts with an agent that sends requests to ask for new tasks, solves them, and sends the results back to the server.

![How does this work?](https://github.com/user-attachments/assets/26364cd8-28cd-490a-81aa-87b52968bbb4)

## Features

- Supports basic arithmetic operations (+, -, *, /, ^) and decimal points
- Handles real numbers and supports parentheses
- Provides comprehensive error handling
- Manages tasks and expressions through various endpoints

## Endpoints

### API Endpoints

- **`/api/v1/calculate`**: Accepts POST requests containing an expression in JSON and returns the result or error in JSON.
- **`/api/v1/expressions`**: Retrieves a list of all expressions evaluated by the server.

### Web Page Endpoints

- **`/calculate`**: Displays the calculator web page where users can input expressions and see results.
- **`/expressions`**: Displays a list of all expressions evaluated by the server.
- **`/expression`**: Displays details of a specific expression by ID.

# Setup
You must have Golang installed.
1. **Clone the Repository:**
   Open your terminal and run the following command to clone the repository:
   
```
git clone https://github.com/Barsenick/calculator.git
```


  
2. **Navigate to the Command Directory:**
  ```
cd ./cmd
  ```

3. **Run the Server:**
Start the calculator server by executing:
```
go run main.go
```

# Usage
To send requests to the API, open PowerShell and run the command:
``` powershell
 Invoke-RestMethod -Method Post -Uri http://localhost:8080/api/v1/calculate -ContentType 'application/json' -Body '{"expression": "5*(22.5+2.5")-2^3"}'
```
Returns: ```{"result":"117"}```.

## Error 422
If the expression is not valid, the server will return error 422.

``` powershell
 Invoke-RestMethod -Method Post -Uri http://localhost:8080/api/v1/calculate -ContentType 'application/json' -Body '{"expression": "2+++2"}'
```
Returns: ```{"error":"invalid expression"}```.

## Error 500
If an internal error occurs during the calculation, the server will return error 500.

``` powershell
 Invoke-RestMethod -Method Post -Uri http://localhost:8080/api/v1/calculate -ContentType 'application/json' -Body '{"expression": "1+9999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999"}'
```
Returns: ```{"error":"internal server error"}``` (the number is too big for float64).
