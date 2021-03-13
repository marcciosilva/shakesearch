package html

import (
	"reflect"
	"testing"
)

func TestAdaptTextForHTML(t *testing.T) {
	type args struct {
		searchText string
		text       []string
	}
	tests := []struct {
		name           string
		args           args
		expectedResult string
	}{
		{
			name: "Text is formatted as expected into HTML",
			args: args{
				searchText: "test",
				text:       []string{"this is a totally ordinary\r\ntest text to test"},
			},
			expectedResult: "<p> this is a totally ordinary <br/> <mark>test</mark> text to <mark>test</mark> </p>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AdaptTextForHTML(tt.args.searchText, tt.args.text)
			if reflect.DeepEqual(tt.expectedResult, tt.args.text) {
				t.Errorf("wanted %v, got %v", tt.expectedResult, tt.args.text)
			}
		})
	}
}
