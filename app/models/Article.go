package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"strings"
)

type Article struct {
	gorm.Model
	Title string `gorm:"size:100;not null"`
	Body string `gorm:"size:100;not null"`
	CreatedBy User `gorm:"foreignKey:UserID;" `
	UserID uint `gorm:"not null"`
}

func (article *Article) Prepare() {
	article.Title = strings.TrimSpace(article.Title)
	article.Body = strings.TrimSpace(article.Body)
	article.CreatedBy = User{}
}

func  (article *Article) Validate() error {
	if article.Title == "" {
		return errors.New("Title is required")
	}
	if article.Body == "" {
		return errors.New("Body is required")
	}
	return nil
}

func (article *Article) Save(db *gorm.DB) (*Article, error) {
	var err error

	// Debug article and show detailed log operation
	err = db.Debug().Create(&article).Error
	if err != nil {
		return &Article{}, err
	}

	return article, nil
}

func GetArticleByID(id int, db *gorm.DB) (*Article, error) {
	article := &Article{}

	if err := db.Debug().Table("articles").Where("id = ?", id).First(article).Error; err != nil {
		return nil, err
	}

	return article, nil
}

func (article *Article) GetArticles(db *gorm.DB) (*[]Article, error) {
	articles := []Article{}

	if err := db.Debug().Table("articles").Find(&articles).Error; err != nil {
		return &[]Article{}, err
	}

	return &articles, nil
}

func (article *Article) UpdateArticle(id int, db *gorm.DB) (*Article, error) {
	if err := db.Debug().Table("articles").Where("id = ? ", id).Updates(Article{
		Title: article.Title,
		Body: article.Body,
	}).Error; err != nil {
		return &Article{}, err
	}

	return article, nil
}

func DeleteArticle(id int, db *gorm.DB) {
	db.Debug().Table("articles").Where("id = ?", id).Unscoped().Delete(&Article{})
}