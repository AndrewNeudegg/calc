package constants

import "github.com/andrewneudegg/calc/pkg/units"

// initUniversal initializes universal and astronomical constants.
// Values are from IAU 2015 Resolution B3 and other standard references.
func (s *System) initUniversal() {
	// Gravitational constant
	s.addConstant(
		"gravitational_constant",
		"G",
		6.67430e-11,
		"m³/(kg·s²)",
		units.DimensionNone,
		"Newtonian constant of gravitation",
		"universal",
	)
	
	// Standard gravity (Earth surface)
	s.addConstant(
		"standard_gravity",
		"g_n",
		9.80665,
		"m/s²",
		units.DimensionNone,
		"Standard acceleration due to gravity",
		"universal",
	)
	
	// Stefan-Boltzmann constant
	s.addConstant(
		"stefan_boltzmann",
		"σ",
		5.670374419e-8,
		"W/(m²·K⁴)",
		units.DimensionNone,
		"Stefan-Boltzmann constant",
		"universal",
	)
	
	// Astronomical unit
	s.addConstant(
		"astronomical_unit",
		"au",
		1.495978707e11,
		"m",
		units.DimensionLength,
		"Astronomical unit (mean Earth-Sun distance)",
		"universal",
	)
	
	// Light year
	s.addConstant(
		"light_year",
		"ly",
		9.4607304725808e15,
		"m",
		units.DimensionLength,
		"Light year",
		"universal",
	)
	
	// Parsec
	s.addConstant(
		"parsec",
		"pc",
		3.0856775814913673e16,
		"m",
		units.DimensionLength,
		"Parsec",
		"universal",
	)
	
	// Solar mass
	s.addConstant(
		"solar_mass",
		"M_☉",
		1.98847e30,
		"kg",
		units.DimensionMass,
		"Solar mass",
		"universal",
	)
	
	// Earth mass
	s.addConstant(
		"earth_mass",
		"M_⊕",
		5.97217e24,
		"kg",
		units.DimensionMass,
		"Earth mass",
		"universal",
	)
	
	// Solar radius
	s.addConstant(
		"solar_radius",
		"R_☉",
		6.957e8,
		"m",
		units.DimensionLength,
		"Solar radius",
		"universal",
	)
	
	// Earth radius (mean)
	s.addConstant(
		"earth_radius",
		"R_⊕",
		6.371e6,
		"m",
		units.DimensionLength,
		"Mean Earth radius",
		"universal",
	)
	
	// Hubble constant (approximate, varies with measurement)
	s.addConstant(
		"hubble",
		"H_0",
		2.2e-18,
		"1/s",
		units.DimensionNone,
		"Hubble constant (≈70 km/s/Mpc)",
		"universal",
	)
}
