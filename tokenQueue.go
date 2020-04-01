package main

type TokenQueue struct {
	tokens []Token
}

func (manager *TokenQueue) size() int {
	return len(manager.tokens)
}

func createTokenQueue(tokens chan Token) TokenQueue {
	return TokenQueue{}
}

func (manager *TokenQueue) pushFront(token Token) {
	manager.tokens = append(manager.tokens, Token{})
	copy(manager.tokens[1:], manager.tokens)
	manager.tokens[0] = token
}

func (manager *TokenQueue) pushBack(token Token) {
	manager.tokens = append(manager.tokens, token)
}

func (manager *TokenQueue) next() Token {
	ret := manager.peek()
	if ret.ty != TT_EOF {
		manager.tokens = manager.tokens[1:]
	}
	return ret
}

func (manager *TokenQueue) peek() Token {
	return manager.tokens[0]
}
