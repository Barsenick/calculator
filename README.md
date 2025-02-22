# Golang Calculator Server
This server is a simple calculator that can evaluate mathematical expressions. It has 1 endpoint, ```/api/v1/calculate```. It gets the expression from POST requests, containing the expression in JSON and returns the result or error in JSON.

# Features
- Supports basic arithmetic operations (+, -, *, /, ^) and decimal points
- Handles real numbers and supports parenthesis
- Provides comprehensive error handling

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
To send requests, open PowerShell and run the command:
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
