package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	base := flag.String("base", envOr("API_BASE", "http://localhost:8999"), "API base URL")
	token := flag.String("token", envOr("OPENVPN_API_TOKEN", "openvpn-ui-api-token"), "API token")
	name := flag.String("name", "", "certificate name for integration tests")
	integration := flag.Bool("integration", false, "run create/renew/revoke/delete/download lifecycle")
	flag.Parse()

	if *name == "" {
		*name = fmt.Sprintf("api-test-%d", time.Now().Unix())
	}

	client := &http.Client{Timeout: 30 * time.Second}
	runner := &runner{
		client: client,
		base:   strings.TrimRight(*base, "/"),
		token:  *token,
		name:   *name,
	}

	runner.runSmokeTests()

	if *integration {
		runner.runIntegrationTests()
	} else {
		fmt.Println("\n(skip integration tests, use -integration to enable)")
	}

	if runner.failed > 0 {
		fmt.Printf("\n%d test(s) failed\n", runner.failed)
		os.Exit(1)
	}
	fmt.Printf("\nall %d test(s) passed\n", runner.passed)
}

type runner struct {
	client *http.Client
	base   string
	token  string
	name   string
	passed int
	failed int
}

func (r *runner) runSmokeTests() {
	fmt.Println("== smoke tests ==")
	r.check("GET without token returns 401", func() error {
		return r.expectStatus("GET", "/api/v1/certificates/smoke-missing", "", nil, 401)
	})
	r.check("GET with bad token returns 401", func() error {
		return r.expectStatus("GET", "/api/v1/certificates/smoke-missing", "bad-token", nil, 401)
	})
	r.check("GET nonexistent cert returns 404", func() error {
		return r.expectStatus("GET", "/api/v1/certificates/smoke-missing", r.token, nil, 404)
	})
	r.check("POST empty body returns 400", func() error {
		return r.expectStatus("POST", "/api/v1/certificates", r.token, []byte(`{}`), 400)
	})
	r.check("POST renew unknown name returns 400", func() error {
		return r.expectStatus("POST", "/api/v1/certificates/smoke-missing/renew", r.token, nil, 400)
	})
}

func (r *runner) runIntegrationTests() {
	fmt.Println("\n== integration tests ==")
	r.check("POST create certificate", func() error {
		body := []byte(fmt.Sprintf(`{"name":%q}`, r.name))
		resp, err := r.do("POST", "/api/v1/certificates", r.token, body)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			b, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("status %d: %s", resp.StatusCode, strings.TrimSpace(string(b)))
		}
		return nil
	})
	if r.failed > 0 {
		fmt.Println("(skip remaining integration tests: create failed, PKI/scripts may be unavailable)")
		return
	}

	r.check("GET download .ovpn", func() error {
		resp, err := r.do("GET", "/api/v1/certificates/"+r.name, r.token, nil)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		data, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != 200 {
			return fmt.Errorf("status %d: %s", resp.StatusCode, strings.TrimSpace(string(data)))
		}
		if !strings.Contains(string(data), "client") && !strings.Contains(string(data), "remote") {
			return fmt.Errorf("response does not look like an .ovpn file")
		}
		return nil
	})

	r.check("POST renew certificate", func() error {
		return r.expectStatus("POST", "/api/v1/certificates/"+r.name+"/renew", r.token, nil, 200)
	})
	r.check("POST revoke certificate", func() error {
		return r.expectStatus("POST", "/api/v1/certificates/"+r.name+"/revoke", r.token, nil, 200)
	})
	r.check("DELETE certificate", func() error {
		return r.expectStatus("DELETE", "/api/v1/certificates/"+r.name, r.token, nil, 200)
	})
}

func (r *runner) check(name string, fn func() error) {
	err := fn()
	if err != nil {
		r.fail(name, err)
		return
	}
	r.pass(name)
}

func (r *runner) pass(name string) {
	r.passed++
	fmt.Printf("  PASS  %s\n", name)
}

func (r *runner) fail(name string, err error) {
	r.failed++
	fmt.Printf("  FAIL  %s: %v\n", name, err)
}

func (r *runner) expectStatus(method, path, token string, body []byte, want int) error {
	resp, err := r.do(method, path, token, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != want {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("want HTTP %d, got %d: %s", want, resp.StatusCode, strings.TrimSpace(string(b)))
	}
	return nil
}

func (r *runner) do(method, path, token string, body []byte) (*http.Response, error) {
	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}
	req, err := http.NewRequest(method, r.base+path, reader)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return r.client.Do(req)
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
