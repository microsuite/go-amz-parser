package product

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/net/html"

	"github.com/antchfx/htmlquery"
	"github.com/microsuite/go-amz-parser/errors"
	"github.com/microsuite/go-amz-parser/utils"
)

type DEProductParser struct{}

func NewDEProductParser() *DEProductParser {
	return &DEProductParser{}
}

func (p *DEProductParser) ParseASIN(doc *html.Node) (string, error) {
	var asin string

	expr := `//div[starts-with(@id, "corePrice") and @data-csa-c-asin]`
	nodes, err := utils.FindNodes(doc, expr, true)
	if err == nil && len(nodes) > 0 {
		asin = htmlquery.SelectAttr(nodes[0], "data-csa-c-asin")
		if asin != "" {
			return asin, nil
		}
	}

	if asin == "" {
		// parse review from top node.
		reviewExpr := `//div[@id="averageCustomerReviews" and @data-asin]`
		nodes, err := utils.FindNodes(doc, reviewExpr, true)
		if err != nil {
			return asin, err
		}

		// parse asin from review.
		asin = htmlquery.SelectAttr(nodes[0], "data-asin")
	}
	return asin, nil
}

func (p *DEProductParser) ParseStar(doc *html.Node) (string, error) {
	expr := `//span[contains(text(),'von 5')]/text()`
	nodes, err := utils.FindNodes(doc, expr, true)
	if err != nil {
		return "unknown", err
	}
	star := strings.TrimSpace(nodes[0].Data)
	return utils.FormatNumberEuro(strings.Split(star, " ")[0]), nil
}

func (p *DEProductParser) ParseRating(doc *html.Node) (string, error) {
	expr := `//a[@id='acrCustomerReviewLink']/span[@id='acrCustomerReviewText']/text()`
	nodes, err := utils.FindNodes(doc, expr, true)
	if err != nil {
		return "unknown", err
	}
	rating := utils.FindNumberHead(strings.TrimSpace(nodes[0].Data))
	return strings.TrimSpace(utils.FormatNumberEuro(rating)), nil
}

func (p *DEProductParser) ParseTitle(doc *html.Node) (string, error) {
	expr := `//span[@id='productTitle']/text()`
	nodes, err := utils.FindNodes(doc, expr, true)
	if err != nil {
		return "unknown", err
	}
	return utils.FormatTitle(nodes[0].Data), nil
}

func (p *DEProductParser) ParseImg(doc *html.Node) (string, error) {
	expr := `//div[@id="imageBlock"]//div[@class="imgTagWrapper"]/img/@src`
	nodes, err := utils.FindNodes(doc, expr, true)
	if err != nil {
		return "unknown", err
	}

	img := htmlquery.SelectAttr(nodes[0], "src")
	if img == "" {
		return "unknown", errors.ErrorNotFoundImgURL
	}
	return img, nil
}

func (p *DEProductParser) ParsePrice(doc *html.Node) (string, error) {
	var err error
	var price string
	exprs := []string{
		`//div[starts-with(@id, "corePrice") and @data-csa-c-asin]/div/span/text()`,
	}

	for _, expr := range exprs {
		nodes, err := utils.FindNodes(doc, expr, true)
		if err == nil {
			splits := strings.Split(strings.TrimSpace(nodes[0].Data), " ")
			if len(splits) > 0 && splits[0] != "" {
				price = splits[0]
			}
			if price == "" && len(splits) > 1 && splits[1] != "" {
				price = splits[1]
			}
			return price, nil
		}
	}
	return "unknown", err
}

func (p *DEProductParser) ParseDispatchFrom(doc *html.Node) (string, error) {
	var err error
	exprs := []string{
		`//div/span[contains(text(), 'Versand')]/following-sibling::span/text()`, // ProductOfCategory
		`//div/span[contains(text(), 'Versand')]/../../following-sibling::div/div/span/text()`,
	}

	for _, expr := range exprs {
		nodes, err := utils.FindNodes(doc, expr, true)
		if err == nil {
			return strings.TrimSpace(nodes[0].Data), nil
		}
	}
	return "unknown", err
}

func (p *DEProductParser) ParseSoldBy(doc *html.Node) (string, error) {
	var err error
	exprs := []string{
		`//div/span[contains(text(), 'Verkäufer')]/../../following-sibling::div/div/span/a/text()`,
		`//div/span[contains(text(), 'Verkäufer')]/../../following-sibling::div/div/span/text()`,
	}

	for _, expr := range exprs {
		nodes, err := utils.FindNodes(doc, expr, true)
		if err == nil {
			return strings.TrimSpace(nodes[0].Data), nil
		}
	}
	return "unknown", err
}

