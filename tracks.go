package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type Track struct {
	Album struct {
		AlbumType string `json:"album_type"`
		Artists   []struct {
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href string `json:"href"`
			ID   string `json:"id"`
			Name string `json:"name"`
			Type string `json:"type"`
			URI  string `json:"uri"`
		} `json:"artists"`
		AvailableMarkets []string `json:"available_markets"`
		ExternalUrls     struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Href   string `json:"href"`
		ID     string `json:"id"`
		Images []struct {
			Height int    `json:"height"`
			URL    string `json:"url"`
			Width  int    `json:"width"`
		} `json:"images"`
		Name                 string `json:"name"`
		ReleaseDate          string `json:"release_date"`
		ReleaseDatePrecision string `json:"release_date_precision"`
		Type                 string `json:"type"`
		URI                  string `json:"uri"`
	} `json:"album"`
	Artists []struct {
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Href string `json:"href"`
		ID   string `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
		URI  string `json:"uri"`
	} `json:"artists"`
	AvailableMarkets []string `json:"available_markets"`
	DiscNumber       int      `json:"disc_number"`
	DurationMs       int      `json:"duration_ms"`
	Explicit         bool     `json:"explicit"`
	ExternalIds      struct {
		Isrc string `json:"isrc"`
	} `json:"external_ids"`
	ExternalUrls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Href        string `json:"href"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Popularity  int    `json:"popularity"`
	PreviewURL  string `json:"preview_url"`
	TrackNumber int    `json:"track_number"`
	Type        string `json:"type"`
	URI         string `json:"uri"`
}

type Item struct {
	AddedAt string `json:"added_at"`
	AddedBy struct {
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Href string `json:"href"`
		ID   string `json:"id"`
		Type string `json:"type"`
		URI  string `json:"uri"`
	} `json:"added_by"`
	IsLocal bool  `json:"is_local"`
	Track   Track `json:"track"`
}

type Tracks struct {
	Href     string      `json:"href"`
	Items    []Item      `json:"items"`
	Limit    int         `json:"limit"`
	Next     string      `json:"next"`
	Offset   int         `json:"offset"`
	Previous interface{} `json:"previous"`
	Total    int         `json:"total"`
}

// GetAllTracks returns all tracks from the provided playlist URI
func GetAllTracks(playlistURI string, accessToken string) ([]Item, error) {
	var trackItems []Item

	resp, err := requestPlaylistTracks(playlistURI, accessToken)
	if err != nil {
		return nil, err
	}

	initialTracks, err := getTracks(resp)
	if err != nil {
		return nil, err
	}
	trackItems = append(trackItems, initialTracks.Items...)

	next := initialTracks.Next
	for next != "" {
		resp, err := requestPlaylistTracksFromFullPath(next, accessToken)
		newTracks, err := getTracks(resp)
		if err != nil {
			return nil, err
		}
		trackItems = append(trackItems, newTracks.Items...)
		next = newTracks.Next
	}

	return trackItems, err
}

func getTracks(resp *http.Response) (Tracks, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Tracks{}, err
	}
	var tracks Tracks
	err = json.Unmarshal(body, &tracks)
	return tracks, nil
}

func requestPlaylistTracks(playlistURI string, accessToken string) (*http.Response, error) {
	queryParams := strings.Split(playlistURI, ":")
	uri := "https://api.spotify.com/v1/users/" + queryParams[2] + "/playlists/" + queryParams[4] + "/tracks"

	client := &http.Client{}
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)
	return client.Do(req)
}

func requestPlaylistTracksFromFullPath(path string, accessToken string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)
	return client.Do(req)
}