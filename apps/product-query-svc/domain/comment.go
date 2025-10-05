package domain

import (
	"strings"
	"time"
	"unicode/utf8"
)

const (
	MaxCommentLength = 2048
)

// Comment represents a user-authored note attached to a product.
type Comment struct {
	ID        int64
	ProductID int64
	UserID    int64
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewComment validates and constructs a comment bound to a product and author.
func NewComment(productID, userID int64, content string) (*Comment, error) {
	c := &Comment{
		ProductID: productID,
		UserID:    userID,
	}
	if err := c.updateContent(strings.TrimSpace(content)); err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	c.CreatedAt = now
	c.UpdatedAt = now
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return c, nil
}

// Validate ensures the comment satisfies domain invariants.
func (c *Comment) Validate() error {
	if c.ProductID <= 0 {
		return ValidationError("product id must be positive")
	}
	if c.UserID <= 0 {
		return ValidationError("user id must be positive")
	}
	trimmed := strings.TrimSpace(c.Content)
	if trimmed == "" {
		return ValidationError("content required")
	}
	if utf8.RuneCountInString(trimmed) > MaxCommentLength {
		return ValidationError("content too long")
	}
	return nil
}

// UpdateContent modifies the body of the comment and bumps the update timestamp.
func (c *Comment) UpdateContent(content string) error {
	if err := c.updateContent(strings.TrimSpace(content)); err != nil {
		return err
	}
	c.UpdatedAt = time.Now().UTC()
	return nil
}

func (c *Comment) updateContent(content string) error {
	if content == "" {
		return ValidationError("content required")
	}
	if utf8.RuneCountInString(content) > MaxCommentLength {
		return ValidationError("content too long")
	}
	c.Content = content
	return nil
}
