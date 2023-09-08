package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Advertisement struct {
	Id               string   `json:"id"`
	CampaignId       string   `json:"campaign_id"`
	Advertiser       string   `json:"advertiser"`
	Title            string   `json:"title"`
	Description      string   `json:"description"`
	ImageURL         string   `json:"image_url"`
	LinkURL          string   `json:"link_url"`
	StartDate        string   `json:"start_date"`
	EndDate          string   `json:"end_date"`
	DisplayFrequency string   `json:"display_frequency"`
	TargetAudience   []string `json:"target_audience"`
	Format           string   `json:"format"`
	MediaWidth       int      `json:"media_width"`
	MediaHeight      int      `json:"media_height"`
	MediaType        string   `json:"media_type"`
	Cost             float32  `json:"cost"`
}

// ===========================================================

func main() {
	http.HandleFunc("/", root)
	http.HandleFunc("/request", request)
	http.HandleFunc("/secondary", secondary)

	port := 8080
	addr := fmt.Sprintf(":%d", port)

	fmt.Printf("Starting server on port %d...\n", port)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Printf("Error starting server: %v", err)
	}
}

// ===========================================================

func root(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	htmlContent := `
	<html><body>
	<h1>Hello, Frameplay!</h1>
	</body></html>
	`
	fmt.Fprintln(w, htmlContent)
}

func request(w http.ResponseWriter, r *http.Request) {
	var payload Advertisement

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&payload); err != nil {
		http.Error(w, "ERR: Problem decoding JSON (data shape).", http.StatusBadRequest)
		return
	}

	if nonempty := advertisementFieldsNonEmpty(payload); nonempty == false {
		http.Error(w, "ERR: JSON must all expected fields non-nil.", http.StatusBadRequest)
		return
	}

	logAdvertisement(payload)

	err := invokeSecondary(payload)
	if err != nil {
		http.Error(w, fmt.Sprintf("ERR: Problem invoking secondary: %v", err), http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("SUCCESS: JSON data received and secondary invoked."))
}

func secondary(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	message := fmt.Sprintf("RECEIVED: %v", r.Body)
	fmt.Println(message)
	w.Write([]byte(message))
}

// ===========================================================

func logAdvertisement(ad Advertisement) {
	fmt.Printf("\n==========\n%v\n==========\n\n", ad)
}

func advertisementFieldsNonEmpty(ad Advertisement) bool {
	// We can have ads with a cost of $0.00 due to potential internal promotions.
	if ad.Id != "" &&
		ad.CampaignId != "" &&
		ad.Advertiser != "" &&
		ad.Title != "" &&
		ad.Description != "" &&
		ad.ImageURL != "" &&
		ad.LinkURL != "" &&
		ad.StartDate != "" &&
		ad.EndDate != "" &&
		ad.DisplayFrequency != "" &&
		len(ad.TargetAudience) != 0 &&
		ad.Format != "" &&
		ad.MediaWidth != 0 &&
		ad.MediaHeight != 0 &&
		ad.MediaType != "" {
		return true
	}
	return false
}

func invokeSecondary(ad Advertisement) error {
	payload, err := json.Marshal(ad)
	if err != nil {
		return fmt.Errorf("ERR: Could not marshal JSON.")
	}

	endpointURL := "http://localhost:8080/secondary"
	req, err := http.NewRequest(
		"POST",
		endpointURL,
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return fmt.Errorf("ERR: Failed to create new request client.")
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("ERR: Failed to invoke secondary through client.")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ERR: Invocation of secondary does not return success.")
	}
	return nil
}
