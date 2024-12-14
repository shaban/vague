package tests

import (
	"fmt"
	"testing"

	"github.com/shaban/vague"
)

func TestParseFor(t *testing.T) {
	tests := []testData{
		{
			name: "invalid: multiple for statements",
			input: fmt.Sprintf(`
				<div>
					<h1>Heading</h1>
					<h1 %[1]s="_,value in values" %[1]s="index in values" >Heading</h1>
				</div>`, vague.FOR_TOKEN),
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrMultiFor,
		},
		{
			name: "invalid: for with v-if condition",
			input: fmt.Sprintf(`
				<div>
					<h1>Heading</h1>
					<h1 %[1]s="show" %[2]s="_,value in values" >Heading</h1>
				</div>`, vague.IF_TOKEN, vague.FOR_TOKEN),
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrForWithCondition,
		},
		{
			name: "invalid: for with v-else-if condition",
			input: fmt.Sprintf(`
				<div>
					<h1 %[1]s="true">Heading</h1>
					<h1 %[3]s="show" %[2]s="_,value in values" >Heading</h1>
				</div>`, vague.IF_TOKEN, vague.FOR_TOKEN, vague.ELSE_IF_TOKEN),
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrForWithCondition,
		},
		{
			name: "invalid: for statement multiple 'in'",
			input: fmt.Sprintf(`
				<div>
					<p %s="k in k in values"></p>
					<h1>Heading</h1>
					<h1>Heading</h1>
				</div>`, vague.FOR_TOKEN),
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrForMultipleIn,
		},
		{
			name: "invalid: for kv statement misplaced ','",
			input: fmt.Sprintf(`
				<div>
					<p %s="k,,v in values"></p>
					<h1>Heading</h1>
					<h1>Heading</h1>
				</div>`, vague.FOR_TOKEN),
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrForExtraComma,
		},
		{
			name: "invalid: for without expression",
			input: fmt.Sprintf(`
				<div>
					<p %[1]s="true"></p>
					<h1 %[2]s>Heading</h1>
				</div>`, vague.IF_TOKEN, vague.FOR_TOKEN),
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrForNoExpr,
		},
		{
			name: "invalid: for no loop variable",
			input: fmt.Sprintf(`
				<div>
					<p %[1]s="true"></p>
					<h1 %[2]s="  in values">Heading</h1>
				</div>`, vague.IF_TOKEN, vague.FOR_TOKEN),
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrForVarInvalid,
		},
		{
			name: "invalid: for no key variable",
			input: fmt.Sprintf(`
				<div>
					<p %[1]s="true"></p>
					<h1 %[2]s=" ,value in values">Heading</h1>
				</div>`, vague.IF_TOKEN, vague.FOR_TOKEN),
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrForKeyVarInvalid,
		},
		{
			name: "invalid: for no loop variable in key, value",
			input: fmt.Sprintf(`
				<div>
					<p %[1]s="true"></p>
					<h1 %[2]s="key,  in values">Heading</h1>
				</div>`, vague.IF_TOKEN, vague.FOR_TOKEN),
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrForKeyVarInvalid,
		},
		{
			name: "valid: mix of correct for statements",
			input: fmt.Sprintf(`
				<ul>
					<li %[1]s="key,_ in rows">{{row.Name}}</li>
					<span %[1]s="value in rows">{{key}}</span>
					<p %[1]s="key,val in vals">{{val.Summary}}</p>
				</ul>`, vague.FOR_TOKEN),
			expectErr:       false,
			expectNilNode:   false,
			expectedErrCode: vague.ErrNone,
			expectedNode: Tag(vague.ELEMENT, "ul").Child(
				Tag(vague.ELEMENT, "li").For("key", "_", "rows").Text("{{row.Name}}").Close(),
			).Child(
				Tag(vague.ELEMENT, "span").For("", "value", "rows").Text("{{key}}").Close(),
			).Child(
				Tag(vague.ELEMENT, "p").For("key", "val", "vals").Text("{{val.Summary}}").Close(),
			).Close(),
		},
		{
			name: "valid: whitespace variation of correct for statements",
			input: fmt.Sprintf(`
				<ul>
					<li %[1]s=   "  key,  _    in   rows">{{row.Name}}</li>
					<span %[1]s="  value   in   rows"  >{{key}}</span>
					<p   %[1]s=  "   key  ,  val   in   vals   "  >{{val.Summary}}</p>
				</ul>`, vague.FOR_TOKEN),
			expectErr:       false,
			expectNilNode:   false,
			expectedErrCode: vague.ErrNone,
			expectedNode: Tag(vague.ELEMENT, "ul").Child(
				Tag(vague.ELEMENT, "li").For("key", "_", "rows").Text("{{row.Name}}").Close(),
			).Child(
				Tag(vague.ELEMENT, "span").For("", "value", "rows").Text("{{key}}").Close(),
			).Child(
				Tag(vague.ELEMENT, "p").For("key", "val", "vals").Text("{{val.Summary}}").Close(),
			).Close(),
		},
		{
			name:          "valid: commata in expression",
			input:         fmt.Sprintf(`<div %s="key,value in getItems(2,4,6) "></div>`, vague.FOR_TOKEN),
			expectErr:     false,
			expectNilNode: false,
			expectedNode: &vague.Node{
				Type:    vague.ELEMENT,
				TagName: "div",
				LoopInfo: &vague.LoopInfo{
					Var:    "value",
					KeyVar: "key",
					Expr:   "getItems(2,4,6)",
				},
			},
		},
	}

	for _, TestRow = range tests {
		t.Run(TestRow.name, testFunc)
	}
}
