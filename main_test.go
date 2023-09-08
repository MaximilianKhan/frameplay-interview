package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSONRequestSuccess(t *testing.T) {
	jsonRequest := []byte(
		`{
			"id": "123", 
			"campaign_id": "123", 
			"advertiser": "Frameplay", 
			"title": "The coolest ad you've ever seen", 
			"description": "Doesn't get cooler than this", 
			"image_url": "https://unsplash.com/photos/a-large-iceberg-in-the-middle-of-the-ocean-n2qV323Fitc", 
			"link_url": "https://www.frameplay.gg/", 
			"start_date": "2023-09-07", 
			"end_date": "2023-09-09", 
			"display_frequency": "hourly", 
			"target_audience": ["gamers", "players", "developers"], 
			"format": "banner", 
			"media_width": 1920, 
			"media_height": 1080, 
			"media_type": "image", 
			"cost": 999.99
		}`,
	)
	req, err := http.NewRequest(
		"POST",
		"/request",
		bytes.NewBuffer(jsonRequest),
	)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	request(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	expectedResponse := "SUCCESS: JSON data received and secondary invoked."
	if rr.Body.String() != expectedResponse {
		t.Errorf("ERR: Expected response '%s', got '%s'", expectedResponse, rr.Body.String())
	}
}

func TestJSONRequesFailureNotEnoughFields(t *testing.T) {
	jsonRequest := []byte(
		`{
			"id": "123", 
			"campaign_id": "123", 
			"advertiser": "Frameplay", 
			"title": "The coolest ad you've ever seen", 
			"description": "Doesn't get cooler than this", 
			"image_url": "https://unsplash.com/photos/a-large-iceberg-in-the-middle-of-the-ocean-n2qV323Fitc", 
			"link_url": "https://www.frameplay.gg/"
		}`,
	)
	req, err := http.NewRequest(
		"POST",
		"/request",
		bytes.NewBuffer(jsonRequest),
	)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	request(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestJSONRequesFailureUnknownField(t *testing.T) {
	jsonRequest := []byte(
		`{
			"id": "123", 
			"campaign_id": "123", 
			"advertiser": "Frameplay", 
			"title": "The coolest ad you've ever seen", 
			"description": "Doesn't get cooler than this", 
			"image_url": "https://unsplash.com/photos/a-large-iceberg-in-the-middle-of-the-ocean-n2qV323Fitc", 
			"link_url": "https://www.frameplay.gg/", 
			"start_date": "2023-09-07", 
			"end_date": "2023-09-09", 
			"display_frequency": "hourly", 
			"target_audience": ["gamers", "players", "developers"], 
			"format": "banner", 
			"media_width": 1920, 
			"media_height": 1080, 
			"media_type": "image", 
			"cost": 999.99,
			"unknown_field": "made you look"
		}`,
	)
	req, err := http.NewRequest(
		"POST",
		"/request",
		bytes.NewBuffer(jsonRequest),
	)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	request(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rr.Code)
	}
}
