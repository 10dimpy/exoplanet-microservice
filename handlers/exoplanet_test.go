package handlers

import (
	"bytes"
	"encoding/json"
	"exoplanet-microservice/models"
	"exoplanet-microservice/services"
	"github.com/gorilla/mux"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type MockService struct {
	exoplanets map[int]models.Exoplanet
}

func (m *MockService) AddExoplanet(exoplanet models.Exoplanet) int {
	id := len(m.exoplanets) + 1
	exoplanet.ID = id
	m.exoplanets[id] = exoplanet
	return id
}

func (m *MockService) ListExoplanets() []models.Exoplanet {
	exoplanets := make([]models.Exoplanet, 0, len(m.exoplanets))
	for _, exoplanet := range m.exoplanets {
		exoplanets = append(exoplanets, exoplanet)
	}
	return exoplanets
}

func (m *MockService) GetExoplanetByID(id int) (models.Exoplanet, error) {
	exoplanet, exists := m.exoplanets[id]
	if !exists {
		return models.Exoplanet{}, services.ErrNotFound
	}
	return exoplanet, nil
}

func (m *MockService) UpdateExoplanet(updatedExoplanet models.Exoplanet) (models.Exoplanet, error) {
	if _, exists := m.exoplanets[updatedExoplanet.ID]; !exists {
		return models.Exoplanet{}, services.ErrNotFound
	}
	m.exoplanets[updatedExoplanet.ID] = updatedExoplanet
	return updatedExoplanet, nil
}

func (m *MockService) DeleteExoplanet(id int) error {
	if _, exists := m.exoplanets[id]; !exists {
		return services.ErrNotFound
	}
	delete(m.exoplanets, id)
	return nil
}

func (m *MockService) EstimateFuel(id int, crewCapacity int) (float64, error) {
	exoplanet, exists := m.exoplanets[id]
	if !exists {
		return 0, services.ErrNotFound
	}

	var gravity float64
	switch exoplanet.Type {
	case "GasGiant":
		gravity = 0.5 / math.Pow(exoplanet.Radius, 2)
	case "Terrestrial":
		gravity = *exoplanet.Mass / math.Pow(exoplanet.Radius, 2)
	default:
		return 0, services.ErrInvalid
	}

	fuel := float64(exoplanet.Distance) / math.Pow(gravity, 2) * float64(crewCapacity)
	return fuel, nil
}

func TestAddExoplanet(t *testing.T) {
	mockService := &MockService{exoplanets: make(map[int]models.Exoplanet)}
	handler := &Handler{Service: mockService}

	// Define an exoplanet with all required fields
	exoplanet := models.Exoplanet{Name: "TerraNova", Description: "A rocky planet similar to Earth.", Type: "GasGiant", Radius: 10, Distance: 1000}
	body, _ := json.Marshal(exoplanet)

	// Create a new POST request to add an exoplanet
	req := httptest.NewRequest("POST", "/exoplanets", bytes.NewReader(body))
	w := httptest.NewRecorder()

	// Call the AddExoplanet handler
	handler.AddExoplanet(w, req)

	// Check the response status code
	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", res.StatusCode)
	}

	// Log the response body for debugging
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	t.Logf("Response body: %s", bodyBytes)

	// Decode the response body
	var result models.Exoplanet
	if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Check if the type matches
	if result.Type != exoplanet.Type {
		t.Errorf("expected type %s, got %s", exoplanet.Type, result.Type)
	}
}

func TestListExoplanets(t *testing.T) {
	mockService := &MockService{exoplanets: make(map[int]models.Exoplanet)}
	handler := &Handler{Service: mockService}

	exoplanet1 := models.Exoplanet{Type: "GasGiant", Radius: 10, Distance: 1000}
	exoplanet2 := models.Exoplanet{Type: "Terrestrial", Radius: 5, Distance: 2000}
	mockService.AddExoplanet(exoplanet1)
	mockService.AddExoplanet(exoplanet2)

	req := httptest.NewRequest("GET", "/exoplanets", nil)
	w := httptest.NewRecorder()

	handler.ListExoplanets(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", res.StatusCode)
	}

	var results []models.Exoplanet
	if err := json.NewDecoder(res.Body).Decode(&results); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("expected 2 exoplanets, got %d", len(results))
	}
}

