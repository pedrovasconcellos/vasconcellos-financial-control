package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestLoadConfigFromFile garante que o carregamento via arquivo YAML popula os campos básicos
func TestLoadConfigFromFile(t *testing.T) {
	t.Cleanup(func() {
		reset()
		os.Unsetenv("CONFIG_FILE")
	})

	dir := t.TempDir()
	filePath := filepath.Join(dir, "cfg.yaml")
	if err := os.WriteFile(filePath, []byte(`
app:
  name: test-app
  environment: homolog
  port: 9000
mongo:
  uri: mongodb://localhost:27017
  database: financial
`), 0o600); err != nil {
		t.Fatalf("falha ao criar arquivo temporário: %v", err)
	}
	if err := os.Setenv("CONFIG_FILE", filePath); err != nil {
		t.Fatalf("falha ao setar variável de ambiente: %v", err)
	}

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig retornou erro: %v", err)
	}
	if cfg.App.Name != "test-app" {
		t.Errorf("esperado app.name 'test-app', obtido '%s'", cfg.App.Name)
	}
	if cfg.App.Port != 9000 {
		t.Errorf("esperado app.port 9000, obtido %d", cfg.App.Port)
	}
	if cfg.Mongo.Database != "financial" {
		t.Errorf("esperado mongo.database 'financial', obtido '%s'", cfg.Mongo.Database)
	}
}

// TestLoadConfigEnvOverride valida que variáveis de ambiente sobrescrevem valores do arquivo
func TestLoadConfigEnvOverride(t *testing.T) {
	t.Cleanup(func() {
		reset()
		os.Unsetenv("CONFIG_FILE")
		os.Unsetenv("APP_PORT")
	})

	dir := t.TempDir()
	filePath := filepath.Join(dir, "cfg.yaml")
	if err := os.WriteFile(filePath, []byte(`app: { port: 8000 }`), 0o600); err != nil {
		t.Fatalf("falha ao criar arquivo temporário: %v", err)
	}
	if err := os.Setenv("CONFIG_FILE", filePath); err != nil {
		t.Fatalf("falha ao setar variável de ambiente: %v", err)
	}
	if err := os.Setenv("APP_PORT", "1234"); err != nil {
		t.Fatalf("falha ao setar APP_PORT: %v", err)
	}

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig retornou erro: %v", err)
	}
	if cfg.App.Port != 1234 {
		t.Errorf("esperado app.port 1234, obtido %d", cfg.App.Port)
	}
}
