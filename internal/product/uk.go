package product

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/net/html"

	"github.com/antchfx/htmlquery"
	"github.com/microsuite/go-amz-parser/errors"
	"github.com/microsuite/go-amz-parser/utils"
)

type UKProductParser struct{}

func (p *UKProductParser) ParseASIN(doc *html.Node) (string, error) {
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

func (p *UKProductParser) ParseStar(doc *html.Node) (string, error) {
	expr := `//span[contains(text(),'out of 5 stars')]/text()`
	nodes, err := utils.FindNodes(doc, expr, true)
	if err != nil || len(nodes) == 0 {
		return "unknown", err
	}
	star := nodes[0].Data
	star = strings.TrimSpace(star)
	star = utils.FindNumberHead(star)
	return utils.FormatNumber(star), nil
}

func (p *UKProductParser) ParseRating(doc *html.Node) (string, error) {
	exprs := []string{
		`//a[@id='acrCustomerReviewLink']/span[@id='acrCustomerReviewText']/text()`,
		`//a[contains(text(),'customer ratings')]/text()`,
	}

	for _, expr := range exprs {
		nodes, err := utils.FindNodes(doc, expr, true)
		if err == nil {
			rating := utils.FindNumberHead(strings.TrimSpace(nodes[0].Data))
			return utils.FormatNumber(rating), nil
		}
	}
	return "unknown", errors.ErrorNotFoundRating
}

func (p *UKProductParser) ParseTitle(doc *html.Node) (string, error) {
	expr := `//span[@id="productTitle"]/text()`
	nodes, err := utils.FindNodes(doc, expr, true)
	if err != nil {
		return "unknown", err
	}
	return utils.FormatTitle(nodes[0].Data), nil
}

func (p *UKProductParser) ParseImg(doc *html.Node) (string, error) {
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

func (p *UKProductParser) ParsePrice(doc *html.Node) (string, error) {
	var price string
	exprs := []string{
		`//span[@class='a-price' and @data-a-color="base"]/span/text()`,
		`//span[starts-with(@class, 'a-price') and @data-a-color="price"]/span/text()`,
		`//span[starts-with(@class, 'a-price') and @data-a-color="base"]/span/text()`,
		`//div[starts-with(@id, "corePrice") and @data-csa-c-asin]/div/span/text()`,
		`//span[starts-with(@id, "a-price") and @data-a-color="price"]/span/text()`,
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
	return "unknown", errors.ErrorNotFoundPrice
}

func (p *UKProductParser) ParseDispatchFrom(doc *html.Node) (string, error) {
	exprs := []string{
		`//div/span[contains(text(), "Dispatches from")]/following-sibling::span/text()`,               // ProductOfCategory
		`//div/span[contains(text(), "Dispatches from")]/../../following-sibling::div/div/span/text()`, // ProductOfSeller
	}

	for _, expr := range exprs {
		nodes, err := utils.FindNodes(doc, expr, true)
		if err == nil {
			return strings.TrimSpace(nodes[0].Data), nil
		}
	}
	return "unknown", errors.ErrorNotFoundDispatchFrom
}

func (p *UKProductParser) ParseSoldBy(doc *html.Node) (string, error) {
	exprs := []string{
		`//div/span[contains(text(), "Sold by")]/following-sibling::span/text()`,                 // ProductOfCategory
		`//div/span[contains(text(), "Sold by")]/../../following-sibling::div/div/span/a/text()`, // ProductOfSeller
	}

	for _, expr := range exprs {
		nodes, err := utils.FindNodes(doc, expr, true)
		if err == nil {
			return strings.TrimSpace(nodes[0].Data), nil
		}
	}
	return "unknown", errors.ErrorNotFoundSoldBy
}

func (p *UKProductParser) ParseProductDimensions(doc *html.Node) (string, error) {
	exprs := []string{
		`//tbody/tr/th[contains(text(), 'Product Dimensions')]/following-sibling::td/text()`,
		`//span[contains(text(), 'Product Dimensions')]/following-sibling::span/text()`,
	}
	for _, expr := range exprs {
		nodes, err := utils.FindNodes(doc, expr, true)
		if err == nil && len(nodes) > 0 {
			return strings.TrimSpace(nodes[0].Data), nil
		}
	}
	return "unknown", errors.ErrorNotFoundDimensions
}

func (p *UKProductParser) ParsePackageDimensions(doc *html.Node) (string, error) {
	exprs := []string{
		`//tbody/tr/th[contains(text(), 'Package Dimensions')]/following-sibling::td/text()`,
		`//span[contains(text(), 'Package Dimensions')]/following-sibling::span/text()`,
	}
	for _, expr := range exprs {
		nodes, err := utils.FindNodes(doc, expr, true)
		if err == nil {
			return strings.TrimSpace(nodes[0].Data), nil
		}
	}
	return "unknown", errors.ErrorNotFoundPackageDimensions
}

func (p *UKProductParser) ParsePackageWeight(doc *html.Node) (string, error) {
	expr := `//tbody/tr/th[contains(text(), "Package Weight")]/following-sibling::td/text()`
	nodes, err := utils.FindNodes(doc, expr, true)
	if err != nil {
		return "unknown", err
	}
	return strings.TrimSpace(nodes[0].Data), nil
}

func (p *UKProductParser) ParseProductWeight(doc *html.Node) (string, error) {
	exprs := []string{
		`//th[contains(text(), 'Item Weight')]/following-sibling::td/text()`,
		`//span[contains(text(), 'Item Weight')]/following-sibling::span/text()`,
	}
	for _, expr := range exprs {
		nodes, err := utils.FindNodes(doc, expr, true)
		if err == nil && len(nodes) > 0 {
			return strings.TrimSpace(nodes[0].Data), nil
		}
	}
	return "unknown", errors.ErrorNotFoundWeight
}

func (p *UKProductParser) ParseFirstAvailDate(doc *html.Node) (string, error) {
	exprs := []string{
		`//tbody/tr/th[contains(text(), 'Date First Available')]/following-sibling::td/text()`,
		`//span[contains(text(), 'Date First Available')]/following-sibling::span/text()`,
	}
	for _, expr := range exprs {
		nodes, err := utils.FindNodes(doc, expr, true)
		if err == nil {
			return strings.TrimSpace(nodes[0].Data), nil
		}
	}
	return "unknown", errors.ErrorNotFoundFirstDate
}

func (p *UKProductParser) ParseSellerId(node *html.Node) (string, error) {
	expr := `//input[@id='deliveryBlockSelectMerchant']/@value`
	nodes, err := utils.FindNodes(node, expr, true)
	if err != nil {
		return "unknown", err
	}

	sellerId := htmlquery.SelectAttr(nodes[0], "value")
	if sellerId == "" {
		return "unknown", errors.ErrorNotFoundSellerId
	}
	return sellerId, nil
}

func (p *UKProductParser) ParseCategoryId(node *html.Node) (string, error) {
	var categoryId string
	var nodes []*html.Node
	var err error

	exprs := []string{
		`//tbody/tr/th[contains(text(), 'Best Sellers Rank')]/following-sibling::td/span/span/a/@href`,
		`//span[contains(text(), 'Best Sellers Rank')]/../ul/li/span/a/@href`,
	}

	for _, expr := range exprs {
		nodes, err = utils.FindNodes(node, expr, true)
		if err == nil {
			break
		}
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

func (p *UKProductParser) ParseHasCart(doc *html.Node) (string, error) {
	nodes, err := utils.FindNodes(doc, `//input[@id='add-to-cart-button']`, true)
	if err != nil || len(nodes) == 0 {
		return "false", err
	}
	return "true", nil
}

func (p *UKProductParser) ParseCoupon(doc *html.Node) (string, error) {
	return "unknown", nil
}

func (p *UKProductParser) ParseColor(doc *html.Node) (string, error) {
	exprs := []string{
		`//label[contains(text(),'Colour Name')]/following-sibling::span/text()`,
	}

	for _, expr := range exprs {
		nodes, err := utils.FindNodes(doc, expr, true)
		if err == nil || len(nodes) != 0 {
			return strings.TrimSpace(nodes[0].Data), nil
		}
	}
	return "unknown", errors.ErrorNotFoundColor
}

func (p *UKProductParser) ParseSize(doc *html.Node) (string, error) {
	exprs := []string{
		`//label[contains(text(),'Size Name')]/following-sibling::span/text()`,
	}

	for _, expr := range exprs {
		nodes, err := utils.FindNodes(doc, expr, true)
		if err == nil || len(nodes) != 0 {
			return strings.TrimSpace(nodes[0].Data), nil
		}
	}
	return "unknown", errors.ErrorNotFoundSize
}

func (p *UKProductParser) ParseSpecs(doc *html.Node) ([]string, error) {
	var specs []string

	expr := `//select[@id='native_dropdown_selected_color_name']`
	_, err := utils.FindNodes(doc, expr, false)
	if err == nil {
		colorExpr := `//option[contains(@id, 'native_color_name')]`
		nodes, err := utils.FindNodes(doc, colorExpr, true)
		if err == nil || len(nodes) != 0 {
			for _, node := range nodes {
				asinStr := htmlquery.SelectAttr(node, "value")
				if asinStr != "" {
					strs := strings.Split(asinStr, ",")
					if len(strs) > 1 {
						specs = append(specs, strs[1])
					}
				}
			}
		}
	} else {
		colorExpr := `//li[contains(@id, 'color_name')]`
		nodes, err := utils.FindNodes(doc, colorExpr, true)
		if err == nil || len(nodes) != 0 {
			for _, node := range nodes {
				asin := htmlquery.SelectAttr(node, "data-csa-c-item-id")
				if asin != "" {
					specs = append(specs, asin)
				}
			}
		}
	}

	// parse sizes
	expr = `//select[@id='native_dropdown_selected_size_name']`
	_, err = utils.FindNodes(doc, expr, false)
	if err == nil {
		sizeExpr := `//option[contains(@id, 'native_size_name')]`
		nodes, err := utils.FindNodes(doc, sizeExpr, true)
		if err == nil || len(nodes) != 0 {
			for _, node := range nodes {
				asinStr := htmlquery.SelectAttr(node, "value")
				if asinStr != "" {
					strs := strings.Split(asinStr, ",")
					if len(strs) > 1 {
						specs = append(specs, strs[1])
					}
				}
			}
		}
	} else {
		sizeExpr := `//li[contains(@id, 'size_name')]`
		nodes, err := utils.FindNodes(doc, sizeExpr, true)
		if err == nil || len(nodes) != 0 {
			for _, node := range nodes {
				asin := htmlquery.SelectAttr(node, "data-csa-c-item-id")
				if asin != "" {
					specs = append(specs, asin)
				}
			}
		}
	}
	return removeDuplicates(specs), nil
}

func (p *UKProductParser) ParseDescription(doc *html.Node) (string, error) {
	expr := `//h1[contains(text(), 'About this item')]/following-sibling::ul[1]/li/span/text()`
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

func (p *UKProductParser) ParseDeliveryTime(doc *html.Node) (string, error) {
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

func (p *UKProductParser) ParseFastestDelivery(doc *html.Node) (string, error) {
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

func (p *UKProductParser) ParsePrimePrice(doc *html.Node) (string, error) {
	return "unknown", nil
}

func (p *UKProductParser) ParseBrand(doc *html.Node) (string, error) {
	return "unknown", nil
}

func (p *UKProductParser) ParseCategoryHierarchy(doc *html.Node) ([]string, error) {
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