func TestGetExoplanetByID(t *testing.T) {
	mockService := &MockService{exoplanets: make(map[int]models.Exoplanet)}
	handler := &Handler{Service: mockService}

	exoplanet := models.Exoplanet{Type: "GasGiant", Radius: 10, Distance: 1000}
	id := mockService.AddExoplanet(exoplanet)

	req := httptest.NewRequest("GET", "/exoplanets/{id:[0-9]+}", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.GetExoplanetByID(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", res.StatusCode)
	}

	var result models.Exoplanet
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if result.ID != id {
		t.Errorf("expected ID %d, got %d", id, result.ID)
	}
}

func TestUpdateExoplanet(t *testing.T) {
	// Initialize the mock service and handler
	mockService := &MockService{exoplanets: make(map[int]models.Exoplanet)}
	handler := &Handler{Service: mockService}

	// Create and add an exoplanet to the mock service
	exoplanet := models.Exoplanet{Type: "GasGiant", Radius: 10, Distance: 1000}
	id := mockService.AddExoplanet(exoplanet)

	// Prepare the updated exoplanet data
	updatedExoplanet := models.Exoplanet{ID: id, Type: "GasGiant", Radius: 5, Distance: 800}
	body, _ := json.Marshal(updatedExoplanet)
	req := httptest.NewRequest("PUT", "/exoplanets/{id}", bytes.NewReader(body))
	req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(id)})
	w := httptest.NewRecorder()

	// Call the UpdateExoplanet handler
	handler.UpdateExoplanet(w, req)

	// Check the response status code
	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", res.StatusCode)
	}

	// Log the response body for debugging
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	t.Logf("Response body: %s", bodyBytes)

	// Decode the response body
	var result models.Exoplanet
	if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Check if the type matches
	if result.Type != updatedExoplanet.Type {
		t.Errorf("expected type %s, got %s", updatedExoplanet.Type, result.Type)
	}

	// Optionally, check other fields if necessary
	if result.Radius != updatedExoplanet.Radius {
		t.Errorf("expected radius %f, got %f", updatedExoplanet.Radius, result.Radius)
	}
	if result.Distance != updatedExoplanet.Distance {
		t.Errorf("expected distance %d, got %d", updatedExoplanet.Distance, result.Distance)
	}
}

func TestDeleteExoplanet(t *testing.T) {
	mockService := &MockService{exoplanets: make(map[int]models.Exoplanet)}
	handler := &Handler{Service: mockService}

	exoplanet := models.Exoplanet{Type: "GasGiant", Radius: 10, Distance: 1000}
	id := mockService.AddExoplanet(exoplanet)

	req := httptest.NewRequest("DELETE", "/exoplanets/{id:[0-9]+}", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.DeleteExoplanet(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204 No Content, got %d", res.StatusCode)
	}

	_, err := mockService.GetExoplanetByID(id)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestEstimateFuel(t *testing.T) {
	// Initialize mock service and handler
	mockService := &MockService{exoplanets: make(map[int]models.Exoplanet)}
	handler := &Handler{Service: mockService}

	// Create an exoplanet and add it to the mock service
	exoplanet := models.Exoplanet{Type: "GasGiant", Radius: 10, Distance: 1000}
	id := mockService.AddExoplanet(exoplanet) // Ensure the exoplanet is added with a valid ID

	// Create a GET request to estimate fuel for the exoplanet with ID 1 and crew capacity 5
	req := httptest.NewRequest("GET", "/exoplanets/{id:[0-9]+}/fuel?crew=5", nil)
	req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(id)})
	w := httptest.NewRecorder()

	// Call the EstimateFuel handler
	handler.EstimateFuel(w, req)

	// Check the response status code
	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", res.StatusCode)
	}

	// Log the response body for debugging
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	t.Logf("Response body: %s", bodyBytes)

	// Decode the response body
	var result map[string]float64
	if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Calculate expected fuel
	gravity := 0.5 / math.Pow(exoplanet.Radius, 2)
	expectedFuel := float64(exoplanet.Distance) / math.Pow(gravity, 2) * 5

	// Compare the calculated fuel with the result
	if math.Abs(result["fuel"]-expectedFuel) > 1e-6 {
		t.Errorf("expected fuel %f, got %f", expectedFuel, result["fuel"])
	}
}
