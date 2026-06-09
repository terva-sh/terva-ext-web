package fetch

import (
	"fmt"
	"net/url"
	"strings"

	xhtml "golang.org/x/net/html"
)

// imgSentinel is the token left in the DOM in place of an <img>; it survives
// Markdown conversion unescaped and is swapped for the `[image:N]` placeholder
// afterward (so neither the brackets nor the alt text get Markdown-escaped).
func imgSentinel(id int) string { return fmt.Sprintf(" {{IMG:%d}} ", id) }

// placeholder is what the model sees inline for an image: a short, stable
// handle plus the alt text when present.
func (im Image) placeholder() string {
	if im.Alt != "" {
		return fmt.Sprintf("[image:%d: %s]", im.ID, im.Alt)
	}
	return fmt.Sprintf("[image:%d]", im.ID)
}

// FormatImages renders the image list for web_images as a compact, model-
// readable list keyed by the same `[image:N]` handles that appear in the
// web_fetch output.
func FormatImages(pageURL string, imgs []Image) string {
	if len(imgs) == 0 {
		return fmt.Sprintf("No images found on %s.", pageURL)
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d image(s) on %s:\n", len(imgs), pageURL)
	for _, im := range imgs {
		fmt.Fprintf(&b, "\n[image:%d] %s", im.ID, im.URL)
		if im.Alt != "" {
			fmt.Fprintf(&b, "\n   %s", strings.Join(strings.Fields(im.Alt), " "))
		}
	}
	return strings.TrimSpace(b.String())
}

// indexImages walks root, replacing each <img> with a sentinel text node and
// collecting the (absolute) image URLs. Identical URLs share one id. base is
// the page URL, used to resolve relative srcs.
func indexImages(root *xhtml.Node, base *url.URL) []Image {
	var images []Image
	byURL := map[string]int{}

	var walk func(n *xhtml.Node)
	walk = func(n *xhtml.Node) {
		// Capture NextSibling before touching the node, since replacing an
		// <img> detaches it from the sibling chain.
		for ch := n.FirstChild; ch != nil; {
			next := ch.NextSibling
			if ch.Type == xhtml.ElementNode && ch.Data == "img" {
				if abs := imgURL(ch, base); abs != "" {
					id, ok := byURL[abs]
					if !ok {
						id = len(images) + 1
						byURL[abs] = id
						images = append(images, Image{ID: id, URL: abs, Alt: strings.TrimSpace(attrVal(ch, "alt"))})
					}
					replaceWithSentinel(ch, id)
				}
				// <img> is void; nothing to recurse into.
			} else {
				walk(ch)
			}
			ch = next
		}
	}
	walk(root)
	return images
}

// applyPlaceholders swaps each image's sentinel for its `[image:N]` placeholder.
func applyPlaceholders(md string, images []Image) string {
	for _, im := range images {
		md = strings.ReplaceAll(md, strings.TrimSpace(imgSentinel(im.ID)), im.placeholder())
	}
	return md
}

// imgURL resolves an <img>'s best source to an absolute URL, skipping inline
// data: URIs. Falls back to the first srcset candidate when src is absent.
func imgURL(n *xhtml.Node, base *url.URL) string {
	src := strings.TrimSpace(attrVal(n, "src"))
	if src == "" {
		src = firstSrcset(attrVal(n, "srcset"))
	}
	if src == "" || strings.HasPrefix(strings.ToLower(src), "data:") {
		return ""
	}
	ref, err := url.Parse(src)
	if err != nil {
		return ""
	}
	return base.ResolveReference(ref).String()
}

// firstSrcset returns the first URL from a srcset attribute (ignoring its
// width/density descriptor).
func firstSrcset(s string) string {
	if s = strings.TrimSpace(s); s == "" {
		return ""
	}
	first := strings.SplitN(s, ",", 2)[0]
	if fields := strings.Fields(first); len(fields) > 0 {
		return fields[0]
	}
	return ""
}

func attrVal(n *xhtml.Node, key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

// replaceWithSentinel swaps node for a text node carrying its image sentinel.
func replaceWithSentinel(node *xhtml.Node, id int) {
	parent := node.Parent
	if parent == nil {
		return
	}
	parent.InsertBefore(&xhtml.Node{Type: xhtml.TextNode, Data: imgSentinel(id)}, node)
	parent.RemoveChild(node)
}
