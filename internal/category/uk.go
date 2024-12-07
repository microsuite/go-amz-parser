package category

import (
	"strings"

	"golang.org/x/net/html"

	"github.com/antchfx/htmlquery"
	"github.com/microsuite/go-amz-parser/errors"
	"github.com/microsuite/go-amz-parser/utils"
)

type UKCategoryParser struct{}

func NewUKCategoryParser() *UKCategoryParser {
	return &UKCategoryParser{}
}

// ParseAllProducts parses all products from the given HTML document.
func (p *UKCategoryParser) ParseAllProducts(doc *html.Node) ([]*html.Node, error) {
	expr := "//div[@class and @data-asin and string-length(@data-asin) > 0 and @data-index and @data-uuid]"
	nodes, err := utils.FindNodes(doc, expr, true)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func (p *UKCategoryParser) ParseMaxPageNum(doc *html.Node) (string, error) {
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

func (p *UKCategoryParser) ParseCurrentPageIndex(doc *html.Node) (string, error) {
	expr := `//span[contains(@aria-label, 'Current page')]/text()`
	nodes, err := utils.FindNodes(doc, expr, true)
	if err != nil {
		return "unknown", err
	}
	return strings.TrimSpace(nodes[0].Data), nil
}

func (p *UKCategoryParser) ParseNextPageURL(doc *html.Node) (string, error) {
	expr := `//a[contains(@aria-label, "Go to next page")]`
	nodes, err := utils.FindNodes(doc, expr, false)
	if err == nil && len(nodes) > 0 {
		nextRef := htmlquery.SelectAttr(nodes[0], "href")
		return nextRef, nil
	}
	return "unknown", errors.ErrorNotFoundNextPage
}

func (p *UKCategoryParser) ParseContentId(doc *html.Node) (string, error) {
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

func (p *UKCategoryParser) ParseContentLink(doc *html.Node) (string, error) {
	expr := `//div[@id='reviewsRefinements']//ul/span/span/li/span/a/@href`

	nodes, err := utils.FindNodes(doc, expr, false)
	if err == nil && len(nodes) > 0 {
		nextRef := htmlquery.SelectAttr(nodes[0], "href")
		return nextRef, nil
	}
	return "unknown", errors.ErrorNotFoundContentLink
}

func (p *UKCategoryParser) ParsePagination(doc *html.Node) (string, error) {
	expr := `//span[contains(text(), 'results for')]/text()`

	nodes, err := utils.FindNodes(doc, expr, false)
	if err != nil {
		return "unknown", nil
	}
	return nodes[0].Data, nil
}

func (p *UKCategoryParser) ParseCategoryName(doc *html.Node) (string, error) {
	expr := `//form//span[@id='nav-search-label-id']//text()`

	nodes, err := utils.FindNodes(doc, expr, false)
	if err != nil {
		return "unknown", nil
	}
	return nodes[0].Data, nil
}

func (p *UKCategoryParser) ParseASIN(node *html.Node) (string, error) {
	expr := `@data-asin`
	nodes, err := utils.FindNodes(node, expr, true)
	if err != nil {
		return "unknown", err
	}
	return htmlquery.SelectAttr(nodes[0], "data-asin"), nil
}

func (p *UKCategoryParser) ParsePrice(node *html.Node) (string, error) {
	expr := `//div//span[@class="a-price"]/span/text()`
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

func (p *UKCategoryParser) ParseStar(node *html.Node) (string, error) {
	expr := `//div//span[contains(@aria-label,'out of 5 stars')]`

	var star string
	nodes, err := utils.FindNodes(node, expr, true)
	if err == nil && len(nodes) > 0 {
		stars := htmlquery.SelectAttr(nodes[0], "aria-label")
		star = utils.FormatNumber(strings.Split(stars, " ")[0])
		if star == "" {
			return "unknown", nil
		}
	} else {
		expr := `//span[contains(text(),'out of 5 stars')]/text()`

		nodes, err := utils.FindNodes(node, expr, true)
		if err != nil {
			return "unknown", err
		}
		stars := strings.TrimSpace(nodes[0].Data)
		star = utils.FormatNumber(strings.Split(stars, " ")[0])
		if star == "" {
			return "unknown", nil
		}
	}
	return star, nil
}

// ParseImg parses the image url from the html document
func (p *UKCategoryParser) ParseImg(node *html.Node) (string, error) {
	expr := `//div//img[contains(@class,"image")]/@src`
	nodes, err := utils.FindNodes(node, expr, true)
	if err != nil {
		return "unknown", err
	}
	return htmlquery.SelectAttr(nodes[0], "src"), nil
}

// ParseTitle parses the title from the html document
func (p *UKCategoryParser) ParseTitle(node *html.Node) (string, error) {
	expr := `//div//span[contains(@class, "text-normal")]/text()`
	nodes, err := utils.FindNodes(node, expr, true)
	if err != nil {
		return "unknown", err
	}
	return utils.FormatTitle(nodes[0].Data), nil
}
