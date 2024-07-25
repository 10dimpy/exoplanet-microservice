package models

type Exoplanet struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Distance    int      `json:"distance" validate:"required,min=10,max=1000"`
	Radius      float64  `json:"radius" validate:"required,min=0.1,max=10"`
	Mass        *float64 `json:"mass,omitempty"`
	Type        string   `json:"type" validate:"required,oneof=GasGiant Terrestrial"`
}
