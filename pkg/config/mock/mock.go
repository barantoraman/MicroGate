package mock

import (
	"encoding/json"
	"fmt"

	loaderContract "github.com/barantoraman/microgate/pkg/config/contract"
)

type mockLoader struct {
	Data map[string]any
}

func NewMockLoader(data map[string]any) loaderContract.Loader {
	return &mockLoader{
		Data: data,
	}
}

func (m *mockLoader) GetConfigByKey(key string, config any) error {
	val, ok := m.Data[key]
	if !ok {
		return fmt.Errorf("key not found: %s", key)
	}

	bytes, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("failed to marshal mock data: %w", err)
	}

	if err := json.Unmarshal(bytes, config); err != nil {
		return fmt.Errorf("failed to unmarshal mock data: %w", err)
	}
	return nil
}
