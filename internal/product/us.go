package product

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"

	"github.com/microsuite/go-amz-parser/errors"
	"github.com/microsuite/go-amz-parser/utils"
)

type USProductParser struct{}

func (p *USProductParser) ParseASIN(doc *html.Node) (string, error) {
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

func (p *USProductParser) ParseStar(doc *html.Node) (string, error) {
	expr := `//span[contains(text(),'out of 5 stars')]/text()`
	nodes, err := utils.FindNodes(doc, expr, true)
	if err != nil {
		return "unknown", err
	}
	star := nodes[0].Data
	star = strings.TrimSpace(star)
	star = utils.FindNumberHead(star)
	return utils.FormatNumber(star), nil
}

func (p *USProductParser) ParseRating(doc *html.Node) (string, error) {
	expr := `//a[@id='acrCustomerReviewLink']/span[@id='acrCustomerReviewText']/text()`
	nodes, err := utils.FindNodes(doc, expr, true)
	if err != nil {
		return "unknown", err
	}
	rating := utils.FindNumberHead(strings.TrimSpace(nodes[0].Data))
	return utils.FormatNumber(rating), nil
}

func (p *USProductParser) ParseTitle(doc *html.Node) (string, error) {
	expr := `//span[@id="productTitle"]/text()`
	nodes, err := utils.FindNodes(doc, expr, true)
	if err != nil {
		return "unknown", err
	}
	return utils.FormatTitle(nodes[0].Data), nil
}

func (p *USProductParser) ParseImg(doc *html.Node) (string, error) {
	expr := `//div[@id="imageBlock"]//div[@class="imgTagWrapper"]/img/@src`
	nodes, err := utils.FindNodes(doc, expr, true)
	if err != nil {
		return "unknown", err
	}

	img := htmlquery.SelectAttr(nodes[0], "src")
	if img == "" {
		return "unknown", fmt.Errorf("img url is empty")
	}
	return img, nil
}

func (p *USProductParser) ParsePrice(doc *html.Node) (string, error) {
	var price string
	exprs := []string{
		`//span[starts-with(@class, 'a-price') and @data-a-color="price"]/span/text()`,
		`//div[starts-with(@id, "corePrice") and @data-csa-c-asin]/div/span/text()`,
		`//span[starts-with(@id, "a-price") and @data-a-color="price"]/span/text()`,
	}

	for _, expr := range exprs {
		nodes, err := utils.FindNodes(doc, expr, true)
		if err == nil {
			splits := strings.Split(strings.TrimSpace(nodes[0].Data), " ")
			if len(splits) > 0 && splits[0] != "" {
				price = utils.FormatNumber(utils.DropMoneySym(splits[0]))
			}
			if price == "" && len(splits) > 1 && splits[1] != "" {
				price = utils.FormatNumber(utils.DropMoneySym(splits[1]))
			}
			return price, nil
		}
	}
	return "0.0", errors.ErrorNotFoundPrice
}

func (p *USProductParser) ParseDispatchFrom(doc *html.Node) (string, error) {
	exprs := []string{
		`//div/span[contains(text(), 'Ships from')]/following-sibling::span/text()`,               // ProductOfCategory
		`//div/span[contains(text(), 'Ships from')]/../../following-sibling::div/div/span/text()`, // ProductOfSeller
	}

	for _, expr := range exprs {
		nodes, err := utils.FindNodes(doc, expr, true)
		if err == nil {
			return strings.TrimSpace(nodes[0].Data), nil
		}
	}
	return "unknown", errors.ErrorNotFoundDispatchFrom
}

func (p *USProductParser) ParseSoldBy(doc *html.Node) (string, error) {
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

func (p *USProductParser) ParseProductDimensions(doc *html.Node) (string, error) {
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

func (p *USProductParser) ParsePackageDimensions(doc *html.Node) (string, error) {
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

func (p *USProductParser) ParsePackageWeight(doc *html.Node) (string, error) {
	expr := `//tbody/tr/th[contains(text(), "Package Weight")]/following-sibling::td/text()`
	nodes, err := utils.FindNodes(doc, expr, true)
	if err != nil {
		return "unknown", err
	}
	return strings.TrimSpace(nodes[0].Data), nil
}

func (p *USProductParser) ParseProductWeight(doc *html.Node) (string, error) {
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

func (p *USProductParser) ParseFirstAvailDate(doc *html.Node) (string, error) {
	exprs := []string{
		`//th[contains(text(), 'Date First Available')]/following-sibling::td/text()`,
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

func (p *USProductParser) ParseSellerId(node *html.Node) (string, error) {
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

func (p *USProductParser) ParseCategoryId(node *html.Node) (string, error) {
	var categoryId string
	var nodes []*html.Node
	var err error

	exprs := []string{
		`//th[contains(text(), 'Best Sellers Rank')]/following-sibling::td/span/span/a/@href`,
		`//span[contains(text(), 'Best Sellers Rank')]/../ul/li/span/a/@href`,
	}

	for _, expr := range exprs {
		nodes, err = utils.FindNodes(node, expr, true)
		if err == nil {
			break
		}
		fmt.Println(err)
	}

	for _, node := range nodes {
		categoryId = htmlquery.SelectAttr(node, "href")
		if categoryId == "" {
			return "unknown", errors.ErrorNotFoundContentId
		}
	}

	regex := regexp.MustCompile(`\d+`)
	mathches := regex.FindStringSubmatch(categoryId)
	for _, match := range mathches {
		return match, nil
	}
	return "unknown", nil
}

func (p *USProductParser) ParseHasCart(doc *html.Node) (string, error) {
	nodes, err := utils.FindNodes(doc, `//input[contains(@id, 'add-to-cart-button')]`, true)
	if err != nil || len(nodes) == 0 {
		return "false", err
	}
	return "true", nil
}

func (p *USProductParser) ParseCoupon(doc *html.Node) (string, error) {
	expr := `//i[contains(text(),'Coupon')]/following-sibling::span/label/text()`

	nodes, err := utils.FindNodes(doc, expr, true)
	if err != nil || len(nodes) == 0 {
		return "unknown", err
	}
	return strings.TrimSpace(nodes[0].Data), nil
}

func (p *USProductParser) ParseColor(doc *html.Node) (string, error) {
	expr := `//label[contains(text(),'Color:')]/following-sibling::span/text()`

	nodes, err := utils.FindNodes(doc, expr, true)
	if err != nil || len(nodes) == 0 {
		return "unknown", err
	}
	return strings.TrimSpace(nodes[0].Data), nil
}

func (p *USProductParser) ParseSize(doc *html.Node) (string, error) {
	exprs := []string{
		`//label[contains(text(),'Size:')]/following-sibling::span/text()`,
		`//span[contains(text(),'Size')]/../following-sibling::td/span/text()`,
	}

	for _, expr := range exprs {
		nodes, err := utils.FindNodes(doc, expr, true)
		if err == nil || len(nodes) != 0 {
			return strings.TrimSpace(nodes[0].Data), nil
		}
	}
	return "unknown", fmt.Errorf("no size found")
}

func (p *USProductParser) ParseSpecs(doc *html.Node) ([]string, error) {
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

func removeDuplicates(strs []string) []string {
	m := make(map[string]bool)
	result := []string{}

	for _, str := range strs {
		if !m[str] {
			m[str] = true
			result = append(result, str)
		}
	}
	return result
}

func (p *USProductParser) ParseDescription(doc *html.Node) (string, error) {
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
		return "unknown", fmt.Errorf("desc is empty")
	}
	return desc, nil
}

func (p *USProductParser) ParseDeliveryTime(doc *html.Node) (string, error) {
	expr := `//div[@id='mir-layout-DELIVERY_BLOCK-slot-PRIMARY_DELIVERY_MESSAGE_LARGE']/span/@data-csa-c-delivery-time`
	nodes, err := utils.FindNodes(doc, expr, true)
	if err != nil {
		return "unknown", err
	}

	fastestDelivery := htmlquery.SelectAttr(nodes[0], "data-csa-c-delivery-time")
	if fastestDelivery == "" {
		return "unknown", fmt.Errorf("delivery time is empty")
	}
	return fastestDelivery, nil
}

func (p *USProductParser) ParseFastestDelivery(doc *html.Node) (string, error) {
	expr := `//div[@id='mir-layout-DELIVERY_BLOCK-slot-SECONDARY_DELIVERY_MESSAGE_LARGE']/span/@data-csa-c-delivery-time`
	nodes, err := utils.FindNodes(doc, expr, true)
	if err != nil {
		return "unknown", err
	}

	fastestDelivery := htmlquery.SelectAttr(nodes[0], "data-csa-c-delivery-time")
	if fastestDelivery == "" {
		return "unknown", fmt.Errorf("fastest delivery is empty")
	}
	return fastestDelivery, nil
}

func (p *USProductParser) ParsePrimePrice(doc *html.Node) (string, error) {
	expr := `//span[contains(text(), 'Join Prime to buy this item at')]/following-sibling::span/text()`
	nodes, err := utils.FindNodes(doc, expr, true)
	if err != nil {
		return "unknown", err
	}
	return strings.TrimSpace(nodes[0].Data), nil
}

func (p *USProductParser) ParseBrand(doc *html.Node) (string, error) {
	exprs := []string{
		`//span[contains(text(),'Brand')]/../following-sibling::td/span/text()`,
	}

	for _, expr := range exprs {
		nodes, err := utils.FindNodes(doc, expr, true)
		if err == nil || len(nodes) != 0 {
			return strings.TrimSpace(nodes[0].Data), nil
		}
	}
	return "unknown", fmt.Errorf("no size found")
}

func (p *USProductParser) ParseCategoryHierarchy(doc *html.Node) ([]string, error) {
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
