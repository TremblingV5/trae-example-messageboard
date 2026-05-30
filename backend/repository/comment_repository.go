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

// GetCommentDepth 获取评论的深度
func (r *CommentRepository) GetCommentDepth(commentID uint) (int, error) {
	var comment model.Comment
	if err := r.db.First(&comment, commentID).Error; err != nil {
		return 0, err
	}

	depth := 0
	currentComment := &comment

	for currentComment.ParentID != nil {
		depth++
		if depth >= model.MaxCommentDepth {
			return depth, nil
		}
		var parent model.Comment
		if err := r.db.First(&parent, *currentComment.ParentID).Error; err != nil {
			return 0, err
		}
		currentComment = &parent
	}

	return depth, nil
}

func (r *CommentRepository) Delete(id uint) error {
	return r.db.Delete(&model.Comment{}, id).Error
}