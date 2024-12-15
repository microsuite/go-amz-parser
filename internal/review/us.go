package review

import (
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/microsuite/go-amz-parser/errors"
	"github.com/microsuite/go-amz-parser/utils"
	"golang.org/x/net/html"
)

type USReviewParser struct{}

func NewUSReviewParser() *USReviewParser {
	return &USReviewParser{}
}

func (p *USReviewParser) ParseAllReviews(doc *html.Node) ([]*html.Node, error) {
	expr := "//div[@id='cm-cr-dp-review-list']/div"
	nodes, err := utils.FindNodes(doc, expr, true)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func (p *USReviewParser) ParseReviewer(node *html.Node) (string, error) {
	expr := `//div[contains(@id, 'customer_review')]/div/a/div/span/text()`
	nodes, err := utils.FindNodes(node, expr, true)
	if err != nil {
		return "unknown", err
	}
	return strings.TrimSpace(nodes[0].Data), nil
}

func (p *USReviewParser) ParseReviewerLink(node *html.Node) (string, error) {
	expr := `//div[contains(@id, 'customer_review')]/div/a/@href`
	nodes, err := utils.FindNodes(node, expr, true)
	if err != nil {
		return "unknown", err
	}

	link := htmlquery.SelectAttr(nodes[0], "href")
	if link == "" {
		return "unknown", errors.ErrorNotFoundReviewerLink
	}
	return link, nil
}

func (p *USReviewParser) ParseStar(node *html.Node) (string, error) {
	expr := `//span[contains(text(),'out of 5 stars')]/text()`
	nodes, err := utils.FindNodes(node, expr, true)
	if err != nil {
		return "unknown", err
	}
	star := strings.TrimSpace(nodes[0].Data)
	star = utils.FindNumberHead(star)
	return utils.FormatNumber(star), nil
}

func (p *USReviewParser) ParseTitle(node *html.Node) (string, error) {
	expr := `//a[contains(@class, 'review-title')]/span/text()`
	nodes, err := utils.FindNodes(node, expr, true)
	if err != nil {
		return "unknown", err
	}
	return strings.TrimSpace(nodes[0].Data), nil
}

func (p *USReviewParser) ParseDate(node *html.Node) (string, error) {
	expr := `//span[contains(@class, 'review-date')]/text()`
	nodes, err := utils.FindNodes(node, expr, true)
	if err != nil {
		return "unknown", err
	}

	date := strings.TrimSpace(nodes[0].Data)
	dates := strings.Split(date, "Bewertet in Deutschland am")
	if len(dates) < 1 {
		return "unknown", errors.ErrorNotFoundReviewerLink
	}
	return strings.TrimSpace(dates[1]), nil
}

func (p *USReviewParser) ParsePurchase(node *html.Node) (string, error) {
	expr := `//div/a[contains(@class, 'a-link-normal')]/span/text()`
	nodes, err := utils.FindNodes(node, expr, true)
	if err != nil {
		return "unknown", err
	}
	return strings.TrimSpace(nodes[0].Data), nil
}

func (p *USReviewParser) ParseContent(node *html.Node) (string, error) {
	expr := `//div[contains(@class, 'review-text-content')]/span/text()`
	nodes, err := utils.FindNodes(node, expr, true)
	if err != nil {
		return "unknown", err
	}
	return strings.TrimSpace(nodes[0].Data), nil
}
