# Calculator in Golang
This package provides a simple calculator that can evaluate mathematical expressions. The `Calc` function takes a string as input, which represents the expression to be evaluated. The function returns a float64 value representing the result of the evaluation, along with any errors that may occur during the process.

## Usage
To use the `Calc` function, import the `calc` package and call the `Calc` function with the desired expression as an argument. Here's an example:

```go
package main

import (
   "fmt"
   "github.com/Barsenick/calc"
)

func main() {
   result, err := calc.Calc("10 + 20 / 4")
   if err != nil {
       fmt.Println("Error:", err)
       return
   }
   fmt.Println("Result:", result)
}
```
## How does this work?
The function converts the equation into Reverse Polish Notation, and then solves it.