func (p *DEProductParser) ParsePackageDimensions(doc *html.Node) (string, error) {
	exprs := []string{
		`//tbody/tr/th[contains(text(), 'Verpackungsabmessungen')]/following-sibling::td/text()`,
		`//tbody/tr/th[contains(text(), 'Paket-Abmessungen')]/following-sibling::td/text()`,
		`//span[contains(text(), 'Verpackungsabmessungen')]/following-sibling::span/text()`,
	}
	for _, expr := range exprs {
		nodes, err := utils.FindNodes(doc, expr, true)
		if err == nil {
			return strings.TrimSpace(nodes[0].Data), nil
		}
	}
	return "unknown", errors.ErrorNotFoundPackageDimensions
}

func (p *DEProductParser) ParsePackageWeight(doc *html.Node) (string, error) {
	return "unknown", nil
}

func (p *DEProductParser) ParseFirstAvailDate(doc *html.Node) (string, error) {
	exprs := []string{
		`//tbody/tr/th[contains(text(), 'Im Angebot von Amazon.de seit')]/following-sibling::td/text()`,
		`//span[contains(text(), 'm Angebot von Amazon.de seit')]/following-sibling::span/text()`,
	}
	for _, expr := range exprs {
		nodes, err := utils.FindNodes(doc, expr, true)
		if err == nil {
			return strings.TrimSpace(nodes[0].Data), nil
		}
	}
	return "unknown", errors.ErrorNotFoundFirstDate
}

func (p *DEProductParser) ParseSellerId(node *html.Node) (string, error) {
	expr := `//input[@id='deliveryBlockSelectMerchant']/@value`
	nodes, err := utils.FindNodes(node, expr, true)
	if err != nil {
		return "unknown", err
	}

	sellerId := htmlquery.SelectAttr(nodes[0], "value")
	if sellerId == "" {
		return "unknown", errors.ErrorNotFoundImgURL
	}
	return sellerId, nil
}

func (p *DEProductParser) ParseCategoryId(node *html.Node) (string, error) {
	var categoryId string

	expr := `//tbody/tr/th[contains(text(), 'Amazon Bestseller-Rang')]/following-sibling::td/span/span/a/@href`
	nodes, err := utils.FindNodes(node, expr, true)
	if err != nil {
		return "unknown", err
	}

	for _, node := range nodes {
		categoryId = htmlquery.SelectAttr(node, "href")
		if categoryId == "" {
			return "unknown", errors.ErrorNotFoundCategoryId
		}
	}

	regex := regexp.MustCompile(`\d+`)
	mathches := regex.FindStringSubmatch(categoryId)

	for _, match := range mathches {
		return match, nil
	}
	return "unknown", nil
}

func (p *DEProductParser) ParseProductDimensions(doc *html.Node) (string, error) {
	exprs := []string{
		`//tbody/tr/th[contains(text(), 'Produktabmessungen')]/following-sibling::td/text()`,
		`//span[contains(text(), 'Produktabmessungen')]/following-sibling::span/text()`,
	}

	for _, expr := range exprs {
		nodes, err := utils.FindNodes(doc, expr, true)
		if err == nil && len(nodes) > 0 {
			return strings.TrimSpace(nodes[0].Data), nil
		}
	}
	return "unknown", errors.ErrorNotFoundDimensions
}

func (p *DEProductParser) ParseProductWeight(doc *html.Node) (string, error) {
	exprs := []string{
		`//tbody/tr/th[contains(text(), 'Artikelgewicht')]/following-sibling::td/text()`,
		`//span[contains(text(), 'Artikelgewicht')]/following-sibling::span/text()`,
	}

	for _, expr := range exprs {
		nodes, err := utils.FindNodes(doc, expr, true)
		if err == nil && len(nodes) > 0 {
			return strings.TrimSpace(nodes[0].Data), nil
		}
	}
	return "unknown", errors.ErrorNotFoundWeight
}

func (p *DEProductParser) ParseHasCart(doc *html.Node) (string, error) {
	nodes, err := utils.FindNodes(doc, `//input[@id='add-to-cart-button']`, true)
	if err != nil || len(nodes) == 0 {
		return "false", err
	}
	return "true", nil
}

func (p *DEProductParser) ParseCoupon(doc *html.Node) (string, error) {
	expr := `//i[contains(text(),'Coupon')]/following-sibling::label/text()`

	nodes, err := utils.FindNodes(doc, expr, true)
	if err != nil || len(nodes) == 0 {
		return "unknown", err
	}
	return strings.TrimSpace(nodes[0].Data), nil
}

func (p *DEProductParser) ParseColor(doc *html.Node) (string, error) {
	exprs := []string{
		`//label[contains(text(),'Farbe:')]/following-sibling::span/text()`,
	}

	for _, expr := range exprs {
		nodes, err := utils.FindNodes(doc, expr, true)
		if err == nil || len(nodes) != 0 {
			return strings.TrimSpace(nodes[0].Data), nil
		}
	}
	return "unknown", errors.ErrorNotFoundColor
}

