package madels

type Country struct {
	Name string
	Slug string
	Capital string
	Population int64 
	Region string
	Subregion string
	Flag string // URL to flag image (SVG)
	Languages []string
	Currencies []string
}