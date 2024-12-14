package tests

import (
	"fmt"
	"testing"

	"github.com/shaban/vague"
)

func TestParseIf(t *testing.T) {
	tests := []testData{
		{
			name: "invalid: if no expression",
			input: fmt.Sprintf(`
				<div %s>
					<h1>Heading</h1>
				</div>`, vague.IF_TOKEN),
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrIfNoExpr,
		},
		{
			name: "invalid: multiple if statements",
			input: fmt.Sprintf(`
				<div>
					<h1 %[1]s="Show" %[1]s="Valid">Heading</h1>
				</div>`, vague.IF_TOKEN),
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrMultiIf,
		},
		{
			name: "invalid: composite else-if",
			input: fmt.Sprintf(`
				<div %[1]s="isValid(x,y int)">
					<h1 %[1]s="Show && Valid">Heading</h1>
					<h2 %[2]s %[1]s="Tabular">Heading</h2>
					<h2 %[2]s>Heading</h2>
				</div>`, vague.IF_TOKEN, vague.ELSE_TOKEN),
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrCompositeElseIf,
		},
		{
			name: "valid: mix of correct if statements",
			input: fmt.Sprintf(`
				<div %[1]s="isValid(x,y int)">
					<h1 %[1]s="Show && Valid">Heading</h1>
				</div>`, vague.IF_TOKEN),
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
