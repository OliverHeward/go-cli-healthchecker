package cmd

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

// Flags
var (
	timeout int
	urls    []string
	verbose bool
)

// Endpoint represents a service to health check
type Endpoint struct {
	Name string
	URL  string
}

// HealthResult contains detailed results from a health check
type HealthResult struct {
	Endpoint   Endpoint
	IsHealthy  bool
	StatusCode int
	Duration   time.Duration
	Error      error
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check the health of configured endpoints",
	Long: `Performs health checks on multiple endpoints concurrently.
	Each endpoint is checked via HTTP GET request, and results include
	status code, response time, and overall health status.
	
	Examples:
	  healthcheck check
	  healthcheck check --timeout 5
	  healthcheck check --urls https://api.github.com,https://dog.ceo/api/breeds/list/all
	  healthcheck check -t 3 -v`,
	Run: runCheck,
}

func init() {
	// Register this command with root
	rootCmd.AddCommand(checkCmd)

	// Define flags
	checkCmd.Flags().IntVarP(&timeout, "timeout", "t", 10, "Request timeout in seconds")
	checkCmd.Flags().StringSliceVarP(&urls, "urls", "u", []string{}, "Comma-separated list of endpoints to check")
	checkCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
}

func runCheck(cmd *cobra.Command, args []string) {
	fmt.Println("Health Checker v0.1")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━")

	if verbose {
		fmt.Printf("⚙️ Timeout: %ds\n", timeout)
	}
	fmt.Println()

	start := time.Now()

	// Use custom URLs if provided, otherwise use defaults
	var endpoints []Endpoint

	if len(urls) > 0 {
		for i, url := range urls {
			endpoints = append(endpoints, Endpoint{
				Name: fmt.Sprintf("Custom-%d", i+1),
				URL:  url,
			})
		}
	} else {
		// Use default endpoints
		endpoints = []Endpoint{
			{Name: "Github API", URL: "https://api.github.com"},
			{Name: "JSONPlaceholder", URL: "https://jsonplaceholder.typicode.com/posts/1"},
			{Name: "Dog Breeds API", URL: "https://dog.ceo/api/breeds/list/all"},
		}
	}

	var wg sync.WaitGroup

	for _, endpoint := range endpoints {
		wg.Add(1)

		go func(ep Endpoint) {
			defer wg.Done()
			result := checkEndpoint(ep)
			printResult(result)
		}(endpoint)
	}

	wg.Wait()

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("✓ Health check complete", len(endpoints), time.Since(start))
}

func checkEndpoint(endpoint Endpoint) HealthResult {
	start := time.Now()

	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	resp, err := client.Get(endpoint.URL)
	duration := time.Since(start)

	if err != nil {
		return HealthResult{
			Endpoint:  endpoint,
			IsHealthy: false,
			Duration:  duration,
			Error:     err,
		}
	}
	defer resp.Body.Close()

	isHealthy := resp.StatusCode >= 200 && resp.StatusCode < 400
	return HealthResult{
		Endpoint:   endpoint,
		IsHealthy:  isHealthy,
		StatusCode: resp.StatusCode,
		Duration:   duration,
		Error:      nil,
	}
}

func printResult(result HealthResult) {
	status := "✓ HEALTHY"
	if !result.IsHealthy {
		status = "✗ UNHEALTHY"
	}

	fmt.Printf("%s [%s]\n", status, result.Endpoint.Name)
	fmt.Printf("  URL: %s\n", result.Endpoint.URL)

	if result.Error != nil {
		fmt.Printf("  Error: %v\n", result.Error)
	} else {
		fmt.Printf("  Status: %d\n", result.StatusCode)
		fmt.Printf("  Response Time: %v\n", result.Duration)
	}
	fmt.Println()
}
