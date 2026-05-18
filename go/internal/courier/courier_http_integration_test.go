package courier

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"courier-technical-test/go/internal/response"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func TestCourierHTTPIntegrationCRUD(t *testing.T) {
	app, _ := newTestApp()

	createResponse := requestJSON(t, app, http.MethodPost, "/api/couriers", map[string]any{
		"name":  "Budiono Hadi Agung",
		"email": "budi-go-integration@example.com",
		"level": 2,
	})
	assertStatus(t, createResponse, http.StatusCreated)
	id := stringPath(t, createResponse.Body, "data.id")

	listResponse := requestJSON(t, app, http.MethodGet, "/api/couriers", nil)
	assertStatus(t, listResponse, http.StatusOK)
	assertNames(t, listResponse.Body, []string{"Budiono Hadi Agung"})

	showResponse := requestJSON(t, app, http.MethodGet, "/api/couriers/"+id, nil)
	assertStatus(t, showResponse, http.StatusOK)
	assertPath(t, showResponse.Body, "data.name", "Budiono Hadi Agung")

	updateResponse := requestJSON(t, app, http.MethodPut, "/api/couriers/"+id, map[string]any{
		"name":   "Budi Updated",
		"level":  3,
		"status": "active",
	})
	assertStatus(t, updateResponse, http.StatusOK)
	assertPath(t, updateResponse.Body, "data.level", float64(3))

	deleteResponse := requestJSON(t, app, http.MethodDelete, "/api/couriers/"+id, nil)
	assertStatus(t, deleteResponse, http.StatusOK)

	deletedShowResponse := requestJSON(t, app, http.MethodGet, "/api/couriers/"+id, nil)
	assertStatus(t, deletedShowResponse, http.StatusNotFound)

	deletedListResponse := requestJSON(t, app, http.MethodGet, "/api/couriers", nil)
	assertStatus(t, deletedListResponse, http.StatusOK)
	assertNames(t, deletedListResponse.Body, []string{})
}

func TestCourierHTTPIntegrationListBehavior(t *testing.T) {
	app, store := newTestApp()
	store.seed(
		courierFixture("Budiono Hadi Agung", 2, "2026-05-18"),
		courierFixture("Budi Santoso", 2, "2026-05-17"),
		courierFixture("Agung Prasetyo", 3, "2026-05-16"),
		courierFixture("Rudi Hartono", 4, "2026-05-15"),
	)

	defaultSortResponse := requestJSON(t, app, http.MethodGet, "/api/couriers", nil)
	assertStatus(t, defaultSortResponse, http.StatusOK)
	assertNames(t, defaultSortResponse.Body, []string{"Agung Prasetyo", "Budi Santoso", "Budiono Hadi Agung", "Rudi Hartono"})

	registeredSortResponse := requestJSON(t, app, http.MethodGet, "/api/couriers?sort=registered_at", nil)
	assertStatus(t, registeredSortResponse, http.StatusOK)
	assertNames(t, registeredSortResponse.Body, []string{"Rudi Hartono", "Agung Prasetyo", "Budi Santoso", "Budiono Hadi Agung"})

	searchResponse := requestJSON(t, app, http.MethodGet, "/api/couriers?search=budi+agung", nil)
	assertStatus(t, searchResponse, http.StatusOK)
	assertNames(t, searchResponse.Body, []string{"Budiono Hadi Agung"})

	singleLevelResponse := requestJSON(t, app, http.MethodGet, "/api/couriers?level=2", nil)
	assertStatus(t, singleLevelResponse, http.StatusOK)
	assertNames(t, singleLevelResponse.Body, []string{"Budi Santoso", "Budiono Hadi Agung"})

	multipleLevelResponse := requestJSON(t, app, http.MethodGet, "/api/couriers?level=2,3", nil)
	assertStatus(t, multipleLevelResponse, http.StatusOK)
	assertNames(t, multipleLevelResponse.Body, []string{"Agung Prasetyo", "Budi Santoso", "Budiono Hadi Agung"})

	firstPage := requestJSON(t, app, http.MethodGet, "/api/couriers?page=1&per_page=2", nil)
	secondPage := requestJSON(t, app, http.MethodGet, "/api/couriers?page=2&per_page=2", nil)
	assertStatus(t, firstPage, http.StatusOK)
	assertStatus(t, secondPage, http.StatusOK)
	assertPath(t, firstPage.Body, "pagination.page", float64(1))
	assertPath(t, firstPage.Body, "pagination.per_page", float64(2))
	assertPath(t, firstPage.Body, "pagination.total", float64(4))
	assertPath(t, firstPage.Body, "pagination.total_pages", float64(2))
	assertNoDuplicateIDs(t, firstPage.Body, secondPage.Body)
}

