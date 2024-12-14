package vague

import (
	"io"
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"
)

var (
	fileName string
	pos      = new(position)
)

func ParseTemplateFile(fn string) (*Node, *parseError) {
	f, err := os.Open(fn)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer f.Close()
	fileName = fn
	return ParseTemplate(f)
}
func ParseTemplate(f io.Reader) (*Node, *parseError) {
	var (
		tokenizer   *html.Tokenizer
		rootNode    *Node
		currentNode *Node
		nodeStack   []*Node
		err         *parseError
		parentNode  *Node
	)

	tokenizer = html.NewTokenizer(f)
	pos.line = 1
	for {
		tokenType := tokenizer.Next()

		//handle lineNumbers
		for _, b := range tokenizer.Raw() {
			if b == '\n' {
				pos.line++
				pos.offset = 0
			} else {
				pos.offset++ // Increment column number for each character
			}
		}
		//start parsing
		switch tokenType {
		case html.ErrorToken:
			//if there is an error
			goErr := tokenizer.Err()
			// and its the end of the file
			if goErr == io.EOF {
				// Check if the root node is closed properly
				if len(nodeStack) > 0 {
					return nil, unclosedRootTag(rootNode.TagName)
				}
				// return the rootNode that holds the rest of the parsed stuff
				return rootNode, nil // End of file
			}
			// if not return an error
			return nil, err // Parsing error

		case html.StartTagToken:
			// get the opening tags name
			// and create a new Node
			token := tokenizer.Token()
			newNode := new(Node)

			// is RootNode
			if currentNode == nil {
				//we don't have a single root throw an ErrMultiRoot
				if rootNode != nil {
					//move to start of tag
					pos.col = pos.offset - len(tokenizer.Raw()) + 1
					return nil, multipleRoots()
				}
				rootNode = newNode // First node is the root
			} else {
				currentNode.Children = append(currentNode.Children, newNode)
			}
			parentNode = currentNode
			// Add to tree and stack
			nodeStack = append(nodeStack, newNode)
			// from here currentNode is the way to go
			currentNode = newNode

			if token.Data == VIRTUAL_TOKEN {
				currentNode.Type = VIRTUAL
			} else {
				currentNode.Type = ELEMENT
			}
			currentNode.TagName = token.Data
			currentNode.Attributes = make(map[string]*Attribute, 0)
			//var condInfo *ConditionInfo

			for _, attr := range token.Attr {
				pos.col = strings.Index(string(tokenizer.Raw()), attr.Key) // Adjust for whitespace or quotes if needed
				pos.col = pos.offset - len(tokenizer.Raw()) + pos.col + 1
				err = parseDirectives(currentNode, parentNode, attr)
				if err != nil {
					return nil, err
				}
			}

		case html.TextToken:
			text := strings.TrimSpace(tokenizer.Token().Data)
			if text == "" {
				break
			}
			newNode := &Node{Type: TEXT, Content: text}
			currentNode.Children = append(currentNode.Children, newNode)

		case html.EndTagToken:
			// check if there is any nodes left to end
			// if not we must have a stray end tag and this is malformed html
			endTagName := strings.TrimSpace(tokenizer.Token().Data)
			if len(nodeStack) == 0 {
				return nil, strayEndTags(endTagName)
			}
			// Get the last opened tag from the node stack
			openingNode := nodeStack[len(nodeStack)-1]
			// Pop from stack
			nodeStack = nodeStack[:len(nodeStack)-1]

			if openingNode.TagName != endTagName {
				return nil, mismatchedTags(openingNode.TagName, endTagName)
			}

			if len(nodeStack) > 0 {
				currentNode = nodeStack[len(nodeStack)-1]
			} else {
				parentNode = nil
				currentNode = nil // Back to root level
			}

			// self closing Tags
		case html.SelfClosingTagToken:
			token := tokenizer.Token()
			newNode := &Node{Type: ELEMENT, TagName: token.Data}
			newNode.Attributes = make(map[string]*Attribute, 0)

			//parse non directive html attributes
			for _, attr := range token.Attr {
				err = parseHTMLAttribute(newNode, attr)
				if err != nil {
					return nil, err
				}
			}
			currentNode.Children = append(currentNode.Children, newNode) // Add to parent
		}
	}
}
