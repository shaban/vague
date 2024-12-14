package vague

import (
	"strings"

	"golang.org/x/net/html"
)

func parseHTMLAttribute(node *Node, attr html.Attribute) *parseError {
	// maybe handle unknown attribute
	if node.Type == VIRTUAL {
		return &parseError{
			Message: "%s nodes can't have attributes like [%s], just [%s, %s, %s, %s]",
			Extra:   []any{VIRTUAL_TOKEN, attr.Key, FOR_TOKEN, IF_TOKEN, ELSE_TOKEN, ELSE_IF_TOKEN},
			Code:    ErrVirtAttr,
		}
	}
	node.Attributes[attr.Key] = &Attribute{Value: attr.Val}
	return nil
}

func parseIf(currentNode *Node, attr html.Attribute) (err *parseError) {
	// avoid overwriting previous condition attributes
	if currentNode.ConditionInfo == nil {
		currentNode.ConditionInfo = new(ConditionInfo)
	} else {
		if currentNode.ConditionInfo.IsElse {
			//we got a composite v-else-if
			return compositeElseIfMessage
		}
	}
	if currentNode.ConditionInfo.Condition != "" {
		return multipleIfMessage
	}
	currentNode.ConditionInfo.Condition = strings.TrimSpace(attr.Val)
	if currentNode.ConditionInfo.Condition == "" {
		return ifNoExpressionMessage
	}
	if currentNode.LoopInfo != nil {
		//we have v-if and v-for on same element
		return forWithConditionMessage
	}
	return nil
}
func parseFor(currentNode *Node, attr html.Attribute) *parseError {
	var (
		keyVar   string
		loopVar  string
		iterable string
	)

	if currentNode.LoopInfo != nil {
		return multipleForStatementsMessage
	}

	expr := strings.TrimSpace(attr.Val)
	if expr == "" {
		return forNoExpressionMessage
	}

	if currentNode.ConditionInfo != nil && currentNode.ConditionInfo.Condition != "" {
		return forWithConditionMessage
	}
	// Check for multiple "in" in the expression
	if strings.Count(expr, " in ") > 1 {
		return forMultipleInMessage // New error message for multiple "in"
	}
	// Extract iterable using strings.Cut
	expr, iterable, ok := strings.Cut(expr, " in ")
	if !ok {
		//since the value was absent the ' in ' was trimmed to 'in '
		if strings.Contains(expr, "in ") {
			return forVarInvalidMessage
		}
		return forNoInMessage
	}
	iterable = strings.TrimSpace(iterable)
	if iterable == "" {
		return forNoIterableMessage
	}

	// Handle key, value or just value
	switch strings.Count(expr, ",") {
	case 0:
		break
	case 1:
		break
	default:
		return forExtraCommaMessage
	}
	if keyVal, val, ok := strings.Cut(expr, ","); ok { // key, value case
		keyVar = strings.TrimSpace(keyVal)
		loopVar = strings.TrimSpace(val)
		if keyVar == "" || loopVar == "" {
			return forKeyVarInvalidMessage
		}
	} else { // Just value case
		loopVar = strings.TrimSpace(expr)
		if loopVar == "" {
			return forVarInvalidMessage
		}
	}

	currentNode.LoopInfo = &LoopInfo{Var: loopVar, KeyVar: keyVar, Expr: iterable}
	return nil
}
func parseElse(currentNode, parentNode, previousSibling *Node, attr html.Attribute) *parseError {
	// if there is no parentNode we must be in the rootNode and that can't have a preceeding if
	if parentNode == nil {
		return elseWithoutIfMessage
	}
	// if we are the first childNode there can't be a preceeding if
	if previousSibling == nil {
		return elseWithoutIfMessage
	}
	// if the previous sibling doesn't have any conditions then there can't be an if
	if previousSibling.ConditionInfo == nil {
		return elseWithoutIfMessage
	}
	// if there is a condition but no expression then previous sibling is an else too
	if previousSibling.ConditionInfo.Condition == "" && previousSibling.ConditionInfo.IsElse {
		return elseAfterElseMessage
	}
	// for the sake of completeness we have that but it will probably never get reached
	if previousSibling.ConditionInfo.IsElse && previousSibling.ConditionInfo.Condition == "" {
		return elseAfterElseMessage
	}
	if attr.Val != "" {
		return elseWithConditionMessage(attr.Val)
	}
	if currentNode.ConditionInfo != nil {
		if currentNode.ConditionInfo.IsElse {
			return multipleElseMessage
		}
		if currentNode.ConditionInfo.Condition != "" {
			//we got a composite v-else-if show an error
			return compositeElseIfMessage
		}

	}
	// only create a new condition struct when it is nil to avoid overwriting other condition information
	if currentNode.ConditionInfo == nil {
		currentNode.ConditionInfo = new(ConditionInfo)
	}

	// else looks good
	currentNode.ConditionInfo.IsElse = true
	return nil
}

func parseElseIf(currentNode, parentNode, previousSibling *Node, attr html.Attribute) *parseError {
	// if there is no parentNode we must be in the rootNode and that can't have a preceeding if
	if parentNode == nil {
		return elseWithoutIfMessage
	}
	// if we are the first childNode there can't be a preceeding if
	if previousSibling == nil {
		return elseWithoutIfMessage
	}
	// if the previous sibling doesn't have any conditions then there can't be an if
	if previousSibling.ConditionInfo == nil {
		return elseWithoutIfMessage
	}
	// if the previous sibling IsElse we can't have an if or elseIf
	if previousSibling.ConditionInfo.IsElse {
		return elseIfAfterElseMessage
	}
	if currentNode.ConditionInfo != nil {
		if currentNode.ConditionInfo.IsElse {
			return mixedElseElseIfMessage
		}
		if currentNode.ConditionInfo.IsElseIf {
			return multipleElseIfMessage
		}
	}
	// only create a new condition struct when it is nil to avoid overwriting other condition information
	if currentNode.ConditionInfo == nil {
		currentNode.ConditionInfo = new(ConditionInfo)
	}
	currentNode.ConditionInfo.Condition = strings.TrimSpace(attr.Val)
	if currentNode.ConditionInfo.Condition == "" {
		return elseIfNoExpressionMessage
	}
	// elseif looks good
	currentNode.ConditionInfo.IsElseIf = true
	return nil
}

func parseDirectives(currentNode *Node, parentNode *Node, attr html.Attribute) (err *parseError) {
	siblings := currentNode.getSiblings(parentNode)
	previousSibling := currentNode.previousSibling(siblings)
	attr.Key = strings.TrimSpace(attr.Key)
	attr.Val = strings.TrimSpace(attr.Val)
	switch attr.Key {
	case IF_TOKEN:
		return parseIf(currentNode, attr)
	case ELSE_IF_TOKEN:
		return parseElseIf(currentNode, parentNode, previousSibling, attr)
	case ELSE_TOKEN:
		return parseElse(currentNode, parentNode, previousSibling, attr)
	case FOR_TOKEN:
		return parseFor(currentNode, attr)
	default:
		return parseHTMLAttribute(currentNode, attr)
	}
}
