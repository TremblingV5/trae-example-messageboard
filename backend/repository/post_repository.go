package repository

import (
	"messageboard/model"

	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(post *model.Post) error {
	return r.db.Create(post).Error
}

func (r *PostRepository) FindByID(id uint) (*model.Post, error) {
	var post model.Post
	err := r.db.Preload("Author").First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostRepository) FindAll(page, pageSize int) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	r.db.Model(&model.Post{}).Count(&total)

	offset := (page - 1) * pageSize
	err := r.db.Preload("Author").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&posts).Error

	return posts, total, err
}

func (r *PostRepository) Search(keyword string, page, pageSize int) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	searchPattern := "%" + keyword + "%"
	r.db.Model(&model.Post{}).
		Where("title LIKE ? OR content LIKE ?", searchPattern, searchPattern).
		Count(&total)

	offset := (page - 1) * pageSize
	err := r.db.Preload("Author").
		Where("title LIKE ? OR content LIKE ?", searchPattern, searchPattern).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&posts).Error

	return posts, total, err
}

func (r *PostRepository) Update(post *model.Post) error {
	return r.db.Save(post).Error
}

func (r *PostRepository) Delete(id uint) error {
	return r.db.Delete(&model.Post{}, id).Error
}