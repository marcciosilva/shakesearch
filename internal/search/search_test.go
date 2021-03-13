package search

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

const shakespeareCompleteWorksPathEnvVariableKey = "TEST_SHAKESPEARE_COMPLETE_WORKS_PATH"

func TestCreateNewSearcher(t *testing.T) {
	type args struct {
		shakespeareCompleteWorksPathEnvVariableKey string
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			name: "Non-nil searcher is returned when filepath in environment variable can be loaded, error is nil",
			args: args{
				shakespeareCompleteWorksPathEnvVariableKey: shakespeareCompleteWorksPathEnvVariableKey,
			},
			wantNil: false,
			wantErr: false,
		},
		{
			name: "Error is returned when filepath in environment variable cannot be loaded, searcher is nil",
			args: args{
				shakespeareCompleteWorksPathEnvVariableKey: "",
			},
			wantNil: true,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loadTestCompleteWorks(t)
			defer cleanUpEnvVariables(t)
			got, err := CreateNewSearcher(tt.args.shakespeareCompleteWorksPathEnvVariableKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateNewSearcher() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantNil != (got == nil) {
				t.Errorf("CreateNewSearcher() got = %v, wantNil %v", got, tt.wantNil)
			}
		})
	}
}

func loadTestCompleteWorks(t *testing.T) {
	workingDirectory, err := os.Getwd()
	if err != nil {
		t.Errorf("failed to get working directory: %v", err)
	}
	testCompleteWorksPath := fmt.Sprintf("%s/resources/test_completeworks.txt", workingDirectory)
	err = os.Setenv(shakespeareCompleteWorksPathEnvVariableKey, testCompleteWorksPath)
	if err != nil {
		t.Errorf("failed to get working directory: %v", err)
	}
}

func cleanUpEnvVariables(t *testing.T) {
	err := os.Unsetenv(shakespeareCompleteWorksFilename)
	if err != nil {
		t.Errorf("failed to clean up env variable %s: %v", shakespeareCompleteWorksFilename, err)
	}
}

func TestShakespeareSearcher_Search(t *testing.T) {
	loadTestCompleteWorks(t)
	defer cleanUpEnvVariables(t)
	searcher, err := CreateNewSearcher(shakespeareCompleteWorksPathEnvVariableKey)
	if err != nil {
		t.Errorf("failed to create new searcher: %v", err)
	}
	expectedResults := map[string][]string{
		"majesty": {"to his new-appearing sight,\r\nServing with looks his sacred majesty,\r\nAnd having climbed the steep-up heavenly hill,\r\nResembling strong youth in his middle age,\r\nYet mortal looks adore his beauty still,\r\nAttending on his golden pilgrimage:\r\nBut when from highmost pitch with weary car,\r\nLike feeble age he reeleth from the day,\r\nThe eyes (fore duteous) now converted are\r\nFrom his low tract and look another way:\r\nSo thou, thy self out-going in thy noon:\r\nUnlooked on diest unless thou get a son.\n"},
		"majsty":  {"sty,\r\nAnd having climbed the steep-up heavenly hill,\r\nResembling strong youth in his middle age,\r\nYet mortal looks adore his beauty still,\r\nAttending on his golden pilgrimage:\r\nBut when from highmost pitch with weary car,\r\nLike feeble age he reeleth from the day,\r\nThe eyes (fore duteous) now converted are\r\nFrom his low tract and look another way:\r\nSo thou, thy self out-going in thy noon:\r\nUnlooked on diest unless thou get a son."},
	}

	resultsByMatchedToken, _ := searcher.Search("from")

	for token, result := range resultsByMatchedToken {
		expectedResultsForToken := expectedResults[token]
		trimCutSet := " \n\r"
		for i := 0; i < len(expectedResultsForToken); i++ {
			if strings.Trim(result[i], trimCutSet) != strings.Trim(expectedResultsForToken[i], trimCutSet) {
				t.Errorf("expected %v, got %v", expectedResultsForToken, result)
			}
		}
	}
}
