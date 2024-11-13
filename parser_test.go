package goamzparser

import (
	"fmt"
	"testing"

	"github.com/antchfx/htmlquery"
)

func TestCategoryParser(t *testing.T) {
	p := NewParser()

	doc, err := htmlquery.LoadDoc("./input.html")
	if err != nil {
		t.Fatalf("Error loading document: %s\n", err.Error())
	}

	region, err := ParseRegion(doc)
	if err != nil {
		t.Fatalf("Error parsing region: %s\n", err.Error())
	}

	parser := p.GetCategoryParser(region)
	if parser == nil {
		t.Fatal("No parser found for region: " + region)
	}

	pageNum, err := parser.ParseMaxPageNum(doc)
	if err != nil {
		t.Errorf("Error parsing max page num: %s\n", err.Error())
	} else {
		fmt.Printf("Max page num: %s\n", pageNum)
	}

	nextPage, err := parser.ParseNextPageURL(doc)
	if err != nil {
		t.Errorf("Error parsing next page: %s\n", err.Error())
	} else {
		fmt.Printf("Next page url: %s\n", nextPage)
	}

	contentId, err := parser.ParseContentId(doc)
	if err != nil {
		t.Errorf("Error parsing content id: %s\n", err.Error())
	} else {
		fmt.Printf("Content id: %s\n", contentId)
	}

	contentLink, err := parser.ParseContentLink(doc)
	if err != nil {
		t.Errorf("Error parsing content link: %s\n", err.Error())
	} else {
		fmt.Printf("Content link: %s\n", contentLink)
	}

	pagination, err := parser.ParsePagination(doc)
	if err != nil {
		t.Errorf("Error parsing pagination: %s\n", err.Error())
	} else {
		fmt.Printf("Pagination: %s\n", pagination)
	}

	nodes, err := parser.ParseAllProducts(doc)
	if err != nil {
		t.Fatalf("Error parsing products: %s\n", err.Error())
	}

	for _, node := range nodes {
		asin, err := parser.ParseASIN(node)
		if err != nil {
			t.Errorf("Error parsing asin: %s\n", err.Error())
		} else {
			fmt.Printf("ASIN: %s\n", asin)
		}

		price, err := parser.ParsePrice(node)
		if err != nil {
			t.Errorf("Error parsing price: %s\n", err.Error())
		} else {
			fmt.Printf("Price: %v\n", price)
		}

		star, err := parser.ParseStar(node)
		if err != nil {
			t.Errorf("Error parsing star: %s\n", err.Error())
		} else {
			fmt.Printf("Star: %v\n", star)
		}
	}
}

func TestSellerParser(t *testing.T) {
	p := NewParser()

	doc, err := htmlquery.LoadDoc("./input.html")
	if err != nil {
		t.Fatalf("Error loading document: %s\n", err.Error())
	}

	region, err := ParseRegion(doc)
	if err != nil {
		t.Fatalf("Error parsing region: %s\n", err.Error())
	}

	parser := p.GetSellerParser(region)
	if parser == nil {
		t.Fatal("No parser found for region: " + region)
	}

	pageNum, err := parser.ParseMaxPageNum(doc)
	if err != nil {
		t.Errorf("Error parsing max page num: %s\n", err.Error())
	} else {
		fmt.Printf("Max page num: %s\n", pageNum)
	}

	nextPage, err := parser.ParseNextPageURL(doc)
	if err != nil {
		t.Errorf("Error parsing next page: %s\n", err.Error())
	} else {
		fmt.Printf("Next page url: %s\n", nextPage)
	}

	contentId, err := parser.ParseContentId(doc)
	if err != nil {
		t.Errorf("Error parsing content id: %s\n", err.Error())
	} else {
		fmt.Printf("Content id: %s\n", contentId)
	}

	contentLink, err := parser.ParseContentLink(doc)
	if err != nil {
		t.Errorf("Error parsing content link: %s\n", err.Error())
	} else {
		fmt.Printf("Content link: %s\n", contentLink)
	}

	pagination, err := parser.ParsePagination(doc)
	if err != nil {
		t.Errorf("Error parsing pagination: %s\n", err.Error())
	} else {
		fmt.Printf("pagination: %s\n", pagination)
	}

	nodes, err := parser.ParseAllProducts(doc)
	if err != nil {
		t.Fatalf("Error parsing products: %s\n", err.Error())
	}

	for _, node := range nodes {
		asin, err := parser.ParseASIN(node)
		if err != nil {
			t.Errorf("Error parsing asin: %s\n", err.Error())
		} else {
			fmt.Printf("ASIN: %s\n", asin)
		}

		price, err := parser.ParsePrice(node)
		if err != nil {
			t.Errorf("Error parsing price: %s\n", err.Error())
		} else {
			fmt.Printf("Price: %v\n", price)
		}

		star, err := parser.ParseStar(node)
		if err != nil {
			t.Errorf("Error parsing star: %s\n", err.Error())
		} else {
			fmt.Printf("Star: %v\n", star)
		}
	}
}

