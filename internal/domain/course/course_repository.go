package course

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

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
	ResolveCourses(params CourseQueryParameters) (courses []Course, err error)
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

func (r *CourseRepositoryMySQL) ResolveCourses(params CourseQueryParameters) (courses []Course, err error) {
	var args []interface{}

	query := courseQueries.selectCourses

	if params.Role != "" {
		query += " WHERE role = ?"
		args = append(args, params.Role)
	}

	if params.Sort != "" {
		isValid, err := r.isValidColumnName(params.Sort)
		if err != nil {
			return nil, err
		}
		if !isValid || !r.isSortableColumn(params.Sort) {
			return nil, errors.New("Invalid sort parameter")
		}

		params.Order, err = r.validateAndCorrectOrder(params.Order)
		if err != nil {
			return nil, err
		}

		query += fmt.Sprintf(" ORDER BY %s %s", params.Sort, params.Order)
	}

	offset := params.Page * params.Limit
	query += " LIMIT ? OFFSET ?"
	args = append(args, params.Limit, offset)

	err = r.DB.Read.Select(&courses, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			err = failure.NotFound("courses")
			logger.ErrorWithStack(err)
			return nil, err
		}

		logger.ErrorWithStack(err)
		return nil, err
	}

	return courses, nil
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

// Internal Functions
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

func (r *CourseRepositoryMySQL) isValidColumnName(columnName string) (bool, error) {
	var columns []string
	const query = `SELECT COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'courses' AND TABLE_SCHEMA = DATABASE()`

	if err := r.DB.Read.Select(&columns, query); err != nil {
		return false, err
	}

	for _, column := range columns {
		if column == columnName {
			return true, nil
		}
	}

	return false, nil
}

func (r *CourseRepositoryMySQL) isSortableColumn(columnName string) bool {
	sortableColumns := map[string]bool{
		"id":         true,
		"user_id":    true,
		"title":      true,
		"created_at": true,
		"created_by": true,
		"updated_at": true,
		"updated_by": true,
		"deleted_at": true,
		"deleted_by": true,
	}

	return sortableColumns[columnName]
}

func (r *CourseRepositoryMySQL) validateAndCorrectOrder(order string) (string, error) {
	order = strings.ToLower(order)
	if order != "asc" && order != "desc" {
		return "", errors.New("Invalid order parameter")
	}
	return order, nil
}
