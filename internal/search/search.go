package search

import (
	"fmt"
	"index/suffixarray"
	"io/ioutil"
	"os"

	"pulley.com/shakesearch/internal/math"
)

const shakespeareCompleteWorksFilename = "resources/completeworks.txt"

type Searcher interface {
	Load(filename string) error
	Search(query string) []string
}

type ShakespeareSearcher struct {
	CompleteWorks string
	SuffixArray   *suffixarray.Index
}

func (s *ShakespeareSearcher) Load(filename string) error {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	s.CompleteWorks = string(dat)
	s.SuffixArray = suffixarray.New(dat)
	return nil
}

func (s *ShakespeareSearcher) Search(query string) []string {
	var results []string
	idxs := s.SuffixArray.Lookup([]byte(query), -1)
	completeWorksLength := len(s.CompleteWorks)
	for _, idx := range idxs {
		maxTextIndex := math.Min(completeWorksLength, idx+250)
		minTextIndex := math.Max(idx-250, 0)
		excerpt := s.CompleteWorks[minTextIndex:maxTextIndex]
		results = append(results, excerpt)
	}

	return results
}

func CreateNewSearcher(shakespeareCompleteWorksPathEnvVariableKey string) (*ShakespeareSearcher, error) {
	searcher := &ShakespeareSearcher{}
	searcherSourcePath := os.Getenv(shakespeareCompleteWorksPathEnvVariableKey)
	err := searcher.Load(searcherSourcePath)
	if err != nil {
		return nil, err
	}
	return searcher, nil
}
