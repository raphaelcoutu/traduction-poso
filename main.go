package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type Dosage struct {
	Text        string `json:"text"`
	Dose        string `json:"dose"`
	DoseUnit    string `json:"dose_unit"`
	FrequencyId int    `json:"frequency_id"`
	Frequency   string `json:"frequency"`
}

func main() {
	file, err := os.Open("in_sample.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	var dosages []Dosage
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		oLine := scanner.Text()

		dosage, err := MapAll(oLine)
		if err != nil {
			log.Fatal(err)
		}

		dosages = append(dosages, dosage)

		if err != nil {
			log.Fatal(err)
			return
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	PrintToJson(dosages)
	PrintToText(dosages)
}

func PrintToJson(dosages []Dosage) {
	jsonData, err := json.MarshalIndent(dosages, "", "  ")
	if err != nil {
		log.Fatal(err)
		return
	}

	err = os.WriteFile("out.json", jsonData, 0644)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func PrintToText(dosages []Dosage) {
	// write to text file
	outFile, err := os.Create("out.txt")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer outFile.Close()

	for _, dosage := range dosages {
		sDosage := fmt.Sprintf("%s, %s, %s, %s", dosage.Text, dosage.Dose, dosage.DoseUnit, dosage.Frequency)
		_, err = outFile.WriteString(sDosage + "\n")
		if err != nil {
			log.Fatal(err)
			return
		}
	}

}

func MapAll(line string) (Dosage, error) {
	dosage := Dosage{}

	dosage.Text = line

	line = strings.ToUpper(line)
	line, err := RemoveAccents(line)
	if err != nil {
		log.Fatal(err)
	}

	dosage.Dose, dosage.DoseUnit = MapDose(line)

	dosage.FrequencyId, dosage.Frequency = MapFrequency(line)

	return dosage, nil
}

func MapDose(line string) (string, string) {

	line = RemoveFraction(line)

	re := regexp.MustCompile(`((((\d(\.|,))?\d+)( A |\-| \- ))?((\d(\.|,))?\d+)) (COMPRIMES|COMPRIME|TABLETS|TABLET)`)
	if re.MatchString(line) {
		dose := strings.Replace(re.FindStringSubmatch(line)[1], "A", "-", -1)
		dose = strings.Replace(dose, " ", "", -1)
		dose = strings.Replace(dose, ",", ".", -1)
		return dose, "comprimé"
	}

	re = regexp.MustCompile(`(((\d+)( A |\-| \- ))?((\d(\.|,))?\d+)) (CAPSULES|CAPSULE)`)
	if re.MatchString(line) {
		dose := strings.Replace(re.FindStringSubmatch(line)[1], "A", "-", -1)
		dose = strings.Replace(dose, " ", "", -1)
		dose = strings.Replace(dose, ",", ".", -1)
		return dose, "capsule"
	}

	re = regexp.MustCompile(`([0-9]+) (VAPORISATIONS|VAPORISATION)`)
	if re.MatchString(line) {
		return re.FindStringSubmatch(line)[1], "vaporisation"
	}

	re = regexp.MustCompile(`([0-9]+) (INHALATIONS|INHALATION)`)
	if re.MatchString(line) {
		return re.FindStringSubmatch(line)[1], "bouffée"
	}

	re = regexp.MustCompile(`([0-9]+) (GRAMMES|GRAMME|G)`)
	if re.MatchString(line) {
		return re.FindStringSubmatch(line)[1], "g"
	}

	re = regexp.MustCompile(`([0-9]+)(GRAMMES|GRAMME|G|G. )`)
	if re.MatchString(line) {
		return re.FindStringSubmatch(line)[1], "g"
	}
	return "", ""
}

func MapFrequency(line string) (int, string) {
	isPrn := false
	withFood := false

	if strings.Contains(line, "PRN") || strings.Contains(line, "AU BESOIN") || strings.Contains(line, "SI BESOIN") || strings.Contains(line, "AS NEEDED") {
		isPrn = true
	}

	if strings.Contains(line, "EN MANGEANT") || strings.Contains(line, "AVEC NOURRITURE") {
		withFood = true
	}

	// # FOIS PAR JOUR
	re := regexp.MustCompile(`([0-9]+) FOIS PAR JOUR`)
	if re.MatchString(line) {
		freq := re.FindStringSubmatch(line)[1]

		if isPrn {
			// prn = true
			if freq == "1" {
				if regexp.MustCompile(`(AU|AVEC LE) DEJEUNER`).MatchString(line) {
					return 0, "1 fois par jour au déjeuner PRN"
				} else if regexp.MustCompile(`AVANT LE DEJEUNER`).MatchString(line) {
					return 0, "1 fois par jour avant le déjeuner PRN"
				} else if strings.Contains(line, "LE MATIN") {
					return 0, "1 fois par jour le matin PRN"
				} else if strings.Contains(line, "AU DINER") {
					return 0, "1 fois par jour au dîner PRN"
				} else if strings.Contains(line, "AU SOUPER") {
					return 0, "1 fois par jour au souper PRN"
				} else if strings.Contains(line, "AU COUCHER") {
					return 0, "1 fois par jour au coucher PRN"
				}

				return 0, "1 fois par jour PRN"
			} else if freq == "2" {
				return 0, "2 fois par jour PRN"
			} else if freq == "3" {
				return 0, "3 fois par jour PRN"
			} else if freq == "4" {
				return 0, "4 fois par jour PRN"
			}
		} else {
			// prn = false
			if freq == "1" {
				if regexp.MustCompile(`(AU|AVEC LE) DEJEUNER`).MatchString(line) {
					return 0, "1 fois par jour au déjeuner"
				} else if regexp.MustCompile(`AVANT LE DEJEUNER`).MatchString(line) {
					return 0, "1 fois par jour avant le déjeuner"
				} else if strings.Contains(line, "LE MATIN") {
					return 0, "1 fois par jour le matin"
				} else if strings.Contains(line, "AU DINER") {
					return 0, "1 fois par jour au dîner"
				} else if strings.Contains(line, "AU SOUPER") {
					return 0, "1 fois par jour au souper"
				} else if strings.Contains(line, "AU COUCHER") {
					return 0, "1 fois par jour au coucher"
				}

				return 0, "1 fois par jour"
			} else if freq == "2" {

				if strings.Contains(line, "DEJEUNER") && strings.Contains(line, "SOUPER") {
					return 0, "2 fois par jour au déjeuner et au souper"
				} else if strings.Contains(line, "MATIN") && strings.Contains(line, "SOIR") && withFood {
					return 0, "2 fois par jour au déjeuner et au souper "
				}

				return 0, "2 fois par jour"
			} else if freq == "3" {
				return 0, "3 fois par jour"
			} else if freq == "4" {
				return 0, "4 fois par jour"
			}
		}
	}

	// # FOIS PAR SEMAINE
	re = regexp.MustCompile(`([0-9]+) FOIS PAR SEMAINE`)
	if re.MatchString(line) {
		freq := re.FindStringSubmatch(line)[1]

		if isPrn {
			// prn = true
			if freq == "1" {
				return 0, "1 fois par semaine PRN"
			} else if freq == "2" {
				return 0, "2 fois par semaine PRN"
			} else if freq == "3" {
				return 0, "3 fois par semaine PRN"
			}
		} else {
			// prn = false
			if freq == "1" {
				return 0, "1 fois par semaine"
			} else if freq == "2" {
				return 0, "2 fois par semaine"
			} else if freq == "3" {
				return 0, "3 fois par semaine"
			}
		}
	}

	re = regexp.MustCompile(`(?i)^(?:[^M]+|M(?:AX(?:IMUM|\.))) [0-9]+ (COMPRIMES?) PAR JOUR`)
	if re.MatchString(line) {
		if isPrn {
			return 0, "1 fois par jour PRN"
		} else {
			return 0, "1 fois par jour"
		}
	}

	re = regexp.MustCompile(`UNE SEULE DOSE`)
	if re.MatchString(line) {
		if isPrn {
			return 0, "1 fois PRN"
		} else {
			return 0, "1 fois"
		}
	}

	return 0, ""
}

func RemoveAccents(text string) (string, error) {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, err := transform.String(t, text)
	if err != nil {
		return "", err
	}
	return result, nil
}

func RemoveFraction(text string) string {
	text = strings.Replace(text, "½", "1/2", -1)
	text = strings.Replace(text, "¼", "1/4", -1)
	text = strings.Replace(text, "¾", "3/4", -1)

	if strings.Contains(text, "1 1/2") {
		text = strings.Replace(text, "1 1/2", "1.5", -1)
	} else if strings.Contains(text, "1 1/4") {
		text = strings.Replace(text, "1 1/4", "1.25", -1)
	} else if strings.Contains(text, "1/2") {
		text = strings.Replace(text, "1/2", "0.5", -1)
	} else if strings.Contains(text, "1/4") {
		text = strings.Replace(text, "1/4", "0.25", -1)
	} else if strings.Contains(text, "3/4") {
		text = strings.Replace(text, "3/4", "0.75", -1)
	}

	return text
}
