package vague

import "fmt"

// Error Types
const (
	ErrNone = iota
	ErrIfNoExpr
	ErrMultiIf

	ErrElseIfNoIf      // 'v-else-if' without preceding 'v-if' or 'v-else-if'
	ErrElseIfAfterElse // 'v-else-if' after 'v-else'
	ErrMultipleElseIf  // Multiple 'v-else-if' on the same element
	ErrCompositeElseIf // trying to create v-else-if from v-else and v-if
	ErrMixedElseElseIf // a mix of v-else and v-else-if on one element
	ErrElseIfNoExpr

	ErrElseAfterElse
	ErrMultiElse
	ErrElseNoIf
	ErrElseCondition

	ErrForNoExpr
	ErrMultiFor
	//ErrInvalidFor
	ErrForKeyVarInvalid
	ErrForVarInvalid
	ErrForNoIterable
	ErrForNoIn
	ErrForWithCondition
	ErrForMultipleIn
	ErrForExtraComma

	ErrVirtAttr  //virtual tag that is neither component nor template envelope, can't have attributes only directives
	ErrMultiRoot // multiple tags on top level
	ErrStrayTags // TODO we have Tags without counterpart
	ErrUnclosedRoot

	ErrMismatchedTag
)

func (e *parseError) Error() string {
	if len(e.Extra) > 0 {
		e.Message = fmt.Sprintf(e.Message, e.Extra...)
	}
	if fileName != "" {
		return fmt.Sprintf("Code: [%d] %s:%v:%v %s", e.Code, fileName, pos.line, pos.col, e.Message)
	}
	return fmt.Sprintf("Code: [%d] %v:%v %s", e.Code, pos.line, pos.col, e.Message)
}

// attribute parsing errors
// if errors
var (
	multipleIfMessage = &parseError{
		Message: `multiple '%[1]s' blocks found for the same element. 
		Only one '%[1]s' block allowed per element`,
		Extra: []any{IF_TOKEN},
		Code:  ErrMultiIf,
	}
	ifNoExpressionMessage = &parseError{
		Message: `'%[1]s' block found without expression. 
			'%[1]s' blocks syntax is [%[1]s=isValid]`,
		Extra: []any{IF_TOKEN},
		Code:  ErrIfNoExpr,
	}
)

// else errors
var (
	elseWithoutIfMessage = &parseError{
		Message: `Error: Encountered a '%[1]s' without a preceding '%[2]s'.
		An '%[1]s' block must always follow a '%[2]s' or '%[3]s' block.`,
		Extra: []any{ELSE_TOKEN, IF_TOKEN, ELSE_IF_TOKEN},
		Code:  ErrElseNoIf,
	}
	elseAfterElseMessage = &parseError{
		Message: `Error: '%[1]s' block found after preceeding '%[1]s' block.
		Only one '%[1]s' block is allowed per '%[2]s' '%[3]s' condition.`,
		Extra: []any{ELSE_TOKEN, IF_TOKEN, ELSE_IF_TOKEN},
		Code:  ErrElseAfterElse,
	}
	multipleElseMessage = &parseError{
		Message: `Error: multiple '%[1]s' directives found on one element.
		Only one '%[1]s' block is allowed per element.`,
		Extra: []any{ELSE_TOKEN},
		Code:  ErrMultiElse,
	}
	elseWithConditionMessage = func(val string) *parseError {
		return &parseError{Message: `Error: '%[1]s' directive found containing condition=%[2]s.
		Either use '%[3]s' or remove the condition.`,
			Extra: []any{ELSE_TOKEN, val, ELSE_IF_TOKEN},
			Code:  ErrElseCondition,
		}
	}
)

// elseif errors
var (
	elseIfAfterElseMessage = &parseError{
		Message: `Error: '%[1]s' block found after preceeding '%[1]s' block.
		Only one '%[1]s' block is allowed per '%[2]s' '%[3]s' condition.`,
		Extra: []any{ELSE_IF_TOKEN, IF_TOKEN, ELSE_IF_TOKEN},
		Code:  ErrElseIfAfterElse,
	}
	elseIfNoExpressionMessage = &parseError{
		Message: `'%[1]s' block found without expression.
		'%[1]s' blocks syntax is [%[1]s='isValid]`,
		Extra: []any{ELSE_IF_TOKEN},
		Code:  ErrElseIfNoExpr,
	}
	multipleElseIfMessage = &parseError{
		Message: `Error: multiple '%[1]s' directives found on one element.
		Only one '%[1]s' block is allowed per element.`,
		Extra: []any{ELSE_IF_TOKEN},
		Code:  ErrMultipleElseIf,
	}
)

