package expressions

import (
	"reflect"
)

// Symbol represents an identifier in the template expressions
type Symbol struct {
	Name     string         // Name of the symbol
	Type     reflect.Type   // Type of the symbol
	Package  string         // Package where the symbol is defined (if applicable)
	Constant bool           // Is the symbol a constant?
	Value    interface{}    // Value of the constant (if Constant is true)
	Callable bool           // Is the symbol callable (a function)?
	Params   []reflect.Type // Parameter types (if Callable is true)
	Variadic bool           // Is the function variadic?
}

// IsVariadic tells if the function is variadic
func (s Symbol) IsVariadic() bool {
	return s.Variadic
}

// ParameterCount returns the amount of parameters the function expects
func (s Symbol) ParameterCount() int {
	return len(s.Params)
}

// ParameterTypeAt returns the type of the parameter at the given index
func (s Symbol) ParameterTypeAt(index int) reflect.Type {
	if index < 0 || index >= len(s.Params) {
		return nil
	}
	return s.Params[index]
}

type SymbolTable map[string]Symbol
