package api

type SiteStat struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type SocialLink struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

type VisitorPayload struct {
	DisplayName string `json:"displayName"`
}

type SitePayload struct {
	Name        string       `json:"name"`
	Headline    string       `json:"headline"`
	Intro       string       `json:"intro"`
	Location    string       `json:"location"`
	Domain      string       `json:"domain"`
	Email       string       `json:"email"`
	Motto       string       `json:"motto"`
	TechStack   []string     `json:"techStack"`
	Stats       []SiteStat   `json:"stats"`
	SocialLinks []SocialLink `json:"socialLinks"`
}

type PostSummaryPayload struct {
	Slug        string   `json:"slug"`
	Title       string   `json:"title"`
	Summary     string   `json:"summary"`
	Category    string   `json:"category"`
	ReadTime    string   `json:"readTime"`
	HeroNote    string   `json:"heroNote"`
	CoverLabel  string   `json:"coverLabel"`
	Tags        []string `json:"tags"`
	Featured    bool     `json:"featured"`
	PublishedAt string   `json:"publishedAt"`
	LikeCount   int      `json:"likeCount"`
	CommentCount int     `json:"commentCount"`
}

type ContentBlock struct {
	Kind  string   `json:"kind"`
	Title string   `json:"title,omitempty"`
	Text  string   `json:"text,omitempty"`
	URL   string   `json:"url,omitempty"`
	Items []string `json:"items,omitempty"`
}

type PostDetailPayload struct {
	PostSummaryPayload
	Blocks         []ContentBlock   `json:"blocks"`
	LikedByVisitor bool             `json:"likedByVisitor"`
	Comments       []CommentPayload `json:"comments"`
}

type ProjectPayload struct {
	Name      string   `json:"name"`
	Summary   string   `json:"summary"`
	Status    string   `json:"status"`
	Link      string   `json:"link"`
	ImageURL  string   `json:"imageUrl"`
	Accent    string   `json:"accent"`
	TechStack []string `json:"techStack"`
}

type VideoPayload struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	ID          int64  `json:"id"`
	ThumbnailURL string `json:"thumbnailUrl"`
	PublishedAt string `json:"publishedAt"`
}

type TimelineEntryPayload struct {
	Period      string `json:"period"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CommentPayload struct {
	ID         int64  `json:"id"`
	AuthorName string `json:"authorName"`
	Content    string `json:"content"`
	CreatedAt  string `json:"createdAt"`
}

type HomeResponse struct {
	Site          SitePayload            `json:"site"`
	FeaturedPosts []PostSummaryPayload   `json:"featuredPosts"`
	RecentPosts   []PostSummaryPayload   `json:"recentPosts"`
	Projects      []ProjectPayload       `json:"projects"`
	Videos        []VideoPayload         `json:"videos"`
	Timeline      []TimelineEntryPayload `json:"timeline"`
}

type PostsResponse struct {
	Posts []PostSummaryPayload `json:"posts"`
}

type PostResponse struct {
	Post    PostDetailPayload `json:"post"`
	Visitor VisitorPayload    `json:"visitor"`
}

type CreateCommentRequest struct {
	Content string `json:"content"`
}

type CommentResponse struct {
	Comment      CommentPayload `json:"comment"`
	CommentCount int            `json:"commentCount"`
	Visitor      VisitorPayload `json:"visitor"`
}

type LikeResponse struct {
	LikeCount int            `json:"likeCount"`
	Liked     bool           `json:"liked"`
	Visitor   VisitorPayload `json:"visitor"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
