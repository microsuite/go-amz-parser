package goamzparser

import (
	"github.com/antchfx/htmlquery"
	"github.com/microsuite/go-amz-parser/errors"
	"github.com/microsuite/go-amz-parser/internal/board"
	"github.com/microsuite/go-amz-parser/internal/category"
	"github.com/microsuite/go-amz-parser/internal/keyword"
	"github.com/microsuite/go-amz-parser/internal/product"
	"github.com/microsuite/go-amz-parser/internal/seller"
	"github.com/microsuite/go-amz-parser/utils"
	"golang.org/x/net/html"
)

type ProductParser interface {
	ParseASIN(doc *html.Node) (string, error)
	ParseStar(doc *html.Node) (string, error)
	ParseRating(doc *html.Node) (string, error)
	ParseTitle(doc *html.Node) (string, error)
	ParseImg(doc *html.Node) (string, error)
	ParsePrice(doc *html.Node) (string, error)
	ParseDispatchFrom(doc *html.Node) (string, error)
	ParseSoldBy(doc *html.Node) (string, error)
	ParseProductDimensions(doc *html.Node) (string, error)
	ParsePackageDimensions(doc *html.Node) (string, error)
	ParsePackageWeight(doc *html.Node) (string, error)
	ParseProductWeight(doc *html.Node) (string, error)
	ParseFirstAvailDate(doc *html.Node) (string, error)
	ParseSellerId(node *html.Node) (string, error)
	ParseCategoryId(node *html.Node) (string, error)
	ParseHasCart(doc *html.Node) (string, error)
	ParseCoupon(doc *html.Node) (string, error)
	ParseColor(doc *html.Node) (string, error)
	ParseSize(doc *html.Node) (string, error)
	ParseSpecs(doc *html.Node) ([]string, error)
	ParseDescription(doc *html.Node) (string, error)
	ParseDeliveryTime(doc *html.Node) (string, error)
	ParseFastestDelivery(doc *html.Node) (string, error)
	ParsePrimePrice(doc *html.Node) (string, error)
	ParseBrand(doc *html.Node) (string, error)
	ParseCategoryHierarchy(doc *html.Node) ([]string, error)
}

type KeywordParser interface {
	ParseAllProducts(doc *html.Node) ([]*html.Node, error)
	ParseCurrentPageIndex(doc *html.Node) (string, error)
	ParseNextPageURL(doc *html.Node) (string, error)
	ParseASIN(node *html.Node) (string, error)
	ParsePrice(node *html.Node) (string, error)
	ParseStar(node *html.Node) (string, error)
	ParseRating(node *html.Node) (string, error)
	ParseSponsered(node *html.Node) (string, error)
	ParsePrime(node *html.Node) (string, error)
	ParseSales(node *html.Node) (string, error)
	ParseImg(node *html.Node) (string, error)
	ParseTitle(node *html.Node) (string, error)
}

type CategoryParser interface {
	ParseAllProducts(doc *html.Node) ([]*html.Node, error)
	ParseCurrentPageIndex(doc *html.Node) (string, error)
	ParseMaxPageNum(doc *html.Node) (string, error)
	ParseNextPageURL(doc *html.Node) (string, error)
	ParseContentId(doc *html.Node) (string, error)
	ParseContentLink(doc *html.Node) (string, error)
	ParsePagination(doc *html.Node) (string, error)
	ParseCategoryName(doc *html.Node) (string, error)
	ParseASIN(node *html.Node) (string, error)
	ParsePrice(node *html.Node) (string, error)
	ParseStar(node *html.Node) (string, error)
	ParseImg(node *html.Node) (string, error)
}

type SellerParser interface {
	ParseAllProducts(doc *html.Node) ([]*html.Node, error)
	ParseCurrentPageIndex(doc *html.Node) (string, error)
	ParseMaxPageNum(doc *html.Node) (string, error)
	ParseNextPageURL(doc *html.Node) (string, error)
	ParseContentId(doc *html.Node) (string, error)
	ParseContentLink(doc *html.Node) (string, error)
	ParsePagination(doc *html.Node) (string, error)
	ParseASIN(node *html.Node) (string, error)
	ParsePrice(node *html.Node) (string, error)
	ParseStar(node *html.Node) (string, error)
	ParseImg(node *html.Node) (string, error)
}

type BoardParser interface {
	ParseNextPageURL(doc *html.Node) (string, error)
}

type Parser struct {
	productParserMap  map[string]ProductParser
	keywordParserMap  map[string]KeywordParser
	categoryParserMap map[string]CategoryParser
	sellerParserMap   map[string]SellerParser
	boardSellerMap    map[string]BoardParser
}

func NewParser() *Parser {
	p := &Parser{
		productParserMap:  make(map[string]ProductParser),
		keywordParserMap:  make(map[string]KeywordParser),
		categoryParserMap: make(map[string]CategoryParser),
		sellerParserMap:   make(map[string]SellerParser),
		boardSellerMap:    make(map[string]BoardParser),
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

func (p *Parser) registerParsers() {
	// Register product parsers.
	p.registerProductParser(US, &product.USProductParser{})
	p.registerProductParser(UK, &product.UKProductParser{})
	p.registerProductParser(DE, &product.DEProductParser{})
	p.registerProductParser(FR, &product.FRProductParser{})

	// Register keyword parsers.
	p.registerKeywordParser(US, &keyword.USKeywordParser{})
	p.registerKeywordParser(UK, &keyword.UKKeywordParser{})
	p.registerKeywordParser(DE, &keyword.DEKeywordParser{})
	p.registerKeywordParser(FR, &keyword.FRKeywordParser{})

	// Register category parsers.
	p.registerCategoryParser(US, &category.USCategoryParser{})
	p.registerCategoryParser(UK, &category.UKCategoryParser{})
	p.registerCategoryParser(DE, &category.DECategoryParser{})
	p.registerCategoryParser(FR, &category.FRCategoryParser{})

	// Register seller parsers.
	p.registerSellerParser(US, &seller.USSellerParser{})
	p.registerSellerParser(UK, &seller.UKSellerParser{})
	p.registerSellerParser(DE, &seller.DESellerParser{})
	p.registerSellerParser(FR, &seller.FRSellerParser{})

	// Register seller parsers.
	p.registerBoardParser(US, &board.USBoardParser{})
	p.registerBoardParser(UK, &board.UKBoardParser{})
	p.registerBoardParser(DE, &board.DEBoardParser{})
	p.registerBoardParser(FR, &board.FRBoardParser{})
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