func TestCourierHTTPIntegrationValidation(t *testing.T) {
	app, _ := newTestApp()

	assertStatus(t, requestJSON(t, app, http.MethodGet, "/api/couriers?level=9", nil), http.StatusBadRequest)
	assertStatus(t, requestJSON(t, app, http.MethodGet, "/api/couriers?page=abc", nil), http.StatusBadRequest)
	assertStatus(t, requestJSON(t, app, http.MethodGet, "/api/couriers?per_page=101", nil), http.StatusBadRequest)
	assertStatus(t, requestJSON(t, app, http.MethodGet, "/api/couriers?sort=deleted_at", nil), http.StatusBadRequest)
	assertStatus(t, requestJSON(t, app, http.MethodGet, "/api/couriers/not-a-valid-object-id", nil), http.StatusNotFound)

	invalidBodies := []map[string]any{
		{"level": 2},
		{"name": "", "level": 2},
		{"name": "Invalid Level", "level": 9},
		{"name": "Invalid Email", "level": 2, "email": "invalid"},
		{"name": "Invalid Status", "level": 2, "status": "unknown"},
		{"name": strings.Repeat("a", 151), "level": 2},
	}
	for _, body := range invalidBodies {
		assertStatus(t, requestJSON(t, app, http.MethodPost, "/api/couriers", body), http.StatusUnprocessableEntity)
	}
}

type testStore struct {
	mutex    sync.Mutex
	couriers map[bson.ObjectID]Courier
}

func newTestApp() (testApp, *testStore) {
	store := &testStore{couriers: map[bson.ObjectID]Courier{}}
	handler := NewHandler(NewService(store))
	app := fiber.New(fiber.Config{AppName: "Courier Technical Test Go"})
	app.Use(func(c *fiber.Ctx) error {
		response.RequestID(c)
		return c.Next()
	})
	couriers := app.Group("/api/couriers")
	couriers.Get("/", handler.Index)
	couriers.Post("/", handler.Store)
	couriers.Get("/:id", handler.Show)
	couriers.Put("/:id", handler.Update)
	couriers.Delete("/:id", handler.Destroy)

	return app, store
}

type testApp interface {
	Test(req *http.Request, msTimeout ...int) (*http.Response, error)
}

func (store *testStore) seed(couriers ...Courier) {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	for _, courier := range couriers {
		if courier.ID.IsZero() {
			courier.ID = bson.NewObjectID()
		}
		store.couriers[courier.ID] = courier
	}
}

func (store *testStore) Create(_ context.Context, courier Courier) (Courier, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	courier.ID = bson.NewObjectID()
	store.couriers[courier.ID] = courier
	return courier, nil
}

func (store *testStore) List(_ context.Context, params QueryParams) (ListResult, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	filtered := make([]Courier, 0, len(store.couriers))
	for _, courier := range store.couriers {
		if courier.DeletedAt != nil || !matchesLevels(courier, params.Levels) || !matchesSearch(courier, params.Search) {
			continue
		}
		filtered = append(filtered, courier)
	}
	sortCouriers(filtered, params.Sort)

	total := int64(len(filtered))
	start := (params.Page - 1) * params.PerPage
	if start > len(filtered) {
		start = len(filtered)
	}
	end := start + params.PerPage
	if end > len(filtered) {
		end = len(filtered)
	}
	totalPages := total / int64(params.PerPage)
	if total%int64(params.PerPage) != 0 {
		totalPages++
	}

	return ListResult{Couriers: filtered[start:end], Page: params.Page, PerPage: params.PerPage, Total: total, TotalPages: totalPages}, nil
}

func (store *testStore) FindByID(_ context.Context, id bson.ObjectID) (Courier, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	courier, ok := store.couriers[id]
	if !ok || courier.DeletedAt != nil {
		return Courier{}, mongo.ErrNoDocuments
	}
	return courier, nil
}

func (store *testStore) Update(_ context.Context, id bson.ObjectID, courier Courier) (Courier, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	current, ok := store.couriers[id]
	if !ok || current.DeletedAt != nil {
		return Courier{}, mongo.ErrNoDocuments
	}
	courier.ID = id
	courier.CreatedAt = current.CreatedAt
	courier.DeletedAt = current.DeletedAt
	store.couriers[id] = courier
	return courier, nil
}

