package tests

import (
	"testing"
	"weld/vague"
)

func TestParseFor(t *testing.T) {
	tests := []testData{
		{
			name: "invalid: multiple for statements",
			input: `
			<div>
				<h1>Heading</h1>
				<h1 v-for="_,value in values" v-for="index in values" >Heading</h1>
			</div>`,
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrMultiFor,
		},
		{
			name: "invalid: for with v-if condition",
			input: `
			<div>
				<h1>Heading</h1>
				<h1 v-if="show" v-for="_,value in values" >Heading</h1>
			</div>`,
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrForWithCondition,
		},
		{
			name: "invalid: for with v-else-if condition",
			input: `
			<div>
				<h1 v-if="true">Heading</h1>
				<h1 v-else-if="show" v-for="_,value in values" >Heading</h1>
			</div>`,
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrForWithCondition,
		},
		{
			name: "invalid: for statement multiple 'in'",
			input: `
			<div>
				<p v-for="k in k in values"></p>
				<h1>Heading</h1>
				<h1>Heading</h1>
			</div>`,
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrForMultipleIn,
		},
		{
			name: "invalid: for kv statement misplaced ','",
			input: `
			<div>
				<p v-for="k,,v in values"></p>
				<h1>Heading</h1>
				<h1>Heading</h1>
			</div>`,
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrForExtraComma,
		},
		{
			name: "invalid: for without expression",
			input: `
			<div>
				<p v-if="true"></p>
				<h1 v-for>Heading</h1>
			</div>`,
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrForNoExpr,
		},
		{
			name: "invalid: for no loop variable",
			input: `
			<div>
				<p v-if="true"></p>
				<h1 v-for="  in values">Heading</h1>
			</div>`,
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrForVarInvalid,
		},
		{
			name: "invalid: for no key variable",
			input: `
			<div>
				<p v-if="true"></p>
				<h1 v-for=" ,value in values">Heading</h1>
			</div>`,
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrForKeyVarInvalid,
		},
		{
			name: "invalid: for no loop variable in key, value",
			input: `
			<div>
				<p v-if="true"></p>
				<h1 v-for="key,  in values">Heading</h1>
			</div>`,
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrForKeyVarInvalid,
		},
		{
			name: "valid: mix of correct for statements",
			input: `
			<ul>
				<li v-for="key,_ in rows">{{row.Name}}</li>
				<span v-for="value in rows">{{key}}</span>
				<p v-for="key,val in vals">{{val.Summary}}</p>
			</ul>`,
			expectErr:       false,
			expectNilNode:   false,
			expectedErrCode: vague.ErrNone,
			expectedNode: Tag(vague.ELEMENT, "ul").Child(
				Tag(vague.ELEMENT, "li").For("key", "_", "rows").Text("{{row.Name}}").Close()).Child(
				Tag(vague.ELEMENT, "span").For("", "value", "rows").Text("{{key}}").Close()).Child(
				Tag(vague.ELEMENT, "p").For("key", "val", "vals").Text("{{val.Summary}}").Close(),
			).Close(),
		}, {
			name: "valid: whitespace variation of correct for statements",
			input: `
			<ul>
				<li v-for=   "  key,  _    in   rows">{{row.Name}}</li>
				<span v-for="  value   in   rows"  >{{key}}</span>
				<p   v-for=  "   key  ,  val   in   vals   "  >{{val.Summary}}</p>
			</ul>`,
			expectErr:       false,
			expectNilNode:   false,
			expectedErrCode: vague.ErrNone,
			expectedNode: Tag(vague.ELEMENT, "ul").Child(
				Tag(vague.ELEMENT, "li").For("key", "_", "rows").Text("{{row.Name}}").Close()).Child(
				Tag(vague.ELEMENT, "span").For("", "value", "rows").Text("{{key}}").Close()).Child(
				Tag(vague.ELEMENT, "p").For("key", "val", "vals").Text("{{val.Summary}}").Close(),
			).Close(),
		},
		{
			name:          "valid: commata in expression",
			input:         `<div v-for="key,value in getItems(2,4,6) "></div>`,
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
