package request_limit

import (
	"log"
	"sync/atomic"
	"time"
)

type Manager struct {
	errorChan    chan error
	releaseChan  chan *Token
	outChan      chan *Token
	inChan       chan struct{}
	needToken    int64
	activeTokens map[string]*Token
	limit        int
	makeToken    tokenFactory
}

func NewManager(conf *Config) *Manager {
	m := &Manager{
		errorChan:    make(chan error),
		outChan:      make(chan *Token),
		inChan:       make(chan struct{}),
		activeTokens: make(map[string]*Token),
		releaseChan:  make(chan *Token),
		needToken:    0,
		limit:        conf.Limit,
		makeToken:    NewToken,
	}

	return m
}

func (m *Manager) Acquire() (*Token, error) {
	go func() {
		m.inChan <- struct{}{}
	}()

	select {
		case t := <-m.outChan:
			return t, nil
		case err := <-m.errorChan:
			return nil, err
	}
}

func (m *Manager) Release(t *Token) {
	go func() {
		m.releaseChan <- t
	}()
}

func (m *Manager) incNeedToken() {
	atomic.AddInt64(&m.needToken, 1)
}

func (m *Manager) decNeedToken() {
	atomic.AddInt64(&m.needToken, -1)
}

func (m *Manager) awaitingToken() bool {
	return atomic.LoadInt64(&m.needToken) > 0
}

func (m *Manager) tryGenerateToken() string {
	if m.makeToken != nil {
		if m.IsLimitExceeded() {
			m.incNeedToken()
			return "Too many requests"
		}

		token := m.makeToken()

		m.activeTokens[token.ID] = token

		go func() {
			m.outChan <- token
		}()
	}
	return ""
}

func (m *Manager) IsLimitExceeded() bool {
	if len(m.activeTokens) >= m.limit {
		return true
	}

	return false
}

func (m *Manager) releaseToken(token *Token) {
	if token == nil {
		log.Print("unable to relase nil token")
		return
	}

	if _, ok := m.activeTokens[token.ID]; !ok {
		log.Printf("unable to relase token %s - not in use", token)
		return
	}

	delete(m.activeTokens, token.ID)

	// process anything waiting for a rate limit
	if m.awaitingToken() {
		m.decNeedToken()
		go m.tryGenerateToken()
	}
}

//func (m *Manager) runResetTokenTask(resetAfter time.Duration) {
//	go func() {
//		ticker := time.NewTicker(resetAfter)
//		for range ticker.C {
//			for _, token := range m.activeTokens {
//				if token.NeedReset(resetAfter) {
//					go func(t *Token) {
//						m.releaseChan <- t
//					}(token)
//				}
//			}
//		}
//	}()
//}

func (t *Token) NeedReset(resetAfter time.Duration) bool {
	if time.Since(t.CreatedAt) >= resetAfter {
		return true
	}
	return false
}