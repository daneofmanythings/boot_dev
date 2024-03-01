package workers

type RSS struct {
	Channel   Channel `xml:"channel"`
	Version   string  `xml:"version"`
	XmlnsAtom string  `xml:"xmlns:atom"`
}

type Channel struct {
	Title         string        `xml:"title"`
	Link          []LinkElement `xml:"link"`
	Description   string        `xml:"description"`
	Generator     string        `xml:"generator"`
	Language      string        `xml:"language"`
	LastBuildDate string        `xml:"lastBuildDate"`
	Item          []Item        `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
	Description string `xml:"description"`
}

type LinkClass struct {
	Href   string `xml:"href"`
	Rel    string `xml:"rel"`
	Type   string `xml:"type"`
	Prefix string `xml:"prefix"`
}

type LinkElement struct {
	LinkClass *LinkClass
	String    *string
}
