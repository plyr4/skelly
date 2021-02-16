package types

// Emoji struct representation for emoji data
// fetched from https://raw.githubusercontent.com/iamcal/emoji-data/master/emoji.json
type Emoji struct {
	Name           string      `json:"name"`
	Unified        string      `json:"unified"`
	NonQualified   string      `json:"non_qualified"`
	Docomo         string      `json:"docomo"`
	Au             string      `json:"au"`
	Softbank       string      `json:"softbank"`
	Google         string      `json:"google"`
	Image          string      `json:"image"`
	SheetX         int         `json:"sheet_x"`
	SheetY         int         `json:"sheet_y"`
	ShortName      string      `json:"short_name"`
	ShortNames     []string    `json:"short_names"`
	Text           interface{} `json:"text"`
	Texts          interface{} `json:"texts"`
	Category       string      `json:"category"`
	SortOrder      int         `json:"sort_order"`
	AddedIn        string      `json:"added_in"`
	HasImgApple    bool        `json:"has_img_apple"`
	HasImgGoogle   bool        `json:"has_img_google"`
	HasImgTwitter  bool        `json:"has_img_twitter"`
	HasImgFacebook bool        `json:"has_img_facebook"`
}
