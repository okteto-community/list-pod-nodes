package main

import (
	"encoding/csv"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/okteto-community/list-pod-nodes/api"
)

func main() {
	token := os.Getenv("OKTETO_TOKEN")
	oktetoURL := os.Getenv("OKTETO_URL")

	logLevel := &slog.LevelVar{} // INFO
	opts := &slog.HandlerOptions{
		Level: logLevel,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))

	if token == "" || oktetoURL == "" {
		logger.Error("OKTETO_TOKEN and OKTETO_URL environment variables are required")
		os.Exit(1)
	}

	u, err := url.Parse(oktetoURL)
	if err != nil {
		logger.Error(fmt.Sprintf("Invalid OKTETO_URL %s", err))
		os.Exit(1)
	}

	nsList, err := api.GetDevelopmentNamespaces(u.Host, token, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("Error getting namespaces: %s", err))
		os.Exit(1)
	}
	// Create CSV File
	currentTime := time.Now()
	filename := fmt.Sprintf("pod-nodes_%s.csv", currentTime.Format("2006-01-02_15-04-05"))
	file, err := os.Create(filename)
	if err != nil {
		logger.Error(fmt.Sprintf("Error creating CSV file: %s", err))
		os.Exit(1)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers
	headers := []string{"Pod Name", "Namespace", "Node Name"}
	if err := writer.Write(headers); err != nil {
		logger.Error(fmt.Sprintf("Error writing headers to CSV: %s", err))
		os.Exit(1)
	}

	// Write data rows
	for _, ns := range nsList {
		cmdStr := fmt.Sprintf(`
		kubectl get pod -n %s -o json | jq -r '
		.items[]
		| "\(.metadata.name)\t\(.metadata.namespace)\t\(.spec.nodeName)"
		'`, ns.Name)
		cmd := exec.Command("bash", "-c", cmdStr)
		output, err := cmd.CombinedOutput()

		if err != nil {
			fmt.Printf("Error executing command: %v\n", err)
			return
		}

		outputStr := strings.TrimSpace(string(output))
		if outputStr != "" {
			lines := strings.Split(string(outputStr), "\n")
			for _, line := range lines {
				parts := strings.Split(line, "\t")
				podName := parts[0]
				namespaceName := parts[1]
				nodeName := parts[2]

				// Write the row to the CSV file
				row := []string{podName, namespaceName, nodeName}
				if err := writer.Write(row); err != nil {
					logger.Error(fmt.Sprintf("Error writing row to CSV: %s", err))
					os.Exit(1)
				}
				fmt.Println(line)
			}
		}
	}

	if err != nil {
		logger.Error(fmt.Sprintf("There was an error requesting the namespaces: %s", err))
		os.Exit(1)
	}
}
