package yahoo

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

func getCookieAndCrumb() (*http.Client, string, error) {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	// Create request to get crumb
	req1, _ := http.NewRequest("GET", "https://fc.yahoo.com", nil)
	req1.Header.Set("User-Agent", "Mozilla/5.0")
	resp1, err := client.Do(req1)
	if err != nil {
		return nil, "", fmt.Errorf("warm cookies: %w", err)
	}
	io.Copy(io.Discard, resp1.Body)
	resp1.Body.Close()

	// Get crumb with cookie
	req2, _ := http.NewRequest("GET", "https://query2.finance.yahoo.com/v1/test/getcrumb", nil)
	req2.Header.Set("User-Agent", "Mozilla/5.0")
	resp2, err := client.Do(req2)
	if err != nil {
		return nil, "", fmt.Errorf("get crumb: %w", err)
	}
	defer resp2.Body.Close()

	crumb, err := io.ReadAll(resp2.Body)
	if err != nil {
		fmt.Printf("Error occurred while reading response body: %s", err)
		return nil, "", fmt.Errorf("read crumb: %w", err)
	}

	// Trim crumb if necessary
	trimmedCrumb := strings.Trim(string(crumb), "\" \n\r")
	if trimmedCrumb == "" {
		return nil, "", fmt.Errorf("empty crumb")
	}
	return client, trimmedCrumb, nil
}
