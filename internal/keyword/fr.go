package keyword

import (
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/microsuite/go-amz-parser/errors"
	"github.com/microsuite/go-amz-parser/utils"
	"golang.org/x/net/html"
)

type FRKeywordParser struct{}

func (p *FRKeywordParser) ParseAllProducts(doc *html.Node) ([]*html.Node, error) {
	expr := "//div[@class and @data-asin and string-length(@data-asin) > 0 and @data-index and @data-uuid]"
	nodes, err := utils.FindNodes(doc, expr, true)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func (p *FRKeywordParser) ParseCurrentPageIndex(doc *html.Node) (string, error) {
	expr := `//span[contains(@aria-label, 'Page actuelle')]/text()`
	nodes, err := utils.FindNodes(doc, expr, true)
	if err != nil {
		return "unknown", err
	}
	return strings.TrimSpace(nodes[0].Data), nil
}

func (p *FRKeywordParser) ParseNextPageURL(doc *html.Node) (string, error) {
	expr := `//a[contains(@aria-label, 'Accéder à la page suivant')]`
	nodes, err := utils.FindNodes(doc, expr, false)
	if err == nil && len(nodes) > 0 {
		nextRef := htmlquery.SelectAttr(nodes[0], "href")
		return nextRef, nil
	}
	return "unknown", errors.ErrorNotFoundNextPage
}

func (p *FRKeywordParser) ParseASIN(node *html.Node) (string, error) {
	expr := `@data-asin`
	nodes, err := utils.FindNodes(node, expr, true)
	if err != nil {
		return "unknown", err
	}
	return htmlquery.SelectAttr(nodes[0], "data-asin"), nil
}

func (p *FRKeywordParser) ParsePrice(node *html.Node) (string, error) {
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

func (p *FRKeywordParser) ParseStar(node *html.Node) (string, error) {
	expr := `//div//span[contains(@aria-label,'sur 5')]`
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

func (p *FRKeywordParser) ParseRating(node *html.Node) (string, error) {
	expr := `//div//span/div/span[contains(@aria-label, 'évaluations')]/a/span/text()`
	nodes, err := utils.FindNodes(node, expr, true)
	if err != nil {
		return "unknown", err
	}
	return utils.FormatRating(nodes[0].Data), nil
}

// ParseSponsered parses the sponsered from the html document
func (p *FRKeywordParser) ParseSponsered(node *html.Node) (string, error) {
	expr := `//div//span[text()='Sponsorisé']`
	_, err := utils.FindNodes(node, expr, true)
	if err != nil {
		return "0", err
	}
	return "1", nil
}

// ParsePrime parses the prime from the html document
func (p *FRKeywordParser) ParsePrime(node *html.Node) (string, error) {
	expr := `//div//i[@aria-label="Amazon Prime"]`
	_, err := utils.FindNodes(node, expr, true)
	if err != nil {
		return "false", err
	}
	return "true", nil
}

// ParseSales parses the sales from the html document
func (p *FRKeywordParser) ParseSales(node *html.Node) (string, error) {
	expr := `//div//span[contains(text(), "achetés au cours du mois dernier")]/text()`
	nodes, err := utils.FindNodes(node, expr, true)
	if err != nil {
		return "unknown", err
	}
	sales := strings.Trim(nodes[0].Data, "achetés au cours du mois dernier")
	return utils.FormatNumber(sales), nil
}

// ParseImg parses the image url from the html document
func (p *FRKeywordParser) ParseImg(node *html.Node) (string, error) {
	expr := `//div//img[contains(@class,"image")]/@src`
	nodes, err := utils.FindNodes(node, expr, true)
	if err != nil {
		return "unknown", err
	}
	return htmlquery.SelectAttr(nodes[0], "src"), nil
}

// ParseTitle parses the title from the html document
func (p *FRKeywordParser) ParseTitle(node *html.Node) (string, error) {
	expr := `//div//span[contains(@class, "text-normal")]/text()`
	nodes, err := utils.FindNodes(node, expr, true)
	if err != nil {
		return "unknown", err
	}
	return utils.FormatTitle(nodes[0].Data), nil
}
