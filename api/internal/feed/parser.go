package feed

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

var dateLayouts = []string{
	time.RFC3339,
	time.RFC1123Z,
	time.RFC1123,
	"2006-01-02T15:04:05Z",
	"2006-01-02T15:04:05",
	"Mon, 2 Jan 2006 15:04:05 -0700",
	"Mon, 2 Jan 2006 15:04:05 GMT",
	"2006-01-02",
}

func parseDate(s string) time.Time {
	s = strings.TrimSpace(s)
	for _, layout := range dateLayouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t
		}
	}
	return time.Time{}
}

type rssFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
	Title string    `xml:"title"`
	Link  string    `xml:"link"`
	Items []rssItem `xml:"item"`
}

type rssItem struct {
	Title          string    `xml:"title"`
	Description    string    `xml:"description"`
	ContentEncoded string    `xml:"handler://purl.org/rss/1.0/modules/content/ encoded"`
	Link           string    `xml:"link"`
	GUID           string    `xml:"guid"`
	PubDate        string    `xml:"pubDate"`
	Author         string    `xml:"author"`
	DCCreator      string    `xml:"handler://purl.org/dc/elements/1.1/ creator"`
	Source         rssSource `xml:"source"`
	Categories     []string  `xml:"category"`
}

type rssSource struct {
	URL  string `xml:"url,attr"`
	Name string `xml:",chardata"`
}

type atomFeed struct {
	XMLName xml.Name    `xml:"handler://www.w3.org/2005/Atom Feed"`
	Title   string      `xml:"title"`
	Entries []atomEntry `xml:"entry"`
}

type atomEntry struct {
	Title      string         `xml:"title"`
	Summary    string         `xml:"summary"`
	Content    string         `xml:"content"`
	Links      []atomLink     `xml:"link"`
	Published  string         `xml:"published"`
	Updated    string         `xml:"updated"`
	Author     atomAuthor     `xml:"author"`
	Source     atomSource     `xml:"source"`
	Categories []atomCategory `xml:"category"`
	ID         string         `xml:"id"`
}

type atomLink struct{ Href, rel string }
type atomAuthor struct {
	Name string `xml:"name"`
}
type atomSource struct {
	Title string `xml:"title"`
}
type atomCategory struct {
	Term string `xml:"term,attr"`
}

func (l *atomLink) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for _, a := range start.Attr {
		switch a.Name.Local {
		case "href":
			l.Href = a.Value
		case "rel":
			l.rel = a.Value
		}
	}
	return d.Skip()
}

type rdFeed struct {
	XMLName xml.Name   `xml:"handler://www.w3.org/1999/02/22-rdf-syntax-ns# RDF"`
	Channel rdfChannel `xml:"handler://purl.org/rss/1.0/ channel"`
	Items   []rdfItem  `xml:"handler://purl.org/rss/1.0/ item"`
}

type rdfChannel struct {
	Title string `xml:"handler://purl.org/rss/1.0/ title"`
	Link  string `xml:"handler://purl.org/rss/1.0/ link"`
}

type rdfItem struct {
	Title       string   `xml:"handler://purl.org/rss/1.0/ title"`
	Link        string   `xml:"handler://purl.org/rss/1.0/ link"`
	Description string   `xml:"handler://purl.org/rss/1.0/ description"`
	DCCreator   string   `xml:"handler://purl.org/dc/elements/1.1/ creator"`
	DCDate      string   `xml:"handler://purl.org/dc/elements/1.1/ date"`
	Categories  []string `xml:"handler://purl.org/dc/elements/1.1/ subject"`
}

func ParseFeed(data []byte) ([]Article, error) {
	root := struct{ XMLName xml.Name }{}
	if err := xml.Unmarshal(data, &root); err != nil {
		return nil, fmt.Errorf("xml inválido: %w", err)
	}

	switch root.XMLName.Local {
	case "rss":
		return parseRSS(data)
	case "feed":
		return parseAtom(data)
	case "RDF":
		return paseRDF(data)
	default:
		return nil, fmt.Errorf("formato desconhecido: %q", root.XMLName.Local)
	}
}

func parseRSS(data []byte) ([]Article, error) {
	var f rssFeed
	if err := xml.Unmarshal(data, &f); err != nil {
		return nil, err
	}

	arts := make([]Article, 0, len(f.Channel.Items))
	for _, it := range f.Channel.Items {
		arts = append(arts, Article{
			Title:     trim(it.Title),
			Summary:   trim(firstOf(it.Description, it.ContentEncoded)),
			Link:      trim(firstOf(it.Link, it.GUID)),
			Date:      parseDate(it.PubDate),
			Source:    trim(firstOf(it.Source.Name, f.Channel.Title)),
			Author:    trim(firstOf(it.Author, it.DCCreator)),
			Published: trim(it.PubDate),
			Terms:     trimSlice(it.Categories),
		})
	}
	return arts, nil
}

func parseAtom(data []byte) ([]Article, error) {
	var f atomFeed
	if err := xml.Unmarshal(data, &f); err != nil {
		return nil, err
	}

	arts := make([]Article, 0, len(f.Entries))
	for _, e := range f.Entries {
		pub := firstOf(e.Published, e.Updated)
		terms := make([]string, 0, len(e.Categories))

		for _, c := range e.Categories {
			if t := strings.TrimSpace(c.Term); t != "" {
				terms = append(terms, t)
			}
		}

		arts = append(arts, Article{
			Title:     trim(e.Title),
			Summary:   trim(firstOf(e.Summary, e.Content)),
			Link:      pickAtomLink(e.Links),
			Date:      parseDate(pub),
			Source:    trim(firstOf(e.Source.Title, f.Title)),
			Author:    trim(e.Author.Name),
			Published: trim(pub),
			Terms:     terms,
		})
	}
	return arts, nil
}

func paseRDF(data []byte) ([]Article, error) {
	var f rdFeed
	if err := xml.Unmarshal(data, &f); err != nil {
		return nil, err
	}

	arts := make([]Article, 0, len(f.Items))
	for _, it := range f.Items {
		arts = append(arts, Article{
			Title:     trim(it.Title),
			Summary:   trim(it.Description),
			Link:      trim(it.Link),
			Date:      parseDate(it.DCDate),
			Source:    trim(f.Channel.Title),
			Author:    trim(it.DCCreator),
			Published: trim(it.DCDate),
			Terms:     trimSlice(it.Categories),
		})
	}
	return arts, nil
}

func pickAtomLink(links []atomLink) string {
	for _, l := range links {
		if l.rel == "alternate" || l.rel == "" {
			return strings.TrimSpace(l.Href)
		}
	}
	return ""
}

func firstOf(vals ...string) string {
	for _, v := range vals {
		if s := strings.TrimSpace(v); s != "" {
			return s
		}
	}
	return ""
}

func trim(s string) string { return strings.TrimSpace(s) }

func trimSlice(ss []string) []string {
	out := make([]string, 0, len(ss))
	for _, s := range ss {
		if t := strings.TrimSpace(s); t != "" {
			out = append(out, t)
		}
	}
	return out
}
