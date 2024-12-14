package expressions

import (
	"fmt"
	"reflect"
)

var currentSymbolTable SymbolTable // Global symbol table (adjust visibility as needed)

// need to flesh this one out later comment out for now
/*func evaluateExpression(expr string) (interface{}, error) {
	// ... expression parsing logic ...

	// Access and update the symbol table during evaluation
	symbol, ok := currentSymbolTable[symbolName]
	if !ok {
		// Handle symbol not found error
		return nil, errors.New("symbol not found: " + symbolName)
	}

	// Use the symbol's type information and other details for evaluation

	// ... logic to use the symbol information ...

	return result, nil
}*/
//AddFunc adds a function to the symbol table
func AddFunc(name string, function interface{}) {
	funcValue := reflect.ValueOf(function)
	funcType := funcValue.Type()

	if funcType.Kind() != reflect.Func {
		panic(fmt.Errorf("provided interface is not a function but a %v", funcType.Kind()))
	}

	numParams := funcType.NumIn()
	params := make([]reflect.Type, numParams)

	for i := 0; i < numParams; i++ {
		params[i] = funcType.In(i)
	}

	currentSymbolTable[name] = Symbol{
		Name:     name,
		Type:     funcType,
		Callable: true,
		Params:   params,
		Variadic: funcType.IsVariadic(),
	}
}

func init() {
	currentSymbolTable = make(SymbolTable)
}

// Example function
func MyFunction(a int, b string) string {
	return fmt.Sprintf("%d: %s", a, b)
}

func main() {
	AddFunc("myFunc", MyFunction)
	fmt.Println(currentSymbolTable["myFunc"])
}
