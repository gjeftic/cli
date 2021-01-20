// "Package main" is the namespace declaration
package main

// importing standard libraries
import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// constants
const (
	githugAPIURL = "https://api.github.com"
	userEndpoint = "/users/"
	toRepos      = "/repos"
)

// User struct represents the JSON data from GitHub API: https://api.github.com/users/defunct
// This struct was generated via a JSON-to-GO utility by Matt Holt: https://mholt.github.io/json-to-go/
type User struct {
	Login             string      `json:"login"`
	ID                int         `json:"id"`
	AvatarURL         string      `json:"avatar_url"`
	GravatarID        string      `json:"gravatar_id"`
	URL               string      `json:"url"`
	HTMLURL           string      `json:"html_url"`
	FollowersURL      string      `json:"followers_url"`
	FollowingURL      string      `json:"following_url"`
	GistsURL          string      `json:"gists_url"`
	StarredURL        string      `json:"starred_url"`
	SubscriptionsURL  string      `json:"subscriptions_url"`
	OrganizationsURL  string      `json:"organizations_url"`
	ReposURL          string      `json:"repos_url"`
	EventsURL         string      `json:"events_url"`
	ReceivedEventsURL string      `json:"received_events_url"`
	Type              string      `json:"type"`
	SiteAdmin         bool        `json:"site_admin"`
	Name              string      `json:"name"`
	Company           string      `json:"company"`
	Blog              string      `json:"blog"`
	Location          string      `json:"location"`
	Email             string      `json:"email"`
	Hireable          interface{} `json:"hireable"`
	Bio               string      `json:"bio"`
	PublicRepos       int         `json:"public_repos"`
	PublicGists       int         `json:"public_gists"`
	Followers         int         `json:"followers"`
	Following         int         `json:"following"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
	Stats             StatsGithub `json:"stats,omitempty"`
}

// Repo struct represents the JSON data
type Repo struct {
	Name        string `json:"name"`
	Private     bool   `json:"private"`
	HTMLURL     string `json:"html_url"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	GitURL      string `json:"git_url"`
	Size        int    `json:"size"`
	Language    string `json:"language"`
	// Open_issues_count int    `json:"open_issues_count"`
	// Forks             int    `json:"forks"`
	// Watchers          int    `json:"watchers"`
	DefaultBranch string `json:"default_branch"`
	ID            int    `json:"id"`
}

// Repos struct
type Repos struct {
	Repos []Repo
}

// StatsGithub map data
type StatsGithub map[string]int

// getUsers queries GitHub API for a given user
func getUsers(name string) User {
	// send GET request to GitHub API with the requested user "name"
	resp, err := http.Get(githugAPIURL + userEndpoint + name)
	// if err occurs during GET request, then throw error and quit application
	check(err)

	// Always good practice to defer closing the response body.
	// If application crashes or function finishes successfully, GO will always execute this "defer" statement
	defer resp.Body.Close()

	// read the response body and handle any errors during reading.
	body, err := ioutil.ReadAll(resp.Body)
	check(err)

	// create a user variable of type "User" struct to store the "Unmarshal"-ed (aka parsed JSON) data, then return the user
	var user User
	json.Unmarshal(body, &user)

	user.Stats = make(map[string]int)
	res := getRepos(name)
	total := len(res.Repos)
	for _, result := range res.Repos {
		if result.Language == "" {
			result.Language = "unknown"
		}
		user.Stats[result.Language]++
	}
	for i, agg := range user.Stats {
		user.Stats[i] = agg * 100 / total
	}
	return user
}

// getRepos queries GitHub API for a given user repositories
func getRepos(name string) Repos {
	// send GET request to GitHub API with the requested user "name"
	resp, err := http.Get(githugAPIURL + userEndpoint + name + toRepos)
	// if err occurs during GET request, then throw error and quit application
	check(err)

	// Always good practice to defer closing the response body.
	// If application crashes or function finishes successfully, GO will always execute this "defer" statement
	defer resp.Body.Close()

	// read the response body and handle any errors during reading.
	body, err := ioutil.ReadAll(resp.Body)
	check(err)

	// create a user variable of type "User" struct to store the "Unmarshal"-ed (aka parsed JSON) data, then return the user
	var repos Repos
	json.Unmarshal(body, &repos.Repos)
	return repos
}
