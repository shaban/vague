package tests

import (
	"encoding/json"
	"log"
	"reflect"
	"strings"
	"testing"
	"weld/vague"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type testData struct {
	name            string
	input           string
	expectErr       bool
	expectNilNode   bool
	expectedErrCode int
	expectedNode    *vague.Node
}

/*
	func deserialize(jsString string, t *testing.T) *vague.Node {
		n := new(vague.Node)
		err := json.Unmarshal([]byte(jsString), &n)
		if err != nil {
			t.Fatal("Test file for: ", t.Name(), "corrupt")
		}
		return n
	}
*/
func SprintNodeTree(node *vague.Node) string {
	nodeJson, err := json.MarshalIndent(node, "", "\n")
	if err != nil {
		log.Fatal(err)
	}
	return string(nodeJson)
}
func compareNodes(actual, expected *vague.Node) (debugInfo []string) {
	debugInfo = append(debugInfo, SprintNodeTree(actual))
	debugInfo = append(debugInfo, SprintNodeTree(expected))
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(debugInfo[0], debugInfo[1], false)
	debugInfo = append(debugInfo, dmp.DiffPrettyText(diffs))
	return debugInfo
}
func testFunc(t *testing.T) {
	reader := strings.NewReader(TestRow.input)
	node, err := vague.ParseTemplate(reader)
	if (err != nil) != TestRow.expectErr {
		t.Errorf("ParseTemplate() error = %v, expectErr %v", err, TestRow.expectErr)
		return
	}

	if TestRow.expectNilNode != (node == nil) {
		t.Errorf("Expected node==nil to be %t, but got node==nil: %t", TestRow.expectNilNode, node == nil)
		return
	}
	if err != nil {
		if err.Code != TestRow.expectedErrCode {
			t.Errorf("Expected error code %d, but got %d", TestRow.expectedErrCode, err.Code)
		}
	}

	if TestRow.expectedNode != nil {
		assertNodeValues(t, node, TestRow.expectedNode)
	}
}
func assertConditionInfoEqual(t *testing.T, actual, expected *vague.ConditionInfo) {
	if actual == nil && expected == nil {
		return // Both nil, considered equal
	}
	if actual == nil || expected == nil {
		t.Error("One ConditionInfo is nil while the other is not")
		return
	}
	if actual.Condition != expected.Condition {
		t.Errorf("Expected Condition '%s', got '%s'", expected.Condition, actual.Condition)
		return
	}
	if actual.IsElse != expected.IsElse {
		t.Errorf("Expected IsElse to be %t, got %t", expected.IsElse, actual.IsElse)
		return
	}
	if actual.IsElseIf != expected.IsElseIf {
		t.Errorf("Expected IsElseIf to be %t, got %t", expected.IsElseIf, actual.IsElseIf)
		return
	}
	// add comparisons for other fields as needed
}
func assertAttributesEqual(t *testing.T, actual, expected map[string]*vague.Attribute) {
	if len(actual) != len(expected) {
		t.Errorf("Expected %d attributes, got %d", len(expected), len(actual))
		return
	}

	for key, expectedAttribute := range expected {
		actualAttribute, ok := actual[key]
		if !ok {
			t.Errorf("Missing attribute '%s'", key)
			return
		}

		if actualAttribute.Value != expectedAttribute.Value {
			t.Errorf("ATestRowribute '%s': Expected Value '%s', got '%s'", key, expectedAttribute.Value, actualAttribute.Value)
			return
		}

		if !reflect.DeepEqual(actualAttribute.Values, expectedAttribute.Values) {
			t.Errorf("ATestRowribute '%s': Expected Values %v, got %v", key, expectedAttribute.Values, actualAttribute.Values)
			return
		}
	}
}

func assertLoopInfoEqual(t *testing.T, actual, expected *vague.LoopInfo) {
	if actual == nil && expected == nil {
		return // Both nil, considered equal
	}
	if actual == nil || expected == nil {
		t.Error("One LoopInfo is nil while the other is not")
		return
	}
	if actual.Expr != expected.Expr {
		t.Errorf("Expected Expr '%s', got '%s'", expected.Expr, actual.Expr)
		return
	}
	if actual.Var != expected.Var {
		t.Errorf("Expected Var '%s', got '%s'", expected.Var, actual.Var)
		return
	}
	if actual.KeyVar != expected.KeyVar {
		t.Errorf("Expected KeyVar '%s', got '%s'", expected.KeyVar, actual.KeyVar)
		return
	}
}
func assertNodeValues(t *testing.T, actual *vague.Node, expected *vague.Node) {
	// Type
	if actual.Type != expected.Type {
		t.Errorf("Expected node type %d, got %d", expected.Type, actual.Type)
		debugAnalysis(actual, expected, t)
		return
	}

	// TagName
	if actual.TagName != expected.TagName {
		t.Errorf("Expected tag name '%s', got '%s'", expected.TagName, actual.TagName)
		debugAnalysis(actual, expected, t)
		return
	}
	//Info structs
	assertLoopInfoEqual(t, actual.LoopInfo, expected.LoopInfo)
	assertConditionInfoEqual(t, actual.ConditionInfo, expected.ConditionInfo)
	assertAttributesEqual(t, actual.Attributes, expected.Attributes)

	// Content (for text nodes)
	if actual.Type == vague.TEXT && actual.Content != expected.Content {
		t.Errorf("Expected text content '%s', got '%s'", expected.Content, actual.Content)
		debugAnalysis(actual, expected, t)
	}

	// Children (recursively check)
	if len(actual.Children) != len(expected.Children) {
		t.Errorf(`Expected %d children, got %d\n`, len(expected.Children), len(actual.Children))
		debugAnalysis(actual, expected, t)
	} else {
		for i, child := range actual.Children {
			assertNodeValues(t, child, expected.Children[i])
		}
	}
}
func debugAnalysis(actual, expected *vague.Node, t *testing.T) {
	debugInfo := compareNodes(actual, expected)
	t.Errorf(`actual: %s\n
		expected: %s\n
		diff: %s`, debugInfo[0], debugInfo[1], debugInfo[2])
}

type NodeBuilder struct {
	node *vague.Node
}

func Tag(nodeType int, tagName string) *NodeBuilder {
	return &NodeBuilder{
		node: &vague.Node{
			Type:     nodeType,
			TagName:  tagName,
			Children: make([]*vague.Node, 0), // Initialize Children slice
		},
	}
}

func (nb *NodeBuilder) Attr(key, value string) *NodeBuilder {
	if nb.node.Attributes == nil {
		nb.node.Attributes = make(map[string]*vague.Attribute)
	}
	nb.node.Attributes[key] = &vague.Attribute{Value: value}
	return nb
}

func (nb *NodeBuilder) Child(child *vague.Node) *NodeBuilder {
	nb.node.Children = append(nb.node.Children, child)
	return nb
}

func (nb *NodeBuilder) Text(content string) *NodeBuilder {
	tn := &vague.Node{
		Type:     vague.TEXT,
		Content:  content,
		Children: make([]*vague.Node, 0), // Initialize Children slice
	}
	nb.node.Children = append(nb.node.Children, tn)
	return nb
}

func (nb *NodeBuilder) If(condition string) *NodeBuilder {
	nb.node.ConditionInfo = &vague.ConditionInfo{Condition: condition}
	return nb
}
func (nb *NodeBuilder) Else() *NodeBuilder {
	nb.node.ConditionInfo = &vague.ConditionInfo{IsElse: true}
	return nb
}

func (nb *NodeBuilder) ElseIf(condition string) *NodeBuilder {
	nb.node.ConditionInfo = &vague.ConditionInfo{Condition: condition, IsElseIf: true}
	return nb
}
func (nb *NodeBuilder) For(keyVar, variable, expression string) *NodeBuilder {
	nb.node.LoopInfo = &vague.LoopInfo{KeyVar: keyVar, Var: variable, Expr: expression}
	return nb
}

func (nb *NodeBuilder) Close() *vague.Node {
	return nb.node
}
