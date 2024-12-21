package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type Dosage struct {
	Id          int    `json:"id"`
	Text        string `json:"text"`
	Dose        string `json:"dose"`
	DoseUnit    string `json:"dose_unit"`
	Route       string `json:"route"`
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
	i := 1
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		oLine := scanner.Text()

		dosage, err := MapAll(oLine)
		if err != nil {
			log.Fatal(err)
		}

		dosage.Id = i

		dosages = append(dosages, dosage)

		i++
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
	dosage.Route = MapRoute(line, dosage)

	dosage.FrequencyId, dosage.Frequency = MapFrequency(line)

	return dosage, nil
}

func MapDose(line string) (string, string) {

	if isComplexDosage(line) {
		return "", ""
	}

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

	if regexp.MustCompile(`(LA|UNE) CAPSULE`).MatchString(line) {
		return "1", "capsule"
	}

	re = regexp.MustCompile(`([0-9]+) (VAPORISATIONS|VAPORISATION)`)
	if re.MatchString(line) {
		return re.FindStringSubmatch(line)[1], "vaporisation"
	}

	re = regexp.MustCompile(`VAPORISE(?:R|Z) ([0-9]+) FOIS`)
	if re.MatchString(line) {
		return re.FindStringSubmatch(line)[1], "vaporisation"
	}

	re = regexp.MustCompile(`(VAPORISER|VAPORISEZ)`)
	if re.MatchString(line) {
		return "1", "vaporisation"
	}

	re = regexp.MustCompile(`([0-9]+) (INHALATIONS|INHALATION)`)
	if re.MatchString(line) {
		return re.FindStringSubmatch(line)[1], "bouffée"
	}

	re = regexp.MustCompile(`([0-9]+) (GOUTTES?|G )`)
	if re.MatchString(line) {
		return re.FindStringSubmatch(line)[1], "goutte"
	}

	re = regexp.MustCompile(`([0-9]+) (GRAMMES?|G )`)
	if re.MatchString(line) {
		return re.FindStringSubmatch(line)[1], "g"
	}

	re = regexp.MustCompile(`([0-9]+)(GRAMMES?|G|G\. )`)
	if re.MatchString(line) {
		return re.FindStringSubmatch(line)[1], "g"
	}

	re = regexp.MustCompile(`([0-9]+|UN) (TIMBRES?|PATCHS?)`)
	if re.MatchString(line) {
		if re.FindStringSubmatch(line)[1] == "UN" {
			return "1", "timbre"
		}
		return re.FindStringSubmatch(line)[1], "timbre"
	}

	return "", ""
}

func MapRoute(line string, dosage Dosage) string {

	if strings.Contains(line, "CHAQUE NARINE") || strings.Contains(line, "DANS LES NARINES") {
		return "nasale"
	}

	if regexp.MustCompile(`INTRA(-|\s)?MUSCULAIRE`).MatchString(line) {
		return "intramusculaire"
	}

	if regexp.MustCompile(`SOUS(-|\s)?CUTANEE`).MatchString(line) {
		return "sous-cutané"
	}

	if strings.Contains(line, "OEIL") || strings.Contains(line, "ŒIL") || strings.Contains(line, "YEUX") {
		if strings.Contains(line, "GAUCHE") {
			return "oeil gauche"
		} else if strings.Contains(line, "DROIT") {
			return "oeil droit"
		} else if strings.Contains(line, "YEUX") {
			return "dans les 2 yeux"
		}
		return "oculaire"
	}

	if regexp.MustCompile(`OREILLES?`).MatchString(line) {
		if strings.Contains(line, "GAUCHE") {
			return "oreille gauche"
		} else if strings.Contains(line, "DROIT") {
			return "oreille droit"
		} else if strings.Contains(line, "OREILLES") {
			return "dans les 2 oreilles"
		}
		return "otique"
	}

	if regexp.MustCompile(`APPLIQUE(R|Z)|APPLICATION LOCALE`).MatchString(line) {
		return "topique"
	}

	if strings.Contains(line, "BOIRE") || regexp.MustCompile(`DISSOU.*\sVERRE D'EAU`).MatchString(line) {
		return "oral"
	}

	if strings.Contains(line, "SOUS LA LANGUE") {
		return "sublingual"
	}

	if regexp.MustCompile(`INHALE(R|Z) (LE CONTENU D'UNE )?CAPSULE`).MatchString(line) {
		return "inhalation"
	}

	if slices.Contains([]string{"comprimé", "capsule"}, dosage.DoseUnit) {
		return "oral"
	} else if dosage.DoseUnit == "bouffée" {
		return "inhalation"
	} else if dosage.DoseUnit == "timbre" {
		return "topique"
	}

	return ""
}

