package mock

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMockDBConnection(t *testing.T) {
	db, err := NewMockConnection()
	require.NoError(t, err)
	require.NotNil(t, db)

	db.Close()
	conn := db.DB()
	require.Nil(t, conn)
}
