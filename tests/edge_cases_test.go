package tests

import (
	"fmt"
	"testing"

	"github.com/shaban/vague"
)

func TestParseConditionalsEdgeCases(t *testing.T) {
	tests := []testData{
		{
			name: "valid: deeply nested conditionals",
			input: fmt.Sprintf(`
				<div %[1]s="outer">
					<p %[1]s="inner1">Inner 1</p>
					<div %[3]s="inner2">
						<span %[1]s="innermost">Innermost</span>
						<span %[2]s>Not Innermost</span>
					</div>
					<p %[2]s>Inner 2</p>
				</div>
			`, vague.IF_TOKEN, vague.ELSE_TOKEN, vague.ELSE_IF_TOKEN),
			expectErr:     false,
			expectNilNode: false,
			expectedNode: Tag(vague.ELEMENT, "div").If("outer").Child(
				Tag(vague.ELEMENT, "p").If("inner1").Text("Inner 1").Close(),
			).Child(
				Tag(vague.ELEMENT, "div").ElseIf("inner2").Child(
					Tag(vague.ELEMENT, "span").If("innermost").Text("Innermost").Close(),
				).Child(
					Tag(vague.ELEMENT, "span").Else().Text("Not Innermost").Close(),
				).Close(),
			).Child(
				Tag(vague.ELEMENT, "p").Else().Text("Inner 2").Close(),
			).Close(),
		},
		{
			name: "invalid: attributes on virtual tag",
			input: fmt.Sprintf(`
				<div %[1]s="outer">
					<%[4]s disabled %[1]s="inner1">Inner 1</%[4]s>
					<div %[3]s="inner2">
						<span %[1]s="innermost">Innermost</span>
						<span %[2]s>Not Innermost</span>
					</div>
					<p %[2]s>Inner 2</p>
				</div>
			`, vague.IF_TOKEN, vague.ELSE_TOKEN, vague.ELSE_IF_TOKEN, vague.VIRTUAL_TOKEN),
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrVirtAttr,
		},
		{
			name: "valid: whitespace variations",
			input: fmt.Sprintf(`
				<div>
					<p %[1]s="show">Some Text</p>
					<p %[3]s  = "  isValid  " >Valid</p>
				</div>
			`, vague.IF_TOKEN, vague.ELSE_TOKEN, vague.ELSE_IF_TOKEN),
			expectErr:     false,
			expectNilNode: false,
			expectedNode: Tag(vague.ELEMENT, "div").Child(
				Tag(vague.ELEMENT, "p").If("show").Text("Some Text").Close(),
			).Child(
				Tag(vague.ELEMENT, "p").ElseIf("isValid").Text("Valid").Close(),
			).Close(),
		},
	}
	for _, TestRow = range tests {
		t.Run(TestRow.name, testFunc)
	}
}
