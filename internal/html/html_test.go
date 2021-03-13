package html

import (
	"reflect"
	"testing"
)

func TestAdaptTextForHTML(t *testing.T) {
	type args struct {
		textByToken map[string][]string
		tokens      []string
	}
	tests := []struct {
		name            string
		args            args
		expectedResults []string
	}{
		{
			name: "Text is formatted as expected into HTML",
			args: args{
				textByToken: map[string][]string{
					"test": {"this is a totally ordinary\r\ntest textByToken to test"},
					"tust": {"this\r\nis tust"},
				},
				tokens: []string{"test", "tust"},
			},
			expectedResults: []string{
				"<p> this is a totally ordinary <br/> <mark>test</mark> textByToken to <mark>test</mark> </p>",
				"<p> this <br/> is <mark>tust</mark> </p>",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := AdaptTextForHTML(tt.args.textByToken, tt.args.tokens)
			if reflect.DeepEqual(tt.expectedResults, results) {
				t.Errorf("wanted %v, got %v", tt.expectedResults, tt.args.textByToken)
			}
		})
	}
}
