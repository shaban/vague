package expressions

import (
	"fmt"
	"reflect"
)

var currentSymbolTable symbolTable

func init() {
	currentSymbolTable = make(symbolTable)
}

// AddFunc adds a function to the symbol table
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
