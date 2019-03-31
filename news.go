package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
)

// constants
const (
	newsURL = "https://newsapi.org/v2/top-headlines?apiKey=f99aa135983b46be95358b8d9da1018e"
)

type News struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
	UrlToImage  string `json:"urlToImage"`
	PublishedAt string `json:"publishedAt"`
	Content     string `json:"content"`
	Source      Source `json:"source"`
}

type Source struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Articles struct {
	Status       string
	TotalResults int
	Articles     []News
}

func getNews(name, category string) Articles {
	var cat string = ""
	// send GET request to GitHub API with the requested user "name"
	if category != "" {
		cat = "&category=" + category
	}
	resp, err := http.Get(newsURL + "&country=" + name + cat)
	// if err occurs during GET request, then throw error and quit application
	check(err)

	// Always good practice to defer closing the response body.
	// If application crashes or function finishes successfully, GO will always execute this "defer" statement
	defer resp.Body.Close()

	// read the response body and handle any errors during reading.
	body, err := ioutil.ReadAll(resp.Body)
	check(err)
	// fmt.Println(string(body))

	// create a user variable of type "User" struct to store the "Unmarshal"-ed (aka parsed JSON) data, then return the user
	var news Articles
	json.Unmarshal(body, &news)
	// json.Unmarshal(&news.Articles, &news.Articles)
	// fmt.Println(news.Articles)
	return news
}

func DisplayNews(news, category, x string) {
	fmt.Printf("Getting %s news: %s\n", category, news)
	results := getNews(news, category)

	var size int
	var err error
	if x != "" {
		if size, err = strconv.Atoi(x); err != nil {
			check(err)
		}
	} else {
		size = 80
	}

	var wg sync.WaitGroup
	for _, res := range results.Articles {
		wg.Add(1)
		go func() {
			fmt.Println("**********************************************************")
			fmt.Println(`Source:             `, res.Source.Name)
			fmt.Println(`Publishing date:    `, res.PublishedAt)
			fmt.Println(`Title:              `, res.Title)
			// fmt.Println(`Description:        `, res.Description)
			fmt.Println(`Content:            `, res.Content)
			fmt.Println(`Url:                `, res.Url)
			// fmt.Println(`UrlToImage:         `, res.UrlToImage)
			fmt.Println()
			if res.UrlToImage != "" {
				asciiArt := Convert2Ascii(res.UrlToImage, size)
				fmt.Println(string(asciiArt))
			}
			wg.Done()
		}()
		wg.Wait()
	}
}
