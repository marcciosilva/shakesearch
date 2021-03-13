package search

import (
	"fmt"
	"index/suffixarray"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"sync"
	"unicode"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"pulley.com/shakesearch/internal/math"
)

const (
	shakespeareCompleteWorksFilename = "resources/completeworks.txt"
	invisibleUnicodeCharacterCutSet  = "\uFEFF\u200B\u200D\u200C"
)

type Searcher interface {
	Load(filename string) error
	Search(query string) (map[string][]string, []string)
}

type ShakespeareSearcher struct {
	CompleteWorks string
	SuffixArray   *suffixarray.Index
	Tokens        []string
}

func (s *ShakespeareSearcher) Load(filename string) error {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	s.CompleteWorks = string(dat)
	tokenizedCompleteWorks := strings.Split(
		strings.Replace(s.CompleteWorks, "\r\n", " ", -1), " ",
	)
	s.Tokens = tokenizedCompleteWorks
	s.SuffixArray = suffixarray.New(dat)
	return nil
}

func (s *ShakespeareSearcher) Search(query string) (map[string][]string, []string) {
	prioritizedTokens := s.getPrioritizedFuzzyMatchesForSearchQuery(query)

	resultsByToken := make(map[string][]string)

	wg := sync.WaitGroup{}
	resultsMutex := sync.RWMutex{}
	for _, tokenToSearchFor := range prioritizedTokens {
		wg.Add(1)
		go s.writeMatchingExcerptsForTokenToResult(tokenToSearchFor, resultsByToken, &wg, &resultsMutex)
	}
	wg.Wait()

	return resultsByToken, prioritizedTokens
}

func (s *ShakespeareSearcher) getPrioritizedFuzzyMatchesForSearchQuery(query string) []string {
	// matches the query (case insensitive) to a token in the text only if said token includes all of the token's letters.
	// then ranks matches using Levenshtein distance.
	matches := fuzzy.RankFindFold(query, s.Tokens)
	sort.Sort(matches)
	wasMatchUsedByToken := make(map[string]bool)
	prioritizedTokens := make([]string, 0)
	for _, match := range matches {
		tokenWithoutInvisibleCharacters := strings.Trim(match.Target, invisibleUnicodeCharacterCutSet)
		_, wasMatchUsed := wasMatchUsedByToken[tokenWithoutInvisibleCharacters]
		if !wasMatchUsed {
			wasMatchUsedByToken[tokenWithoutInvisibleCharacters] = true
			prioritizedTokens = append(prioritizedTokens, tokenWithoutInvisibleCharacters)
		}
	}
	return prioritizedTokens
}

func (s *ShakespeareSearcher) writeMatchingExcerptsForTokenToResult(tokenToSearchFor string, resultsByToken map[string][]string,
	wg *sync.WaitGroup, resultsMutex *sync.RWMutex) {

	defer wg.Done()
	tokenAppearanceIndexes := s.SuffixArray.Lookup([]byte(tokenToSearchFor), -1)
	completeWorksLength := len(s.CompleteWorks)
	resultsForToken := make([]string, 0)
	for _, idx := range tokenAppearanceIndexes {
		maxTextIndex := math.Min(completeWorksLength, idx+250)
		minTextIndex := math.Max(idx-250, 0)
		excerpt := s.CompleteWorks[minTextIndex:maxTextIndex]
		for index, rune := range excerpt {
			if unicode.IsUpper(rune) {
				excerpt = "..." + excerpt[index:] + "..."
				break
			}
		}
		resultsForToken = append(resultsForToken, excerpt)
	}
	resultsMutex.Lock()
	resultsByToken[tokenToSearchFor] = resultsForToken
	resultsMutex.Unlock()
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
