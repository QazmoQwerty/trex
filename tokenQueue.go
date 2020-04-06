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
	if len(manager.tokens) == 0 {
		return Token{TT_UNKNOWN, "", Position{0, 0, 0}}
	}
	return manager.tokens[0]
}

func (manager *TokenQueue) peekBack() Token {
	if len(manager.tokens) == 0 {
		return Token{TT_UNKNOWN, "", Position{0, 0, 0}}
	}
	return manager.tokens[len(manager.tokens)-1]
}

func (manager *TokenQueue) peekBeforeBack() Token {
	if len(manager.tokens) <= 1 {
		return Token{TT_UNKNOWN, "", Position{0, 0, 0}}
	}
	return manager.tokens[len(manager.tokens)-2]
}

func (manager *TokenQueue) popBack() Token {
	if len(manager.tokens) == 0 {
		return Token{TT_UNKNOWN, "", Position{0, 0, 0}}
	}
	defer func() { manager.tokens = manager.tokens[:len(manager.tokens)-1] }()
	return manager.peekBack()
}
