package tests

import (
	"testing"
	"weld/vague"
)

func TestParseIf(t *testing.T) {
	tests := []testData{
		{
			name: "invalid: if no expression",
			input: `
			<div v-if>
				<h1>Heading</h1>
			</div>`,
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrIfNoExpr,
		},
		{
			name: "invalid: multiple if statements",
			input: `
			<div>
				<h1 v-if="Show" v-if="Valid">Heading</h1>
			</div>`,
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrMultiIf,
		},
		{
			name: "invalid: composite else-if",
			input: `
			<div v-if="isValid(x,y int)">
				<h1 v-if="Show && Valid">Heading</h1>
				<h2 v-else v-if="Tabular">Heading</h2>
				<h2 v-else>Heading</h2>
			</div>`,
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrCompositeElseIf,
		},
		{
			name: "valid: mix of correct if statements",
			input: `
			<div v-if="isValid(x,y int)">
				<h1 v-if="Show && Valid">Heading</h1>
			</div>`,
			expectErr:     false,
			expectNilNode: false,
			expectedNode: Tag(vague.ELEMENT, "div").If("isValid(x,y int)").Child(
				Tag(vague.ELEMENT, "h1").If("Show && Valid").Text("Heading").Close(),
			).Close(),
		},
	}
	for _, TestRow = range tests {
		t.Run(TestRow.name, testFunc)
	}

}
