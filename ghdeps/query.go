package ghdeps

import (
	"strings"

	"golang.org/x/net/html"
)

func getAttribute(node *html.Node, name string) string {
	if node == nil || node.Attr == nil {
		return ""
	}
	for _, attr := range node.Attr {
		if attr.Key == name {
			return attr.Val
		}
	}
	return ""
}

func queryBox(node *html.Node) *html.Node {
	if node.Type == html.ElementNode {
		if node.Data == "div" {
			if getAttribute(node, "id") == "dependents" {
				for child := node.FirstChild; child != nil; child = child.NextSibling {
					if getAttribute(child, "class") == "Box" {
						return child
					}
				}
			}
		}
	}
	return nil
}

func queryNextRow(row *html.Node) *html.Node {
	for ; row != nil; row = row.NextSibling {
		if row.Type != html.ElementNode || row.Data != "div" {
			continue
		}
		if strings.Contains(getAttribute(row, "class"), "Box-row") {
			return row
		}
	}
	return nil
}

func queryNextPageButton(box *html.Node) *html.Node {
	for next := box.NextSibling; next != nil; next = next.NextSibling {
		if next.Type != html.ElementNode {
			continue
		}
		if next.Data != "div" {
			continue
		}
		if !strings.Contains(getAttribute(next, "class"), "paginate-container") {
			continue
		}
		for group := next.FirstChild; group != nil; group = group.NextSibling {
			if group.Type == html.ElementNode &&
				group.Data == "div" &&
				strings.Contains(getAttribute(group, "class"), "BtnGroup") {
				for btn := group.FirstChild; btn != nil; btn = btn.NextSibling {
					if btn.Type == html.ElementNode &&
						btn.Data == "a" &&
						strings.Contains(getAttribute(btn, "href"), "dependents_after") {
						return btn
					}
				}
			}
		}
	}
	return nil
}
