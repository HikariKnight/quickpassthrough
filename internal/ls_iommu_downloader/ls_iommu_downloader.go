package ls_iommu_downloader

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/HikariKnight/ls-iommu/pkg/errorcheck"
	"github.com/HikariKnight/quickpassthrough/internal/untar"
	"github.com/cavaliergopher/grab/v3"
)

// Generated from github API response using https://mholt.github.io/json-to-go/
type Response struct {
	URL       string `json:"url"`
	AssetsURL string `json:"assets_url"`
	UploadURL string `json:"upload_url"`
	HTMLURL   string `json:"html_url"`
	ID        int    `json:"id"`
	Author    struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"author"`
	NodeID          string    `json:"node_id"`
	TagName         string    `json:"tag_name"`
	TargetCommitish string    `json:"target_commitish"`
	Name            string    `json:"name"`
	Draft           bool      `json:"draft"`
	Prerelease      bool      `json:"prerelease"`
	CreatedAt       time.Time `json:"created_at"`
	PublishedAt     time.Time `json:"published_at"`
	Assets          []struct {
		URL      string `json:"url"`
		ID       int    `json:"id"`
		NodeID   string `json:"node_id"`
		Name     string `json:"name"`
		Label    string `json:"label"`
		Uploader struct {
			Login             string `json:"login"`
			ID                int    `json:"id"`
			NodeID            string `json:"node_id"`
			AvatarURL         string `json:"avatar_url"`
			GravatarID        string `json:"gravatar_id"`
			URL               string `json:"url"`
			HTMLURL           string `json:"html_url"`
			FollowersURL      string `json:"followers_url"`
			FollowingURL      string `json:"following_url"`
			GistsURL          string `json:"gists_url"`
			StarredURL        string `json:"starred_url"`
			SubscriptionsURL  string `json:"subscriptions_url"`
			OrganizationsURL  string `json:"organizations_url"`
			ReposURL          string `json:"repos_url"`
			EventsURL         string `json:"events_url"`
			ReceivedEventsURL string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"uploader"`
		ContentType        string    `json:"content_type"`
		State              string    `json:"state"`
		Size               int       `json:"size"`
		DownloadCount      int       `json:"download_count"`
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_at"`
		BrowserDownloadURL string    `json:"browser_download_url"`
	} `json:"assets"`
	TarballURL string `json:"tarball_url"`
	ZipballURL string `json:"zipball_url"`
	Body       string `json:"body"`
}

func GetLsIOMMU() {
	// Check the API for releases
	resp, err := http.Get("https://api.github.com/repos/hikariknight/ls-iommu/releases/latest")
	errorcheck.ErrorCheck(err)

	// Close the response when function ends
	defer resp.Body.Close()

	// Get the response body
	body, err := io.ReadAll(resp.Body)
	errorcheck.ErrorCheck(err)

	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Cant decode JSON")
	}

	// Make the directory for ls-iommu if it does not exist
	path := "utils"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		errorcheck.ErrorCheck(err)
	}

	// Generate the download url
	downloadUrl := fmt.Sprintf(
		"https://github.com/HikariKnight/ls-iommu/releases/download/%s/ls-iommu_%s_Linux_x86_64.tar.gz",
		result.TagName,
		result.TagName,
	)

	// Create blank file
	fileName := fmt.Sprintf("%s/ls-iommu_Linux_x86_64.tar.gz", path)
	grabClient := grab.NewClient()
	req, _ := grab.NewRequest(fileName, downloadUrl)

	// Download ls-iommu
	download := grabClient.Do(req)

	// check for errors
	if err := download.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Download saved to ./%v \n", download.Filename)

	r, err := os.Open(fileName)
	errorcheck.ErrorCheck(err)
	err = untar.Untar(fmt.Sprintf("%s/", path), r)
	errorcheck.ErrorCheck(err)

}
