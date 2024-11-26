package keyword

import (
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/microsuite/go-amz-parser/errors"
	"github.com/microsuite/go-amz-parser/utils"
	"golang.org/x/net/html"
)

type DEKeywordParser struct{}

func NewDEKeywordParser() *DEKeywordParser {
	return &DEKeywordParser{}
}

// ParseAllProducts parses all products from the given HTML document.
func (p *DEKeywordParser) ParseAllProducts(doc *html.Node) ([]*html.Node, error) {
	expr := "//div[@class and @data-asin and string-length(@data-asin) > 0 and @data-index and @data-uuid]"
	nodes, err := utils.FindNodes(doc, expr, true)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func (p *DEKeywordParser) ParseCurrentPageIndex(doc *html.Node) (string, error) {
	expr := `//span[contains(@aria-label, 'Aktuelle Seite')]/text()`
	nodes, err := utils.FindNodes(doc, expr, true)
	if err != nil {
		return "unknown", err
	}
	return strings.TrimSpace(nodes[0].Data), nil
}

func (p *DEKeywordParser) ParseNextPageURL(doc *html.Node) (string, error) {
	expr := `//a[contains(@aria-label, 'Zur nächsten Seite')]`
	nodes, err := utils.FindNodes(doc, expr, false)
	if err == nil && len(nodes) > 0 {
		nextRef := htmlquery.SelectAttr(nodes[0], "href")
		return nextRef, nil
	}
	return "unknown", errors.ErrorNotFoundNextPage
}

func (p *DEKeywordParser) ParseKeyword(doc *html.Node) (string, error) {
	expr := `//input[@id='twotabsearchtextbox']/@value`
	nodes, err := utils.FindNodes(doc, expr, true)
	if err != nil {
		return "unknown", err
	}

	keyword := htmlquery.SelectAttr(nodes[0], "value")
	if keyword == "" {
		return "unknown", errors.ErrorNotFoundImgURL
	}
	return keyword, nil
}

func (p *DEKeywordParser) ParseASIN(node *html.Node) (string, error) {
	expr := `@data-asin`
	nodes, err := utils.FindNodes(node, expr, true)
	if err != nil {
		return "unknown", err
	}
	return htmlquery.SelectAttr(nodes[0], "data-asin"), nil
}

func (p *DEKeywordParser) ParsePrice(node *html.Node) (string, error) {
	expr := `//div//span[@class='a-price']/span/text()`
	nodes, err := utils.FindNodes(node, expr, true)
	if err != nil {
		return "unknown", err
	}

	price := strings.TrimSpace(nodes[0].Data)
	if price == "" {
		return "unknown", nil
	}
	return price, nil
}

func (p *DEKeywordParser) ParseStar(node *html.Node) (string, error) {
	expr := `//div//span[contains(@aria-label,'von 5 Sternen')]`
	nodes, err := utils.FindNodes(node, expr, true)
	if err != nil {
		return "unknown", err
	}
	stars := htmlquery.SelectAttr(nodes[0], "aria-label")
	star := utils.FormatNumberEuro(strings.Split(stars, " ")[0])
	if star == "" {
		return "unknown", nil
	}
	return star, nil
}

func (p *DEKeywordParser) ParseRating(node *html.Node) (string, error) {
	expr := `//span[contains(@aria-label, 'Gesponsert')]/a/span/text()`
	nodes, err := utils.FindNodes(node, expr, true)
	if err != nil {
		return "0", err
	}
	return utils.FormatRating(nodes[0].Data), nil
}

// ParseSponsered parses the sponsered from the html document
func (p *DEKeywordParser) ParseSponsered(node *html.Node) (string, error) {
	expr := `//div//span[text()="Sponsored"]`
	_, err := utils.FindNodes(node, expr, true)
	if err != nil {
		return "0", err
	}
	return "1", nil
}

// ParsePrime parses the prime from the html document
func (p *DEKeywordParser) ParsePrime(node *html.Node) (string, error) {
	expr := `//div//i[@aria-label="Amazon Prime"]`
	_, err := utils.FindNodes(node, expr, true)
	if err != nil {
		return "false", err
	}
	return "true", nil
}

// ParseSales parses the sales from the html document
func (p *DEKeywordParser) ParseSales(node *html.Node) (string, error) {
	expr := `//div//span[contains(text(), "Mal im letzten Monat gekauft")]/text()`
	nodes, err := utils.FindNodes(node, expr, true)
	if err != nil {
		return "unknown", err
	}
	sales := strings.Trim(nodes[0].Data, "Mal im letzten Monat gekauft")
	return utils.FormatNumber(sales), nil
}

// ParseImg parses the image url from the html document
func (p *DEKeywordParser) ParseImg(node *html.Node) (string, error) {
	expr := `//div//img[contains(@class,"image")]/@src`
	nodes, err := utils.FindNodes(node, expr, true)
	if err != nil {
		return "unknown", err
	}
	return htmlquery.SelectAttr(nodes[0], "src"), nil
}

// ParseTitle parses the title from the html document
func (p *DEKeywordParser) ParseTitle(node *html.Node) (string, error) {
	expr := `//div//span[contains(@class, "text-normal")]/text()`
	nodes, err := utils.FindNodes(node, expr, true)
	if err != nil {
		return "unknown", err
	}
	return utils.FormatTitle(nodes[0].Data), nil
}
