package nut

import (
	"fmt"
	"strings"
)

type stackItem struct {
	indentation int
	node_t *node_t
}

type parser struct {  
	cursor *cursor
	tokens []token 
	line uint
	stack *stack[stackItem]
	tree
}

func (parser *parser) expect(kind int, tokenType string) (*token,error) {
	currenttoken := parser.peek()
	if currenttoken == nil { 
		return nil, fmt.Errorf("invalid or unexpected token")
	}
	
	if currenttoken.kind != kind {
		return nil, fmt.Errorf("expected %s, got %s on line %v", tokenType, currenttoken.value, parser.line)
	}

	if kind == newLine {
		parser.line++
	}

	parser.cursor.next()
	return currenttoken, nil
}
 
func (parser *parser) peek() *token {
	if parser.isEndOftokens() {
		return nil
	}
	token :=  &parser.tokens[parser.cursor.getValue()]
	return token
}

func (parser *parser) check(kind int) bool {
	currenttoken := parser.peek()
	if(currenttoken == nil || currenttoken.kind != kind) {
		return false
	}
	return true
}

func (parser *parser) isEndOftokens() bool {
	return parser.cursor.getValue() >= len(parser.tokens)
}

func (parser *parser) parseAttribute() (*struct{attribute string;value string}, error) {
	attribute, err := parser.expect(node, "Attribute")
	if err != nil {
		return nil, err
	}
	parser.expect(equals, "Equals")
	value, err := parser.expect(qString, "Value")
	if err != nil {
		return nil, err
	}
	return &struct{attribute string; value string}{attribute: attribute.value, value: value.value}, nil
}

func (parser *parser) parseAttributes() ([]*attribute, error) {
	var attributes []*attribute
	for !parser.isEndOftokens() {
		attr, err := parser.parseAttribute()
		if err != nil {
			return nil,err
		}
		attributes = append(attributes, &attribute{Key: attr.attribute, Value: attr.value})
		if(parser.check(rParen)) {
			break
		}
		
		parser.expect(comma, "Comma")
		
	}

	return attributes, nil
} 

func (parser *parser) ignoreNewLine() {
	for !parser.isEndOftokens() {
		token := parser.peek()
		if token.kind != newLine {
			break
		}
		parser.line++
		parser.cursor.next()
	}
}

func (parser *parser) validateTextNode(tag string, attributes []*attribute) error {
	if len(attributes) > 1 || attributes[0].Key != "value" {
		return fmt.Errorf("%s is a custom element and should only contain the value attribute on line %v", tag, parser.line)
	}
	return nil
}

func (parser *parser) parseNode() (int, *node_t, error) {
	var indentation int
	var attributes []*attribute
	if parser.check(indent) {
		token := parser.peek()
		indentation = len(token.value)
		parser.cursor.next()	
	} 
	if parser.check(variable) {
		token, err := parser.expect(variable, "Variable")
		if err != nil {
			return 0, nil, err
		}
		parser.line += uint(strings.Count(token.value, "\n"))
		return indentation, &node_t{tag: "text", type_t:block, attributes: append(attributes, &attribute{Key: "value", Value: token.value}), child: nil, sibling: nil}, nil
	}

	token, err := parser.expect(node, "Element")

	if err != nil {
		return 0, nil, err
	}

	if !parser.check(newLine) {
		parser.expect(lParen, "Left Parenthesis")
		if parser.check(node) {
			attributes, err = parser.parseAttributes()

			if err != nil {
				return 0, nil, err
			}

			if len(attributes) > 0 && token.value == "text" {
				parser.validateTextNode(token.value, attributes)
			}
			
		}
		parser.expect(rParen, "Right Parenthesis")
	}

	node_t := &node_t{tag: token.value, type_t: block, attributes: attributes, child: nil, sibling: nil}
	elementType, isAnHTMLElement := elements[node_t.tag]
	if !isAnHTMLElement  {
		return 0, nil, fmt.Errorf("%s is not an HTML element on line %v", node_t.tag, parser.line)
	}
	node_t.type_t = elementType
	return indentation, node_t, nil
}

func (parser *parser) parseNodes() (*node_t, error) {
	var root *node_t
	for !parser.isEndOftokens() {	
		indentation, node_t, err := parser.parseNode()
		if err != nil {
			return nil, err
		}
		if !parser.isEndOftokens() {
			parser.expect(newLine, "New Line")
			parser.ignoreNewLine()
		}
		if parser.stack.empty() {
			root = node_t
		}
		
		for !parser.stack.empty() && parser.stack.top().indentation > indentation {
			parser.stack.pop()
		}
		if !parser.stack.empty() && indentation > parser.stack.top().indentation {
			parent := parser.stack.top()

			if  _, ok := metaElements[node_t.tag]; ok && parent.node_t.tag != "head" {
				return nil, fmt.Errorf("%s should be a child of <head> on line %v", parent.node_t.tag, parser.line)
			}

			if parent.node_t.tag == "text" || parent.node_t.type_t == void {
				return nil, fmt.Errorf("%s is a void element and cannot have children on line %v", parent.node_t.tag, parser.line)
			}

			parser.addChild(parent.node_t, node_t)
		}
		if parser.stack.size() > 0 && indentation == parser.stack.top().indentation {	
			sibling := parser.stack.top()
			parser.addSibling(sibling.node_t,node_t)
		}
		parser.stack.push(stackItem{indentation, node_t})
	}
	return root, nil
}

func (parser *parser) addChild(parent *node_t, node *node_t) {
	if parent.child == nil {
		parent.child = node
		return
	}
	parser.addSibling(parent.child, node)	
}

func (parser *parser) addSibling(node *node_t, sibling *node_t) {	
	if node.sibling == nil {
		node.sibling = sibling
		return
	}
	parser.addSibling(node.sibling, sibling)
}

func (parser *parser) parse() error {
	parser.ignoreNewLine()
	nodes, err := parser.parseNodes()
	if err != nil {
		return err
	}
	parser.tree.rootNode.child = nodes
	return nil
}

func newParser(tokens []token) *parser {
	return &parser{newCursor(0), tokens,1, newStack[stackItem](), tree{rootNode: &node_t{tag: "ROOT"}}}
} 