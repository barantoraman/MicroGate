package mock

import (
	"context"
	"crypto/sha256"
	"testing"
	"time"

	tokenPkg "github.com/barantoraman/microgate/pkg/token"
	"github.com/stretchr/testify/require"
)

func TestNewMockStore(t *testing.T) {
	t.Run("New mock store", func(t *testing.T) {
		store := NewMockStore()
		require.NotNil(t, store)
	})
}

func TestMockStore_SetGetDelete(t *testing.T) {
	ctx := context.Background()
	store := NewMockStore()

	plainText := "example-token"
	hash := sha256.Sum256([]byte(plainText))

	token := &tokenPkg.Token{
		PlainText: plainText,
		Hash:      hash[:],
		UserID:    123,
		Expiry:    time.Now().Add(10 * time.Minute),
		Scope:     "auth",
	}

	// SET
	err := store.Set(ctx, token)
	require.NoError(t, err, "set should not return an error")

	// GET (existing)
	gotToken, err := store.Get(ctx, plainText)
	require.NoError(t, err, "get should not return an error for existing token")
	require.Equal(t, token.UserID, gotToken.UserID, "UserID should match")
	require.Equal(t, token.PlainText, gotToken.PlainText, "PlainText should match")
	require.WithinDuration(t, token.Expiry, gotToken.Expiry, time.Second, "Expiry should be close enough")

	// DELETE
	err = store.Delete(ctx, string(token.Hash))
	require.NoError(t, err, "Delete should not return an error")

	// GET (after delete)
	_, err = store.Get(ctx, plainText)
	require.Error(t, err, "get should return error for deleted token")
	require.Equal(t, "session not found", err.Error(), "Expected session not found error")
}
