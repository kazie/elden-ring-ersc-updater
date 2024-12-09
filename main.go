package main

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/google/go-github/v67/github"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

type Config struct {
	CurrentVersion    string `toml:"current_version"`
	EldenRingGamePath string `toml:"elden_ring_game_path"`
	GithubReadToken   string `toml:"github_read_token"`
	IgnoreIniFile     bool   `toml:"ignore_ini_file"`
}

func readConfig() *Config {
	var conf Config
	if _, err := toml.DecodeFile("./config.toml", &conf); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Could not find or read the config file, [./config.toml\n %s", err)
		os.Exit(1)
	}
	return &conf
}

func getLatestVersion(client *github.Client) *github.RepositoryRelease {
	release, _, err := client.Repositories.GetLatestRelease(context.Background(), "LukeYui", "EldenRingSeamlessCoopRelease")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not fetch latest version, [%s]\n", err)
		os.Exit(1)
	}
	return release
}

func getZipFile(zipFileURL string, auth string) []byte {
	client := &http.Client{}
	req, err := http.NewRequest("GET", zipFileURL, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create GET request, [%s]\n", err)
		os.Exit(1)
	}
	req.Header.Set("Authorization", "Bearer "+auth)
	req.Header.Set("Accept", "application/zip")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not download zip file, [%s]\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading response body, [%s]\n", err)
		os.Exit(1)
	}

	return data
}

func unzipDataIntoFolder(data []byte, outputDir string, ignoreIni bool) {
	reader := bytes.NewReader(data)
	zipReader, err := zip.NewReader(reader, int64(len(data)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading zip file, [%s]\n", err)
		os.Exit(1)
	}

	for _, file := range zipReader.File {
		// don't want to override settings, not handling writing for first time.
		if ignoreIni && file.Name == "SeamlessCoop/ersc_settings.ini" {
			fmt.Printf("Ignoring [%s] from zip file\n", file.Name)
			continue
		}

		fullPath := filepath.Join(outputDir, file.Name)
		directory := filepath.Dir(fullPath)
		err := os.MkdirAll(directory, os.ModePerm)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not create directories for %s [%s]\n", directory, err)
			os.Exit(1)
		}

		fileReader, err := file.Open()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not open file %s [%s]\n", file.Name, err)
			os.Exit(1)
		}

		outFile, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not open file %s [%s]\n", fullPath, err)
			os.Exit(1)
		}

		_, err = io.Copy(outFile, fileReader)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not write file %s [%s]\n", fullPath, err)
			os.Exit(1)
		}
		outFile.Close()
		fileReader.Close() // close reader after using it
	}

}

func updateVersionInFile(version string, filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read file %s [%s]\n", filename, err)
		os.Exit(1)
	}

	re := regexp.MustCompile(`current_version\s*=\s*".*"`)
	newContent := re.ReplaceAllString(string(content), `current_version = "`+version+`"`)

	err = os.WriteFile(filename, []byte(newContent), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not update file %s with latest version [%s]\n", filename, err)
		os.Exit(1)
	}
}

func main() {
	conf := readConfig()
	client := github.NewClient(nil).WithAuthToken(conf.GithubReadToken)
	release := getLatestVersion(client)

	if *release.TagName != conf.CurrentVersion {
		fmt.Printf("New version available! Upgrade: [%s → %s]\n", conf.CurrentVersion, *release.TagName)
		if len(release.Assets) == 0 {
			fmt.Fprintf(os.Stderr, "No assets for release [%s] is missing\n", *release.TagName)
			os.Exit(1)
		}
		if *release.Assets[0].Name != "ersc.zip" {
			fmt.Fprintf(os.Stderr, "Only expecting one asset named ersc.zip, but found [%s] instead\n", *release.Assets[0].Name)
			os.Exit(1)
		}
		if *release.Assets[0].BrowserDownloadURL == "" {
			fmt.Fprintf(os.Stderr, "Missing download URL for ersc.zip \n")
			os.Exit(1)
		}
		fmt.Printf("Downloading zip from: %s\n", *release.ZipballURL)
		data := getZipFile(*release.Assets[0].BrowserDownloadURL, conf.GithubReadToken)
		unzipDataIntoFolder(data, conf.EldenRingGamePath, conf.IgnoreIniFile)
		updateVersionInFile(*release.TagName, "./config.toml")
		fmt.Printf("Updated to latest version [%s]\n", *release.TagName)
	} else {
		fmt.Printf("Up to date! [%s]\n", conf.CurrentVersion)
		os.Exit(0)
	}

}
