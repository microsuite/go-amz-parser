package review

import (
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/microsuite/go-amz-parser/errors"
	"github.com/microsuite/go-amz-parser/utils"
	"golang.org/x/net/html"
)

type FRReviewParser struct{}

func NewFRReviewParser() *FRReviewParser {
	return &FRReviewParser{}
}

func (p *FRReviewParser) ParseAllReviews(doc *html.Node) ([]*html.Node, error) {
	expr := "//body/div[@review]"
	nodes, err := utils.FindNodes(doc, expr, true)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func (p *FRReviewParser) ParseReviewer(node *html.Node) (string, error) {
	expr := `//div[contains(@id, 'customer_review')]/div/a/div/span/text()`
	nodes, err := utils.FindNodes(node, expr, true)
	if err != nil {
		return "unknown", err
	}
	return strings.TrimSpace(nodes[0].Data), nil
}

func (p *FRReviewParser) ParseReviewerLink(node *html.Node) (string, error) {
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

func (p *FRReviewParser) ParseStar(node *html.Node) (string, error) {
	expr := `//span[contains(text(),'sur 5')]/text()`
	nodes, err := utils.FindNodes(node, expr, true)
	if err != nil {
		return "unknown", err
	}
	star := strings.TrimSpace(nodes[0].Data)
	return utils.FormatNumberEuro(strings.Split(star, " ")[0]), nil
}

func (p *FRReviewParser) ParseTitle(node *html.Node) (string, error) {
	expr := `//a[contains(@class, 'review-title')]/span/text()`
	nodes, err := utils.FindNodes(node, expr, true)
	if err != nil {
		return "unknown", err
	}
	return strings.TrimSpace(nodes[0].Data), nil
}

func (p *FRReviewParser) ParseDate(node *html.Node) (string, error) {
	expr := `//span[contains(@data-hook, 'review-date')]/text()`
	nodes, err := utils.FindNodes(node, expr, true)
	if err != nil {
		return "unknown", err
	}

	date := strings.TrimSpace(nodes[0].Data)
	dates := strings.Split(date, "Avis laissé en France le")
	if len(dates) < 1 {
		return "unknown", errors.ErrorNotFoundReviewerLink
	}
	return strings.TrimSpace(dates[1]), nil
}

func (p *FRReviewParser) ParsePurchase(node *html.Node) (string, error) {
	expr := `//div/a[contains(@class, 'a-link-normal')]/span/text()`
	nodes, err := utils.FindNodes(node, expr, true)
	if err != nil {
		return "unknown", err
	}
	return strings.TrimSpace(nodes[0].Data), nil
}

func (p *FRReviewParser) ParseContent(node *html.Node) (string, error) {
	expr := `//div[@review-text-content]/span/text()`
	nodes, err := utils.FindNodes(node, expr, true)
	if err != nil {
		return "unknown", err
	}
	return strings.TrimSpace(nodes[0].Data), nil
}
