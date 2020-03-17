package main

type TokenChanManager struct {
	inserted Token
	curr     Token
	tokens   chan Token
}

func createTokenChanManager(tokens chan Token) TokenChanManager {
	return TokenChanManager{Token{}, <-tokens, tokens}
}

func (this *TokenChanManager) insertAtFront(token Token) {
	this.inserted = token
}

func (this *TokenChanManager) next() Token {
	if this.inserted.ty != TT_UNKNOWN {
		ret := this.inserted
		this.inserted = Token{}
		return ret
	}
	ret := this.curr
	if this.curr.ty != TT_EOF {
		this.curr = <-this.tokens
	}
	return ret
}

func (this *TokenChanManager) peek() Token {
	if this.inserted.ty != TT_UNKNOWN {
		return this.inserted
	}
	return this.curr
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

// func (this TokenChanManager) nextToken() Token {
// 	this.curr = this.next
// 	if this.curr.ty != TT_EOF {
// 		this.next = <-this.tokens
// 	}
// 	return this.curr
// }

// func (this TokenChanManager) peekToken() Token {
// 	return this.curr
// }
