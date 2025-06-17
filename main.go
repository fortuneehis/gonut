package nut

func Run(input []byte, variables map[string]string) (string, error) {
	lexer := newlexer(input)
	lexer.scan()
	parser := newParser(lexer.tokens)
	
	if err := parser.parse(); err != nil {
		return "", err
	}
	output, err := generate(&parser.tree, variables)
	if err != nil {
		return "", err
	}
	return output, nil
}


