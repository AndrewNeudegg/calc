package constants

import "github.com/andrewneudegg/calc/pkg/units"

// initFundamental initializes fundamental physical constants.
// Values are from CODATA 2018 recommended values.
func (s *System) initFundamental() {
	// Speed of light in vacuum
	s.addConstant(
		"speed_of_light",
		"c",
		299792458.0,
		"m/s",
		units.DimensionSpeed,
		"Speed of light in vacuum",
		"fundamental",
	)
	
	// Planck constant
	s.addConstant(
		"planck",
		"h",
		6.62607015e-34,
		"J·s",
		units.DimensionNone, // Energy × Time (composite)
		"Planck constant",
		"fundamental",
	)
	
	// Reduced Planck constant (h-bar)
	s.addConstant(
		"planck_reduced",
		"ℏ",
		1.054571817e-34,
		"J·s",
		units.DimensionNone,
		"Reduced Planck constant (ℏ = h/2π)",
		"fundamental",
	)
	
	// Elementary charge
	s.addConstant(
		"elementary_charge",
		"e",
		1.602176634e-19,
		"C",
		units.DimensionNone, // Coulombs (charge)
		"Elementary charge",
		"fundamental",
	)
	
	// Electron mass
	s.addConstant(
		"electron_mass",
		"m_e",
		9.1093837015e-31,
		"kg",
		units.DimensionMass,
		"Electron rest mass",
		"fundamental",
	)
	
	// Proton mass
	s.addConstant(
		"proton_mass",
		"m_p",
		1.67262192369e-27,
		"kg",
		units.DimensionMass,
		"Proton rest mass",
		"fundamental",
	)
	
	// Neutron mass
	s.addConstant(
		"neutron_mass",
		"m_n",
		1.67492749804e-27,
		"kg",
		units.DimensionMass,
		"Neutron rest mass",
		"fundamental",
	)
	
	// Fine-structure constant (dimensionless)
	s.addConstant(
		"fine_structure",
		"α",
		7.2973525693e-3,
		"",
		units.DimensionNone,
		"Fine-structure constant (α ≈ 1/137)",
		"fundamental",
	)
	
	// Rydberg constant
	s.addConstant(
		"rydberg",
		"R_∞",
		10973731.568160,
		"1/m",
		units.DimensionNone,
		"Rydberg constant",
		"fundamental",
	)
	
	// Avogadro constant
	s.addConstant(
		"avogadro",
		"N_A",
		6.02214076e23,
		"1/mol",
		units.DimensionNone,
		"Avogadro constant",
		"fundamental",
	)
	
	// Boltzmann constant
	s.addConstant(
		"boltzmann",
		"k_B",
		1.380649e-23,
		"J/K",
		units.DimensionNone,
		"Boltzmann constant",
		"fundamental",
	)
	
	// Gas constant
	s.addConstant(
		"gas_constant",
		"R",
		8.314462618,
		"J/(mol·K)",
		units.DimensionNone,
		"Molar gas constant (R = N_A × k_B)",
		"fundamental",
	)
}
