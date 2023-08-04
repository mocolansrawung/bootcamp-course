package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/evermos/boilerplate-go/internal/domain/course"
	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/transport/http/middleware"
	"github.com/evermos/boilerplate-go/transport/http/response"
	"github.com/go-chi/chi"
)

type CourseHandler struct {
	CourseService  course.CourseService
	AuthMiddleware *middleware.Authentication
}

func ProvideCourseHandler(courseService course.CourseService, authMiddleware *middleware.Authentication) CourseHandler {
	return CourseHandler{
		CourseService:  courseService,
		AuthMiddleware: authMiddleware,
	}
}

func (h *CourseHandler) Router(r chi.Router) {
	r.Route("/courses", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(h.AuthMiddleware.ValidateAuth)
			r.Post("/", h.CreateCourse)
			r.Get("/", h.ResolveCourses)
		})
	})
}

func (h *CourseHandler) CreateCourse(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var requestFormat course.CourseRequestFormat
	err := decoder.Decode(&requestFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
	}

	err = shared.GetValidator().Struct(requestFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	resp, ok := r.Context().Value("responseBody").(shared.Claims)
	if !ok {
		response.WithError(w, failure.Unauthorized("User not authorized"))
	}

	course, err := h.CourseService.CreateCourse(requestFormat, resp.UserID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	fmt.Println(course)
	response.WithJSON(w, http.StatusCreated, course)
}

func (h *CourseHandler) ResolveCourses(w http.ResponseWriter, r *http.Request) {
	pageString := r.URL.Query().Get("page")
	page, err := convertIdParamsToInt(pageString)
	if err != nil || page < 0 {
		page = 0
	}

	limitString := r.URL.Query().Get("limit")
	limit, err := convertIdParamsToInt(limitString)
	if err != nil || limit <= 0 {
		limit = 10
	}

	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	courses, err := h.CourseService.ResolveCourses(page, limit, sort, order)
	if err != nil {
		response.WithError(w, err)
		return
	}

	fmt.Println(courses)

	response.WithJSON(w, http.StatusOK, courses)
}

func convertIdParamsToInt(idStr string) (int, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("error converting ID parameter to integer: %w", err)
	}

	return id, nil
}
