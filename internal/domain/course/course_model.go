package course

import (
	"encoding/json"
	"time"

	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/nuuid"
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
)

type Course struct {
	ID        uuid.UUID   `db:"id" validate:"required"`
	UserID    uuid.UUID   `db:"user_id" validate:"required"`
	Title     string      `db:"title" validate:"required"`
	Content   string      `db:"content" validate:"required"`
	CreatedAt time.Time   `db:"created_at" validate:"required"`
	CreatedBy uuid.UUID   `db:"created_by" validate:"required"`
	UpdatedAt null.Time   `db:"updated_at"`
	UpdatedBy nuuid.NUUID `db:"updated_by"`
	DeletedAt null.Time   `db:"deleted_at"`
	DeletedBy nuuid.NUUID `db:"deleted_by"`
}

type CourseQueryParameters struct {
	Page  int
	Limit int
	Sort  string
	Order string
	Role  string
}

func (c Course) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.ToResponseFormat())
}

func (c Course) NewCourseFromRequestFormat(req CourseRequestFormat, userID uuid.UUID) (newCourse Course, err error) {
	courseID, _ := uuid.NewV4()

	newCourse = Course{
		ID:        courseID,
		UserID:    userID,
		Title:     req.Title,
		Content:   req.Content,
		CreatedAt: time.Now(),
		CreatedBy: userID,
	}

	err = newCourse.Validate()

	return
}

func (c *Course) Validate() (err error) {
	validator := shared.GetValidator()
	return validator.Struct(c)
}

func (c Course) ToResponseFormat() CourseResponseFormat {
	return CourseResponseFormat{
		ID:        c.ID,
		UserID:    c.UserID,
		Title:     c.Title,
		Content:   c.Content,
		CreatedBy: c.CreatedBy,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		UpdatedBy: c.UpdatedBy.Ptr(),
		DeletedAt: c.DeletedAt,
		DeletedBy: c.DeletedBy.Ptr(),
	}
}

type CourseRequestFormat struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}

type CourseResponseFormat struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"userID"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"createdAt"`
	CreatedBy uuid.UUID  `json:"createdBy"`
	UpdatedAt null.Time  `json:"updatedAt"`
	UpdatedBy *uuid.UUID `json:"updatedBy"`
	DeletedAt null.Time  `json:"deletedAt,omitempty"`
	DeletedBy *uuid.UUID `json:"deletedBy,omitempty"`
}
