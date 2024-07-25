package utils

import (
	"errors"
	"exoplanet-microservice/models"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateExoplanet(exoplanet models.Exoplanet) error {
	// Basic validation using tags
	if err := validate.Struct(exoplanet); err != nil {
		return err
	}

	// Additional custom validation
	if exoplanet.Type == "Terrestrial" {
		if exoplanet.Mass == nil {
			return errors.New("mass must be specified for Terrestrial planets")
		}
	}

	return nil
}
