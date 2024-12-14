package tests

import (
	"testing"

	"github.com/shaban/vague"
)

func TestParseElseIf(t *testing.T) {
	tests := []testData{
		{
			name: "invalid: else-if without if",
			input: `
                <div>
                    <h1>Heading</h1>
                    <h1 v-else-if="show">Heading</h1>
                </div>`,
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrElseNoIf,
		},
		{
			name: "invalid: else-if after else",
			input: `
                <div>
                    <h1 v-if="true">Heading</h1>
					<h1 v-else>Heading</h1>
                    <h1 v-else-if="show">Heading</h1>
                </div>`,
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrElseIfAfterElse,
		},
		{
			name: "invalid: mixed else-if and else",
			input: `
                <div>
                    <h1 v-if="true">Heading</h1>
                    <h1 v-else v-else-if="show">Heading</h1>
                </div>`,
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrMixedElseElseIf,
		},
		{
			name: "invalid: mutliple else-if on element",
			input: `
                <div>
                    <h1 v-if="true">Heading</h1>
                    <h1 v-else-if="of_course" v-else-if="show">Heading</h1>
                </div>`,
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrMultipleElseIf,
		},
		{
			name: "invalid: else-if no condition",
			input: `
                <div>
                    <h1 v-if="true">Heading</h1>
                    <h1 v-else-if=>Heading</h1>
                </div>`,
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrElseIfNoExpr,
		},
		// test for mixed else and else-if
		{
			name: "valid: else-if after if",
			input: `
                <div>
                    <p v-if="yes"></p>
                    <h1 v-else-if="show">Heading</h1>
                </div>`,
			expectErr:     false,
			expectNilNode: false,
			expectedNode: Tag(vague.ELEMENT, "div").Child(
				Tag(vague.ELEMENT, "p").If("yes").Close()).Child(
				Tag(vague.ELEMENT, "h1").ElseIf("show").Text("Heading").Close(),
			).Close(),
		},
		{
			name: "valid: else-if after else-if",
			input: `
                <div>
                    <p v-if="yes"></p>
                    <h1 v-else-if="show">Heading</h1>
                    <h2 v-else-if="isValid">Other Heading</h2>
                </div>`,
			expectErr:     false,
			expectNilNode: false,
			expectedNode: Tag(vague.ELEMENT, "div").Child(
				Tag(vague.ELEMENT, "p").If("yes").Close()).Child(
				Tag(vague.ELEMENT, "h1").ElseIf("show").Text("Heading").Close()).Child(
				Tag(vague.ELEMENT, "h2").ElseIf("isValid").Text("Other Heading").Close(),
			).Close(),
		},
		{
			name: "valid: else after else-if with condition",
			input: `
                <div>
                    <p v-if="yes"></p>
                    <h1 v-else-if="show">Heading</h1>
                    <h2 v-else>Sub Heading</h2>
                </div>`,
			expectErr:     false,
			expectNilNode: false,
			expectedNode: Tag(vague.ELEMENT, "div").Child(
				Tag(vague.ELEMENT, "p").If("yes").Close()).Child(
				Tag(vague.ELEMENT, "h1").ElseIf("show").Text("Heading").Close()).Child(
				Tag(vague.ELEMENT, "h2").Else().Text("Sub Heading").Close(),
			).Close(),
		},
	}
	for _, TestRow = range tests {
		t.Run(TestRow.name, testFunc)
	}
}
