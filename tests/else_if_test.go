package tests

import (
	"fmt"
	"testing"

	"github.com/shaban/vague"
)

func TestParseElseIf(t *testing.T) {
	tests := []testData{
		{
			name: "invalid: else-if without if",
			input: fmt.Sprintf(`
				<div>
					<h1>Heading</h1>
					<h1 %s="show">Heading</h1>
				</div>`, vague.ELSE_IF_TOKEN),
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrElseNoIf,
		},
		{
			name: "invalid: else-if after else",
			input: fmt.Sprintf(`
				<div>
					<h1 %[1]s="true">Heading</h1>
					<h1 %[2]s>Heading</h1>
					<h1 %[3]s="show">Heading</h1>
				</div>`, vague.IF_TOKEN, vague.ELSE_TOKEN, vague.ELSE_IF_TOKEN),
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrElseIfAfterElse,
		},
		{
			name: "invalid: mixed else-if and else",
			input: fmt.Sprintf(`
				<div>
					<h1 %[1]s="true">Heading</h1>
					<h1 %[2]s %[3]s="show">Heading</h1>
				</div>`, vague.IF_TOKEN, vague.ELSE_TOKEN, vague.ELSE_IF_TOKEN),
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrMixedElseElseIf,
		},
		{
			name: "invalid: multiple else-if on element",
			input: fmt.Sprintf(`
				<div>
					<h1 %[1]s="true">Heading</h1>
					<h1 %[3]s="of_course" %[3]s="show">Heading</h1>
				</div>`, vague.IF_TOKEN, vague.ELSE_TOKEN, vague.ELSE_IF_TOKEN),
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrMultipleElseIf,
		},
		{
			name: "invalid: else-if no condition",
			input: fmt.Sprintf(`
				<div>
					<h1 %[1]s="true">Heading</h1>
					<h1 %[3]s>Heading</h1>
				</div>`, vague.IF_TOKEN, vague.ELSE_TOKEN, vague.ELSE_IF_TOKEN),
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrElseIfNoExpr,
		},
		{
			name: "valid: else-if after if",
			input: fmt.Sprintf(`
				<div>
					<p %[1]s="yes"></p>
					<h1 %[3]s="show">Heading</h1>
				</div>`, vague.IF_TOKEN, vague.ELSE_TOKEN, vague.ELSE_IF_TOKEN),
			expectErr:     false,
			expectNilNode: false,
			expectedNode: Tag(vague.ELEMENT, "div").Child(
				Tag(vague.ELEMENT, "p").If("yes").Close(),
			).Child(
				Tag(vague.ELEMENT, "h1").ElseIf("show").Text("Heading").Close(),
			).Close(),
		},
		{
			name: "valid: else-if after else-if",
			input: fmt.Sprintf(`
				<div>
					<p %[1]s="yes"></p>
					<h1 %[3]s="show">Heading</h1>
					<h2 %[3]s="isValid">Other Heading</h2>
				</div>`, vague.IF_TOKEN, vague.ELSE_TOKEN, vague.ELSE_IF_TOKEN),
			expectErr:     false,
			expectNilNode: false,
			expectedNode: Tag(vague.ELEMENT, "div").Child(
				Tag(vague.ELEMENT, "p").If("yes").Close(),
			).Child(
				Tag(vague.ELEMENT, "h1").ElseIf("show").Text("Heading").Close(),
			).Child(
				Tag(vague.ELEMENT, "h2").ElseIf("isValid").Text("Other Heading").Close(),
			).Close(),
		},
		{
			name: "valid: else after else-if with condition",
			input: fmt.Sprintf(`
				<div>
					<p %[1]s="yes"></p>
					<h1 %[3]s="show">Heading</h1>
					<h2 %[2]s>Sub Heading</h2>
				</div>`, vague.IF_TOKEN, vague.ELSE_TOKEN, vague.ELSE_IF_TOKEN),
			expectErr:     false,
			expectNilNode: false,
			expectedNode: Tag(vague.ELEMENT, "div").Child(
				Tag(vague.ELEMENT, "p").If("yes").Close(),
			).Child(
				Tag(vague.ELEMENT, "h1").ElseIf("show").Text("Heading").Close(),
			).Child(
				Tag(vague.ELEMENT, "h2").Else().Text("Sub Heading").Close(),
			).Close(),
		},
	}

	for _, TestRow = range tests {
		t.Run(TestRow.name, testFunc)
	}
}
