package nut

import (
	"fmt"
	"regexp"
	"strings"
)

type generator struct {
	tree *tree
	variables map[string]string
} 


func (generator *generator) generateBlockElement(tag, children, attributes string, indentationLength int) string {
	startBlockIndentation := strings.Repeat(" ", indentationLength)
	endBlockIndentation := ""
	if children != "" {
		endBlockIndentation = startBlockIndentation
	}

	if attributes != "" {
		attributes = " "+attributes
	}

	return fmt.Sprintf("\n%s<%s%s>%s%s</%s>\n", startBlockIndentation, tag, attributes, children,endBlockIndentation, tag)
}

func (generator *generator) generateAttributes(attributes []*attribute) (string, error) {
	var generatedAttributes string
	for _, attribute := range attributes {
		expr := regexp.MustCompile(`\{[a-zA-Z_]+([a-zA-Z0-9_]*)\}`)
		for _,variable := range expr.FindAllString(attribute.Value, -1) {
			variableName :=variable[1:len(variable) - 1]
			value, variableExists := generator.variables[variableName];
			if !variableExists {
				return "", fmt.Errorf("%s is not defined", variable[1:len(variable) - 1])
			} 
			attribute.Value = expr.ReplaceAllString(attribute.Value, value)

		}
		generatedAttributes = fmt.Sprintf("%s%s=\"%s\" ", generatedAttributes,attribute.Key, attribute.Value)
	}
	attributesLength := len(generatedAttributes)

	if attributesLength == 0 {
		return generatedAttributes, nil
	}

	return generatedAttributes[:attributesLength - 1], nil
}

func (generator *generator) generateVoidElement(tag, attributes string, indentationLength int) string {
	return fmt.Sprintf("\n%s<%s %s/>\n", strings.Repeat(" ", indentationLength), tag, attributes)
}

func (generator *generator) generateTextElement(value string, indentationLength int) string {
	return fmt.Sprintf("\n%s%s\n", strings.Repeat(" ", indentationLength), value)
}

func (generator *generator) generate(node *node_t, level int) (string, error) {
	if node == nil {
		return "", nil
	}

	var generatedElement string = strings.Repeat(" ", level) + "\n"
	var generatedAttributes string

	if len(node.attributes) > 0 {
		_generatedAttributes, err := generator.generateAttributes(node.attributes)
		if err != nil {
			return "", err
		}
		generatedAttributes = _generatedAttributes
	}

	if node.type_t == block {
		var childrenElement string
		if node.child != nil {
			_childrenElement, err := generator.generate(node.child, level+1)
			if err != nil {
				return "", err
			}
			childrenElement = _childrenElement
		}
		generatedElement = generator.generateBlockElement(node.tag, childrenElement, generatedAttributes, level * 4)
	}
	indentationLength := level * 4
	if node.tag == "text" && len(node.attributes) > 0 {
		generatedElement = generator.generateTextElement(node.attributes[0].Value, indentationLength)
	}

	if node.tag != "text" && node.type_t == void {
		generatedElement = generator.generateVoidElement(node.tag, generatedAttributes, indentationLength)
	}

	if node.sibling != nil {
		sibling, err := generator.generate(node.sibling, level)
		if err != nil {
			return "", err
		}
		generatedElement = fmt.Sprintf("%s%s", generatedElement[:len(generatedElement) - 1], sibling)
	}
	
	return generatedElement, nil
}

func generate(tree *tree, variables map[string]string) (string, error) {
	generator := &generator{tree, variables}
	return generator.generate(generator.tree.rootNode.child, 0)
}