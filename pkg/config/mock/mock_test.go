package mock

import (
	"testing"

	"github.com/barantoraman/microgate/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestMockLoader_GetConfigByKey(t *testing.T) {
	const configKey = "auth_service"

	type testCase struct {
		name         string
		input        map[string]any
		valid        bool
		expectErrMsg string
	}

	cases := []testCase{
		{
			name: "Valid auth service configurations",
			input: map[string]any{
				configKey: map[string]any{
					"db_type": "postgres",
					"db_user": "admin",
					"db_pass": "secret",
					"db_host": "localhost",
					"db_port": 5432,
					"db_name": "authdb",
				},
			},
			valid: true,
		},
		{
			name:         "Key not found",
			input:        map[string]any{},
			valid:        false,
			expectErrMsg: "key not found",
		},
		{
			name: "Marshal failure",
			input: map[string]any{
				configKey: func() {},
			},
			valid:        false,
			expectErrMsg: "failed to marshal mock data",
		},
		{
			name: "Unmarshal failure",
			input: map[string]any{
				configKey: "this is not a struct",
			},
			valid:        false,
			expectErrMsg: "failed to unmarshal mock data",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			loader := NewMockLoader(tt.input)
			var cfg config.AuthServiceConfigurations
			err := loader.GetConfigByKey(configKey, &cfg)

			if tt.valid {
				require.NoError(t, err)
				assertValidConfig(t, cfg)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectErrMsg)
			}
		})
	}
}

func assertValidConfig(t *testing.T, cfg config.AuthServiceConfigurations) {
	require.Equal(t, "postgres", cfg.DBType)
	require.Equal(t, "admin", cfg.DBUser)
	require.Equal(t, 5432, cfg.DBPort)
	require.Equal(t, "localhost", cfg.DBHost)
	require.Equal(t, "authdb", cfg.DBName)
}
