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

type FRProductParser struct{}

func (p *FRProductParser) ParseASIN(doc *html.Node) (string, error) {
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

func (p *FRProductParser) ParseStar(doc *html.Node) (string, error) {
	expr := `//span[contains(text(),'sur 5')]/text()`
	nodes, err := utils.FindNodes(doc, expr, true)
	if err != nil {
		return "unknown", err
	}
	star := nodes[0].Data
	star = strings.TrimSpace(star)
	return utils.FormatNumberEuro(strings.Split(star, " ")[0]), nil
}

func (p *FRProductParser) ParseRating(doc *html.Node) (string, error) {
	expr := `//a[@id='acrCustomerReviewLink']/span[@id='acrCustomerReviewText']/text()`
	nodes, err := utils.FindNodes(doc, expr, true)
	if err != nil {
		return "unknown", err
	}
	rating := utils.FindNumberHead(strings.TrimSpace(nodes[0].Data))
	return strings.TrimSpace(utils.FormatNumberEuro(rating)), nil
}

func (p *FRProductParser) ParseTitle(doc *html.Node) (string, error) {
	expr := `//span[@id='productTitle']/text()`
	nodes, err := utils.FindNodes(doc, expr, true)
	if err != nil {
		return "unknown", err
	}
	return utils.FormatTitle(nodes[0].Data), nil
}

func (p *FRProductParser) ParseImg(doc *html.Node) (string, error) {
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

func (p *FRProductParser) ParsePrice(doc *html.Node) (string, error) {
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

func (p *FRProductParser) ParseDispatchFrom(doc *html.Node) (string, error) {
	var err error
	exprs := []string{
		`//div/span[contains(text(), 'Expédié par')]/following-sibling::span/text()`, // ProductOfCategory
		`//div/span[contains(text(), 'Expédié par')]/../../following-sibling::div/div/span/text()`,
	}

	for _, expr := range exprs {
		nodes, err := utils.FindNodes(doc, expr, true)
		if err == nil {
			return strings.TrimSpace(nodes[0].Data), nil
		}
	}
	return "unknown", err
}

func (p *FRProductParser) ParseSoldBy(doc *html.Node) (string, error) {
	exprs := []string{
		`//div/span[contains(text(), 'Vendu par')]/../../following-sibling::div/div/span/a/text()`,
		`//div/span[contains(text(), 'Vendu par')]/../../following-sibling::div/div/span/text()`,
	}

	for _, expr := range exprs {
		nodes, err := utils.FindNodes(doc, expr, true)
		if err == nil {
			return strings.TrimSpace(nodes[0].Data), nil
		}
	}
	return "unknown", errors.ErrorNotFoundSoldBy
}

func (p *FRProductParser) ParsePackageDimensions(doc *html.Node) (string, error) {
	exprs := []string{
		`//th[contains(text(), 'Dimensions du colis')]/following-sibling::td/text()`,
		`//span[contains(text(), 'Dimensions du colis')]/following-sibling::span/text()`,
	}
	for _, expr := range exprs {
		nodes, err := utils.FindNodes(doc, expr, true)
		if err == nil && len(nodes) > 0 {
			return strings.TrimSpace(nodes[0].Data), nil
		}
	}
	return "unknown", errors.ErrorNotFoundPackageDimensions
}

func (p *FRProductParser) ParsePackageWeight(doc *html.Node) (string, error) {
	exprs := []string{
		`//th[contains(text(), 'Artikelgewicht')]/following-sibling::td/text()`,
		`//span[contains(text(), 'Artikelgewicht')]/following-sibling::span/text()`,
	}
	for _, expr := range exprs {
		nodes, err := utils.FindNodes(doc, expr, true)
		if err == nil && len(nodes) > 0 {
			return strings.TrimSpace(nodes[0].Data), nil
		}
	}
	return "unknown", errors.ErrorNotFoundPackageWeight
}

func (p *FRProductParser) ParseFirstAvailDate(doc *html.Node) (string, error) {
	exprs := []string{
		`//tbody/tr/th[contains(text(), 'Date de mise en ligne sur Amazon.fr')]/following-sibling::td/text()`,
		`//span[contains(text(), 'Date de mise en ligne sur Amazon.fr')]/following-sibling::span/text()`,
	}
	for _, expr := range exprs {
		nodes, err := utils.FindNodes(doc, expr, true)
		if err == nil && len(nodes) > 0 {
			return strings.TrimSpace(nodes[0].Data), nil
		}
	}
	return "unknown", errors.ErrorNotFoundFirstDate
}

func (p *FRProductParser) ParseSellerId(node *html.Node) (string, error) {
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

func (p *FRProductParser) ParseCategoryId(node *html.Node) (string, error) {
	var err error
	var categoryId string
	var nodes []*html.Node

	exprs := []string{
		`//th[contains(text(), "Classement des meilleures ventes d'Amazon")]/following-sibling::td/span/span/a/@href`,
		`//span[contains(text(), "Classement des meilleures ventes d'Amazon")]/following-sibling::ul//a/@href`,
	}

	for _, expr := range exprs {
		nodes, err = utils.FindNodes(node, expr, true)
		if err == nil {
			continue
		}
	}
	if len(nodes) == 0 {
		return "unknown", fmt.Errorf("no nodes found")
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

func (p *FRProductParser) ParseProductDimensions(doc *html.Node) (string, error) {
	exprs := []string{
		`//tbody/tr/th[contains(text(), 'Dimensions du produit')]/following-sibling::td/text()`,
		`//span[contains(text(), 'Dimensions du produit')]/following-sibling::span/text()`,
	}
	for _, expr := range exprs {
		nodes, err := utils.FindNodes(doc, expr, true)
		if err == nil && len(nodes) > 0 {
			return strings.TrimSpace(nodes[0].Data), nil
		}
	}
	return "unknown", errors.ErrorNotFoundDimensions
}

func (p *FRProductParser) ParseProductWeight(doc *html.Node) (string, error) {
	exprs := []string{
		`//th[contains(text(), "Poids")]/following-sibling::td/text()`,
		`//span[contains(text(), "Poids")]/following-sibling::span/text()`,
	}
	for _, expr := range exprs {
		nodes, err := utils.FindNodes(doc, expr, true)
		if err == nil && len(nodes) > 0 {
			return strings.TrimSpace(nodes[0].Data), nil
		}
	}
	return "unknown", errors.ErrorNotFoundWeight
}

func (p *FRProductParser) ParseHasCart(doc *html.Node) (string, error) {
	nodes, err := utils.FindNodes(doc, `//input[@id='add-to-cart-button']`, true)
	if err != nil || len(nodes) == 0 {
		return "false", err
	}
	return "true", nil
}

func (p *FRProductParser) ParseCoupon(doc *html.Node) (string, error) {
	return "unknown", nil
}

func (p *FRProductParser) ParseColor(doc *html.Node) (string, error) {
	exprs := []string{
		`//label[contains(text(),'Couleur:')]/following-sibling::span/text()`,
		`//span[contains(text(), 'Couleur')]/../following-sibling::td/span/text()`,
	}

	for _, expr := range exprs {
		nodes, err := utils.FindNodes(doc, expr, true)
		if err == nil || len(nodes) != 0 {
			return strings.TrimSpace(nodes[0].Data), nil
		}
	}
	return "unknown", errors.ErrorNotFoundColor
}

func (p *FRProductParser) ParseSize(doc *html.Node) (string, error) {
	exprs := []string{
		`//label[contains(text(),'Taille:')]/following-sibling::span/text()`,
		`//span[contains(text(), 'Taille')]/../following-sibling::td/span/text()`,
	}

	for _, expr := range exprs {
		nodes, err := utils.FindNodes(doc, expr, true)
		if err == nil || len(nodes) != 0 {
			return strings.TrimSpace(nodes[0].Data), nil
		}
	}
	return "unknown", errors.ErrorNotFoundSize
}

func (p *FRProductParser) ParseSpecs(doc *html.Node) ([]string, error) {
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

func (p *FRProductParser) ParseDescription(doc *html.Node) (string, error) {
	expr := `//h1[contains(text(), 'À propos de cet article')]/following-sibling::ul[1]/li/span/text()`
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

func (p *FRProductParser) ParseDeliveryTime(doc *html.Node) (string, error) {
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

func (p *FRProductParser) ParseFastestDelivery(doc *html.Node) (string, error) {
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

func (p *FRProductParser) ParsePrimePrice(doc *html.Node) (string, error) {
	return "unknown", nil
}

func (p *FRProductParser) ParseBrand(doc *html.Node) (string, error) {
	return "unknown", nil
}

func (p *FRProductParser) ParseCategoryHierarchy(doc *html.Node) ([]string, error) {
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