// for errors
var (
	/*invalidForMessage = func(value string) *parseError {
		return &parseError{
			Message: `Error: invalid '%[1]s' expression [%[2]s].
			'for' expressions can be of the form [%[1]s='value in values'] or [%[1]s='key, value in values']`,
			Extra: []any{FOR_TOKEN, value},
			Code:  ErrInvalidFor,
		}
	}*/
	forNoExpressionMessage = &parseError{
		Message: `Error: '%[1]s' directive without expression.
		'%[1]s' expressions can be of the form [%[1]s='value in values'] or [%[1]s='key, value in values']`,
		Extra: []any{FOR_TOKEN},
		Code:  ErrForNoExpr,
	}
	multipleForStatementsMessage = &parseError{
		Message: `Error: Multiple '%[1]s' directives found on one element. Only one '%[1]s' per tag allowed`,
		Extra:   []any{FOR_TOKEN},
		Code:    ErrMultiFor,
	}
	forWithConditionMessage = &parseError{
		Message: `Error: '%[1]s' directive found together with a [%[2]s or %[3]s] directive on one element. 
		Either have the condition inside the '%[1]s' loop or outside.`,
		Extra: []any{FOR_TOKEN, IF_TOKEN, ELSE_IF_TOKEN},
		Code:  ErrForWithCondition,
	}
	forNoInMessage = &parseError{
		Message: `Error: '%[1]s' directive needs an 'in' separator between the item and the iterable. 
		if you have an 'in' separator make sure that the value before it is not empty,
		like [%[1]s=" in values"]
		'for' expressions can be of the form [%[1]s='value in values'] or [%[1]s='key, value in values']`,
		Extra: []any{FOR_TOKEN},
		Code:  ErrForNoIn,
	}
	forVarInvalidMessage = &parseError{
		Message: `Error: '%[1]s' directive of the form [%[1]s='value in values'] can't have an empty value.
		if you want to rather have the key use [%[1]s='key,_ in values']`,
		Extra: []any{FOR_TOKEN},
		Code:  ErrForVarInvalid,
	}
	forKeyVarInvalidMessage = &parseError{
		Message: `Error: '%[1]s' directive of the form [%[1]s='key, value in values'] can't have an empty value nor key.
		if you want to omit the value use [%[1]s='key,_ in values']
		if you rather want to omit the key use [%[1]s='value in values'] without the comma`,
		Extra: []any{FOR_TOKEN},
		Code:  ErrForKeyVarInvalid,
	}
	forNoIterableMessage = &parseError{
		Message: `Error: '%[1]s' directive must have an iterable after the 'in' separator.
		'for' expressions can be of the form [%[1]s='value in values'] or [%[1]s='key, value in values']`,
		Extra: []any{FOR_TOKEN},
		Code:  ErrForKeyVarInvalid,
	}
	forMultipleInMessage = &parseError{
		Message: `Error: '%[1]s' directive can have only one 'in' separator.
		'for' expressions can be of the form [%[1]s='value in values'] or [%[1]s='key, value in values']`,
		Extra: []any{FOR_TOKEN},
		Code:  ErrForMultipleIn,
	}
	forExtraCommaMessage = &parseError{
		Message: `Error: '%[1]s' directive can have only one ',' to separate key and value.
		'for' expressions can be of the form [%[1]s='value in values'] or [%[1]s='key, value in values']`,
		Extra: []any{FOR_TOKEN},
		Code:  ErrForExtraComma,
	}
)

// composite errors
var (
	mixedElseElseIfMessage = &parseError{
		Message: `Error: '%[1]s' and %[2]s directive found on one element.
		Either use '%[1]s' or %[2]s. Can't have it both ways.`,
		Extra: []any{ELSE_IF_TOKEN, ELSE_TOKEN},
		Code:  ErrMixedElseElseIf,
	}

	compositeElseIfMessage = &parseError{Message: `Error: Composite '%[1]s'.
		Use '%[1]s' instead of '%[2]s' and '%[3]s'.`,
		Extra: []any{ELSE_IF_TOKEN, IF_TOKEN, ELSE_TOKEN},
		Code:  ErrCompositeElseIf,
	}
)

//dom parsing errors (well formedness, matching tags, multiple root tags etc.)

func unclosedRootTag(tagName string) *parseError {
	return &parseError{
		Code:    ErrUnclosedRoot,
		Message: "unclosed root tag <%s>",
		Extra:   []any{tagName},
	}
}
func multipleRoots() *parseError {
	return &parseError{
		Message: `multiple 'root' elements found. 
		A Template can have only one 'root' element`,
		Code: ErrMultiRoot,
	}
}
func strayEndTags(tagName string) *parseError {
	return &parseError{
		Code:    ErrStrayTags,
		Message: "unexpected closing tag </%s>",
		Extra:   []any{tagName},
	}
}
func mismatchedTags(openingTag, closingTag string) *parseError {
	return &parseError{
		Code:    ErrMismatchedTag,
		Message: "mismatched closing tag: expected </%s>, got </%s>",
		Extra:   []any{openingTag, closingTag},
	}
}
