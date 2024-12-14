package tests

import (
	"testing"
	"weld/vague"
)

/*
expectedNode: Tag(vague.ELEMENT,"div").If("isValid(x,y int)").Child(

	Tag(vague.ELEMENT,"h1").If("Show && Valid").Text("Heading").Close().Child(

Tag(vague.ELEMENT,"h2").If("Tabular")

		)
	),
*/
func TestParseElse(t *testing.T) {
	tests := []testData{
		{
			name: "invalid: else without if",
			input: `
			<div>
				<h1>Heading</h1>
				<h1 v-else>Heading</h1>
			</div>`,
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrElseNoIf,
		},
		{
			name: "invalid: else after else without if",
			input: `
			<div>
				<p v-if="yes"></p>
				<h1 v-else>Heading</h1>
				<h1 v-else>Heading</h1>
			</div>`,
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrElseAfterElse,
		},
		{
			name: "invalid: multiple else statements",
			input: `
			<div>
				<p v-if="true"></p>
				<h1 v-else v-else>Heading</h1>
			</div>`,
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrMultiElse,
		},
		{
			name: "invalid: composite else-if",
			input: `
			<div v-if="isValid(x,y int)">
				<h1 v-if="Show && Valid">Heading</h1>
				<h2 v-if="Tabular" v-else>Heading</h2>
				<h2 v-else>Heading</h2>
			</div>`,
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrCompositeElseIf,
		},
		{
			name: "invalid: v-else with condition",
			input: `
			<div>
				<div v-if="show">Show</div>
				<div v-else="hide">Hide</div>
			</div>
	`,
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrElseCondition,
		},
		//need one sensible valid test
		//also must check for else with condition
	}
	for _, TestRow = range tests {
		t.Run(TestRow.name, testFunc)
	}

}