func (store *testStore) SoftDelete(_ context.Context, id bson.ObjectID, deletedAt time.Time) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	courier, ok := store.couriers[id]
	if !ok || courier.DeletedAt != nil {
		return mongo.ErrNoDocuments
	}
	courier.DeletedAt = &deletedAt
	courier.UpdatedAt = deletedAt
	store.couriers[id] = courier
	return nil
}

func courierFixture(name string, level int, registeredAt string) Courier {
	date, err := time.Parse("2006-01-02", registeredAt)
	if err != nil {
		panic(err)
	}
	now := time.Now().UTC()
	return Courier{ID: bson.NewObjectID(), Name: name, Level: level, Status: "active", RegisteredAt: date, CreatedAt: now, UpdatedAt: now, DeletedAt: nil}
}

func matchesLevels(courier Courier, levels []int) bool {
	if len(levels) == 0 {
		return true
	}
	for _, level := range levels {
		if courier.Level == level {
			return true
		}
	}
	return false
}

func matchesSearch(courier Courier, search string) bool {
	for _, word := range strings.Fields(search) {
		if !strings.Contains(strings.ToLower(courier.Name), strings.ToLower(word)) {
			return false
		}
	}
	return true
}

func sortCouriers(couriers []Courier, sortParam string) {
	sort.Slice(couriers, func(leftIndex int, rightIndex int) bool {
		left := couriers[leftIndex]
		right := couriers[rightIndex]
		switch sortParam {
		case "registered_at":
			return left.RegisteredAt.Before(right.RegisteredAt)
		case "-registered_at":
			return right.RegisteredAt.Before(left.RegisteredAt)
		case "created_at":
			return left.CreatedAt.Before(right.CreatedAt)
		case "-created_at":
			return right.CreatedAt.Before(left.CreatedAt)
		case "-name":
			return left.Name > right.Name
		default:
			return left.Name < right.Name
		}
	})
}

type jsonResponse struct {
	StatusCode int
	Body       map[string]any
}

func requestJSON(t *testing.T, app testApp, method string, target string, body map[string]any) jsonResponse {
	t.Helper()
	var requestBody io.Reader
	if body != nil {
		encodedBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		requestBody = bytes.NewReader(encodedBody)
	}
	request := httptestRequest(t, method, target, requestBody)
	if body != nil {
		request.Header.Set("content-type", "application/json")
	}
	response, err := app.Test(request, -1)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	var decoded map[string]any
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil {
		t.Fatal(err)
	}
	return jsonResponse{StatusCode: response.StatusCode, Body: decoded}
}

func httptestRequest(t *testing.T, method string, target string, body io.Reader) *http.Request {
	t.Helper()
	request, err := http.NewRequest(method, target, body)
	if err != nil {
		t.Fatal(err)
	}
	return request
}

func assertStatus(t *testing.T, response jsonResponse, status int) {
	t.Helper()
	if response.StatusCode != status {
		t.Fatalf("expected status %d, got %d with body %#v", status, response.StatusCode, response.Body)
	}
}

func assertPath(t *testing.T, body map[string]any, path string, expected any) {
	t.Helper()
	if actual := pathValue(t, body, path); actual != expected {
		t.Fatalf("expected %s to be %#v, got %#v", path, expected, actual)
	}
}

func stringPath(t *testing.T, body map[string]any, path string) string {
	t.Helper()
	value, ok := pathValue(t, body, path).(string)
	if !ok {
		t.Fatalf("expected %s to be string", path)
	}
	return value
}

func pathValue(t *testing.T, body map[string]any, path string) any {
	t.Helper()
	current := any(body)
	for _, part := range strings.Split(path, ".") {
		object, ok := current.(map[string]any)
		if !ok {
			t.Fatalf("expected object at %s", part)
		}
		current = object[part]
	}
	return current
}

func assertNames(t *testing.T, body map[string]any, expected []string) {
	t.Helper()
	data := body["data"].([]any)
	actual := make([]string, 0, len(data))
	for _, rawCourier := range data {
		actual = append(actual, rawCourier.(map[string]any)["name"].(string))
	}
	if strings.Join(actual, "|") != strings.Join(expected, "|") {
		t.Fatalf("expected names %#v, got %#v", expected, actual)
	}
}

func assertNoDuplicateIDs(t *testing.T, first map[string]any, second map[string]any) {
	t.Helper()
	seen := map[string]bool{}
	for _, page := range []map[string]any{first, second} {
		for _, rawCourier := range page["data"].([]any) {
			id := rawCourier.(map[string]any)["id"].(string)
			if seen[id] {
				t.Fatalf("duplicate id across pages: %s", id)
			}
			seen[id] = true
		}
	}
}
