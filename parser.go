package goamzparser

import (
	"github.com/antchfx/htmlquery"
	"github.com/microsuite/go-amz-parser/errors"
	"github.com/microsuite/go-amz-parser/internal/board"
	"github.com/microsuite/go-amz-parser/internal/category"
	"github.com/microsuite/go-amz-parser/internal/keyword"
	"github.com/microsuite/go-amz-parser/internal/product"
	"github.com/microsuite/go-amz-parser/internal/review"
	"github.com/microsuite/go-amz-parser/internal/seller"
	"github.com/microsuite/go-amz-parser/utils"
	"golang.org/x/net/html"
)

type ProductParser interface {
	// ParseASIN parses the ASIN from the given HTML document.
	ParseASIN(doc *html.Node) (string, error)

	// ParseStar parses the star from the given HTML document.
	ParseStar(doc *html.Node) (string, error)

	// ParseRating parses the rating from the given HTML document.
	ParseRating(doc *html.Node) (string, error)

	// ParseTitle parses the title from the given HTML document.
	ParseTitle(doc *html.Node) (string, error)

	// ParseImg parses the image from the given HTML document.
	ParseImg(doc *html.Node) (string, error)

	// ParsePrice parses the price from the given HTML document.
	ParsePrice(doc *html.Node) (string, error)

	// ParseDispatchFrom parses the dispatch from the given HTML document.
	ParseDispatchFrom(doc *html.Node) (string, error)

	// ParseSoldBy parses the sold by from the given HTML document.
	ParseSoldBy(doc *html.Node) (string, error)

	// ParseProductDimensions parses the product dimensions from the given HTML document.
	ParseProductDimensions(doc *html.Node) (string, error)

	// ParsePackageDimensions parses the package dimensions from the given HTML document.
	ParsePackageDimensions(doc *html.Node) (string, error)

	// ParsePackageWeight parses the package weight from the given HTML document.
	ParsePackageWeight(doc *html.Node) (string, error)

	// ParseProductWeight parses the product weight from the given HTML document.
	ParseProductWeight(doc *html.Node) (string, error)

	// ParseFirstAvailDate parses the first available date from the given HTML document.
	ParseFirstAvailDate(doc *html.Node) (string, error)

	// ParseSellerId parses the seller id from the given HTML document.
	ParseSellerId(node *html.Node) (string, error)

	// ParseCategoryId parses the category id from the given HTML document.
	ParseCategoryId(node *html.Node) (string, error)

	// ParseHasCart parses the cart from the given HTML document.
	ParseHasCart(doc *html.Node) (string, error)

	// ParseCoupon parses the coupon from the given HTML document.
	ParseCoupon(doc *html.Node) (string, error)

	// ParseColor parses the color from the given HTML document.
	ParseColor(doc *html.Node) (string, error)

	// ParseSize parses the size from the given HTML document.
	ParseSize(doc *html.Node) (string, error)

	// ParseSpecs parses the specs from the given HTML document.
	ParseSpecs(doc *html.Node) ([]string, error)

	// ParseDescription parses the description from the given HTML document.
	ParseDescription(doc *html.Node) (string, error)

	// ParseDeliveryTime parses the delivery time from the given HTML document.
	ParseDeliveryTime(doc *html.Node) (string, error)

	// ParseFastestDelivery parses the fastest delivery from the given HTML document.
	ParseFastestDelivery(doc *html.Node) (string, error)

	// ParsePrimePrice parses the prime price from the given HTML document.
	ParsePrimePrice(doc *html.Node) (string, error)

	// ParseBrand parses the brand from the given HTML document.
	ParseBrand(doc *html.Node) (string, error)

	// ParseCategoryHierarchy parses the category hierarchy from the given HTML document.
	ParseCategoryHierarchy(doc *html.Node) ([]string, error)

	// ParseCustomerReviews parses the customer reviews from the given HTML document.
	ParseCustomerReviews(doc *html.Node) (map[string]string, error)
}

