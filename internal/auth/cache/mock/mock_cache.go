package mock

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	cacheContract "github.com/barantoraman/microgate/internal/auth/cache/contract"
	tokenPkg "github.com/barantoraman/microgate/pkg/token"
)

type mockStore struct {
	msMap map[string][]byte
	msMu  sync.RWMutex
}

func NewMockStore() cacheContract.Store {
	return &mockStore{
		msMap: make(map[string][]byte),
		msMu:  sync.RWMutex{},
	}
}

func (m *mockStore) Delete(ctx context.Context, token string) error {
	m.msMu.Lock()
	defer m.msMu.Unlock()

	delete(m.msMap, token)
	return nil
}

func (m *mockStore) Get(ctx context.Context, sessionToken string) (tokenPkg.Token, error) {
	hash := sha256.Sum256([]byte(sessionToken))
	tHash := hash[:]

	m.msMu.RLock()
	defer m.msMu.RUnlock()

	val, ok := m.msMap[string(tHash)]
	if !ok {
		return tokenPkg.Token{}, errors.New("session not found")
	}
	var session tokenPkg.Token
	if err := json.Unmarshal(val, &session); err != nil {
		return tokenPkg.Token{}, errors.New("failed to unmarshal session")
	}
	return session, nil
}

func (m *mockStore) Set(ctx context.Context, sessionToken *tokenPkg.Token) error {
	sessionByteArr, err := json.Marshal(sessionToken)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %v", err)
	}

	m.msMu.Lock()
	defer m.msMu.Unlock()

	m.msMap[string(sessionToken.Hash)] = sessionByteArr
	return nil
}
