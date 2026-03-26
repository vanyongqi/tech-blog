package api

type AdminLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AdminSessionPayload struct {
	Authenticated bool   `json:"authenticated"`
	Username      string `json:"username"`
}

type AdminLoginResponse struct {
	Session AdminSessionPayload `json:"session"`
}

type AdminContentBlock struct {
	Kind  string   `json:"kind"`
	Title string   `json:"title,omitempty"`
	Text  string   `json:"text,omitempty"`
	URL   string   `json:"url,omitempty"`
	Items []string `json:"items,omitempty"`
}

type AdminPostSummaryPayload struct {
	Slug         string   `json:"slug"`
	Title        string   `json:"title"`
	Summary      string   `json:"summary"`
	Category     string   `json:"category"`
	ReadTime     string   `json:"readTime"`
	CoverLabel   string   `json:"coverLabel"`
	Tags         []string `json:"tags"`
	Featured     bool     `json:"featured"`
	PublishedAt  string   `json:"publishedAt"`
	LikeCount    int      `json:"likeCount"`
	CommentCount int      `json:"commentCount"`
}

type AdminPostPayload struct {
	Slug         string              `json:"slug"`
	Title        string              `json:"title"`
	Summary      string              `json:"summary"`
	Category     string              `json:"category"`
	ReadTime     string              `json:"readTime"`
	HeroNote     string              `json:"heroNote"`
	CoverLabel   string              `json:"coverLabel"`
	Tags         []string            `json:"tags"`
	Featured     bool                `json:"featured"`
	PublishedAt  string              `json:"publishedAt"`
	Blocks       []AdminContentBlock `json:"blocks"`
	LikeCount    int                 `json:"likeCount"`
	CommentCount int                 `json:"commentCount"`
}

type AdminSavePostRequest struct {
	Slug        string              `json:"slug"`
	Title       string              `json:"title"`
	Summary     string              `json:"summary"`
	Category    string              `json:"category"`
	ReadTime    string              `json:"readTime"`
	HeroNote    string              `json:"heroNote"`
	CoverLabel  string              `json:"coverLabel"`
	Tags        []string            `json:"tags"`
	Featured    bool                `json:"featured"`
	PublishedAt string              `json:"publishedAt"`
	Blocks      []AdminContentBlock `json:"blocks"`
}

type AdminPostsResponse struct {
	Posts []AdminPostSummaryPayload `json:"posts"`
}

type AdminPostResponse struct {
	Post AdminPostPayload `json:"post"`
}

type AdminProjectPayload struct {
	ID        int64    `json:"id"`
	Name      string   `json:"name"`
	Summary   string   `json:"summary"`
	Status    string   `json:"status"`
	Link      string   `json:"link"`
	ImageURL  string   `json:"imageUrl"`
	Accent    string   `json:"accent"`
	TechStack []string `json:"techStack"`
}

type AdminSaveProjectRequest struct {
	Name      string   `json:"name"`
	Summary   string   `json:"summary"`
	Status    string   `json:"status"`
	Link      string   `json:"link"`
	ImageURL  string   `json:"imageUrl"`
	Accent    string   `json:"accent"`
	TechStack []string `json:"techStack"`
}

type AdminProjectsResponse struct {
	Projects []AdminProjectPayload `json:"projects"`
}

type AdminProjectResponse struct {
	Project AdminProjectPayload `json:"project"`
}

type AdminVideoPayload struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	ThumbnailURL string `json:"thumbnailUrl"`
	PublishedAt string `json:"publishedAt"`
}

type AdminSaveVideoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	ThumbnailURL string `json:"thumbnailUrl"`
	PublishedAt string `json:"publishedAt"`
}

type AdminVideosResponse struct {
	Videos []AdminVideoPayload `json:"videos"`
}

type AdminVideoResponse struct {
	Video AdminVideoPayload `json:"video"`
}

type AdminSuggestThumbnailRequest struct {
	URL string `json:"url"`
}

type AdminSuggestThumbnailResponse struct {
	ThumbnailURL string `json:"thumbnailUrl"`
}
