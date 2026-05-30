package response

import "time"

// Response 统一响应格式
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SuccessResponse 成功响应
func SuccessResponse(data interface{}) Response {
	return Response{
		Code:    200,
		Message: "success",
		Data:    data,
	}
}

// SuccessWithMessage 带消息的成功响应
func SuccessWithMessage(message string, data interface{}) Response {
	return Response{
		Code:    200,
		Message: message,
		Data:    data,
	}
}

// ErrorResponse 错误响应
func ErrorResponse(code int, message string) Response {
	return Response{
		Code:    code,
		Message: message,
	}
}

// 常用错误响应
var (
	BadRequestResponse       = ErrorResponse(400, "bad request")
	UnauthorizedResponse     = ErrorResponse(401, "unauthorized")
	ForbiddenResponse        = ErrorResponse(403, "forbidden")
	NotFoundResponse         = ErrorResponse(404, "not found")
	InternalServerErrorResponse = ErrorResponse(500, "internal server error")
)

// LoginResponse 登录响应
type LoginResponse struct {
	Token    string    `json:"token"`
	ExpireAt time.Time `json:"expire_at"`
}

// UserResponse 用户信息响应
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Nickname  string    `json:"nickname"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PostResponse 帖子响应
type PostResponse struct {
	ID        uint          `json:"id"`
	Title     string        `json:"title"`
	Content   string        `json:"content"`
	Author    *UserResponse `json:"author,omitempty"`
	ImageURL  string        `json:"image_url"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// PostListResponse 帖子列表响应
type PostListResponse struct {
	Posts     []PostResponse `json:"posts"`
	Total     int64          `json:"total"`
	Page      int            `json:"page"`
	PageSize  int            `json:"page_size"`
}

// CommentResponse 评论响应
type CommentResponse struct {
	ID        uint              `json:"id"`
	Content   string            `json:"content"`
	PostID    uint              `json:"post_id"`
	Author    *UserResponse     `json:"author,omitempty"`
	ParentID  *uint             `json:"parent_id"`
	Children  []CommentResponse `json:"children,omitempty"`
	VoteCount int               `json:"vote_count"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// VoteResponse 点赞响应
type VoteResponse struct {
	VoteCount int  `json:"vote_count"`
	Voted     bool `json:"voted"`
	Value     int  `json:"value"` // 1: upvote, -1: downvote, 0: not voted
}