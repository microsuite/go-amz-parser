package board

import (
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/microsuite/go-amz-parser/errors"
	"github.com/microsuite/go-amz-parser/utils"
	"golang.org/x/net/html"
)

type FRBoardParser struct{}

func NewFRBoardParser() *FRBoardParser {
	return &FRBoardParser{}
}

// ParseAllProducts parses all products from the given HTML document.
func (p *FRBoardParser) ParseAllProducts(doc *html.Node) ([]*html.Node, error) {
	expr := `/html/body/div[@id="a-page"]//div[@data-client-recs-list and @data-reftag]`
	nodes, err := utils.FindNodes(doc, expr, false)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

// ParseNextPageURL parses the next page reference from the html document.
func (p *FRBoardParser) ParseNextPageURL(doc *html.Node) (string, error) {
	expr := `//li/a[contains(text(), "Page suivante") and string-length(@href) > 0]`
	nodes, err := utils.FindNodes(doc, expr, false)
	if err == nil && len(nodes) > 0 {
		nextRef := htmlquery.SelectAttr(nodes[0], "href")
		return nextRef, nil
	}
	return "unknown", errors.ErrorNotFoundNextPage
}

// ParseRecsList parses the recs list from the give html node.
func (p *FRBoardParser) ParseRecsList(doc *html.Node) (string, error) {
	expr := `/html/body/div[@id='a-page']//div[@data-client-recs-list and @data-reftag]`

	nodes, err := utils.FindNodes(doc, expr, true)
	if err == nil && len(nodes) > 0 {
		return htmlquery.SelectAttr(nodes[0], "data-client-recs-list"), nil
	}
	return "unknown", errors.ErrorNotFoundRecsList
}

// ParseReftag parses the ref tag from the give html node.
func (p *FRBoardParser) ParseReftag(doc *html.Node) (string, error) {
	expr := `/html/body/div[@id='a-page']//div[@data-client-recs-list and @data-reftag]`

	nodes, err := utils.FindNodes(doc, expr, true)
	if err == nil && len(nodes) > 0 {
		return htmlquery.SelectAttr(nodes[0], "data-reftag"), nil
	}
	return "unknown", errors.ErrorNotFoundReftag
}

// ParseOffset parses the offset from the give html node.
func (p *FRBoardParser) ParseOffset(doc *html.Node) (string, error) {
	expr := `/html/body/div[@id='a-page']//div[@data-client-recs-list and @data-reftag]`

	nodes, err := utils.FindNodes(doc, expr, true)
	if err == nil && len(nodes) > 0 {
		return htmlquery.SelectAttr(nodes[0], "data-index-offset"), nil
	}
	return "unknown", errors.ErrorNotFoundOffset
}

// ParseAcpParam parses the acp param from the give html node.
func (p *FRBoardParser) ParseAcpParam(doc *html.Node) (string, error) {
	expr := `//div[@data-acp-params and @data-acp-path]`

	nodes, err := utils.FindNodes(doc, expr, true)
	if err == nil && len(nodes) > 0 {
		return htmlquery.SelectAttr(nodes[0], "data-acp-params"), nil
	}
	return "unknown", errors.ErrorNotFoundAcpParam
}

// ParseAcpPath parses the acp path from the give html node.
func (p *FRBoardParser) ParseAcpPath(doc *html.Node) (string, error) {
	expr := `//div[@data-acp-params and @data-acp-path]`

	nodes, err := utils.FindNodes(doc, expr, true)
	if err == nil && len(nodes) > 0 {
		return htmlquery.SelectAttr(nodes[0], "data-acp-path"), nil
	}
	return "unknown", errors.ErrorNotFoundAcpPath
}

// ParseBestSellersCategory parses the best seller category from the give html document.
func (p *FRBoardParser) ParseBestSellersCategory(doc *html.Node) (string, error) {
	expr := `//div/div/h1[contains(text(), 'Les meilleures ventes en')]/text()`

	nodes, err := utils.FindNodes(doc, expr, true)
	if err != nil {
		return "unknown", err
	}
	if nodes[0].Data == "" {
		return "unknown", errors.ErrorNotFoundBestSellerCategory
	}
	return nodes[0].Data, nil
}

// ParseNewReleasesCategory parses the new release category from the give html document.
func (p *FRBoardParser) ParseNewReleasesCategory(doc *html.Node) (string, error) {
	expr := `//div/div/h1[contains(text(), 'Dernières nouveautés en')]/text()`

	nodes, err := utils.FindNodes(doc, expr, true)
	if err != nil {
		return "unknown", err
	}
	if nodes[0].Data == "" {
		return "unknown", errors.ErrorNotFoundNewReleasesCategory
	}
	return nodes[0].Data, nil
}

// ParseASIN parses the ASIN from the given html node.
func (p *FRBoardParser) ParseASIN(node *html.Node) (string, error) {
	expr := `//div[@data-asin]`

	nodes, err := utils.FindNodes(node, expr, false)
	if err == nil && len(nodes) > 0 {
		return htmlquery.SelectAttr(nodes[0], "data-asin"), nil
	}
	return "unknown", errors.ErrorNotFoundASIN
}

// ParsePrice parses the price from the give html node.
func (p *FRBoardParser) ParsePrice(node *html.Node) (string, error) {
	exprs := []string{
		`div//span[contains(@class, "price")]/text()`,
		`div//span[contains(@class, "price")]/span/text()`,
	}

	for _, expr := range exprs {
		nodes, err := utils.FindNodes(node, expr, false)
		if err == nil && len(nodes) > 0 {
			return strings.TrimSpace(nodes[0].Data), nil
		}
	}
	return "unknown", errors.ErrorNotFoundPrice
}

// ParseStar parses the star from the give html node.
func (p *FRBoardParser) ParseStar(node *html.Node) (string, error) {
	expr := `//div/a[@title]`

	nodes, err := utils.FindNodes(node, expr, true)
	if err == nil && len(nodes) > 0 {
		stars := htmlquery.SelectAttr(nodes[0], "title")
		return utils.FindNumberHead(strings.TrimSpace(stars)), nil
	}
	return "unknown", errors.ErrorNotFoundStar
}

// ParseRating parses the rating from the give html node.
func (p *FRBoardParser) ParseRating(node *html.Node) (string, error) {
	expr := `//div/a[@title]/span/text()`

	nodes, err := utils.FindNodes(node, expr, true)
	if err == nil && len(nodes) > 0 {
		return nodes[0].Data, nil
	}
	return "unknown", errors.ErrorNotFoundRating
}

// ParseTitle parses the title from the give html node.
func (p *FRBoardParser) ParseTitle(node *html.Node) (string, error) {
	expr := `//a/span/div/text()`

	nodes, err := utils.FindNodes(node, expr, true)
	if err == nil && len(nodes) > 0 {
		return nodes[0].Data, nil
	}
	return "unknown", errors.ErrorNotFoundTitle
}

// ParseRank parses the rank from the give html node.
func (p *FRBoardParser) ParseRank(node *html.Node) (string, error) {
	expr := `//div/span/text()`

	nodes, err := utils.FindNodes(node, expr, true)
	if err == nil && len(nodes) > 0 {
		return strings.Replace(nodes[0].Data, "#", "", -1), nil
	}
	return "unknown", errors.ErrorNotFoundRank
}
