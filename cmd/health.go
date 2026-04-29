package cmd

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func healthCheckCmd(addr string) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("http://" + addr + "/health")
	if err != nil {
		fmt.Printf("health check failed: %s\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("health check failed: status %d\n", resp.StatusCode)
		os.Exit(1)
	}

	fmt.Println("OK")
}
