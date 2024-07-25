package services

import (
	"errors"
	"exoplanet-microservice/models"
	"math"
)

var (
	ErrNotFound = errors.New("exoplanet not found")
	ErrInvalid  = errors.New("invalid exoplanet data")
)

type ExoplanetService interface {
	AddExoplanet(exoplanet models.Exoplanet) int
	ListExoplanets() []models.Exoplanet
	GetExoplanetByID(id int) (models.Exoplanet, error)
	UpdateExoplanet(updatedExoplanet models.Exoplanet) (models.Exoplanet, error)
	DeleteExoplanet(id int) error
	EstimateFuel(id int, crewCapacity int) (float64, error)
}

type Service struct {
	exoplanets map[int]models.Exoplanet
	nextID     int
}

func NewService() *Service {
	return &Service{
		exoplanets: make(map[int]models.Exoplanet),
		nextID:     1,
	}
}

func (s *Service) AddExoplanet(exoplanet models.Exoplanet) int {
	exoplanet.ID = s.nextID
	s.exoplanets[s.nextID] = exoplanet
	s.nextID++
	return exoplanet.ID
}

func (s *Service) ListExoplanets() []models.Exoplanet {
	exoplanets := make([]models.Exoplanet, 0, len(s.exoplanets))
	for _, exoplanet := range s.exoplanets {
		exoplanets = append(exoplanets, exoplanet)
	}
	return exoplanets
}

func (s *Service) GetExoplanetByID(id int) (models.Exoplanet, error) {
	exoplanet, exists := s.exoplanets[id]
	if !exists {
		return models.Exoplanet{}, ErrNotFound
	}
	return exoplanet, nil
}

func (s *Service) UpdateExoplanet(updatedExoplanet models.Exoplanet) (models.Exoplanet, error) {
	_, exists := s.exoplanets[updatedExoplanet.ID]
	if !exists {
		return models.Exoplanet{}, ErrNotFound
	}
	s.exoplanets[updatedExoplanet.ID] = updatedExoplanet
	return updatedExoplanet, nil
}

func (s *Service) DeleteExoplanet(id int) error {
	_, exists := s.exoplanets[id]
	if !exists {
		return ErrNotFound
	}
	delete(s.exoplanets, id)
	return nil
}

func (s *Service) EstimateFuel(id int, crewCapacity int) (float64, error) {
	exoplanet, exists := s.exoplanets[id]
	if !exists {
		return 0, ErrNotFound
	}

	var gravity float64
	switch exoplanet.Type {
	case "GasGiant":
		gravity = 0.5 / math.Pow(exoplanet.Radius, 2)
	case "Terrestrial":
		gravity = *exoplanet.Mass / math.Pow(exoplanet.Radius, 2)
	default:
		return 0, ErrInvalid
	}

	fuel := float64(exoplanet.Distance) / math.Pow(gravity, 2) * float64(crewCapacity)
	return fuel, nil
}
