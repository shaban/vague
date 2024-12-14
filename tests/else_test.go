package tests

import (
	"fmt"
	"testing"

	"github.com/shaban/vague"
)

func TestParseElse(t *testing.T) {
	tests := []testData{
		{
			name: "invalid: else without if",
			input: fmt.Sprintf(`
				<div>
					<h1>Heading</h1>
					<h1 %s>Heading</h1>
				</div>`, vague.ELSE_TOKEN),
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrElseNoIf,
		},
		{
			name: "invalid: else after else without if",
			input: fmt.Sprintf(`
				<div>
					<p %[1]s="yes"></p>
					<h1 %[2]s>Heading</h1>
					<h1 %[2]s>Heading</h1>
				</div>`, vague.IF_TOKEN, vague.ELSE_TOKEN),
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrElseAfterElse,
		},
		{
			name: "invalid: multiple else statements",
			input: fmt.Sprintf(`
				<div>
					<p %[1]s="true"></p>
					<h1 %[2]s %[2]s>Heading</h1>
				</div>`, vague.IF_TOKEN, vague.ELSE_TOKEN),
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrMultiElse,
		},
		{
			name: "invalid: composite else-if",
			input: fmt.Sprintf(`
				<div %[1]s="isValid(x,y int)">
					<h1 %[1]s="Show && Valid">Heading</h1>
					<h2 %[1]s="Tabular" %[2]s>Heading</h2>
					<h2 %[2]s>Heading</h2>
				</div>`, vague.IF_TOKEN, vague.ELSE_TOKEN),
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrCompositeElseIf,
		},
		{
			name: "invalid: v-else with condition",
			input: fmt.Sprintf(`
				<div>
					<div %[1]s="show">Show</div>
					<div %[2]s="hide">Hide</div>
				</div>
			`, vague.IF_TOKEN, vague.ELSE_TOKEN),
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrElseCondition,
		},
		{
			name: "valid: else after if", // Added a valid test case
			input: fmt.Sprintf(`
				<div>
					<div %[1]s="show">Show</div>
					<div %[2]s>Hide</div>
				</div>
			`, vague.IF_TOKEN, vague.ELSE_TOKEN),
			expectErr:     false,
			expectNilNode: false,
			expectedNode: Tag(vague.ELEMENT, "div").Child(
				Tag(vague.ELEMENT, "div").If("show").Text("Show").Close(),
			).Child(
				Tag(vague.ELEMENT, "div").Else().Text("Hide").Close(),
			).Close(),
		},
	}
	for _, TestRow = range tests {
		t.Run(TestRow.name, testFunc)
	}
}
