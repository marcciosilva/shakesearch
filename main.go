package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"pulley.com/shakesearch/internal/html"
	"pulley.com/shakesearch/internal/search"
)

const (
	defaultAppHttpPort                         = "3001"
	portEnvVariableKey                         = "PORT"
	searchQueryParamKey                        = "q"
	shakespeareCompleteWorksPathEnvVariableKey = "SHAKESPEARE_COMPLETE_WORKS_PATH"
)

func main() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get working directory: %v", err)
	}
	shakespeareCompleteWorksPath := fmt.Sprintf("%s/internal/search/resources/completeworks.txt", workingDirectory)
	err = os.Setenv(shakespeareCompleteWorksPathEnvVariableKey, shakespeareCompleteWorksPath)
	if err != nil {
		log.Fatalf("failed to set env variable for Shakespeare's complete works' path: %v", err)
	}
	searcher, err := search.CreateNewSearcher(shakespeareCompleteWorksPathEnvVariableKey)
	if err != nil {
		log.Fatalf("failed to create searcher: %v", err)
	}
	fileServer := http.FileServer(http.Dir("./static"))

	handleApiRoutes(searcher, fileServer)
	startApp()
}

func handleApiRoutes(searcher search.Searcher, fileServer http.Handler) {
	http.HandleFunc("/search", getSearchHandler(searcher))
	http.Handle("/", fileServer)
}

func getSearchHandler(searcher search.Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query, ok := r.URL.Query()[searchQueryParamKey]
		if !ok || len(query[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			writeBytesToResponseWriter(w, []byte("missing search query in URL params"))
			return
		}
		searchedText := strings.Split(query[0], " ")
		// TODO: create package for text manipulation, trimming and whatnot

		resultsBySearchedToken := make(map[string][]string)
		resultsMutex := &sync.RWMutex{}
		wg := sync.WaitGroup{}
		// TODO: support matching entire search?
		//wg.Add(1)
		//go getResultsForSearchedTextToken(searcher, query[0], resultsBySearchedToken, resultsMutex, &wg)
		for _, searchedTextToken := range searchedText {
			wg.Add(1)
			go getResultsForSearchedTextToken(searcher, searchedTextToken, resultsBySearchedToken, resultsMutex, &wg)
		}
		wg.Wait()

		results := make([]string, 0)
		for _, searchedTextToken := range searchedText {
			results = append(results, resultsBySearchedToken[searchedTextToken]...)
		}
		writeSearchResults(w, results)
	}
}

func getResultsForSearchedTextToken(searcher search.Searcher, searchedTextToken string,
	resultsBySearchedToken map[string][]string, resultsMutex *sync.RWMutex, wg *sync.WaitGroup) {

	defer wg.Done()
	resultsByMatchedToken, prioritizedMatchingTokens := searcher.Search(searchedTextToken)
	results := html.AdaptTextForHTML(resultsByMatchedToken, prioritizedMatchingTokens)
	resultsMutex.Lock()
	resultsBySearchedToken[searchedTextToken] = results
	resultsMutex.Unlock()
}

func writeSearchResults(w http.ResponseWriter, results []string) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	err := enc.Encode(results)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeBytesToResponseWriter(w, []byte("encoding failure"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	writeBytesToResponseWriter(w, buf.Bytes())
}

func writeBytesToResponseWriter(w http.ResponseWriter, bytes []byte) {
	_, err := w.Write(bytes)
	if err != nil {
		log.Fatalf("failed to write to response writer: %v", err)
	}
}

func startApp() {
	port := os.Getenv(portEnvVariableKey)
	if port == "" {
		port = defaultAppHttpPort
	}
	fmt.Printf("Listening on port %s...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
