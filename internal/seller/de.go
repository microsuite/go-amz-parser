package seller

import (
	"strings"

	"golang.org/x/net/html"

	"github.com/antchfx/htmlquery"
	"github.com/microsuite/go-amz-parser/errors"
	"github.com/microsuite/go-amz-parser/utils"
)

type DESellerParser struct{}

// ParseAllProducts parses all products from the given HTML document.
func (p *DESellerParser) ParseAllProducts(doc *html.Node) ([]*html.Node, error) {
	expr := "//div[@class and @data-asin and string-length(@data-asin) > 0 and @data-index and @data-uuid]"
	nodes, err := utils.FindNodes(doc, expr, true)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func (p *DESellerParser) ParseMaxPageNum(doc *html.Node) (string, error) {
	expr := `//span[@class='s-pagination-item s-pagination-disabled']/text()`
	nodes, err := utils.FindNodes(doc, expr, true)
	if err != nil {
		expr := `//a[@class='s-pagination-item s-pagination-button']/text()`
		nodes, err := utils.FindNodes(doc, expr, true)
		if err != nil {
			return "unknown", err
		}
		return strings.TrimSpace(nodes[len(nodes)-1].Data), nil
	}
	return strings.TrimSpace(nodes[0].Data), nil
}

func (p *DESellerParser) ParseCurrentPageIndex(doc *html.Node) (string, error) {
	expr := `//span[contains(@aria-label, 'Current page')]/text()`
	nodes, err := utils.FindNodes(doc, expr, true)
	if err != nil {
		return "unknown", err
	}
	return strings.TrimSpace(nodes[0].Data), nil
}

func (p *DESellerParser) ParseNextPageURL(doc *html.Node) (string, error) {
	expr := `//a[contains(@aria-label, 'Zur nächsten Seite')]`
	nodes, err := utils.FindNodes(doc, expr, false)
	if err == nil && len(nodes) > 0 {
		nextRef := htmlquery.SelectAttr(nodes[0], "href")
		return nextRef, nil
	}
	return "unknown", errors.ErrorNotFoundNextPage
}

func (p *DESellerParser) ParseContentId(doc *html.Node) (string, error) {
	var contentId string

	expr := `//div[@id='reviewsRefinements']//span/li/@id`
	nodes, err := utils.FindNodes(doc, expr, true)
	if err != nil {
		return "unknown", err
	}

	for _, node := range nodes {
		contentId = htmlquery.SelectAttr(node, "id")
		if contentId == "" {
			return "unknown", errors.ErrorNotFoundContentId
		}
	}

	ids := strings.Split(contentId, "p_72/")
	if len(ids) >= 2 {
		return ids[1], nil
	}
	return "unknown", errors.ErrorNotFoundContentId
}

func (p *DESellerParser) ParseContentLink(doc *html.Node) (string, error) {
	expr := `//div[@id='reviewsRefinements']//ul/span/span/li/span/a/@href`

	nodes, err := utils.FindNodes(doc, expr, false)
	if err == nil && len(nodes) > 0 {
		nextRef := htmlquery.SelectAttr(nodes[0], "href")
		return nextRef, nil
	}
	return "unknown", errors.ErrorNotFoundContentLink
}

func (p *DESellerParser) ParsePagination(doc *html.Node) (string, error) {
	return "unknown", nil
}

func (p *DESellerParser) ParseASIN(node *html.Node) (string, error) {
	expr := `@data-asin`
	nodes, err := utils.FindNodes(node, expr, true)
	if err != nil {
		return "unknown", err
	}
	return htmlquery.SelectAttr(nodes[0], "data-asin"), nil
}

func (p *DESellerParser) ParsePrice(node *html.Node) (string, error) {
	expr := `//div//span[@class='a-price']/span/text()`
	nodes, err := utils.FindNodes(node, expr, true)
	if err != nil {
		return "unknown", err
	}

	price := utils.FormatNumberEuro(utils.DropMoneySym(nodes[0].Data))
	if price == "" {
		return "unknown", nil
	}
	return price, nil
}

func (p *DESellerParser) ParseStar(node *html.Node) (string, error) {
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

// ParseImg parses the image url from the html document
func (p *DESellerParser) ParseImg(node *html.Node) (string, error) {
	expr := `//div//img[contains(@class,"image")]/@src`
	nodes, err := utils.FindNodes(node, expr, true)
	if err != nil {
		return "unknown", err
	}
	return htmlquery.SelectAttr(nodes[0], "src"), nil
}
