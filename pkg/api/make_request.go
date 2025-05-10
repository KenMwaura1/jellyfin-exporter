package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	activeSessionsGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "jellyfin_active_sessions",
		Help: "Number of active Jellyfin sessions",
	})
	jellyfinToken string // Add this to store the Jellyfin token
)

func init() {
	prometheus.MustRegister(activeSessionsGauge)
	jellyfinToken = "your-jellyfin-token-here" // Replace with your actual Jellyfin token
}

func fetchJellyfinMetrics() {
	for {
		sessions, err := getJellyfinSessions()
		if err != nil {
			log.Printf("Error fetching Jellyfin sessions: %v", err)
		} else {
			activeSessionsGauge.Set(float64(len(sessions)))
		}
		time.Sleep(15 * time.Second)
	}
}

func getJellyfinSessions() ([]interface{}, error) {
	req, err := http.NewRequest("GET", "http://192.168.100.16:8096/Sessions", nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "MediaBrowser Token="+jellyfinToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	var sessions []interface{}
	err = json.NewDecoder(resp.Body).Decode(&sessions)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return sessions, nil
}

func main() {
	go fetchJellyfinMetrics()

	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}