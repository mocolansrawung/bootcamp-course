package course

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/logger"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	courseQueries = struct {
		selectCourses string
		insertCourse  string
	}{
		selectCourses: `
			SELECT
				id, 
				user_id,
				title,
				content,
				created_at,
				created_by,
				updated_at,
				updated_by,
				deleted_at,
				deleted_by
			FROM courses
		`,

		insertCourse: `
			INSERT INTO courses (
				id, 
				user_id,
				title,
				content,
				created_at,
				created_by,
				updated_at,
				updated_by,
				deleted_at,
				deleted_by
			) VALUES (
				:id,
				:user_id,
				:title,
				:content,
				:created_at,
				:created_by,
				:updated_at,
				:updated_by,
				:deleted_at,
				:deleted_by
			)
		`,
	}
)

type CourseRepository interface {
	CreateCourse(course Course) (err error)
	ResolveCourses(page int, limit int, sort string, order string) (courses []Course, err error)
}

type CourseRepositoryMySQL struct {
	DB *infras.MySQLConn
}

func ProvideCourseRepositoryMySQL(db *infras.MySQLConn) *CourseRepositoryMySQL {
	s := new(CourseRepositoryMySQL)
	s.DB = db

	return s
}

func (r *CourseRepositoryMySQL) CreateCourse(course Course) (err error) {
	exists, err := r.ExistsByID(course.ID)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	if exists {
		err = failure.Conflict("create", "userId", "already exists")
		logger.ErrorWithStack(err)
		return
	}

	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := r.txCreate(tx, course); err != nil {
			e <- err
			return
		}

		e <- nil
	})
}

func (r *CourseRepositoryMySQL) ResolveCourses(page int, limit int, sort string, order string) (courses []Course, err error) {
	var args []interface{}

	query := courseQueries.selectCourses

	if sort != "" {
		validColumns := map[string]bool{
			"id":         true,
			"user_id":    true,
			"title":      true,
			"content":    false,
			"created_at": true,
			"created_by": true,
			"updated_at": true,
			"updated_by": true,
			"deleted_at": true,
			"deleted_by": true,
		}
		if !validColumns[sort] {
			return nil, errors.New("Invalid sort parameter")
		}

		validOrders := map[string]bool{
			"asc":  true,
			"desc": true,
		}
		if !validOrders[order] {
			return nil, errors.New("Invalid order parameter")
		}

		if order == "" {
			order = "asc"
		}

		query += fmt.Sprintf(" ORDER BY %s %s", sort, order)
	}

	offset := page * limit
	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	err = r.DB.Read.Select(&courses, query, args...)

	if err != nil {
		if err == sql.ErrNoRows {
			err = failure.NotFound("courses")
			logger.ErrorWithStack(err)
			return
		}

		logger.ErrorWithStack(err)
		return
	}

	return
}

func (r *CourseRepositoryMySQL) ExistsByID(id uuid.UUID) (exists bool, err error) {
	err = r.DB.Read.Get(
		&exists,
		"SELECT COUNT(id) FROM courses WHERE id = ?",
		id.String())

	if err != nil {
		logger.ErrorWithStack(err)
	}

	return
}

// internal functions
func (r *CourseRepositoryMySQL) txCreate(tx *sqlx.Tx, course Course) (err error) {
	stmt, err := tx.PrepareNamed(courseQueries.insertCourse)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(course)
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return
}
