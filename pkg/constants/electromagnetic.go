package constants

import "github.com/andrewneudegg/calc/pkg/units"

// initElectromagnetic initializes electromagnetic constants.
// Values are from CODATA 2018 recommended values.
func (s *System) initElectromagnetic() {
	// Vacuum permeability (magnetic constant)
	s.addConstant(
		"vacuum_permeability",
		"μ_0",
		1.25663706212e-6,
		"N/A²",
		units.DimensionNone,
		"Vacuum permeability (magnetic constant)",
		"electromagnetic",
	)
	
	// Vacuum permittivity (electric constant)
	s.addConstant(
		"vacuum_permittivity",
		"ε_0",
		8.8541878128e-12,
		"F/m",
		units.DimensionNone,
		"Vacuum permittivity (electric constant)",
		"electromagnetic",
	)
	
	// Coulomb constant (ke)
	s.addConstant(
		"coulomb",
		"k_e",
		8.9875517923e9,
		"N·m²/C²",
		units.DimensionNone,
		"Coulomb constant (1/(4πε_0))",
		"electromagnetic",
	)
	
	// Impedance of free space
	s.addConstant(
		"impedance_vacuum",
		"Z_0",
		376.730313668,
		"Ω",
		units.DimensionNone,
		"Characteristic impedance of vacuum",
		"electromagnetic",
	)
	
	// Bohr magneton
	s.addConstant(
		"bohr_magneton",
		"μ_B",
		9.2740100783e-24,
		"J/T",
		units.DimensionNone,
		"Bohr magneton",
		"electromagnetic",
	)
	
	// Nuclear magneton
	s.addConstant(
		"nuclear_magneton",
		"μ_N",
		5.0507837461e-27,
		"J/T",
		units.DimensionNone,
		"Nuclear magneton",
		"electromagnetic",
	)
}
