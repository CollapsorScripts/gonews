package config

import "testing"

func TestInit(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Тест загрузки конфигурации",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Init()
		})
		t.Logf("Кофигурация: %+v", *Cfg)
	}
}
