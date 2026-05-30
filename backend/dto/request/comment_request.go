package request

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	Content  string `json:"content" binding:"required,min=1"`
	ParentID *uint  `json:"parent_id"` // nil for root comment
}

// VoteActionRequest 点赞/取消点赞请求
type VoteActionRequest struct {
	Action string `json:"action" binding:"required,oneof=upvote downvote cancel"`
}
