package services

import (
	"exoplanet-microservice/models"
	"math"
	"testing"
)

func TestAddExoplanet(t *testing.T) {
	s := NewService()
	exoplanet := models.Exoplanet{
		ID:       1,
		Type:     "GasGiant",
		Radius:   10,
		Distance: 1000,
	}

	id := s.AddExoplanet(exoplanet)
	if id != 1 {
		t.Errorf("expected ID 1, got %d", id)
	}

	storedExoplanet, _ := s.GetExoplanetByID(id)
	if storedExoplanet != exoplanet {
		t.Errorf("expected exoplanet %v, got %v", exoplanet, storedExoplanet)
	}
}

func TestListExoplanets(t *testing.T) {
	s := NewService()
	exoplanet1 := models.Exoplanet{Type: "GasGiant", Radius: 10, Distance: 1000}
	exoplanet2 := models.Exoplanet{Type: "Terrestrial", Radius: 5, Distance: 2000}

	s.AddExoplanet(exoplanet1)
	s.AddExoplanet(exoplanet2)

	exoplanets := s.ListExoplanets()
	if len(exoplanets) != 2 {
		t.Errorf("expected 2 exoplanets, got %d", len(exoplanets))
	}
}

func TestGetExoplanetByID(t *testing.T) {
	s := NewService()
	exoplanet := models.Exoplanet{ID: 1, Type: "GasGiant", Radius: 10, Distance: 1000}
	id := s.AddExoplanet(exoplanet)

	storedExoplanet, err := s.GetExoplanetByID(id)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if storedExoplanet != exoplanet {
		t.Errorf("expected exoplanet %v, got %v", exoplanet, storedExoplanet)
	}
}

func TestUpdateExoplanet(t *testing.T) {
	s := NewService()
	exoplanet := models.Exoplanet{Type: "GasGiant", Radius: 10, Distance: 1000}
	id := s.AddExoplanet(exoplanet)

	updatedExoplanet := models.Exoplanet{ID: id, Type: "Terrestrial", Radius: 5, Distance: 2000}
	result, err := s.UpdateExoplanet(updatedExoplanet)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if result != updatedExoplanet {
		t.Errorf("expected exoplanet %v, got %v", updatedExoplanet, result)
	}
}

func TestDeleteExoplanet(t *testing.T) {
	s := NewService()
	exoplanet := models.Exoplanet{Type: "GasGiant", Radius: 10, Distance: 1000}
	id := s.AddExoplanet(exoplanet)

	err := s.DeleteExoplanet(id)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	_, err = s.GetExoplanetByID(id)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestEstimateFuel(t *testing.T) {
	s := NewService()
	exoplanet := models.Exoplanet{Type: "GasGiant", Radius: 10, Distance: 1000}
	id := s.AddExoplanet(exoplanet)

	fuel, err := s.EstimateFuel(id, 5)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expectedFuel := float64(exoplanet.Distance) / math.Pow(0.5/math.Pow(exoplanet.Radius, 2), 2) * 5
	if math.Abs(fuel-expectedFuel) > 1e-6 {
		t.Errorf("expected fuel %f, got %f", expectedFuel, fuel)
	}
}
