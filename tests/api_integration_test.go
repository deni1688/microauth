package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestApi(t *testing.T) {
	token := ""

	t.Run("Login - POST /api/v1/login", func(t *testing.T) {
		token = getToken(t)

		if token == "" {
			t.Fatal("token is empty")
		}

		if len(token) != 256/4 {
			t.Fatalf("expected token length to be 256 bytes, got %d", len(token))
		}
	})

	t.Run("Create Credential - POST /api/v1/dashboard/credentials", func(t *testing.T) {
		payload := `{"name": "app_123", "password": "who281330"}`
		req, err := http.NewRequest(http.MethodPost, "http://localhost:9876/api/v1/dashboard/credentials", strings.NewReader(payload))
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", token)
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}

		var result map[string]any
		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		err = json.Unmarshal(respBytes, &result)
		if err != nil {
			t.Fatal(err)
		}

		if result["error"] != nil {
			t.Fatalf("expected no error, got '%s'", result["error"])
		}

		if result["message"] != "credential saved" {
			t.Fatalf("expected message 'credential saved', got '%s'", result["message"])
		}
	})

	t.Run("Get Credentials - GET /api/v1/dashboard/credentials", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "http://localhost:9876/api/v1/dashboard/credentials", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}

		var result []map[string]any
		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		err = json.Unmarshal(respBytes, &result)
		if err != nil {
			t.Fatal(err)
		}

		if len(result) != 2 {
			t.Fatalf("expected 2 credential, got %d", len(result))
		}
	})

	t.Run("Update Credential - POST /api/v1/dashboard/credentials", func(t *testing.T) {
		payload := `{"id": 2, "name": "app_123"}`
		req, err := http.NewRequest(http.MethodPost, "http://localhost:9876/api/v1/dashboard/credentials", strings.NewReader(payload))
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", token)
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		var result map[string]any
		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		err = json.Unmarshal(respBytes, &result)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}

		if result["error"] != nil {
			t.Fatalf("expected no error, got '%s'", result["error"])
		}

		if result["message"] != "credential saved" {
			t.Fatalf("expected message 'credential saved', got '%s'", result["message"])
		}
	})

	t.Run("Delete Credential - DELETE /api/v1/dashboard/credentials/:id", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, "http://localhost:9876/api/v1/dashboard/credentials/2", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}

		var result map[string]any
		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		err = json.Unmarshal(respBytes, &result)
		if err != nil {
			t.Fatal(err)
		}

		if result["error"] != nil {
			t.Fatalf("expected no error, got '%s'", result["error"])
		}

		if result["message"] != "credential deleted" {
			t.Fatalf("expected message 'credential deleted', got '%s'", result["message"])
		}
	})

	t.Run("Logout - POST /api/v1/logout", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "http://localhost:9876/api/v1/logout", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}

		var result map[string]any
		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		err = json.Unmarshal(respBytes, &result)
		if err != nil {
			t.Fatal(err)
		}

		if result["error"] != nil {
			t.Fatalf("expected no error, got '%s'", result["error"])
		}

		if result["message"] != "logout successful" {
			t.Fatalf("expected message 'logout successful', got '%s'", result["message"])
		}
	})

	t.Run("Get Credentials with Invalid Token - GET /api/v1/dashboard/credentials", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "http://localhost:9876/api/v1/dashboard/credentials", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}
	})
}

func getToken(t *testing.T) string {
	credentials := `{"name": "root", "password": "root1234"}`
	resp, err := http.Post("http://localhost:9876/api/v1/login", "application/json", strings.NewReader(credentials))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	tokenBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]string
	err = json.Unmarshal(tokenBytes, &result)
	if err != nil {
		t.Fatal(err)
	}

	return result["token"]
}
