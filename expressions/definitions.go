package expressions

import "reflect"

type Symbol struct {
	Name    string
	Type    reflect.Type // Use reflect.Type for flexibility
	Package string       // Package name if it's an imported symbol
	// Add other relevant info like constant values, etc. if needed
}

type SymbolTable map[string]Symbol
