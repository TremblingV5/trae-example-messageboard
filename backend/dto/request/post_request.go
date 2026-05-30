package request

// CreatePostRequest 创建帖子请求
type CreatePostRequest struct {
	Title    string `json:"title" binding:"required,min=1,max=200"`
	Content  string `json:"content" binding:"required"`
	ImageURL string `json:"image_url" binding:"max=255"`
}

// PostListRequest 帖子列表请求
type PostListRequest struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
}

// SearchPostRequest 搜索帖子请求
type SearchPostRequest struct {
	Keyword  string `form:"keyword" binding:"required,min=1"`
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
}