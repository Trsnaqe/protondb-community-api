package background_services

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/google/go-github/github"
)

const (
	repoURL    = "https://github.com/bdefore/protondb-data/tree/master/reports"
	rawBaseURL = "https://raw.githubusercontent.com/bdefore/protondb-data/master"
)

func GetLatestProcessedReportFile(lastProcessedFile string) ([]byte, string, string, error) {
	// Get the list of files in the reports directory
	fileList, err := getFileList(lastProcessedFile)
	if err != nil || len(fileList) == 0 {
		log.Println("Error getting file list:", err)
		return nil, "", "", err
	}

	// Find the oldest file that is newer than the last processed file
	var oldestFile string
	if lastProcessedFile == "" {
		// If lastProcessedFile is empty, get the oldest file in the fileList
		oldestFile = fileList[0]
	} else {
		for _, file := range fileList {
			// Check if the file is newer than the last processed file
			if compareFiles(lastProcessedFile, file) {
				oldestFile = file
				break
			}
		}

		// Check if an older file was found
		if oldestFile == "" {
			log.Println("No new file to download. Using the last processed file:", lastProcessedFile)
			return nil, lastProcessedFile, "", nil
		}
	}

	fileURL := fmt.Sprintf("%s/%s", rawBaseURL, oldestFile)
	log.Println("Downloading file:", fileURL)
	resp, err := http.Get(fileURL)
	if err != nil {
		log.Println("Error downloading file:", err)
		return nil, lastProcessedFile, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error downloading file. Status code: %d\n", resp.StatusCode)
		return nil, lastProcessedFile, "", fmt.Errorf("failed to download file. Status code: %d", resp.StatusCode)
	}

	// Create a temporary file to save the downloaded .tar.gz file
	tempFile, err := ioutil.TempFile("", "report")
	if err != nil {
		log.Println("Error creating temporary file:", err)
		return nil, lastProcessedFile, "", err
	}
	defer tempFile.Close()

	// Save the downloaded .tar.gz file to the temporary file
	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		log.Println("Error saving downloaded file:", err)
		return nil, lastProcessedFile, "", err
	}

	// Extract the .tar.gz file
	extractedDir, err := extractTarGz(tempFile.Name())
	if err != nil {
		log.Println("Error extracting .tar.gz file:", err)
		return nil, lastProcessedFile, "", err
	}

	// Remove the .tar.gz file after extraction
	err = os.Remove(tempFile.Name())
	if err != nil {
		log.Println("Error removing .tar.gz file:", err)
	}

	// Find the JSON file within the extracted directory
	jsonFilePath, err := findJSONFile(extractedDir)
	if err != nil {
		log.Println("Error finding JSON file:", err)
		return nil, lastProcessedFile, extractedDir, err
	}

	// Read the JSON file contents
	jsonData, err := ioutil.ReadFile(jsonFilePath)
	if err != nil {
		log.Println("Error reading JSON file:", err)
		return nil, lastProcessedFile, extractedDir, err
	}

	log.Println("Processed JSON file:", jsonFilePath)
	return jsonData, oldestFile, extractedDir, nil
}

func getFileList(latestProcessedFile string) ([]string, error) {
	client := github.NewClient(nil)
	ctx := context.Background()

	tree, _, err := client.Git.GetTree(ctx, "bdefore", "protondb-data", "master", true)
	if err != nil {
		return nil, err
	}

	var fileList []string
	for _, entry := range tree.Entries {
		if entry.GetType() == "blob" && strings.HasPrefix(entry.GetPath(), "reports/") && strings.HasSuffix(entry.GetPath(), ".tar.gz") {
			fileList = append(fileList, entry.GetPath())
		}
	}

	// Sort the file list by month and year
	sort.Slice(fileList, func(i, j int) bool {
		return compareFiles(fileList[i], fileList[j])
	})

	novFile := "reports/reports_nov1_2019.tar.gz"

	// If latestProcessedFile is older than nov1_2019, return nov1_2019
	if compareFiles(latestProcessedFile, novFile) {
		return []string{novFile}, nil
	}

	// If latestProcessedFile is not older than nov1_2019, return the newest file
	newestFile := fileList[len(fileList)-1]
	return []string{newestFile}, nil
}

// extractTarGz extracts the .tar.gz archive and returns the path to the extracted directory
func extractTarGz(tarGzFile string) (string, error) {
	// Open the .tar.gz file for reading
	file, err := os.Open(tarGzFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Create a gzip reader
	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return "", err
	}
	defer gzReader.Close()

	// Create a tar reader
	tarReader := tar.NewReader(gzReader)

	// Create a temporary directory for extraction
	extractedDir, err := ioutil.TempDir("", "extracted")
	if err != nil {
		return "", err
	}

	// Extract each file in the .tar.gz archive
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // Reached the end of the archive
		}
		if err != nil {
			log.Println("Error reading tar entry:", err)
			continue // Skip this entry and continue with the next one
		}

		// Determine the file path for extraction
		targetPath := filepath.Join(extractedDir, header.Name)

		// Create directories as needed
		if header.Typeflag == tar.TypeDir {
			err := os.MkdirAll(targetPath, 0755)
			if err != nil {
				log.Println("Error creating directory:", err)
			}
			continue
		}

		// Create the file and write its contents
		file, err := os.Create(targetPath)
		if err != nil {
			log.Println("Error creating file:", err)
			continue // Skip this file and continue with the next one
		}

		_, err = io.Copy(file, tarReader)
		file.Close() // Close the file after copying

		if err != nil {
			log.Println("Error writing file:", err)
			// If there was an error writing the file, delete it to avoid partial extraction
			if err := os.Remove(targetPath); err != nil {
				log.Println("Error removing partial file:", err)
			}
		}
	}
	return extractedDir, nil
}

// findJSONFile finds the JSON file in the specified directory and returns its path.
func findJSONFile(directory string) (string, error) {
	var jsonFilePath string

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".json") {
			jsonFilePath = path
			return io.EOF // Stop the walk
		}

		return nil
	})

	if err == io.EOF {
		return jsonFilePath, nil
	}

	return "", fmt.Errorf("JSON file not found")
}
