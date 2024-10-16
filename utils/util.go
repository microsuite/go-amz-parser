package utils

import (
	"fmt"
	"net/url"
	"strings"
	"unicode"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

func FindNodes(doc *html.Node, expr string, multi bool) ([]*html.Node, error) {
	nodes, err := htmlquery.QueryAll(doc, expr)
	if err != nil {
		return nil, fmt.Errorf("'%v' error, %v", expr, err)
	}

	if len(nodes) == 0 {
		return nil, fmt.Errorf("'%v' error, no nodes selected", expr)
	}

	if len(nodes) != 1 && !multi {
		return nil, fmt.Errorf("'%v' error, %v nodes selected", expr, len(nodes))
	}
	return nodes, err
}

func FormatNumber(s string) string {
	s = strings.ReplaceAll(s, "k", "000")
	s = strings.ReplaceAll(s, "K", "000")
	s = strings.ReplaceAll(s, ",", "")
	return strings.ReplaceAll(s, "+", "")
}

func FormatNumberEuro(old string) string {
	old = strings.ReplaceAll(old, ".", "")
	old = strings.ReplaceAll(old, ",", ".")
	old = strings.ReplaceAll(old, "k", "000")
	return strings.ReplaceAll(old, "K", "000")
}

func FormatTitle(s string) string {
	return strings.Replace(strings.Replace(strings.TrimSpace(s), "'", " ", -1), "\"", " ", -1)
}

func DropMoneySym(s string) string {
	switch {
	case strings.Contains(s, "$"):
		/* com ca */
		s = strings.ReplaceAll(s, "$", "")
	case strings.Contains(s, "€"):
		/* de es fr it */
		s = strings.ReplaceAll(s, "€", "")
	case strings.Contains(s, "£"):
		/* uk */
		s = strings.ReplaceAll(s, "£", "")
	case strings.Contains(s, "￥"):
		/* jp */
		s = strings.ReplaceAll(s, "￥", "")
	case strings.Contains(s, "EUR"):
		s = strings.ReplaceAll(s, "EUR", "")
	default:
	}
	return strings.TrimSpace(s)
}

func FindNumberHead(s string) string {
	strs := strings.Split(s, " ")
	if len(strs) > 0 {
		s = strs[0]
	}

	lastIndex := strings.LastIndexFunc(s, unicode.IsDigit)
	if lastIndex == -1 {
		return s
	}
	return s[0 : lastIndex+1]
}

func FormatRating(s string) string {
	lastIndex := strings.LastIndexFunc(s, unicode.IsDigit)
	if lastIndex == -1 {
		return s
	}
	return strings.ReplaceAll(strings.ReplaceAll(s[0:lastIndex+1], ",", ""), ".", "")
}

func FormalMerchant(s string) string {
	u, err := url.Parse(html.UnescapeString(s))
	if err != nil {
		return "unkown"
	}

	values := u.Query()

	if seller, ok := values["seller"]; ok {
		return seller[0]
	} else {
		if _, ok := values["nodeId"]; ok {
			return "amazon"
		} else {
			return "unkown"
		}
	}
}
