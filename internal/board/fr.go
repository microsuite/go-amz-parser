package board

import (
	"github.com/antchfx/htmlquery"
	"github.com/microsuite/go-amz-parser/errors"
	"github.com/microsuite/go-amz-parser/utils"
	"golang.org/x/net/html"
)

type FRBoardParser struct{}

// ParseNextPageURL parses the next page reference from the html document.
func (p *DEBoardParser) ParseNextPageURL(doc *html.Node) (string, error) {
	expr := `//li/a[contains(text(), "Next page") and string-length(@href) > 0]`
	nodes, err := utils.FindNodes(doc, expr, false)
	if err == nil && len(nodes) > 0 {
		nextRef := htmlquery.SelectAttr(nodes[0], "href")
		return nextRef, nil
	}
	return "unknown", errors.ErrorNotFoundNextPage
}