func TestProductParser(t *testing.T) {
	p := NewParser()

	doc, err := htmlquery.LoadDoc("./input.html")
	if err != nil {
		t.Fatalf("Error loading document: %s\n", err.Error())
	}

	region, err := ParseRegion(doc)
	if err != nil {
		t.Fatalf("Error parsing region: %s\n", err.Error())
	}

	parser := p.GetProductParser(region)
	if parser == nil {
		t.Fatal("No parser found for region: " + region)
	}

	asin, err := parser.ParseASIN(doc)
	if err != nil {
		t.Errorf("Error parsing product asin: %s\n", err.Error())
	} else {
		fmt.Printf("asin: %s\n", asin)
	}

	star, err := parser.ParseStar(doc)
	if err != nil {
		t.Errorf("Error parsing product star: %s\n", err.Error())
	} else {
		fmt.Printf("star: %s\n", star)
	}

	rating, err := parser.ParseRating(doc)
	if err != nil {
		t.Errorf("Error parsing product rating: %s\n", err.Error())
	} else {
		fmt.Printf("rating: %s\n", rating)
	}

	title, err := parser.ParseTitle(doc)
	if err != nil {
		t.Errorf("Error parsing product title: %s\n", err.Error())
	} else {
		fmt.Printf("Title: %s\n", title)
	}

	price, err := parser.ParsePrice(doc)
	if err != nil {
		t.Errorf("Error parsing product price: %s\n", err.Error())
	} else {
		fmt.Printf("Price: %v\n", price)
	}

	img, err := parser.ParseImg(doc)
	if err != nil {
		t.Errorf("Error parsing product img: %s\n", err.Error())
	} else {
		fmt.Printf("img: %v\n", img)
	}

	dispatchFrom, err := parser.ParseDispatchFrom(doc)
	if err != nil {
		t.Errorf("Error parsing product dispatchFrom: %s\n", err.Error())
	} else {
		fmt.Printf("dispatchFrom: %v\n", dispatchFrom)
	}

	soldBy, err := parser.ParseSoldBy(doc)
	if err != nil {
		t.Errorf("Error parsing product soldBy: %s\n", err.Error())
	} else {
		fmt.Printf("soldBy: %v\n", soldBy)
	}

	pkgDims, err := parser.ParsePackageDimensions(doc)
	if err != nil {
		t.Errorf("Error parsing product pkgDims: %s\n", err.Error())
	} else {
		fmt.Printf("pkgDims: %v\n", pkgDims)
	}

	pkgWeight, err := parser.ParsePackageWeight(doc)
	if err != nil {
		t.Errorf("Error parsing product pkgWeight: %s\n", err.Error())
	} else {
		fmt.Printf("pkgWeight: %v\n", pkgWeight)
	}

	firstDate, err := parser.ParseFirstAvailDate(doc)
	if err != nil {
		t.Errorf("Error parsing product firstDate: %s\n", err.Error())
	} else {
		fmt.Printf("firstDate: %v\n", firstDate)
	}

	sellerId, err := parser.ParseSellerId(doc)
	if err != nil {
		t.Errorf("Error parsing product sellerId: %s\n", err.Error())
	} else {
		fmt.Printf("sellerId: %v\n", sellerId)
	}

	categoryId, err := parser.ParseCategoryId(doc)
	if err != nil {
		t.Errorf("Error parsing product categoryId: %s\n", err.Error())
	} else {
		fmt.Printf("categoryId: %v\n", categoryId)
	}

	productDims, err := parser.ParseProductDimensions(doc)
	if err != nil {
		t.Errorf("Error parsing product productDims: %s\n", err.Error())
	} else {
		fmt.Printf("productDims: %v\n", productDims)
	}

	itemWeight, err := parser.ParseProductWeight(doc)
	if err != nil {
		t.Errorf("Error parsing product itemWeight: %s\n", err.Error())
	} else {
		fmt.Printf("itemWeight: %v\n", itemWeight)
	}

	hasCart, err := parser.ParseHasCart(doc)
	if err != nil {
		t.Errorf("Error parsing product hasCart: %s\n", err.Error())
	} else {
		fmt.Printf("hasCart: %v\n", hasCart)
	}

	coupon, err := parser.ParseCoupon(doc)
	if err != nil {
		t.Errorf("Error parsing product coupon: %s\n", err.Error())
	} else {
		fmt.Printf("coupon: %v\n", coupon)
	}

	color, err := parser.ParseColor(doc)
	if err != nil {
		t.Errorf("Error parsing product color: %s\n", err.Error())
	} else {
		fmt.Printf("color: %v\n", color)
	}

	size, err := parser.ParseSize(doc)
	if err != nil {
		t.Errorf("Error parsing product size: %s\n", err.Error())
	} else {
		fmt.Printf("size: %v\n", size)
	}

	specs, err := parser.ParseSpecs(doc)
	if err != nil {
		t.Errorf("Error parsing product specs: %s\n", err.Error())
	} else {
		fmt.Printf("specs: %v\n", specs)
	}

	desc, err := parser.ParseDescription(doc)
	if err != nil {
		t.Errorf("Error parsing product desc: %s\n", err.Error())
	} else {
		fmt.Printf("desc: %v\n", desc)
	}

	deliveryTime, err := parser.ParseDeliveryTime(doc)
	if err != nil {
		t.Errorf("Error parsing product deliveryTime: %s\n", err.Error())
	} else {
		fmt.Printf("deliveryTime: %v\n", deliveryTime)
	}

	fastestDelivery, err := parser.ParseFastestDelivery(doc)
	if err != nil {
		t.Errorf("Error parsing product fastestDelivery: %s\n", err.Error())
	} else {
		fmt.Printf("fastestDelivery: %v\n", fastestDelivery)
	}

	primePrice, err := parser.ParsePrimePrice(doc)
	if err != nil {
		t.Errorf("Error parsing product primePrice: %s\n", err.Error())
	} else {
		fmt.Printf("primePrice: %v\n", primePrice)
	}

	brand, err := parser.ParseBrand(doc)
	if err != nil {
		t.Errorf("Error parsing product brand: %s\n", err.Error())
	} else {
		fmt.Printf("brand: %v\n", brand)
	}

	categoryHierarchy, err := parser.ParseCategoryHierarchy(doc)
	if err != nil {
		t.Errorf("Error parsing product categoryHierarchy: %s\n", err.Error())
	} else {
		fmt.Printf("categoryHierarchy: %v\n", categoryHierarchy)
	}

	customerReviews, err := parser.ParseCustomerReviews(doc)
	if err != nil {
		t.Errorf("Error parsing product customerReviews: %s\n", err.Error())
	} else {
		fmt.Printf("customerReviews: %v\n", customerReviews)
	}
}

