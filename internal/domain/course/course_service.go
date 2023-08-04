package course

import (
	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/gofrs/uuid"
)

type CourseService interface {
	CreateCourse(requestFormat CourseRequestFormat, userID uuid.UUID) (course Course, err error)
	ResolveCourses(page int, limit int, sort string, order string) (courses []Course, err error)
}

type CourseServiceImpl struct {
	CourseRepository CourseRepository
	Config           *configs.Config
}

func ProvideCourseServiceImpl(courseRepository CourseRepository, config *configs.Config) *CourseServiceImpl {
	s := new(CourseServiceImpl)
	s.CourseRepository = courseRepository
	s.Config = config

	return s
}

func (s *CourseServiceImpl) CreateCourse(requestFormat CourseRequestFormat, userID uuid.UUID) (course Course, err error) {
	course, err = course.NewCourseFromRequestFormat(requestFormat, userID)
	if err != nil {
		return
	}

	if err != nil {
		return course, failure.BadRequest(err)
	}

	err = s.CourseRepository.CreateCourse(course)
	if err != nil {
		return
	}

	return
}

func (s *CourseServiceImpl) ResolveCourses(page int, limit int, sort string, order string) (courses []Course, err error) {
	courses, err = s.CourseRepository.ResolveCourses(page, limit, sort, order)
	if err != nil {
		return courses, failure.BadRequest(err)
	}

	return
}
