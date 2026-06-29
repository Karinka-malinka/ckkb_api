package issue

import (
	"fmt"
	"net/http"

	"github.com/ckkb_api/internal/domain/issue"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type IssueHandler struct {
	actionIssue IIssueAction
}

func NewIssueHandler(actionIssue IIssueAction) *IssueHandler {
	return &IssueHandler{actionIssue}
}

type PaginationRequest struct {
	Statuses      []string `json:"statuses"`
	Assignees     []int    `json:"assignees"`
	CreatedAtFrom string   `json:"created_at_from"`
	PerPage       int      `json:"per_page"`
}

type ListIssueResponse struct {
	Status string     `json:"status"`
	Count  int        `json:"count"`
	Data   []TaskItem `json:"data"`
}

// TaskItem соответствует одному элементу в поле "data"
type TaskItem struct {
	ID           int    `json:"id"`
	PublicID     int    `json:"issue_id"`
	IssueURL     string `json:"issue_url"`
	InspectionID int    `json:"inspection_id"`
	Status       string `json:"status"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	CreatedAt    string `json:"created_at"`
	//TODO добавить другие поля задачи
}

func (ih *IssueHandler) RegisterHandler(apiGroup *echo.Group) {
	apiGroup.GET("/list", ih.ListIssue)
	//restrictedGroup.GET(path+"/settings/:"+queryID+"&"+queryPortal+"&"+queryGBUName, ph.Read)
	//restrictedGroup.GET(path+"/list", ph.List)
	//restrictedGroup.POST("/issue/listpaginated", ih.ListPaginated)
	//restrictedGroup.PUT(path, ph.Update)
}

func (ih *IssueHandler) ListIssue(c echo.Context) error {

	request := PaginationRequest{}

	err := c.Bind(&request)
	if err != nil {
		logrus.Errorf("error in bind Pagination request. error: %v", err)
		return c.JSON(http.StatusBadRequest, "please check request struct")
	}

	if err := c.Validate(request); err != nil {
		logrus.Errorf("error in bind register Pagination request. error: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "please check request struct")
	}

	param := issue.RequestParams{
		Statuses:      request.Statuses,
		Assignees:     request.Assignees,
		CreatedAtFrom: request.CreatedAtFrom,
		PerPage:       request.PerPage,
	}

	listIssue, err := ih.actionIssue.GetListIssueCheckOffice(&param)
	if err != nil {
		return c.String(http.StatusInternalServerError, "error in getting issue list")
	}

	data := make([]TaskItem, len(listIssue))

	for i, task := range listIssue {
		data[i] = TaskItem{
			ID:       task,
			IssueURL: fmt.Sprintf("https://ckkb-mos.online/publicapi/v1/tasks/%d/view", task),
		}
	}

	response := ListIssueResponse{
		Status: request.Statuses[0],
		Count:  len(listIssue),
		Data:   data,
	}

	return c.JSON(http.StatusOK, response)
}
