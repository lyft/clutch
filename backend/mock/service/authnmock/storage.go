package authnmock

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
)

type MockAuthnStorage struct {
	// mapping from provider -> user -> token
	Tokens map[string]map[string]*oauth2.Token
}

func (m MockAuthnStorage) Store(ctx context.Context, userID, provider string, token *oauth2.Token) error {
	if _, ok := m.Tokens[provider]; !ok {
		m.Tokens[provider] = make(map[string]*oauth2.Token)
	}

	m.Tokens[provider][userID] = token

	return nil
}

func (m MockAuthnStorage) Read(ctx context.Context, userID, provider string) (*oauth2.Token, error) {
	if _, ok := m.Tokens[provider]; !ok {
		return nil, fmt.Errorf("token provider '%s' not found for user '%s'", provider, userID)
	}

	token, ok := m.Tokens[provider][userID]
	if !ok {
		return nil, fmt.Errorf("token user '%s' not found for provider '%s'", userID, provider)
	}

	return token, nil
}

func NewMockStorage() *MockAuthnStorage {
	return &MockAuthnStorage{
		Tokens: map[string]map[string]*oauth2.Token{},
	}
}
