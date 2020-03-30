package main

type TokenChanManager struct {
	inserted Token
	curr     Token
	tokens   chan Token
}

func createTokenChanManager(tokens chan Token) TokenChanManager {
	return TokenChanManager{Token{}, <-tokens, tokens}
}

func (manager *TokenChanManager) insertAtFront(token Token) {
	manager.inserted = token
}

func (manager *TokenChanManager) next() Token {
	if manager.inserted.ty != TT_UNKNOWN {
		ret := manager.inserted
		manager.inserted = Token{}
		return ret
	}
	ret := manager.curr
	if manager.curr.ty != TT_EOF {
		manager.curr = <-manager.tokens
	}
	return ret
}

func (manager *TokenChanManager) peek() Token {
	if manager.inserted.ty != TT_UNKNOWN {
		return manager.inserted
	}
	return manager.curr
}

// type TokenChanManager struct {
// 	curr   Token
// 	next   Token
// 	tokens chan Token
// }

// func createTokenChanManager(tokens chan Token) TokenChanManager {
// 	curr := <-tokens
// 	next := <-tokens
// 	return TokenChanManager{curr, next, tokens}
// }

// func (manager TokenChanManager) nextToken() Token {
// 	manager.curr = manager.next
// 	if manager.curr.ty != TT_EOF {
// 		manager.next = <-manager.tokens
// 	}
// 	return manager.curr
// }

// func (manager TokenChanManager) peekToken() Token {
// 	return manager.curr
// }
