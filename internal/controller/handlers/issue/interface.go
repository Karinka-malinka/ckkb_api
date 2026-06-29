package issue

import "github.com/ckkb_api/internal/domain/issue"

type IIssueAction interface {
	GetListIssueCheckOffice(params *issue.RequestParams) ([]issue.TaskItem, error)
}
