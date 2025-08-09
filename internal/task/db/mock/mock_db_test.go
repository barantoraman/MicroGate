package mock

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMockDBConnection(t *testing.T) {
	db, err := NewMockConnection()
	require.NoError(t, err)
	require.NotNil(t, db)

	conn := db.DB()
	require.Nil(t, conn)

	db.Close()
}
