package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// You can specify a tag as a command line argument to generate the changelog for a specific version.
// Example: go run tools/changelog/changelog.go v0.0.33
// If no tag is provided, the latest release will be used.

// Setting repo owner and repo name by generate changelog
const (
	repoOwner = "openimsdk"
	repoName  = "open-im-server"
)

// GitHubRepo struct represents the repo details.
type GitHubRepo struct {
	Owner         string
	Repo          string
	FullChangelog string
}

// ReleaseData represents the JSON structure for release data.
type ReleaseData struct {
	TagName   string `json:"tag_name"`
	Body      string `json:"body"`
	HtmlUrl   string `json:"html_url"`
	Published string `json:"published_at"`
}

// Method to classify and format release notes.
func (g *GitHubRepo) classifyReleaseNotes(body string) map[string][]string {
	result := map[string][]string{
		"feat":     {},
		"fix":      {},
		"chore":    {},
		"refactor": {},
		"build":    {},
		"other":    {},
	}

	// Regular expression to extract PR number and URL (case insensitive)
	rePR := regexp.MustCompile(`(?i)in (https://github\.com/[^\s]+/pull/(\d+))`)

	// Split the body into individual lines.
	lines := strings.Split(body, "\n")

	for _, line := range lines {
		// Skip lines that contain "deps: Merge"
		if strings.Contains(strings.ToLower(line), "deps: merge #") {
			continue
		}

		// Use a regular expression to extract Full Changelog link and its title (case insensitive).
		if strings.Contains(strings.ToLower(line), "**full changelog**") {
			matches := regexp.MustCompile(`(?i)\*\*full changelog\*\*: (https://github\.com/[^\s]+/compare/([^\s]+))`).FindStringSubmatch(line)
			if len(matches) > 2 {
				// Format the Full Changelog link with title
				g.FullChangelog = fmt.Sprintf("[%s](%s)", matches[2], matches[1])
			}
			continue // Skip further processing for this line.
		}

		if strings.HasPrefix(line, "*") {
			var category string

			// Use strings.ToLower to make the matching case insensitive
			lowerLine := strings.ToLower(line)

			// Determine the category based on the prefix (case insensitive).
			if strings.HasPrefix(lowerLine, "* feat") {
				category = "feat"
			} else if strings.HasPrefix(lowerLine, "* fix") {
				category = "fix"
			} else if strings.HasPrefix(lowerLine, "* chore") {
				category = "chore"
			} else if strings.HasPrefix(lowerLine, "* refactor") {
				category = "refactor"
			} else if strings.HasPrefix(lowerLine, "* build") {
				category = "build"
			} else {
				category = "other"
			}

			// Extract PR number and URL (case insensitive)
			matches := rePR.FindStringSubmatch(line)
			if len(matches) == 3 {
				prURL := matches[1]
				prNumber := matches[2]
				// Format the line with the PR link and use original content for the final result
				formattedLine := fmt.Sprintf("* %s [#%s](%s)", strings.Split(line, " by ")[0][2:], prNumber, prURL)
				result[category] = append(result[category], formattedLine)
			} else {
				// If no PR link is found, just add the line as is
				result[category] = append(result[category], line)
			}
		}
	}

	return result
}

// Method to generate the final changelog.
func (g *GitHubRepo) generateChangelog(tag, date, htmlURL, body string) string {
	sections := g.classifyReleaseNotes(body)

	// Convert ISO 8601 date to simpler format (YYYY-MM-DD)
	formattedDate := date[:10]

	// Changelog header with tag, date, and links.
	changelog := fmt.Sprintf("## [%s](%s) \t(%s)\n\n", tag, htmlURL, formattedDate)

	if len(sections["feat"]) > 0 {
		changelog += "### New Features\n" + strings.Join(sections["feat"], "\n") + "\n\n"
	}
	if len(sections["fix"]) > 0 {
		changelog += "### Bug Fixes\n" + strings.Join(sections["fix"], "\n") + "\n\n"
	}
	if len(sections["chore"]) > 0 {
		changelog += "### Chores\n" + strings.Join(sections["chore"], "\n") + "\n\n"
	}
	if len(sections["refactor"]) > 0 {
		changelog += "### Refactors\n" + strings.Join(sections["refactor"], "\n") + "\n\n"
	}
	if len(sections["build"]) > 0 {
		changelog += "### Builds\n" + strings.Join(sections["build"], "\n") + "\n\n"
	}
	if len(sections["other"]) > 0 {
		changelog += "### Others\n" + strings.Join(sections["other"], "\n") + "\n\n"
	}

	if g.FullChangelog != "" {
		changelog += fmt.Sprintf("**Full Changelog**: %s\n", g.FullChangelog)
	}

	return changelog
}

// Method to fetch release data from GitHub API.
func (g *GitHubRepo) fetchReleaseData(version string) (*ReleaseData, error) {
	var apiURL string

	if version == "" {
		// Fetch the latest release.
		apiURL = fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", g.Owner, g.Repo)
	} else {
		// Fetch a specific version.
		apiURL = fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/tags/%s", g.Owner, g.Repo, version)
	}

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var releaseData ReleaseData
	err = json.Unmarshal(body, &releaseData)
	if err != nil {
		return nil, err
	}

	return &releaseData, nil
}

func main() {
	repo := &GitHubRepo{Owner: repoOwner, Repo: repoName}

	// Get the version from command line arguments, if provided
	var version string // Default is use latest

	if len(os.Args) > 1 {
		version = os.Args[1] // Use the provided version
	}

	// Fetch release data (either for latest or specific version)
	releaseData, err := repo.fetchReleaseData(version)
	if err != nil {
		fmt.Println("Error fetching release data:", err)
		return
	}

	// Generate and print the formatted changelog
	changelog := repo.generateChangelog(releaseData.TagName, releaseData.Published, releaseData.HtmlUrl, releaseData.Body)
	fmt.Println(changelog)
}
