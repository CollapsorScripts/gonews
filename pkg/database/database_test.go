package database

import (
	"newsaggr/pkg/config"
	"testing"
)

// preload - загрузка конфигурации
func preload() error {
	if err := config.Init(); err != nil {
		return err
	}

	return nil
}

func TestInit(t *testing.T) {
	if err := preload(); err != nil {
		t.Errorf("%v", err)
		return
	}

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Тест подключения к бд",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Init()
			if (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("Экземпляр подключения к бд: %+v", got)
		})
	}
}
