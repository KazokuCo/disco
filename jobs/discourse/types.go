package discourse

type Post struct {
	ID              int64  `json:"id"`
	Num             int64  `json:"post_number"`
	UserID          int64  `json:"user_id"`
	Username        string `json:"username"`
	UserTitle       string `json:"user_title"`
	DisplayUsername string `json:"display_username"`
	Admin           bool   `json:"admin"`
	Moderator       bool   `json:"moderator"`
	Staff           bool   `json:"staff"`
	TrustLevel      int    `json:"trust_level"`
	Cooked          string `json:"cooked"`
}

type Topic struct {
	ID         int64  `json:"id"`
	Title      string `json:"title"`
	FancyTitle string `json:"fancy_title"`
	Slug       string `json:"slug"`
}

type SearchEnvelope struct {
	Topics []Topic `json:"topics"`
	Posts  []Post  `json:"posts"`
}
