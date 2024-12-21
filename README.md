# Go Calculator Server
This server is a simple calculator that can evaluate mathematical expressions. It supports basic arithmetic operations (+, -, *, /, ^), decimal points and handles parentheses.

# Usage
This server has 1 endpoint, ```/api/v1/calculate```. It gets the expression from POST requests, containing the expression in JSON.
``` powershell
 Invoke-RestMethod -Method Post -Uri http://localhost:8080/api/v1/calculate -ContentType 'application/json' -Body '{"expression": "5*(22.5+2.5")-2^3}'
```
Returns: ```{"result":"117"}```.

## Error 422
If the expression is not valid, the calculator will return the error 422.

``` powershell
 Invoke-RestMethod -Method Post -Uri http://localhost:8080/api/v1/calculate -ContentType 'application/json' -Body '{"expression": "2+++2"}'
```
Returns: ```{"error":"invalid expression"}```.

## Error 500
If an internal error occurs during the calculation, the calculator will return the error 500.

``` powershell
 Invoke-RestMethod -Method Post -Uri http://localhost:8080/api/v1/calculate -ContentType 'application/json' -Body '{"expression": "1+9999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999"}'
```
Returns: ```{"error":"internal server error"}``` (the number is too big for float64).
