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

	t.Run("Login - POST /rest/v1/login", func(t *testing.T) {
		token = getToken(t)

		if token == "" {
			t.Fatal("token is empty")
		}

		if len(token) != 256/4 {
			t.Fatalf("expected token length to be 256 bytes, got %d", len(token))
		}
	})

	t.Run("Create Admin - POST /rest/v1/dashboard/admins", func(t *testing.T) {
		payload := `{"email": "mike.jones@mail.com", "password": "who281330", "firstname": "Mike", "lastname": "Jones"}`
		req, err := http.NewRequest(http.MethodPost, "http://localhost:9876/api/v1/dashboard/admins", strings.NewReader(payload))
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

		if result["message"] != "core saved" {
			t.Fatalf("expected message 'core saved', got '%s'", result["message"])
		}
	})

	t.Run("Get Admins - GET /rest/v1/dashboard/admins", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "http://localhost:9876/api/v1/dashboard/admins", nil)
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
			t.Fatalf("expected 2 core, got %d", len(result))
		}
	})

	t.Run("Update Admin - POST /rest/v1/dashboard/admins", func(t *testing.T) {
		payload := `{"id": 2, "email": "mike-jones@mail.com", "firstname": "Michael", "lastname": "Jones"}`
		req, err := http.NewRequest(http.MethodPost, "http://localhost:9876/api/v1/dashboard/admins", strings.NewReader(payload))
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

		if result["message"] != "core saved" {
			t.Fatalf("expected message 'core saved', got '%s'", result["message"])
		}
	})

	t.Run("Delete Admin - DELETE /rest/v1/dashboard/admins/:id", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, "http://localhost:9876/api/v1/dashboard/admins/2", nil)
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

		if result["message"] != "core deleted" {
			t.Fatalf("expected message 'core deleted', got '%s'", result["message"])
		}
	})

	t.Run("Logout - POST /rest/v1/logout", func(t *testing.T) {
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

	t.Run("Get Admins with Invalid Token - GET /rest/v1/dashboard/admins", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "http://localhost:9876/api/v1/dashboard/admins", nil)
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
	credentials := `{"email": "root@mail.com", "password": "root1234"}`
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
