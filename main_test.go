package main

import (
	"testing"
)

func TestMapDose(t *testing.T) {
	testCases := []struct {
		input            string
		expectedDose     string
		expectedDoseUnit string
	}{
		{
			input:            "PRENDRE 1 COMPRIME",
			expectedDose:     "1",
			expectedDoseUnit: "comprimé",
		},
		{
			input:            "SELON LES DIRECTIVES DU MEDECIN",
			expectedDose:     "",
			expectedDoseUnit: "",
		},
		{
			input:            "2 VAPORISATIONS",
			expectedDose:     "2",
			expectedDoseUnit: "vaporisation",
		},
		{
			input:            "PRENDRE 2 INHALATIONS",
			expectedDose:     "2",
			expectedDoseUnit: "bouffée",
		},
		{
			input:            "PRENDRE 1 A 2 COMPRIMES",
			expectedDose:     "1-2",
			expectedDoseUnit: "comprimé",
		},
		{
			input:            "PRENDRE 1 - 2 COMPRIMES",
			expectedDose:     "1-2",
			expectedDoseUnit: "comprimé",
		},
		{
			input:            "PRENDRE 1-2 COMPRIMES",
			expectedDose:     "1-2",
			expectedDoseUnit: "comprimé",
		},
		{
			input:            "PRENDRE 17 GRAMMES",
			expectedDose:     "17",
			expectedDoseUnit: "g",
		},
		{
			input:            "PRENDRE 0.5 COMPRIME",
			expectedDose:     "0.5",
			expectedDoseUnit: "comprimé",
		},
		{
			input:            "PRENDRE 0.5 A 1 COMPRIME",
			expectedDose:     "0.5-1",
			expectedDoseUnit: "comprimé",
		},
		{
			input:            "PRENDRE 0.5 COMPRIME",
			expectedDose:     "0.5",
			expectedDoseUnit: "comprimé",
		},
		{
			input:            "PRENDRE 0.5 A 2 COMPRIME",
			expectedDose:     "0.5-2",
			expectedDoseUnit: "comprimé",
		},
		{
			input:            "PRENDRE 0,5 A 2 COMPRIME",
			expectedDose:     "0.5-2",
			expectedDoseUnit: "comprimé",
		},
		{
			input:            "PRENDRE 0,5 COMPRIME",
			expectedDose:     "0.5",
			expectedDoseUnit: "comprimé",
		},
		{
			input:            "1 COMPRIME FOIS PAR JOUR 1/2 HEURE AVANT COUCHER",
			expectedDose:     "1",
			expectedDoseUnit: "comprimé",
		},
	}

	for _, tc := range testCases {
		t.Run("TestMapDose", func(t *testing.T) {
			actualDose, actualDoseUnit := MapDose(tc.input)

			if actualDose != tc.expectedDose {
				t.Errorf("Dose\nI: %v\nE: %v\nA: %v", tc.input, tc.expectedDose, actualDose)
				return
			}

			if actualDoseUnit != tc.expectedDoseUnit {
				t.Errorf("DoseUnit\nI: %v\nE: %v\nA: %v", tc.input, tc.expectedDoseUnit, actualDoseUnit)
				return
			}
		})
	}
}

func TestRemoveFraction(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    "PRENDRE 1/2 COMPRIME",
			expected: "PRENDRE 0.5 COMPRIME",
		},
		{
			input:    "PRENDRE 1/4 COMPRIME",
			expected: "PRENDRE 0.25 COMPRIME",
		},
		{
			input:    "PRENDRE 3/4 COMPRIME",
			expected: "PRENDRE 0.75 COMPRIME",
		},
		{
			input:    "PRENDRE 1 1/2 COMPRIME",
			expected: "PRENDRE 1.5 COMPRIME",
		},
		{
			input:    "PRENDRE ½ COMPRIME",
			expected: "PRENDRE 0.5 COMPRIME",
		},
		{
			input:    "PRENDRE ¼ COMPRIME",
			expected: "PRENDRE 0.25 COMPRIME",
		},
		{
			input:    "PRENDRE ¾ COMPRIME",
			expected: "PRENDRE 0.75 COMPRIME",
		},
		{
			input:    "PRENDRE 1 1/4 COMPRIME",
			expected: "PRENDRE 1.25 COMPRIME",
		},
	}

	for _, tc := range testCases {
		t.Run("TestRemoveFraction", func(t *testing.T) {
			actual := RemoveFraction(tc.input)
			if actual != tc.expected {
				t.Errorf("I: %v\nE: %v\nA: %v", tc.input, tc.expected, actual)
				return
			}
		})
	}
}

func TestMapFrequency(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    "PRENDRE 1 COMPRIME 1 FOIS PAR JOUR",
			expected: "1 fois par jour",
		},
		{
			input:    "PRENDRE 1 COMPRIME 2 FOIS PAR JOUR",
			expected: "2 fois par jour",
		},
		{
			input:    "PRENDRE 1 A 2 COMPRIMES AUX 4 A 6 HEURES SI BESOIN (MAXIMUM 8 COMPRIMES PAR JOUR)",
			expected: "",
		},
		{
			input:    "PRENDRE 1 COMPRIME PAR JOUR",
			expected: "1 fois par jour",
		},
		{
			input:    "PRENDRE 1 COMPRIME PAR JOUR MAXIMUM 10 COMPRIME PAR JOUR",
			expected: "1 fois par jour",
		},
	}

	for _, tc := range testCases {
		t.Run("TestMapFrequency", func(t *testing.T) {
			_, actual := MapFrequency(tc.input)
			if actual != tc.expected {
				t.Errorf("I: %v\nE: %v\nA: %v", tc.input, tc.expected, actual)
				return
			}
		})
	}
}
