package rpclient

import "golang.org/x/net/html"

func traverseN(n *html.Node, attrs []string, result map[string]*html.Node) map[string]*html.Node {
	if n.Type == html.ElementNode {
		for _, attr := range n.Attr {
			for _, attr_value := range attrs {
				if attr.Val == attr_value {
					result[attr.Val] = n
				}
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		r := traverseN(c, attrs, result)
		for key := range r {
			result[key] = r[key]
		}

	}
	return result
}

// attrs []attrs_value  -> result  map[attrs_value]*html.Node
func getElementsByAttrs(n *html.Node, attrs []string) map[string]*html.Node {
	var res = make(map[string]*html.Node)
	return traverseN(n, attrs, res)
}

func GetAttribute(n *html.Node, key string) (string, bool) {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val, true
		}
	}
	return "", false
}
