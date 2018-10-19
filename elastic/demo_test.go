package elastic

import (
	"testing"
	"net/url"
)

func TestCanonicalize(t *testing.T) {
	tests := []struct {
		Input  []string
		Output []string
	}{
		// #0
		{
			Input:  []string{"http://58.222.20.252/"},
			Output: []string{"http://58.222.20.252"},
		},
		// #1
		{
			Input:  []string{"http://58.222.20.252:9200/", "gopher://golang.org/", "http://58.222.20.252:9201"},
			Output: []string{"http://58.222.20.252:9200", "http://58.222.20.252:9201"},
		},
		// #2
		{
			Input:  []string{"http://user:secret@58.222.20.252/path?query=1#fragment"},
			Output: []string{"http://user:secret@58.222.20.252/path"},
		},
		// #3
		{
			Input:  []string{"https://somewhere.on.mars:9999/path?query=1#fragment"},
			Output: []string{"https://somewhere.on.mars:9999/path"},
		},
		// #4
		{
			Input:  []string{"https://prod1:9999/one?query=1#fragment", "https://prod2:9998/two?query=1#fragment"},
			Output: []string{"https://prod1:9999/one", "https://prod2:9998/two"},
		},
		// #5
		{
			Input:  []string{"http://58.222.20.252/one/"},
			Output: []string{"http://58.222.20.252/one"},
		},
		// #6
		{
			Input:  []string{"http://58.222.20.252/one///"},
			Output: []string{"http://58.222.20.252/one"},
		},
		// #7: Invalid URL
		{
			Input:  []string{"58.222.20.252/"},
			Output: []string{},
		},
		// #8: Invalid URL
		{
			Input:  []string{"58.222.20.252:9200"},
			Output: []string{},
		},
	}

	for i, test := range tests {
		got := canonicalize(test.Input...)
		if want, have := len(test.Output), len(got); want != have {
			t.Fatalf("#%d: expected %d elements; got: %d", i, want, have)
		}
		for i := 0; i < len(got); i++ {
			if want, have := test.Output[i], got[i]; want != have {
				t.Errorf("#%d: expected %q; got: %q", i, want, have)
			}
		}
	}
}

func canonicalize(rawurls ...string) []string {
	var canonicalized []string
	for _, rawurl := range rawurls {
		u, err := url.Parse(rawurl)
		if err == nil {
			if u.Scheme == "http" || u.Scheme == "https" {
				// Trim trailing slashes
				for len(u.Path) > 0 && u.Path[len(u.Path)-1] == '/' {
					u.Path = u.Path[0 : len(u.Path)-1]
				}
				u.Fragment = ""
				u.RawQuery = ""
				canonicalized = append(canonicalized, u.String())
			}
		}
	}
	return canonicalized
}