func (p *DEProductParser) ParseSize(doc *html.Node) (string, error) {
	exprs := []string{
		`//span[contains(text(), 'Größe')]/../following-sibling::td/span/text()`,
	}

	for _, expr := range exprs {
		nodes, err := utils.FindNodes(doc, expr, true)
		if err == nil || len(nodes) != 0 {
			return strings.TrimSpace(nodes[0].Data), nil
		}
	}
	return "unknown", errors.ErrorNotFoundSize
}

func (p *DEProductParser) ParseSpecs(doc *html.Node) ([]string, error) {
	specs := make([]string, 0)
	var m map[string]interface{}

	pattern := `"asinVariationValues(.*)`
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	match := re.FindString(htmlquery.InnerText(doc))

	values := strings.Split(match, `"asinVariationValues" : `)
	if len(values) > 1 {
		str := strings.Trim(values[1], ",")
		if err := json.Unmarshal([]byte(str), &m); err != nil {
			return nil, err
		}

		for key := range m {
			specs = append(specs, key)
		}
	}
	return specs, nil
}

func (p *DEProductParser) ParseDescription(doc *html.Node) (string, error) {
	expr := `//h1[contains(text(), 'Info zu diesem Artikel')]/following-sibling::ul[1]/li/span/text()`
	nodes, err := utils.FindNodes(doc, expr, true)
	if err != nil {
		return "unknown", err
	}

	var desc string
	for i, n := range nodes {
		desc += fmt.Sprintf("%v. %v ", i+1, n.Data)
	}

	if desc == "" {
		return "unknown", errors.ErrorNotFoundDesc
	}
	return desc, nil
}

func (p *DEProductParser) ParseDeliveryTime(doc *html.Node) (string, error) {
	expr := `//div[@id='mir-layout-DELIVERY_BLOCK-slot-PRIMARY_DELIVERY_MESSAGE_LARGE']/span/@data-csa-c-delivery-time`
	nodes, err := utils.FindNodes(doc, expr, true)
	if err != nil {
		return "unknown", err
	}

	fastestDelivery := htmlquery.SelectAttr(nodes[0], "data-csa-c-delivery-time")
	if fastestDelivery == "" {
		return "unknown", errors.ErrorNotFoundDeliveryTime
	}
	return fastestDelivery, nil
}

func (p *DEProductParser) ParseFastestDelivery(doc *html.Node) (string, error) {
	expr := `//div[@id='mir-layout-DELIVERY_BLOCK-slot-SECONDARY_DELIVERY_MESSAGE_LARGE']/span/@data-csa-c-delivery-time`
	nodes, err := utils.FindNodes(doc, expr, true)
	if err != nil {
		return "unknown", err
	}

	fastestDelivery := htmlquery.SelectAttr(nodes[0], "data-csa-c-delivery-time")
	if fastestDelivery == "" {
		return "unknown", errors.ErrorNotFoundFastestDelivery
	}
	return fastestDelivery, nil
}

func (p *DEProductParser) ParsePrimePrice(doc *html.Node) (string, error) {
	return "unknown", nil
}

func (p *DEProductParser) ParseBrand(doc *html.Node) (string, error) {
	return "unknown", nil
}

func (p *DEProductParser) ParseCategoryHierarchy(doc *html.Node) ([]string, error) {
	expr := `//div[@id='wayfinding-breadcrumbs_feature_div']//li/span/a/text()`
	nodes, err := utils.FindNodes(doc, expr, true)
	if err != nil {
		return nil, err
	}

	var categoryHierarchies []string
	for _, node := range nodes {
		categoryHierarchies = append(categoryHierarchies, strings.TrimSpace(node.Data))
	}
	return categoryHierarchies, err
}

func (p *DEProductParser) ParseCustomerReviews(doc *html.Node) (map[string]string, error) {
	expr := `//li[@class='a-align-center a-spacing-none']`
	nodes, err := utils.FindNodes(doc, expr, true)
	if err != nil {
		return nil, err
	}

	customerReviews := make(map[string]string)
	for _, node := range nodes {
		leftExpr := `//span[@class='a-list-item']//div[contains(@class, 'a-text-left')]/text()`
		rightExpr := `//span[@class='a-list-item']//div[contains(@class, 'a-text-right')]/text()`

		leftNodes, err := utils.FindNodes(node, leftExpr, false)
		if err != nil {
			continue
		}

		var percentage string
		rightNodes, err := utils.FindNodes(node, rightExpr, false)
		if err != nil {
			percentage = "unknown"
		}
		percentage = strings.TrimSpace(rightNodes[0].Data)

		customerReviews[strings.TrimSpace(leftNodes[0].Data)] = strings.TrimSpace(percentage)
	}
	return customerReviews, nil
}