type KeywordParser interface {
	// ParseAllProducts parses all products from the given HTML document.
	ParseAllProducts(doc *html.Node) ([]*html.Node, error)

	// ParseCurrentPageIndex parses the current page index from the given HTML document.
	ParseCurrentPageIndex(doc *html.Node) (string, error)

	// ParseNextPageURL parses the next page url from the given HTML document.
	ParseNextPageURL(doc *html.Node) (string, error)

	// ParseKeyword parses the keyword from the given HTML document.
	ParseKeyword(doc *html.Node) (string, error)

	// ParseASIN parses the ASIN from the given HTML node.
	ParseASIN(node *html.Node) (string, error)

	// ParsePrice parses the price from the given HTML node.
	ParsePrice(node *html.Node) (string, error)

	// ParseStar parses the star from the given HTML node.
	ParseStar(node *html.Node) (string, error)

	// ParseRating parses the rating from the given HTML node.
	ParseRating(node *html.Node) (string, error)

	// ParseSponsered parses the sponsered from the given HTML node.
	ParseSponsered(node *html.Node) (string, error)

	// ParsePrime parses the prime from the given HTML node.
	ParsePrime(node *html.Node) (string, error)

	// ParseSales parses the sales from the given HTML node.
	ParseSales(node *html.Node) (string, error)

	// ParseImg parses the img from the given HTML node.
	ParseImg(node *html.Node) (string, error)

	// ParseTitle parses the title from the given HTML node.
	ParseTitle(node *html.Node) (string, error)
}

type CategoryParser interface {
	// ParseAllProducts parses all products from the given HTML document.
	ParseAllProducts(doc *html.Node) ([]*html.Node, error)

	// ParseCurrentPageIndex parses the current page index from the given HTML document.
	ParseCurrentPageIndex(doc *html.Node) (string, error)

	// ParseMaxPageNum parses the max page number from the given HTML document.
	ParseMaxPageNum(doc *html.Node) (string, error)

	// ParseCurrentPageIndex parses the next page url from the given HTML document.
	ParseNextPageURL(doc *html.Node) (string, error)

	// ParseContentId parses the content id from the given HTML document.
	ParseContentId(doc *html.Node) (string, error)

	// ParseContentLink parses the content link from the given HTML document.
	ParseContentLink(doc *html.Node) (string, error)

	// ParsePagination parses the pagination from the given HTML document.
	ParsePagination(doc *html.Node) (string, error)

	// ParseCategoryName parses the category name from the given HTML document.
	ParseCategoryName(doc *html.Node) (string, error)

	// ParseASIN parses the ASIN from the given HTML node.
	ParseASIN(node *html.Node) (string, error)

	// ParsePrice parses the price from the given HTML node.
	ParsePrice(node *html.Node) (string, error)

	// ParseStar parses the star from the given HTML node.
	ParseStar(node *html.Node) (string, error)

	// ParseImg parses the image from the given HTML node.
	ParseImg(node *html.Node) (string, error)

	// ParseTitle parses the title from the given HTML node.
	ParseTitle(node *html.Node) (string, error)
}

type SellerParser interface {
	// ParseAllProducts parses all products from the given HTML document.
	ParseAllProducts(doc *html.Node) ([]*html.Node, error)

	// ParseCurrentPageIndex parses the current page index from the given HTML document.
	ParseCurrentPageIndex(doc *html.Node) (string, error)

	// ParseMaxPageNum parses the max page number from the given HTML document.
	ParseMaxPageNum(doc *html.Node) (string, error)

	// ParseNextPageURL parses the next page url from the given HTML document.
	ParseNextPageURL(doc *html.Node) (string, error)

	// ParseContentId parses the content id from the given HTML document.
	ParseContentId(doc *html.Node) (string, error)

	// ParseContentLink parses the content link from the given HTML document.
	ParseContentLink(doc *html.Node) (string, error)

	// ParsePagination parses the pagination from the given HTML document.
	ParsePagination(doc *html.Node) (string, error)

	// ParseASIN parses the ASIN from the given html node.
	ParseASIN(node *html.Node) (string, error)

	// ParsePrice parses the price from the give html node.
	ParsePrice(node *html.Node) (string, error)

	// ParseStar parses the star from the give html node.
	ParseStar(node *html.Node) (string, error)

	// ParseImg parses the img from the give html node.
	ParseImg(node *html.Node) (string, error)

	// ParseTitle parses the title from the given HTML node.
	ParseTitle(node *html.Node) (string, error)
}