func MapFrequency(line string) (int, string) {
	isPrn := false
	withFood := false

	if isComplexDosage(line) {
		return 0, ""
	}

	if strings.Contains(line, "PRN") || strings.Contains(line, "AU BESOIN") || strings.Contains(line, "SI BESOIN") || strings.Contains(line, "AS NEEDED") || regexp.MustCompile(`SI (DOULEURS?)`).MatchString(line) {
		isPrn = true
	}

	if strings.Contains(line, "EN MANGEANT") || strings.Contains(line, "AVEC NOURRITURE") {
		withFood = true
	}

	// # FOIS PAR JOUR (FR)
	re := regexp.MustCompile(`([0-9]+|UNE) FOIS PAR JOUR`)
	if re.MatchString(line) {
		freq := re.FindStringSubmatch(line)[1]

		if isPrn {
			// prn = true
			if freq == "1" || freq == "UNE" {
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
				} else if regexp.MustCompile(`AVANT (LE )?COUCHER`).MatchString(line) {
					return 0, "1 fois par jour au coucher PRN"
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
			if freq == "1" || freq == "UNE" {
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
				} else if regexp.MustCompile(`AVANT (LE )?COUCHER`).MatchString(line) {
					return 0, "1 fois par jour au coucher"
				} else if strings.Contains(line, "AU COUCHER") {
					return 0, "1 fois par jour au coucher"
				}

				return 0, "1 fois par jour"
			} else if freq == "2" {

				if strings.Contains(line, "DEJEUNER") && strings.Contains(line, "SOUPER") {
					return 0, "2 fois par jour au déjeuner et au souper"
				} else if strings.Contains(line, "MATIN") && strings.Contains(line, "SOIR") && withFood {
					return 0, "2 fois par jour au déjeuner et au souper"
				}

				return 0, "2 fois par jour"
			} else if freq == "3" {
				return 0, "3 fois par jour"
			} else if freq == "4" {
				return 0, "4 fois par jour"
			}
		}
	}

	// # FOIS PAR JOUR (EN)
	if strings.Contains(line, "DAILY") {
		if strings.Contains(line, "TWICE") {
			//
		} else if strings.Contains(line, "THREE TIMES") {
			//
		} else {
			if isPrn {
				return 0, "1 fois par jour PRN"
			} else {
				if strings.Contains(line, "BREAKFAST") {
					return 0, "1 fois par jour au déjeuner"
				} else if strings.Contains(line, "IN THE MORNING") {
					return 0, "1 fois par jour le matin"
				} else if strings.Contains(line, "AT BEDTIME") {
					return 0, "1 fois par jour au coucher"
				}
				return 0, "1 fois par jour"
			}
		}
	}

	// AUX # HEURES
	re = regexp.MustCompile(`(?:AUX|TOU(?:TE)?S LES) ([0-9]+)( A [0-9]+)?\s+HEURES`)
	if re.MatchString(line) {
		freq := re.FindStringSubmatch(line)[1]

		if isPrn {
			// prn = true
			if freq == "3" {
				return 0, "q3h PRN"
			} else if freq == "4" {
				return 0, "q4h PRN"
			} else if freq == "6" {
				return 0, "q6h PRN"
			} else if freq == "8" {
				return 0, "q8h PRN"
			} else if freq == "12" {
				return 0, "q8h PRN"
			} else if freq == "24" {
				return 0, "q8h PRN"
			}
		} else {
			// prn = false
			if freq == "3" {
				return 0, "q3h"
			} else if freq == "4" {
				return 0, "q4h"
			} else if freq == "6" {
				return 0, "q6h"
			} else if freq == "8" {
				return 0, "q8h"
			} else if freq == "12" {
				return 0, "q8h"
			} else if freq == "24" {
				return 0, "q8h"
			}
		}
	}

	// AUX # HEURES
	re = regexp.MustCompile(`AUX ([0-9]+)( A [0-9]+)?\s+MIN(UTES|UTE)?`)
	if re.MatchString(line) {
		freq := re.FindStringSubmatch(line)[1]

		if isPrn {
			// prn = true
			if freq == "5" {
				return 0, "q5min PRN"
			} else if freq == "10" {
				return 0, "q10min PRN"
			} else if freq == "15" {
				return 0, "q15min PRN"
			} else if freq == "30" {
				return 0, "q30min PRN"
			}
		} else {
			// prn = false
			if freq == "5" {
				return 0, "q5min"
			} else if freq == "10" {
				return 0, "q10min"
			} else if freq == "15" {
				return 0, "q15min"
			} else if freq == "30" {
				return 0, "q30min"
			}
		}
	}

	// # FOIS PAR SEMAINE
	re = regexp.MustCompile(`([0-9]+) FOIS PAR SEMAINE|(ONCE|TWICE) A WEEK`)
	if re.MatchString(line) {
		freq := re.FindStringSubmatch(line)[1]
		freqEn := re.FindStringSubmatch(line)[2]

		if isPrn {
			// prn = true
			if freq == "1" || freqEn == "ONCE" {
				return 0, "1 fois par semaine PRN"
			} else if freq == "2" {
				return 0, "2 fois par semaine PRN"
			} else if freq == "3" {
				return 0, "3 fois par semaine PRN"
			}
		} else {
			// prn = false
			if freq == "1" || freqEn == "ONCE" {
				return 0, "1 fois par semaine"
			} else if freq == "2" {
				return 0, "2 fois par semaine"
			} else if freq == "3" {
				return 0, "3 fois par semaine"
			}
		}
	}

	// PAR JOUR (mais pas MAXIMUM # COMPRIMES PAR JOUR)
	re = regexp.MustCompile(`[0-9]+ (COMPRIMES?|CAPSULES?) PAR JOUR`)
	matches := re.FindAllString(line, -1)

	var filteredMatches []string
	for _, match := range matches {
		if !strings.Contains(line[:strings.Index(line, match)], "MAXIMUM") {
			filteredMatches = append(filteredMatches, match)
		}
	}
	if filteredMatches != nil {
		if isPrn {
			return 0, "1 fois par jour PRN"
		} else {
			if strings.Contains(line, "30 MINUTES AVANT LE COUCHER") {
				return 0, "1 fois par jour au coucher PRN"
			} else if regexp.MustCompile(`(AU|AVEC LE) DEJEUNER`).MatchString(line) {
				return 0, "1 fois par jour au déjeuner"
			}
		}

		return 0, "1 fois par jour"
	}

	// PAR SEMAINE (mais pas MAXIMUM # COMPRIMES PAR SEMAINE)
	re = regexp.MustCompile(`[0-9]+ (COMPRIMES?|CAPSULES?|TIMBRES?) PAR SEMAINE`)
	matches = re.FindAllString(line, -1)

	filteredMatches = nil
	for _, match := range matches {
		if !strings.Contains(line[:strings.Index(line, match)], "MAXIMUM") {
			filteredMatches = append(filteredMatches, match)
		}
	}
	if filteredMatches != nil {
		return 0, "1 fois par semaine"
	}

	if strings.Contains(line, "LE MATIN") {
		if isPrn {
			return 0, "1 fois par jour le matin PRN"
		} else {
			return 0, "1 fois par jour le matin"
		}
	}

	if regexp.MustCompile(`(30 MINUTES|1/2 HEURE) AVANT LE COUCHER`).MatchString(line) {
		if isPrn {
			return 0, "1 fois par jour au coucher PRN"
		} else {
			return 0, "1 fois par jour au coucher"
		}
	}

	re = regexp.MustCompile(`UNE SEULE DOSE|DOSE UNIQUE|UNE SEULE PRISE|IMMEDIATEMENT`)
	if re.MatchString(line) {
		if isPrn {
			return 0, "1 fois PRN"
		} else {
			return 0, "1 fois"
		}
	}

	if regexp.MustCompile(`GARDER 24 HEURES,? ?RETIRER .* CHANGER`).MatchString(line) {
		return 0, "1 fois par jour"
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

func isComplexDosage(line string) bool {
	if regexp.MustCompile(`PUIS ([0-9]+) (?:COMPRIMES?|CAPSULES?)`).MatchString(line) {
		return true
	}

	if regexp.MustCompile(`MAINTENANT.*CHAQUE SELLE`).MatchString(line) {
		return true
	}

	return false
}
