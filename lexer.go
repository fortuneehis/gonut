package nut

import "fmt"


type lexer struct {
	tokens []token
	source []byte
	cursor *cursor
}

func newlexer(source []byte) *lexer {
	return &lexer{tokens: []token{}, source: source, cursor: newCursor(0)}
}

func (lexer *lexer) isEndOfFile() bool {
	return lexer.cursor.getValue() >= len(lexer.source)
}

func (lexer *lexer) seek() byte {
	char := lexer.peek()
	lexer.cursor.next()
	return char
}

func (lexer *lexer) getByteCollection(breakIf func(currentByte byte) bool) []byte {
	bytes := []byte{}
	lexer.cursor.previous()
	for !lexer.isEndOfFile() {
		currentByte := lexer.seek()
		if(breakIf(currentByte)) {
			lexer.cursor.previous()
			break
		}
		bytes = append(bytes, currentByte)
	}
	return bytes
}

func (lexer *lexer) createToken(kind int, value string) token {
	return *newToken(kind, value)
}

func (lexer *lexer) peek() byte {
	if(lexer.isEndOfFile()) {
		return 0
	}
	return lexer.source[lexer.cursor.getValue()]
}

func (lexer *lexer) scan() error {
	for !lexer.isEndOfFile() {
		current := lexer.seek()
		var token token
		switch current {
		case 10:
			token = lexer.createToken(newLine, string(current))
			lexer.tokens = append(lexer.tokens, token)
		
			lexer.getByteCollection(func(currentByte byte) bool {
				return currentByte != 10
			})
				nextByte := lexer.seek();

			if  nextByte != 32 {
				lexer.cursor.previous()
				continue
			}

			whiteSpace := lexer.getByteCollection(func(currentByte byte) bool {
				return currentByte != 32 
			})
			if currentByte := lexer.peek(); currentByte != 10 && currentByte != 35 && currentByte != 0 {
				token = lexer.createToken(indent, string(whiteSpace))
				lexer.tokens = append(lexer.tokens, token)
			}
			continue
		case 32:
			continue
		case 34:
			lexer.cursor.next()
			token = lexer.createToken(qString, string(lexer.getByteCollection(func(currentByte byte) bool {
				return currentByte == 34
			})))
			if lexer.isEndOfFile() || lexer.seek() != 34 {
				return fmt.Errorf("invalid or unexpected token")
			}
		case 39:
			lexer.cursor.next()
			token = lexer.createToken(qString, string(lexer.getByteCollection(func(currentByte byte) bool {
				return currentByte == 39
			})))
			if lexer.isEndOfFile() || lexer.seek() != 39 {
				return fmt.Errorf("invalid or unexpected token")
			}
		case 91:
			lexer.cursor.next()
			token = lexer.createToken(variable, string(lexer.getByteCollection(func(currentByte byte) bool {
				return currentByte == 93
			})))

			if lexer.isEndOfFile() || lexer.seek() != 93 {
				return fmt.Errorf("invalid or unexpected token")
			}
		case 35:
			lexer.getByteCollection(func(currentByte byte) bool {
				return currentByte == 10
			})
			continue
		case 44:
			token = lexer.createToken(comma, string(current))
		case 61:
			token = lexer.createToken(equals, string(current))
		case 40:
			token = lexer.createToken(lParen, string(current))
		case 41:
			token = lexer.createToken(rParen, string(current))
		default:
			token = lexer.createToken(node, string(lexer.getByteCollection(func(currentByte byte) bool {
				return currentByte == 32 || currentByte == 91 || currentByte == 93 || currentByte == 34 || currentByte == 10 || currentByte == 44 || currentByte == 39 || currentByte == 40 || currentByte == 41 || currentByte == 61
			}))) 
		}
		lexer.tokens = append(lexer.tokens, token)
	}
	return nil
}