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
		{
			input:            "PRENDRE 2 COMPRIMES LE 1ER JOUR, PUIS 1 COMPRIME 1 FOIS PAR JOUR AUX 24 HEURES DU 2IÈME AU 5IÈME JOUR",
			expectedDose:     "",
			expectedDoseUnit: "",
		},
		{
			input:            "VAPORISEZ 2 FOIS DANS LES NARINES LE MATIN - REGULIEREMENT (ALLERGIE)",
			expectedDose:     "2",
			expectedDoseUnit: "vaporisation",
		},
		{
			input:            "1 GOUTTE DANS L'OEIL AFFECTE 4 FOIS PAR JOUR POUR 7 JOURS",
			expectedDose:     "1",
			expectedDoseUnit: "goutte",
		},
		{
			input:            "APPLIQUER 2G. SUR LES ZONES DOULOUREUSES",
			expectedDose:     "2",
			expectedDoseUnit: "g",
		},
		{
			input:            "PRENEZ LA CAPSULE EN MANGEANT - DOSE UNIQUE (INFECTION)",
			expectedDose:     "1",
			expectedDoseUnit: "capsule",
		},
		{
			input:            "VAPORISEZ DANS LA BOUCHE AUX 5 MINUTES SI BESOIN",
			expectedDose:     "1",
			expectedDoseUnit: "vaporisation",
		},
		{
			input:            "INHALER LE CONTENU D'UNE CAPSULE 1 FOIS PAR JOUR LE MATIN",
			expectedDose:     "1",
			expectedDoseUnit: "capsule",
		},
		{
			input:            "PRENEZ 2 COMPRIMES MAINTENANT ET 1 COMPRIME APRÈS CHAQUE SELLE LIQUIDE-MAX 8/JR (DIARRHÉE)",
			expectedDose:     "",
			expectedDoseUnit: "",
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
			input:    "PRENDRE 1 À 2 COMPRIMES AUX 4 A 6  HEURES SI BESOIN (MAXIMUM 8 COMPRIMES PAR JOUR)",
			expected: "q4h PRN",
		},
		{
			input:    "PRENDRE 1 COMPRIME PAR JOUR",
			expected: "1 fois par jour",
		},
		{
			input:    "PRENDRE 1 COMPRIME PAR JOUR MAXIMUM 10 COMPRIME PAR JOUR",
			expected: "1 fois par jour",
		},
		{
			input:    "PRENDRE 1 COMPRIME 1 FOIS PAR JOUR AVANT LE DEJEUNER",
			expected: "1 fois par jour avant le déjeuner",
		},
		{
			input:    "PRENEZ 1 COMPRIME 30 MINUTES AVANT LE COUCHER - AU BESOIN (INSOMNIE)",
			expected: "1 fois par jour au coucher PRN",
		},
		{
			input:    "PRENEZ 1 COMPRIME PAR JOUR SANS ARRET",
			expected: "1 fois par jour",
		},
		{
			input:    "PRENEZ 1 COMPRIME PAR JOUR AVEC LE DEJEUNER - RÉGULIÈREMENT (PRESSION)",
			expected: "1 fois par jour au déjeuner",
		},
		{
			input:    "PRENDRE 2 COMPRIMES LE 1ER JOUR, PUIS 1 COMPRIME 1 FOIS PAR JOUR AUX 24 HEURES DU 2IÈME AU 5IÈME JOUR",
			expected: "",
		},
		{
			input:    "VAPORISEZ 2 FOIS DANS LES NARINES LE MATIN - REGULIEREMENT (ALLERGIE)",
			expected: "1 fois par jour le matin",
		},
		{
			input:    "PRENDRE 1 COMPRIME 1 FOIS PAR JOUR 1/2 HEURE AVANT COUCHER",
			expected: "1 fois par jour au coucher",
		},
		{
			input:    "PRENEZ 1 COMPRIME AUX 4 HEURES - AU BESOIN (DOULEUR)",
			expected: "q4h PRN",
		},
		{
			input:    "PRENEZ 2 COMPRIMES AUX 6 HEURES - REGULIEREMENT (DOULEUR - FIEVRE)",
			expected: "q6h",
		},
		{
			input:    "PRENEZ 1 COMPRIME AUX 4 A 6 HEURES - AU BESOIN (DOULEUR)",
			expected: "q4h PRN",
		},
		{
			input:    "PRENDRE 1 COMPRIME PAR SEMAINE AVEC 120 ML D'EAU, LE MATIN, AU MOINS 30 MINUTES AVANT NOURRITURE OU AUTRE MEDICAMENT",
			expected: "1 fois par semaine",
		},
		{
			input:    "PRENDRE 1 CAPSULE UNE FOIS PAR JOUR, A LA MEME HEURE CHAQUE JOUR",
			expected: "1 fois par jour",
		},
		{
			input:    "TAKE 1 TABLET ONCE DAILY",
			expected: "1 fois par jour",
		},
		{
			input:    "TAKE 1 TABLET ONCE DAILY AT BEDTIME",
			expected: "1 fois par jour au coucher",
		},
		{
			input:    "TAKE 1 TABLET ONCE DAILY IN THE MORNING",
			expected: "1 fois par jour le matin",
		},
		{
			input:    "TAKE 1 TABLET DAILY WITH BREAKFAST - REGULARLY (PRESSURE)",
			expected: "1 fois par jour au déjeuner",
		},
		{
			input:    "TAKE 1 TABLET ONCE A WEEK",
			expected: "1 fois par semaine",
		},
		{
			input:    "1 VAPORISATION SOUS LA LANGUE AUX 5 MINUTES SI DOULEUR A LA POITRINE. MAX. 3 VAPORISATIONS SI BESOIN.",
			expected: "q5min PRN",
		},
		{
			input:    "1 COMPRIME TOUTES LES 4 HEURES SI DOULEURS",
			expected: "q4h PRN",
		},
		{
			input:    "COLLER UN TIMBRE, GARDER 24 HEURES, RETIRER ET CHANGER. POURSUIVRE PENDANT 6 SEMAINES ET PASSER A L'ETAPE 2",
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

func TestMapRoute(t *testing.T) {
	testCases := []struct {
		input    string
		doseUnit string
		expected string
	}{
		{
			input:    "INHALER LE CONTENU D'UNE CAPSULE 1 FOIS PAR JOUR LE MATIN",
			doseUnit: "capsule",
			expected: "inhalation",
		},
	}

	for _, tc := range testCases {
		t.Run("TestMapRoute", func(t *testing.T) {
			actual := MapRoute(tc.input, Dosage{DoseUnit: tc.doseUnit})
			if actual != tc.expected {
				t.Errorf("I: %v\nE: %v\nA: %v", tc.input, tc.expected, actual)
				return
			}
		})
	}
}
