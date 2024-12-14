package tests

import (
	"fmt"
	"testing"

	"github.com/shaban/vague"
)

func TestParseWellFormed(t *testing.T) {
	tests := []testData{
		{
			name: "invalid: multiple root elements",
			input: `
				<div>
					<h1>Heading</h1>
					<h1>Heading</h1>
				</div>
				<h1>Heading</h1>
			`,
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrMultiRoot,
		},
		{
			name: "invalid: stray end tags",
			input: `
				<div>
					<h1>Heading</h1>
					<h1>Heading</h1>
				</div>
				</h1>
			`,
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrStrayTags,
		},
		{
			name: "valid: no multiroot",
			input: fmt.Sprintf(`
				<div>
					<p %s="yes"></p>
					<h1 %s>Heading</h1>
					<h1>Heading</h1>
				</div>`, vague.IF_TOKEN, vague.ELSE_TOKEN),
			expectErr:     false,
			expectNilNode: false,
			expectedNode: Tag(vague.ELEMENT, "div").Child(
				Tag(vague.ELEMENT, "p").If("yes").Close(),
			).Child(
				Tag(vague.ELEMENT, "h1").Else().Text("Heading").Close(),
			).Child(
				Tag(vague.ELEMENT, "h1").Text("Heading").Close(),
			).Close(), // Correctly close the div
		},
		{
			name: "invalid: unclosed Root Tag",
			input: `
				<div>
					<h1>Heading</h1>
					<h1>Heading</h1>
					<div>Hello
				</div>
			`,
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrUnclosedRoot,
		},
		{
			name: "invalid: tags don't match",
			input: fmt.Sprintf(`
				<div>
					<p %s="yes"></pt>
					<h1 %s>Heading</h1>
					<h1>Heading</h1>
				</div>`, vague.IF_TOKEN, vague.ELSE_TOKEN),
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrMismatchedTag,
		},
		{
			name: "valid: mix of tags that should match",
			input: fmt.Sprintf(`
				<div>
					<p %s="yes"></p>
					<h1 %s>Heading</h1>
					<span></span>
					<ul>
						<li></li>
					</ul>
					<h1>Heading</h1>
				</div>`, vague.IF_TOKEN, vague.ELSE_TOKEN),
			expectErr:     false,
			expectNilNode: false,
			expectedNode: Tag(vague.ELEMENT, "div").Child(
				Tag(vague.ELEMENT, "p").If("yes").Close(),
			).Child(
				Tag(vague.ELEMENT, "h1").Else().Text("Heading").Close(),
			).Child(
				Tag(vague.ELEMENT, "span").Close(),
			).Child(
				Tag(vague.ELEMENT, "ul").Child(
					Tag(vague.ELEMENT, "li").Close(),
				).Close(),
			).Child(
				Tag(vague.ELEMENT, "h1").Text("Heading").Close(),
			).Close(),
		},
	}
	for _, TestRow = range tests {
		t.Run(TestRow.name, testFunc)
	}
}
