package getlist

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Jokerjoker91/Storage-App/handlers/signer"
)

type ListBucketResult struct {
	XMLName xml.Name `xml:"ListBucketResult"`
	Contents []Content `xml:"Contents"`
}

type Content struct {
	Key string `xml:"Key"`
}

// Folder represents the folder-file structure.
type Folder struct {
	Name     string    `json:"name"`
	Files    []string  `json:"files"`
	SubFolders []*Folder `json:"subFolders"`
}

// Handler to fetch and process object storage data
func GetBucketContentsHandler(w http.ResponseWriter, r *http.Request) {
	// Define Scaleway bucket URL and region
	url := "https://long-term-strg-app.s3.fr-par.scw.cloud/?list-type=2"
	region := "fr-par"

	// Generate the signed request
	req, err := signer.CreateSignedRequest("GET", url, region)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create signed request: %v", err), http.StatusInternalServerError)
		return
	}

	// Make the signed request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to fetch bucket contents", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	// Parse the XML response
	var result ListBucketResult
	err = xml.Unmarshal(body, &result)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse XML: %v", err), http.StatusInternalServerError)
		return
	}

	// Organize files into folder structure
	folderTree := buildFolderTree(result.Contents)

	// Serialize the folder structure to JSON
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(folderTree)
	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}
}


func buildFolderTree(contents []Content) *Folder {
	root := &Folder{Name: "Root", Files: []string{}, SubFolders: []*Folder{}}
	folderMap := map[string]*Folder{"": root} // Map for quick folder access

	for _, content := range contents {
		parts := strings.Split(content.Key, "/")
		currentFolder := root

		// Traverse or create folder structure
		for i, part := range parts {
			if i == len(parts)-1 && part != "" { // Last part is a file
				currentFolder.Files = append(currentFolder.Files, part)
				break
			}

			if part == "" {
				continue
			}

			// Check if folder exists, otherwise create it
			if _, exists := folderMap[part]; !exists {
				newFolder := &Folder{Name: part, Files: []string{}, SubFolders: []*Folder{}}
				currentFolder.SubFolders = append(currentFolder.SubFolders, newFolder)
				folderMap[part] = newFolder
			}

			currentFolder = folderMap[part]
		}
	}

	return root
}
