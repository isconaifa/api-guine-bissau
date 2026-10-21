package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/isconaifa/api-guine-bissau/internal/dados"
)

func init() { gin.SetMode(gin.TestMode) }

func baseTeste(t testing.TB) *dados.Base {
	t.Helper()
	b, err := os.ReadFile("../dados/testdata/valido.json")
	if err != nil {
		t.Fatal(err)
	}
	base, err := dados.Carregar(b)
	if err != nil {
		t.Fatal(err)
	}
	return base
}

func routerTeste(t testing.TB) *gin.Engine {
	t.Helper()
	return NovoRouter(baseTeste(t))
}

// pedir executa um pedido e devolve a resposta e o corpo JSON decodificado.
func pedir(t *testing.T, r http.Handler, metodo, caminho string) (*httptest.ResponseRecorder, map[string]any) {
	t.Helper()
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(metodo, caminho, nil))
	var corpo map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &corpo); err != nil {
		t.Fatalf("%s %s: corpo não é JSON (%v): %q", metodo, caminho, err, w.Body.String())
	}
	return w, corpo
}

// verificarErro confirma status, formato e código de uma resposta de erro.
func verificarErro(t *testing.T, w *httptest.ResponseRecorder, corpo map[string]any, status int, codigo string) string {
	t.Helper()
	if w.Code != status {
		t.Errorf("status = %d; esperado %d", w.Code, status)
	}
	if len(corpo) != 1 {
		t.Errorf("corpo de erro deve ter só a chave 'erro': %v", corpo)
	}
	erro, _ := corpo["erro"].(map[string]any)
	if erro["codigo"] != codigo {
		t.Errorf("codigo = %v; esperado %s", erro["codigo"], codigo)
	}
	mensagem, _ := erro["mensagem"].(string)
	if mensagem == "" {
		t.Error("mensagem de erro vazia")
	}
	return mensagem
}
