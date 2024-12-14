package vague

// Position in the source template file
type position struct {
	line   int
	offset int
	col    int
}

// Node Types
const (
	TEXT = iota
	ELEMENT
	VIRTUAL
)

// Token Types
const (
	IF_TOKEN      = "v-if"      //boolean expression
	ELSE_TOKEN    = "v-else"    //no value
	ELSE_IF_TOKEN = "v-else-if" //boolean expression

	FOR_TOKEN     = "v-for"                     //map, slice, array
	EVENT_TOKEN   = `@(\p{L}+)=`                //function,primitive operation(++,--)
	BINDING_TOKEN = `:([a-zA-Z][a-zA-Z0-9-]*)=` //string,variable for two way data binding
	PROPS_TOKEN   = `:([a-zA-Z][a-zA-Z0-9-]*)=` // complex data, any
	// on special attribute bindings have to check inside BINDING_TOKE for the names and then call functions from there
	CLASS_TOKEN   = ":class" //Add(class,boolean Expression) AddAll(variadic string...)
	STYLE_TOKEN   = ":style" //Add(style,value stringExpressions) AddUnit(style stringExpr,value numericalExpr,Unit stringExpr)
	VIRTUAL_TOKEN = "template"
)

type LoopInfo struct {
	Expr   string `json:"Expr,omitempty"`
	Var    string `json:"Var,omitempty"`
	KeyVar string `json:"KeyVar,omitempty"`
}

// ConditionalInfo stores the information of the if else and elseif nodes
// if it has a condition its an if node if it has none it an else node
// if it has a condition and Else is true it is an elseif node
type ConditionInfo struct {
	Condition string `json:"Condition,omitempty"`
	IsElse    bool   `json:"IsElse,omitempty"`
	IsElseIf  bool   `json:"IsElseIf,omitempty"`
}
type Attribute struct {
	Value  string   `json:"Value,omitempty"`
	Values []string `json:"Values,omitempty"`
}

type EventInfo struct {
	Name       string
	Options    []*EventOption
	Expression string //Call Signature of the event handler
}
type EventOption struct {
	Capture bool
	Once    bool
	Passive bool
	Prevent bool //preventDefault
}

type DefaultHandler struct {
	Arguments []string
}

// Node is a recursive and complete way of describing the DOM and the template logic
// need a function  to evaluate wether this node is dynamic and skip evaluation if it doesn't maybe have
// to handle [Directives] and [Interpolated content and attributes] separately
// as in skip interpolation, skip directives
type Node struct {
	Type       int                   `json:"Type"`                 // "element", "text", virtual
	TagName    string                `json:"TagName,omitempty"`    // For elements
	Attributes map[string]*Attribute `json:"Attributes,omitempty"` // For elements
	//BoolAttrs     map[string]string     `json:"BoolAttrs,omitempty"`     // For Elements
	Content  string  `json:"Content,omitempty"`  // For text and CDATA
	Children []*Node `json:"Children,omitempty"` // For child nodes
	//Events        map[string]string     `json:"Events,omitempty"`        // Event binding value is the whole function with parameters verbatim
	LoopInfo      *LoopInfo      `json:"LoopInfo,omitempty"`      // Pointer to loop information (nil if not a loop)
	ConditionInfo *ConditionInfo `json:"ConditionInfo,omitempty"` // Pointer to conditional info (nil if not conditional)
	//verbatim      string                `json:"-"`                       // For caching stuff that doesn't need rerendering
	//dirty         bool                  `json:"-"`                       // For indicating that this element needs to rerender
}

type parseError struct {
	Code    int
	Message string
	Extra   []any
}
