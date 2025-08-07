package mock

import (
	"testing"

	"github.com/barantoraman/microgate/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestMockDBConnection(t *testing.T) {
	var cfg config.AuthServiceConfigurations

	db, err := NewMockConnection(cfg)
	require.NoError(t, err)
	require.NotNil(t, db)

	db.Close()
	conn := db.DB()
	require.Nil(t, conn)
}
