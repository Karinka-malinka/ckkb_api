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
	CreatedAtFrom string   `json:"created_at_from,omitempty"`
	PerPage       int      `json:"per_page"`
}

type StatisticsResponse struct {
	Statistics Statistics `json:"statistics"`
}

type Statistics struct {
	NewTask New     `json:"new"`
	InWork  Process `json:"in_work"`
	Revise  Revise  `json:"revise"`
	Check   Check   `json:"check"`
}

type New struct {
	Count int        `json:"count"`
	Data  []TaskItem `json:"data"`
}

type Process struct {
	Count int        `json:"count"`
	Data  []TaskItem `json:"data"`
}

type Revise struct {
	Count int        `json:"count"`
	Data  []TaskItem `json:"data"`
}

type Check struct {
	Count int        `json:"count"`
	Data  []TaskItem `json:"data"`
}

// TaskItem соответствует одному элементу в поле "data"
type TaskItem struct {
	ID           int    `json:"id"`
	PublicID     int    `json:"issue_id"`
	IssueURL     string `json:"issue_url"`
	InspectionID int    `json:"inspection_id"`
	Status       string `json:"status"`
	Title        string `json:"topic"`
	Description  string `json:"object_type"`
	CreatedAt    string `json:"created_at"`
	//TODO добавить другие поля задачи
}

func (ih *IssueHandler) RegisterHandler(apiGroup *echo.Group) {
	apiGroup.POST("/list", ih.ListIssue)
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

	dataNew := []TaskItem{}
	dataProcess := []TaskItem{}
	dataRevise := []TaskItem{}
	dataCheck := []TaskItem{}

	for _, task := range listIssue {
		data := TaskItem{
			ID:           task.ID,
			PublicID:     task.PublicID,
			InspectionID: task.InspectionID,
			Title:        task.Title,
			Description:  task.Description,
			IssueURL:     fmt.Sprintf("https://ckkb-mos.online/tasks/%d/view", task.ID),
			Status:       task.Status,
		}

		switch task.Status {
		case "created":
			dataNew = append(dataNew, data)
		case "process":
			dataProcess = append(dataProcess, data)
		case "revise":
			dataRevise = append(dataRevise, data)
		case "validation":
			dataCheck = append(dataCheck, data)
		case "review":
			dataCheck = append(dataCheck, data)
		}
	}

	response := StatisticsResponse{
		Statistics: Statistics{
			NewTask: New{
				Count: len(dataNew),
				Data:  dataNew,
			},
			InWork: Process{
				Count: len(dataProcess),
				Data:  dataProcess,
			},
			Revise: Revise{
				Count: len(dataRevise),
				Data:  dataRevise,
			},
			Check: Check{
				Count: len(dataCheck),
				Data:  dataCheck,
			},
		},
	}

	return c.JSON(http.StatusOK, response)
}
