package main

import (
	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

func queryAll(n *html.Node, query string) []*html.Node {
	sel, err := cascadia.Parse(query)
	if err != nil {
		return []*html.Node{}
	}
	return cascadia.QueryAll(n, sel)
}
