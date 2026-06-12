package prometheus

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Metric struct {
	Name   string
	Labels map[string]string
	Value  float64
}

var (
	lineRe  = regexp.MustCompile(`^(\w+)\{([^}]*)\}\s+([\d.eE+\-]+)$`)
	labelRe = regexp.MustCompile(`(\w+)="([^"]*)"`)
)

func Parse(text string) ([]Metric, error) {
	var metrics []Metric
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		m := lineRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		val, err := strconv.ParseFloat(m[3], 64)
		if err != nil {
			continue
		}
		labels := map[string]string{}
		for _, lm := range labelRe.FindAllStringSubmatch(m[2], -1) {
			labels[lm[1]] = lm[2]
		}
		metrics = append(metrics, Metric{Name: m[1], Labels: labels, Value: val})
	}
	return metrics, nil
}

func Fetch(ctx context.Context, address, port string, timeout time.Duration) ([]Metric, error) {
	url := fmt.Sprintf("http://%s:%s/metrics", address, port)
	reqCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching metrics: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body: %w", err)
	}
	return Parse(string(body))
}
