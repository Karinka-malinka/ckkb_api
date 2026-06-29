package issue

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type IssueAction struct {
	store IIssueStore
}

func NewIssueAction(store IIssueStore) *IssueAction {
	return &IssueAction{store: store}
}

const (
	apiKey       = "jTp6NlejUBRC2lWsuXKHKrDwtSLE3fiL0VRVkc2cpuLGXlHPOv6lc6zJGTT9"
	baseURLTasks = "https://ckkb-mos.online/publicapi/v1/tasks"
)

// RequestParams содержит все параметры для запроса
type RequestParams struct {
	Statuses      []string // например, []string{"created"}
	Assignees     []int    // например, []int{42085, 42089}
	CreatedAtFrom string   // формат "2026-06-29T11:00:00"
	PerPage       int
}

// TaskItem соответствует одному элементу в поле "data"
type TaskItem struct {
	ID           int    `json:"id"`
	PublicID     int    `json:"public_id"`
	InspectionID int    `json:"inspection_id"`
	Status       string `json:"status"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	CreatedAt    string `json:"created_at"`
	ExpireAt     string `json:"expire_at"`
	FinishedAt   string `json:"finished_at"`
}

// PaginatedResponse соответствует полному ответу сервера
type PaginatedResponse struct {
	Data        []TaskItem `json:"data"`
	Total       int        `json:"total"`
	PerPage     int        `json:"per_page"`
	CurrentPage int        `json:"current_page"`
	LastPage    int        `json:"last_page"`
	From        *int       `json:"from"`
	To          *int       `json:"to"`
	Status      string     `json:"status"`
}

func (ia *IssueAction) GetCreatedIssueCheckOffice(params *RequestParams) ([]TaskItem, error) {

	/*
		params_loc := RequestParams{
			Statuses:      []string{"created"},
			Assignees:     []int{42085, 42089},
			CreatedAtFrom: "2026-06-29T11:00:00",
			Page:          1,
		}
	*/

	params.Statuses = append(params.Statuses, "created")

	// Строим полный URL
	fullURL, err := buildURL(baseURLTasks, params, 1)
	if err != nil {
		log.Fatalf("Ошибка построения URL: %v", err)
	}

	// Множество для уникальных ID
	uniqueIDs := make(map[int]TaskItem)

	firstResp, err := ia.fetchPage(fullURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка при запросе первой страницы: %v", err)
	}

	// Добавляем ID из первой страницы
	for _, item := range firstResp.Data {
		uniqueIDs[item.InspectionID] = item
	}

	// Если всего одна страница или данных нет, возвращаем результат
	if firstResp.LastPage <= 1 || firstResp.Total == 0 {
		return mapKeysToSlice(uniqueIDs), nil
	}

	// Последовательно загружаем остальные страницы (со 2-й до последней)
	for page := 2; page <= firstResp.LastPage; page++ {

		pageURL, err := buildURL(baseURLTasks, params, page)
		if err != nil {
			return nil, err
		}
		resp, err := ia.fetchPage(pageURL)
		if err != nil {
			// Логируем ошибку, но продолжаем (можно изменить поведение)
			log.Printf("Предупреждение: ошибка при загрузке страницы %d: %v", page, err)
			continue
		}

		for _, item := range resp.Data {
			uniqueIDs[item.InspectionID] = item
		}
		// Небольшая задержка, чтобы не перегружать сервер
		time.Sleep(200 * time.Millisecond)
	}

	return mapKeysToSlice(uniqueIDs), nil
}

func (ia *IssueAction) GetListIssueCheckOffice(params *RequestParams) ([]TaskItem, error) {

	/*
		params_loc := RequestParams{
			Statuses:      []string{"created"},
			Assignees:     []int{42085, 42089},
			CreatedAtFrom: "2026-06-29T11:00:00",
			Page:          1,
		}
	*/

	//params.Statuses = append(params.Statuses, "process")

	// Строим полный URL
	fullURL, err := buildURL(baseURLTasks, params, 1)
	if err != nil {
		log.Fatalf("Ошибка построения URL: %v", err)
	}

	// Множество для уникальных ID
	uniqueIDs := make(map[int]TaskItem)

	firstResp, err := ia.fetchPage(fullURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка при запросе первой страницы: %v", err)
	}

	// Добавляем ID из первой страницы
	for _, item := range firstResp.Data {
		uniqueIDs[item.ID] = item
	}

	// Если всего одна страница или данных нет, возвращаем результат
	if firstResp.LastPage <= 1 || firstResp.Total == 0 {
		return mapKeysToSlice(uniqueIDs), nil
	}

	// Последовательно загружаем остальные страницы (со 2-й до последней)
	for page := 2; page <= firstResp.LastPage; page++ {

		pageURL, err := buildURL(baseURLTasks, params, page)
		if err != nil {
			return nil, err
		}
		resp, err := ia.fetchPage(pageURL)

		if err != nil {
			// Логируем ошибку, но продолжаем (можно изменить поведение)
			log.Printf("Предупреждение: ошибка при загрузке страницы %d: %v", page, err)
			continue
		}

		for _, item := range resp.Data {
			uniqueIDs[item.ID] = item
		}
		// Небольшая задержка, чтобы не перегружать сервер
		time.Sleep(200 * time.Millisecond)
	}

	return mapKeysToSlice(uniqueIDs), nil
}

// BuildURL собирает URL с правильно закодированными параметрами
func buildURL(base string, params *RequestParams, page int) (string, error) {
	u, err := url.Parse(base)
	if err != nil {
		return "", err
	}

	q := url.Values{}

	// Добавляем статусы (ключ "t:statuses[]")
	for _, status := range params.Statuses {
		q.Add("t:statuses[]", status)
	}

	// Добавляем исполнителей (ключ "t:assignees[]")
	for _, assignee := range params.Assignees {
		q.Add("t:assignees[]", strconv.Itoa(assignee))
	}

	// Дата создания
	if params.CreatedAtFrom != "" {
		q.Set("t:created_at_from", params.CreatedAtFrom)
	}

	// Страница
	q.Set("page", strconv.Itoa(page))

	// Записей на странице
	if params.PerPage > 0 {
		q.Set("per_page", strconv.Itoa(params.PerPage))
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}

// fetchPage выполняет один запрос к указанной странице и возвращает распарсенный ответ
func (ia *IssueAction) fetchPage(url string) (*PaginatedResponse, error) {

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("API-Key", apiKey)

	//добавляю вызов api в log
	ia.store.AddLog(url)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус: %s, тело: %s", resp.Status, string(body))
	}

	var result PaginatedResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON: %v, тело: %s", err, string(body))
	}
	return &result, nil
}

// Вспомогательная функция для преобразования map-ключа в слайс
func mapKeysToSlice(m map[int]TaskItem) []TaskItem {

	result := make([]TaskItem, 0, len(m))
	for _, v := range m {
		result = append(result, v)
	}
	return result
}
