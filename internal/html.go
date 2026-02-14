package internal

import "golang.org/x/net/html"

type HtmlManager struct{}

func NewHtmlManager() *HtmlManager {
	return &HtmlManager{}
}

func (h *HtmlManager) FindHtmlContent(
	n *html.Node, nTargetData, nTargetKey, nTargetVal string,
) string {
	if n.Type == html.ElementNode && n.Data == nTargetData {
		var content string
		hasAttr := false
		for _, attr := range n.Attr {
			if attr.Key == nTargetKey && attr.Val == nTargetVal {
				hasAttr = true
			}
			if attr.Key == "content" {
				content = attr.Val
			}
		}
		if hasAttr {
			return content
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if ch := h.FindHtmlContent(
			c, nTargetData, nTargetKey, nTargetVal,
		); ch != "" {
			return ch
		}
	}

	return ""
}
