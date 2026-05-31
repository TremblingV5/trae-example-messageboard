package repository

import (
	"messageboard/model"

	"gorm.io/gorm"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(comment *model.Comment) error {
	return r.db.Create(comment).Error
}

func (r *CommentRepository) FindByID(id uint) (*model.Comment, error) {
	var comment model.Comment
	err := r.db.Preload("Author").First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// FindByPostID 获取帖子的所有评论
func (r *CommentRepository) FindByPostID(postID uint) ([]model.Comment, error) {
	var comments []model.Comment
	err := r.db.Where("post_id = ?", postID).
		Preload("Author").
		Order("created_at ASC").
		Find(&comments).Error
	return comments, err
}

// FindRootCommentsByPostID 获取帖子的根评论（parent_id 为 null）
func (r *CommentRepository) FindRootCommentsByPostID(postID uint) ([]model.Comment, error) {
	var comments []model.Comment
	err := r.db.Where("post_id = ? AND parent_id IS NULL", postID).
		Preload("Author").
		Order("created_at ASC").
		Find(&comments).Error
	return comments, err
}

// FindChildrenByParentID 获取子评论
func (r *CommentRepository) FindChildrenByParentID(parentID uint) ([]model.Comment, error) {
	var comments []model.Comment
	err := r.db.Where("parent_id = ?", parentID).
		Preload("Author").
		Order("created_at ASC").
		Find(&comments).Error
	return comments, err
}

// GetCommentDepth 获取评论的深度（使用递归 CTE 优化查询）
func (r *CommentRepository) GetCommentDepth(commentID uint) (int, error) {
	// 使用递归 CTE 一次性查询评论深度
	// 如果数据库不支持 CTE，则使用逐级查询（最多 5 层）
	
	// 首先获取评论信息
	var comment model.Comment
	if err := r.db.First(&comment, commentID).Error; err != nil {
		return 0, err
	}

	// 如果没有父评论，深度为 0
	if comment.ParentID == nil {
		return 0, nil
	}

	// 使用逐级查询，但限制最大深度为 MaxCommentDepth
	depth := 0
	currentID := *comment.ParentID
	visited := make(map[uint]bool) // 防止循环引用

	for depth < model.MaxCommentDepth {
		// 检查是否已访问过（防止循环）
		if visited[currentID] {
			return depth, nil
		}
		visited[currentID] = true

		var parent model.Comment
		if err := r.db.First(&parent, currentID).Error; err != nil {
			// 如果父评论不存在，返回当前深度
			return depth, nil
		}

		depth++

		// 如果没有更多父评论，返回深度
		if parent.ParentID == nil {
			return depth, nil
		}

		currentID = *parent.ParentID
	}

	return depth, nil
}

// BatchGetCommentDepths 批量获取多个评论的深度
// 通过一次性加载所有相关评论，在内存中计算深度
func (r *CommentRepository) BatchGetCommentDepths(commentIDs []uint) (map[uint]int, error) {
	if len(commentIDs) == 0 {
		return make(map[uint]int), nil
	}

	// 获取所有评论
	var comments []model.Comment
	if err := r.db.Find(&comments).Error; err != nil {
		return nil, err
	}

	// 构建评论映射
	commentMap := make(map[uint]*model.Comment, len(comments))
	for i := range comments {
		commentMap[comments[i].ID] = &comments[i]
	}

	// 计算每个请求的评论深度
	result := make(map[uint]int, len(commentIDs))
	for _, id := range commentIDs {
		result[id] = calculateDepth(commentMap, id, make(map[uint]int), make(map[uint]bool))
	}

	return result, nil
}

// calculateDepth 递归计算评论深度（带缓存和循环检测）
func calculateDepth(commentMap map[uint]*model.Comment, commentID uint, cache map[uint]int, visiting map[uint]bool) int {
	// 检查缓存
	if depth, ok := cache[commentID]; ok {
		return depth
	}

	// 检查循环引用
	if visiting[commentID] {
		return 0
	}

	comment, ok := commentMap[commentID]
	if !ok || comment.ParentID == nil {
		return 0
	}

	visiting[commentID] = true
	depth := 1 + calculateDepth(commentMap, *comment.ParentID, cache, visiting)
	delete(visiting, commentID)

	// 限制最大深度
	if depth > model.MaxCommentDepth {
		depth = model.MaxCommentDepth
	}

	cache[commentID] = depth
	return depth
}

func (r *CommentRepository) Delete(id uint) error {
	return r.db.Delete(&model.Comment{}, id).Error
}
