package errors

import "fmt"

var (
	ErrorNotFoundPrice             = fmt.Errorf("not found price")
	ErrorNotFoundStar              = fmt.Errorf("not found star")
	ErrorNotFoundRating            = fmt.Errorf("not found rating")
	ErrorNotFoundLanguage          = fmt.Errorf("not found lang")
	ErrorNotFoundNextPage          = fmt.Errorf("not found next page")
	ErrorNotFoundContentId         = fmt.Errorf("not found content id")
	ErrorNotFoundDimensions        = fmt.Errorf("not found dimensions")
	ErrorNotFoundWeight            = fmt.Errorf("not found weight")
	ErrorNotFoundFirstDate         = fmt.Errorf("not found first available date")
	ErrorNotFoundPackageDimensions = fmt.Errorf("not found package dimensions")
	ErrorNotFoundPackageWeight     = fmt.Errorf("not found package weight")
	ErrorNotFoundDispatchFrom      = fmt.Errorf("not found package dispatch from")
	ErrorNotFoundSoldBy            = fmt.Errorf("not found package sold by")
	ErrorNotFoundSellerId          = fmt.Errorf("not found package seller id")
)
