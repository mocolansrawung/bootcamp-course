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
			r.Use(h.AuthMiddleware.UserRoleCheck)
			r.Get("/", h.ResolveCourses)
			r.Post("/", h.CreateCourse)
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

	response.WithJSON(w, http.StatusCreated, course)
}

func (h *CourseHandler) ResolveCourses(w http.ResponseWriter, r *http.Request) {
	pageString := r.URL.Query().Get("page")
	page, err := convertQueryParamsToInt(pageString)
	if err != nil || page < 0 {
		page = 0
	}

	limitString := r.URL.Query().Get("limit")
	limit, err := convertQueryParamsToInt(limitString)
	if err != nil || limit <= 0 {
		limit = 10
	}

	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	resp, ok := r.Context().Value("responseBody").(shared.Claims)
	if !ok {
		response.WithError(w, failure.InternalError(err))
	}

	params := course.CourseQueryParameters{
		Page:  page,
		Limit: limit,
		Sort:  sort,
		Order: order,
		Role:  resp.Role,
	}

	courses, err := h.CourseService.ResolveCourses(params)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, courses)
}

func convertQueryParamsToInt(idStr string) (int, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("error converting ID parameter to integer: %w", err)
	}

	return id, nil
}
