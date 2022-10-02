package utils

import (
	"fmt"
	"net/url"
	"strings"
)

type FlagURL struct {
	URLs []*url.URL
}

func (f *FlagURL) String() string {
	return fmt.Sprint(f.URLs)
}

func (f *FlagURL) Set(s string) error {
	urls := strings.Split(s, ",")
	for _, item := range urls {
		parsedURL, err := url.ParseRequestURI(item)
		if err != nil {
			return err
		}
		f.URLs = append(f.URLs, parsedURL)
	}
	return nil
}
