package vague

// getSiblings retrieves the child Nodes including this Node from
// the parentNode
// if there is no parentNode it return nil
// same if the parentNode has no children because the currentNode has not yet been added
// and also when the children nodes are an empty array
func (n *Node) getSiblings(parentNode *Node) []*Node {
	if parentNode == nil {
		return nil
	}
	siblings := parentNode.Children
	if siblings == nil {
		return nil
	}
	if len(siblings) == 0 {
		return nil
	}
	return siblings
}
func (n *Node) previousSibling(siblings []*Node) *Node {
	for index, sibling := range siblings {
		if sibling == n && index > 0 {
			return siblings[index-1]
		}
	}
	return nil
}