func TestBoardParser(t *testing.T) {
	p := NewParser()

	doc, err := htmlquery.LoadDoc("./input.html")
	if err != nil {
		t.Fatalf("Error loading document: %s\n", err.Error())
	}

	region, err := ParseRegion(doc)
	if err != nil {
		t.Fatalf("Error parsing region: %s\n", err.Error())
	}

	parser := p.GetBoardParser(region)
	if parser == nil {
		t.Fatal("No parser found for region: " + region)
	}

	nextPage, err := parser.ParseNextPageURL(doc)
	if err != nil {
		t.Errorf("Error parsing next page: %s\n", err.Error())
	} else {
		fmt.Printf("Next page url: %s\n", nextPage)
	}

	recsList, err := parser.ParseRecsList(doc)
	if err != nil {
		t.Errorf("Error parsing recsList: %s\n", err.Error())
	} else {
		fmt.Printf("recsList: %s\n", recsList)
	}

	reftag, err := parser.ParseReftag(doc)
	if err != nil {
		t.Errorf("Error parsing reftag: %s\n", err.Error())
	} else {
		fmt.Printf("reftag: %s\n", reftag)
	}

	offset, err := parser.ParseOffset(doc)
	if err != nil {
		t.Errorf("Error parsing offset: %s\n", err.Error())
	} else {
		fmt.Printf("offset: %s\n", offset)
	}

	acpParam, err := parser.ParseAcpParam(doc)
	if err != nil {
		t.Errorf("Error parsing acpParam: %s\n", err.Error())
	} else {
		fmt.Printf("acpParam: %s\n", acpParam)
	}

	acpPath, err := parser.ParseAcpPath(doc)
	if err != nil {
		t.Errorf("Error parsing acpPath: %s\n", err.Error())
	} else {
		fmt.Printf("acpPath: %s\n", acpPath)
	}

	nodes, err := parser.ParseAllProducts(doc)
	if err != nil {
		t.Fatalf("Error parsing products: %s\n", err.Error())
	}

	for _, node := range nodes {
		asin, err := parser.ParseASIN(node)
		if err != nil {
			t.Errorf("Error parsing asin: %s\n", err.Error())
		} else {
			fmt.Printf("asin: %s\n", asin)
		}

		price, err := parser.ParsePrice(node)
		if err != nil {
			t.Errorf("Error parsing price: %s\n", err.Error())
		} else {
			fmt.Printf("price: %s\n", price)
		}

		star, err := parser.ParseStar(node)
		if err != nil {
			t.Errorf("Error parsing star: %s\n", err.Error())
		} else {
			fmt.Printf("star: %s\n", star)
		}

		rating, err := parser.ParseRating(node)
		if err != nil {
			t.Errorf("Error parsing rating: %s\n", err.Error())
		} else {
			fmt.Printf("rating: %s\n", rating)
		}

		title, err := parser.ParseTitle(node)
		if err != nil {
			t.Errorf("Error parsing title: %s\n", err.Error())
		} else {
			fmt.Printf("title: %s\n", title)
		}

		rank, err := parser.ParseRank(node)
		if err != nil {
			t.Errorf("Error parsing rank: %s\n", err.Error())
		} else {
			fmt.Printf("rank: %s\n", rank)
		}
	}
}