type BoardParser interface {
	// ParseAllProducts parses all products from the given HTML document.
	ParseAllProducts(doc *html.Node) ([]*html.Node, error)

	// ParseNextPageURL parses the next page url from the given HTML document.
	ParseNextPageURL(doc *html.Node) (string, error)

	// ParseRecsList parses the recs list from the give html document.
	ParseRecsList(doc *html.Node) (string, error)

	// ParseReftag parses the ref tag from the give html document.
	ParseReftag(doc *html.Node) (string, error)

	// ParseOffset parses the offset from the give html document.
	ParseOffset(doc *html.Node) (string, error)

	// ParseAcpParam parses the acp param from the give html document.
	ParseAcpParam(doc *html.Node) (string, error)

	// ParseAcpPath parses the acp path from the give html document.
	ParseAcpPath(doc *html.Node) (string, error)

	// ParseBestSellersCategory parses the best seller category from the give html document.
	ParseBestSellersCategory(doc *html.Node) (string, error)

	// ParseNewReleasesCategory parses the new release category from the give html document.
	ParseNewReleasesCategory(doc *html.Node) (string, error)

	// ParseASIN parses the ASIN from the given html node.
	ParseASIN(node *html.Node) (string, error)

	// ParsePrice parses the price from the give html node.
	ParsePrice(node *html.Node) (string, error)

	// ParseStar parses the star from the give html node.
	ParseStar(node *html.Node) (string, error)

	// ParseStar parses the rating from the give html node.
	ParseRating(node *html.Node) (string, error)

	// ParseTitle parses the title from the give html node.
	ParseTitle(node *html.Node) (string, error)

	// ParseRank parses the rank from the give html node.
	ParseRank(node *html.Node) (string, error)
}

type AmzReviewParser interface {
	// ParseAllReviews parses all reviews from the given HTML document.
	ParseAllReviews(doc *html.Node) ([]*html.Node, error)

	// ParseReviewer parses reviewer from the given HTML node.
	ParseReviewer(node *html.Node) (string, error)

	// ParseReviewerLink parses reviewer link from the given HTML node.
	ParseReviewerLink(node *html.Node) (string, error)

	// ParseStar parses the star from the give html node.
	ParseStar(node *html.Node) (string, error)

	// ParseTitle parses the title from the give html node.
	ParseTitle(node *html.Node) (string, error)

	// ParseDate parses the date from the give html node.
	ParseDate(node *html.Node) (string, error)

	// ParsePurchase parses whether it has been purchased from the give html node.
	ParsePurchase(node *html.Node) (string, error)

	// ParseContent parses the content from the give html node.
	ParseContent(node *html.Node) (string, error)
}

type Parser struct {
	productParserMap  map[string]ProductParser
	keywordParserMap  map[string]KeywordParser
	categoryParserMap map[string]CategoryParser
	sellerParserMap   map[string]SellerParser
	boardSellerMap    map[string]BoardParser
	reviewParserMap   map[string]AmzReviewParser
}

func NewParser() *Parser {
	p := &Parser{
		productParserMap:  make(map[string]ProductParser),
		keywordParserMap:  make(map[string]KeywordParser),
		categoryParserMap: make(map[string]CategoryParser),
		sellerParserMap:   make(map[string]SellerParser),
		boardSellerMap:    make(map[string]BoardParser),
		reviewParserMap:   make(map[string]AmzReviewParser),
	}
	p.registerParsers()
	return p
}

