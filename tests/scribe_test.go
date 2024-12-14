package tests

import (
	"testing"
)

var (
	TestRow testData
)

func TestParseConditionals(t *testing.T) {
	TestParseIf(t)
	TestParseElseIf(t)
	TestParseElse(t)
}

func TestParseTemplate(t *testing.T) {
	TestParseConditionals(t)
	TestParseFor(t)
	TestParseWellFormed(t)
	TestParseConditionalsEdgeCases(t)
}
