package scraper

import (
	"regexp"

	"golang.org/x/net/html"
)

type HtmlManager struct{}

type HtmlSearchFunc func(*html.Node, string, string, string) (string, bool)

func NewHtmlManager() *HtmlManager {
	return &HtmlManager{}
}

func (h *HtmlManager) searchInHtmlNode(
	n *html.Node, searchFunc HtmlSearchFunc,
	nTargetData, nTargetKey, nTargetVal string,
) string {
	if n == nil {
		return ""
	}

	if ch, ok := searchFunc(n, nTargetData, nTargetKey, nTargetVal); ok {
		return ch
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		res := h.searchInHtmlNode(c, searchFunc, nTargetData, nTargetKey, nTargetVal)
		if res != "" {
			return res
		}
	}
	return ""
}

func (h *HtmlManager) searchListInHtmlNode(
	n *html.Node, searchFunc HtmlSearchFunc,
	nTargetData, nTargetKey, nTargetVal string,
) []string {
	content := []string{}
	if n == nil {
		return nil
	}

	if ch, ok := searchFunc(n, nTargetData, nTargetKey, nTargetVal); ok {
		return []string{ch}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		res := h.searchListInHtmlNode(c, searchFunc, nTargetData, nTargetKey, nTargetVal)
		if res != nil {
			content = append(content, res[:]...)
		}
	}
	return content
}

func (h *HtmlManager) extractNodeData(
	n *html.Node, nTargetData, nTargetKey, nTargetVal string,
) (string, bool) {
	if n.Type == html.ElementNode && n.Data == nTargetData {
		hasAttr := false
		for _, attr := range n.Attr {
			if attr.Key == nTargetKey && attr.Val == nTargetVal {
				hasAttr = true
			}
		}
		if hasAttr {
			return n.LastChild.Data, true
		}
	}
	return "", false
}

func (h *HtmlManager) extractNodeAttrVal(
	n *html.Node, nTargetData, nTargetKey, nTargetVal string,
) (string, bool) {
	if n.Type == html.ElementNode && n.Data == nTargetData {
		var content string
		found := false
		for _, attr := range n.Attr {
			match, err := regexp.MatchString(nTargetVal, attr.Val)
			if err != nil {
				continue
			}
			if match {
				found = true
			}
			if attr.Key == "content" {
				content = attr.Val
			}
			if attr.Key == "src" {
				content = attr.Val
			}
			if attr.Key == "href" {
				content = attr.Val
			}
		}
		if found {
			return content, true
		}
	}
	return "", false
}

func (h *HtmlManager) FindHtmlContentData(
	n *html.Node, nTargetData, nTargetKey, nTargetVal string,
) string {
	return h.searchInHtmlNode(n, h.extractNodeData, nTargetData, nTargetKey, nTargetVal)
}

func (h *HtmlManager) FindHtmlContent(
	n *html.Node, nTargetData, nTargetKey, nTargetVal string,
) string {
	return h.searchInHtmlNode(n, h.extractNodeAttrVal, nTargetData, nTargetKey, nTargetVal)
}

func (h *HtmlManager) ListHtmlContent(
	n *html.Node, nTargetData, nTargetKey, nTargetVal string,
) []string {
	return h.searchListInHtmlNode(n, h.extractNodeAttrVal, nTargetData, nTargetKey, nTargetVal)
}
