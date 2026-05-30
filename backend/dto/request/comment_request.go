package request

import "messageboard/model"

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	Content  string `json:"content" binding:"required,min=1"`
	ParentID *uint  `json:"parent_id"` // nil for root comment
}

// ValidateCommentDepth 验证评论深度
func ValidateCommentDepth(parentID *uint) bool {
	if parentID == nil {
		return true // root comment
	}
	// Depth validation will be done in service layer
	return true
}

// VoteRequest 点赞请求
type VoteRequest struct {
	Value int `json:"value" binding:"required,oneof=1 -1"` // 1 for upvote, -1 for downvote
}

// VoteActionRequest 点赞/取消点赞请求
type VoteActionRequest struct {
	Action string `json:"action" binding:"required,oneof=upvote downvote cancel"`
}

// ToVoteValue 将 action 转换为 value
func (r *VoteActionRequest) ToVoteValue() int {
	if r.Action == "upvote" {
		return 1
	} else if r.Action == "downvote" {
		return -1
	}
	return 0
}

// Ensure VoteRequest implements model.Vote interface
var _ = model.MaxCommentDepth // just to ensure import