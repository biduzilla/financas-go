package handler

import (
	"financas/internal/model"
	"financas/internal/model/filters"
	"financas/internal/service"
	"financas/utils"
	e "financas/utils/errors"
	"financas/utils/validator"
	"fmt"
	"net/http"
)

type CategoryHandler struct {
	ContextGetUser  func(r *http.Request) *model.User
	CategoryService service.CategoryServiceInterface
	errorResponse   e.ErrorResponseInterface
}

type CategoryHandlerInterface interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	GetById(w http.ResponseWriter, r *http.Request)
	Insert(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

func NewCategoryHandler(c service.CategoryServiceInterface, ContextGetUser func(r *http.Request) *model.User, errorResponse e.ErrorResponseInterface) *CategoryHandler {
	return &CategoryHandler{
		ContextGetUser:  ContextGetUser,
		CategoryService: c,
		errorResponse:   errorResponse,
	}
}

func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string
		filters.Filters
	}

	v := validator.New()

	qs := r.URL.Query()
	input.Name = utils.ReadString(qs, "name", "")
	input.Filters.Page = utils.ReadInt(qs, "page", 1, v)
	input.Filters.PageSize = utils.ReadInt(qs, "page_size", 20, v)
	input.Filters.Sort = utils.ReadString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "name", "-id", "-name"}

	if filters.ValidateFilters(v, input.Filters); !v.Valid() {
		h.errorResponse.HandlerErrorResponse(w, r, e.ErrInvalidData, v)
		return
	}

	user := h.ContextGetUser(r)
	categories, metadata, err := h.CategoryService.GetAll(input.Name, user.ID, input.Filters, v)
	if err != nil {
		h.errorResponse.HandlerErrorResponse(w, r, err, v)
		return
	}

	categoriesDTO := []*model.CategoryDTO{}

	for _, c := range categories {
		c.User = user
		categoriesDTO = append(categoriesDTO, c.ToDTO())
	}

	h.respond(w, r, http.StatusOK, utils.Envelope{"categories": categoriesDTO, "metadata": metadata}, nil)
}

func (h *CategoryHandler) GetById(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadIntParam(r, "id")
	if err != nil {
		h.errorResponse.BadRequestResponse(w, r, err)
		return
	}

	user := h.ContextGetUser(r)
	c, err := h.CategoryService.GetByID(id, user.ID)

	if err != nil {
		h.errorResponse.HandlerErrorResponse(w, r, err, nil)
		return
	}

	c.User = user
	h.respond(w, r, http.StatusOK, utils.Envelope{"category": c.ToDTO()}, nil)
}

func (h *CategoryHandler) Insert(w http.ResponseWriter, r *http.Request) {
	var dto model.CategoryDTO
	if err := utils.ReadJSON(w, r, &dto); err != nil {
		h.errorResponse.BadRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	user := h.ContextGetUser(r)
	category := dto.ToModel()

	if err := h.CategoryService.Insert(category, v); err != nil {
		h.errorResponse.HandlerErrorResponse(w, r, err, v)
		return
	}

	headers := http.Header{"Location": {fmt.Sprintf("/v1/categories/%d", category.ID)}}

	category.User = user
	h.respond(w, r, http.StatusCreated, utils.Envelope{"category": category.ToDTO()}, headers)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := h.parseID(w, r)
	if !ok {
		return
	}

	dto := model.CategoryDTO{}
	if err := utils.ReadJSON(w, r, &dto); err != nil {
		h.errorResponse.BadRequestResponse(w, r, err)
		return
	}

	dto.ID = &id

	category := dto.ToModel()
	v := validator.New()
	user := h.ContextGetUser(r)

	err := h.CategoryService.Update(category, user.ID, v)
	if err != nil {
		h.errorResponse.HandlerErrorResponse(w, r, err, v)
		return
	}

	category.User = user
	h.respond(w, r, http.StatusOK, utils.Envelope{"category": category.ToDTO()}, nil)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := h.parseID(w, r)
	if !ok {
		return
	}

	user := h.ContextGetUser(r)

	err := h.CategoryService.Delete(id, user.ID)
	if err != nil {
		h.errorResponse.HandlerErrorResponse(w, r, err, nil)
		return
	}

	h.respond(w, r, http.StatusOK, utils.Envelope{"message": "category successfully deleted"}, nil)
}

func (h *CategoryHandler) parseID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := utils.ReadIntParam(r, "id")
	if err != nil {
		h.errorResponse.BadRequestResponse(w, r, err)
		return 0, false
	}
	return id, true
}

func (h *CategoryHandler) respond(w http.ResponseWriter, r *http.Request, status int, data utils.Envelope, headers http.Header) {
	err := utils.WriteJSON(w, status, data, headers)
	if err != nil {
		h.errorResponse.ServerErrorResponse(w, r, err)
	}
}
