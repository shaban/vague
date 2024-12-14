package tests

import (
	"testing"

	"github.com/shaban/vague"
)

func TestParseConditionalsEdgeCases(t *testing.T) {
	tests := []testData{
		// Deeply nested conditionals
		// add expected Node
		{
			name: "valid: deeply nested conditionals",
			input: `
		        <div v-if="outer">
		            <p v-if="inner1">Inner 1</p>
		            <div v-else-if="inner2">
		                <span v-if="innermost">Innermost</span>
		                <span v-else>Not Innermost</span>
		            </div>
		            <p v-else>Inner 2</p>
		        </div>
		    `,
			expectErr:     false,
			expectNilNode: false,
		},
		{
			name: "invalid: attributes on virtual tag",
			input: `
		        <div v-if="outer">
		            <template disabled v-if="inner1">Inner 1</template>
		            <div v-else-if="inner2">
		                <span v-if="innermost">Innermost</span>
		                <span v-else>Not Innermost</span>
		            </div>
		            <p v-else>Inner 2</p>
		        </div>
		    `,
			expectErr:       true,
			expectNilNode:   true,
			expectedErrCode: vague.ErrVirtAttr,
		},
		// Whitespace variations
		{
			name: "valid: whitespace variations",
			input: `
			<div>
        <p v-if="show">Some Text</p>
        <p v-else-if  = "  isValid  " >Valid</p>
		</div>
    `,
			expectErr: false,
		},
	}
	for _, TestRow = range tests {
		t.Run(TestRow.name, testFunc)
	}
}