// registerProductParser registers a product parser for a given
func (p *Parser) registerProductParser(region string, parser ProductParser) {
	p.productParserMap[region] = parser
}

func (p *Parser) GetProductParser(region string) ProductParser {
	return p.productParserMap[region]
}

func (p *Parser) registerKeywordParser(region string, parser KeywordParser) {
	p.keywordParserMap[region] = parser
}

func (p *Parser) GetKeywordParser(region string) KeywordParser {
	return p.keywordParserMap[region]
}

func (p *Parser) registerCategoryParser(region string, parser CategoryParser) {
	p.categoryParserMap[region] = parser
}

func (p *Parser) GetCategoryParser(region string) CategoryParser {
	return p.categoryParserMap[region]
}

func (p *Parser) registerSellerParser(region string, parser SellerParser) {
	p.sellerParserMap[region] = parser
}

func (p *Parser) GetSellerParser(region string) SellerParser {
	return p.sellerParserMap[region]
}

func (p *Parser) registerBoardParser(region string, parser BoardParser) {
	p.boardSellerMap[region] = parser
}

func (p *Parser) GetBoardParser(region string) BoardParser {
	return p.boardSellerMap[region]
}

func (p *Parser) registerReviewParser(region string, parser AmzReviewParser) {
	p.reviewParserMap[region] = parser
}

func (p *Parser) GetReviewParser(region string) AmzReviewParser {
	return p.reviewParserMap[region]
}

func (p *Parser) registerParsers() {
	// Register product parsers.
	p.registerProductParser(US, product.NewUSProductParser())
	p.registerProductParser(UK, product.NewUKProductParser())
	p.registerProductParser(DE, product.NewDEProductParser())
	p.registerProductParser(FR, product.NewFRProductParser())

	// Register keyword parsers.
	p.registerKeywordParser(US, keyword.NewUSKeywordParser())
	p.registerKeywordParser(UK, keyword.NewUKKeywordParser())
	p.registerKeywordParser(DE, keyword.NewDEKeywordParser())
	p.registerKeywordParser(FR, keyword.NewFRKeywordParser())

	// Register category parsers.
	p.registerCategoryParser(US, category.NewUSCategoryParser())
	p.registerCategoryParser(UK, category.NewUKCategoryParser())
	p.registerCategoryParser(DE, category.NewDECategoryParser())
	p.registerCategoryParser(FR, category.NewFRCategoryParser())

	// Register seller parsers.
	p.registerSellerParser(US, seller.NewUSSellerParser())
	p.registerSellerParser(UK, seller.NewUKSellerParser())
	p.registerSellerParser(DE, seller.NewDESellerParser())
	p.registerSellerParser(FR, seller.NewFRSellerParser())

	// Register board parsers.
	p.registerBoardParser(US, board.NewUSBoardParser())
	p.registerBoardParser(UK, board.NewUKBoardParser())
	p.registerBoardParser(DE, board.NewDEBoardParser())
	p.registerBoardParser(FR, board.NewFRBoardParser())

	// Register review parsers.
	p.registerReviewParser(US, review.NewUSReviewParser())
	p.registerReviewParser(UK, review.NewUKReviewParser())
	p.registerReviewParser(DE, review.NewDEReviewParser())
	p.registerReviewParser(FR, review.NewFRReviewParser())
}

func ParseRegion(doc *html.Node) (string, error) {
	langExpr := "/html/@lang"
	langNodes, err := utils.FindNodes(doc, langExpr, false)
	if err != nil {
		return "unknown", err
	}

	lang := htmlquery.SelectAttr(langNodes[0], "lang")
	if lang == "" {
		return "unknown", errors.ErrorNotFoundLanguage
	}
	return lang, nil
}
