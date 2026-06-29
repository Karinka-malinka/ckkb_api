package issue

type IIssueStore interface {
	AddLog(reqURL string) error
}
