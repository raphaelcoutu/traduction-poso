package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	file, err := os.Open("in.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	outFile, err := os.Create("out.txt")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer outFile.Close()

	writer := bufio.NewWriter(outFile)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		line = Transform(line)

		_, err := writer.WriteString(line + "\n")
		if err != nil {
			log.Fatal(err)
			return
		}

	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	err = writer.Flush()
	if err != nil {
		log.Fatal(err)
		return
	}
}

func Transform(line string) string {
	line = strings.Replace(line, "TAKE", "PRENDRE", -1)
	line = strings.Replace(line, "ADMINISTER", "ADMINISTRER", -1)
	line = strings.Replace(line, "APPLY", "APPLIQUER", -1)
	line = strings.Replace(line, "INJECT", "INJECTER", -1)
	line = strings.Replace(line, "INHALE", "INHALER", -1)
	line = strings.Replace(line, "PLACE RECTALLY", "INSÉRER DANS LE RECTUM", -1)
	line = strings.Replace(line, "PLACE 1 APPLICATOR VAGINALLY", "INSÉRER 1 APPLICATEUR DANS LE VAGIN", -1)
	line = strings.Replace(line, "PLACE", "APPLIQUER", -1)
	line = strings.Replace(line, "INSTILL", "INSTILLER", -1)
	line = strings.Replace(line, "CHEW", "MÂCHER", -1)

	re := regexp.MustCompile(`^SPRAY `)
	line = re.ReplaceAllString(line, "VAPORISER ")

	line = strings.Replace(line, "ONE-HALF", "1/2", -1)

	line = strings.Replace(line, "TABLET", "COMPRIMÉ", -1)
	line = strings.Replace(line, "TAB ", "COMPRIMÉ ", -1)
	line = strings.Replace(line, "DROP", "GOUTTE", -1)
	line = strings.Replace(line, "SPRAY", "VAPORISATION", -1)
	line = strings.Replace(line, "UNIT", "UNITÉ", -1)
	line = strings.Replace(line, "PUFFS", "BOUFFÉES", -1)
	line = strings.Replace(line, "PUFF(S)", "BOUFFÉE(S)", -1)
	line = strings.Replace(line, "PATCH ", "TIMBRE ", -1)

	re = regexp.MustCompile(`\(([\d.]+) (MG|G) TOTAL\)`)
	line = re.ReplaceAllStringFunc(line, replaceDecimalDoseWithComma)

	re = regexp.MustCompile(`\((\d{1,3}(?:,\d{3})*) MG TOTAL\)`)
	line = re.ReplaceAllStringFunc(line, removeCommaInThousand)

	line = strings.Replace(line, "BY MOUTH", "PAR LA BOUCHE", -1)
	line = strings.Replace(line, "ORALLY", "PAR LA BOUCHE", -1)
	line = strings.Replace(line, "INTO AFFECTED EAR(S)", "DANS LES OREILLES AFFECTÉES", -1)
	line = strings.Replace(line, "INTO EACH NOSTRIL", "DANS CHAQUE NARINE", -1)
	line = strings.Replace(line, "IN BOTH NOSTRILS", "DANS CHAQUE NARINE", -1)
	line = strings.Replace(line, "EACH NARE ROUTE", "DANS CHAQUE NARINE", -1)
	line = strings.Replace(line, "INTO NOSE", "DANS LE NEZ", -1)
	line = strings.Replace(line, "BY NASAL ROUTE", "PAR VOIE NASALE", -1)
	line = strings.Replace(line, "BY NEBULIZATION", "EN NÉBULISATION", -1)
	line = strings.Replace(line, "BY NEBULIZER ROUTE", "EN NÉBULISATION", -1)
	line = strings.Replace(line, "INTO BOTH EYES", "DANS LES 2 YEUX", -1)
	line = strings.Replace(line, "UNDER THE SKIN", "SOUS LA PEAU", -1)
	line = strings.Replace(line, "INTO THE SKIN", "SOUS LA PEAU", -1)
	line = strings.Replace(line, "ONTO THE SKIN", "SUR LA PEAU", -1)
	line = strings.Replace(line, "TO SKIN", "SUR LA PEAU", -1)
	line = strings.Replace(line, "TOPICALLY", "LOCALEMENT", -1)
	line = strings.Replace(line, "TO AFFECTED AREA", "SUR LA ZONE AFFECTÉE", -1)
	line = strings.Replace(line, "INTO THE LEFT EYE", "DANS L'OEIL GAUCHE", -1)
	line = strings.Replace(line, "INTO THE RIGHT EYE", "DANS L'OEIL DROIT", -1)
	line = strings.Replace(line, "IN AFFECTED EYE(S)", "DANS LES YEUX AFFECTÉS", -1)
	line = strings.Replace(line, "TO EYE", "DANS L'OEIL", -1)
	line = strings.Replace(line, "INTRAMUSCULARLY", "PAR VOIE INTRAMUSCULAIRE", -1)

	re = regexp.MustCompile(`MAX DAILY AMOUNT: (\d+) (MG|ML)`)
	line = re.ReplaceAllString(line, `(DOSE MAX PAR JOUR: $1 MG)`)

	re = regexp.MustCompile(`MAX DAILY AMOUNT: (\d+)\.(\d+) (MG|ML)`)
	line = re.ReplaceAllString(line, `(DOSE MAX PAR JOUR: $1,$2 MG)`)

	re = regexp.MustCompile(`MAX DAILY AMOUNT: (\d+),(\d+) (MG|ML)`)
	line = re.ReplaceAllString(line, `(DOSE MAX PAR JOUR: $1$2 MG)`)

	re = regexp.MustCompile(`(ONCE|ONE TIME) (DAILY|A DAY)`)
	line = re.ReplaceAllString(line, `1 FOIS PAR JOUR`)

	re = regexp.MustCompile(`(2|2 \(TWO\)|TWO) TIMES (DAILY|A DAY)`)
	line = re.ReplaceAllString(line, `2 FOIS PAR JOUR`)

	re = regexp.MustCompile(`TWICE (DAILY|A DAY)`)
	line = re.ReplaceAllString(line, `2 FOIS PAR JOUR`)

	re = regexp.MustCompile(`(3|3 \(THREE\)|THREE) TIMES (DAILY|A DAY)`)
	line = re.ReplaceAllString(line, `3 FOIS PAR JOUR`)

	re = regexp.MustCompile(`(4|4 \(FOUR\)|FOUR) TIMES (DAILY|A DAY)`)
	line = re.ReplaceAllString(line, `4 FOIS PAR JOUR`)

	re = regexp.MustCompile(`EVERY ([2-9]) DAYS`)
	line = re.ReplaceAllString(line, `À TOUS LES $1 JOURS`)

	line = strings.Replace(line, "BEFORE BREAKFAST", "AVANT LE DÉJEUNER", -1)
	line = strings.Replace(line, "BEFORE MEALS AND NIGHTLY", "AVANT LES REPAS ET AU COUCHER", -1)
	line = strings.Replace(line, "WITH BREAKFAST", "AVEC LE DÉJEUNER", -1)

	line = strings.Replace(line, "Q DAY", "1 FOIS PAR JOUR", -1)
	line = strings.Replace(line, " BID ", " 2 FOIS PAR JOUR ", -1)

	line = strings.Replace(line, "TWICE DAILY", "2 FOIS PAR JOUR", -1)
	line = strings.Replace(line, "THREE TIMES DAILY", "3 FOIS PAR JOUR", -1)
	line = strings.Replace(line, "NIGHTLY", "1 FOIS PAR JOUR AU COUCHER", -1)
	line = strings.Replace(line, "EVERY DAY", "1 FOIS PAR JOUR", -1)
	line = strings.Replace(line, "EVERY 4 (FOUR) TO 6 (SIX) HOURS", "AUX 4 À 6 HEURES", -1)
	line = strings.Replace(line, "EVERY 3 (THREE) HOURS", "AUX 3 HEURES", -1)
	line = strings.Replace(line, "EVERY 3 HOURS", "AUX 3 HEURES", -1)
	line = strings.Replace(line, "EVERY 4 (FOUR) HOURS", "AUX 4 HEURES", -1)
	line = strings.Replace(line, "EVERY 4 HOURS", "AUX 4 HEURES", -1)
	line = strings.Replace(line, "EVERY 6 (SIX) HOURS", "AUX 6 HEURES", -1)
	line = strings.Replace(line, "EVERY 6 HOURS", "AUX 6 HEURES", -1)
	line = strings.Replace(line, "EVERY 8 (EIGHT) HOURS", "AUX 8 HEURES", -1)
	line = strings.Replace(line, "EVERY 8 HOURS", "AUX 8 HEURES", -1)

	re = regexp.MustCompile(`EVERY (12|12 \(TWELVE\)|TWELVE) HOURS`)
	line = re.ReplaceAllString(line, `AUX 12 HEURES`)

	line = strings.Replace(line, "ONCE A WEEK", "1 FOIS PAR SEMAINE", -1)
	line = strings.Replace(line, "WEEKLY", "1 FOIS PAR SEMAINE", -1)

	line = strings.Replace(line, "EVERY MORNING", "CHAQUE MATIN", -1)
	line = strings.Replace(line, "EVERY EVENING", "CHAQUE SOIR", -1)
	line = strings.Replace(line, "EVERY NIGHT", "CHAQUE SOIR", -1)
	line = strings.Replace(line, "AT BEDTIME", "AU COUCHER", -1)
	line = strings.Replace(line, "DAILY", "1 FOIS PAR JOUR", -1)
	line = strings.Replace(line, "EVERY OTHER DAY", "AUX 2 JOURS", -1)
	line = strings.Replace(line, "ONCE FOR 1 DOSE", "POUR 1 DOSE", -1)

	line = strings.Replace(line, "AS NEEDED", "AU BESOIN", -1)
	line = strings.Replace(line, "IF NEEDED", "AU BESOIN", -1)

	re = regexp.MustCompile(`(FOR UP TO|FOR) (\d+) DAYS`)
	line = re.ReplaceAllString(line, `POUR $2 JOURS`)

	re = regexp.MustCompile(`X (\d+) DAYS`)
	line = re.ReplaceAllString(line, `X $1 JOURS`)

	line = strings.Replace(line, "FOR MILD PAIN (1-3)", "POUR DOULEUR LÉGÈRE", -1)
	line = strings.Replace(line, "FOR MODERATE PAIN (PAIN SCALE 4-7)", "POUR DOULEUR MODÉRÉE", -1)

	line = strings.Replace(line, "FOR DEPRESSION", "POUR LA DÉPRESSION", -1)
	line = strings.Replace(line, "FOR COUGH", "POUR LA TOUX", -1)
	line = strings.Replace(line, "FOR FEVER", "POUR LA FIÈVRE", -1)
	line = strings.Replace(line, "FOR SLEEP", "POUR L'INSOMNIE", -1)
	line = strings.Replace(line, "COUGH", "TOUX", -1)
	line = strings.Replace(line, "ITCHING", "PRURIT", -1)
	line = strings.Replace(line, "FOR NAUSEA OR VOMITING", "POUR LES NAUSÉES/VOMISSEMENTS", -1)
	line = strings.Replace(line, "TO THE RASH", "SUR LES ROUGEURS", -1)
	line = strings.Replace(line, "TO RASH", "SUR LES ROUGEURS", -1)
	line = strings.Replace(line, "FOR WHEEZING", "SI RESPIRATION SIFFLANTE", -1)
	line = strings.Replace(line, "ANXIETY", "ANXIÉTÉ", -1)
	line = strings.Replace(line, "HEADACHES", "MAUX DE TÊTE", -1)
	line = strings.Replace(line, "MUSCLE SPASMS", "SPASMES MUSCULAIRE", -1)
	line = strings.Replace(line, "CRAMPING", "CRAMPES", -1)

	line = strings.Replace(line, "PAIN", "DOULEUR", -1)
	line = strings.Replace(line, "FEVER", "FIÈVRE", -1)
	line = strings.Replace(line, "SHORTNESS OF BREATH", "DYSPNÉE", -1)
	line = strings.Replace(line, "SHORTNESS OF AIR", "DYSPNÉE", -1)
	line = strings.Replace(line, "NAUSEA", "NAUSÉES", -1)
	line = strings.Replace(line, "NAUSE/VOMITING", "NAUSÉES/VOMISSEMENTS", -1)

	line = strings.Replace(line, "AVOID GRAPEFRUIT PRODUCTS", "ÉVITER LE PAMPLEMOUSSE", -1)

	re = regexp.MustCompile(`^ONE\s`)
	line = re.ReplaceAllString(line, `UN `)

	re = regexp.MustCompile(`\sONE\s`)
	line = re.ReplaceAllString(line, ` UN `)

	re = regexp.MustCompile(`\sTWO\s`)
	line = re.ReplaceAllString(line, ` DEUX `)

	re = regexp.MustCompile(`^TWO\s`)
	line = re.ReplaceAllString(line, `DEUX `)

	line = strings.Replace(line, "TWO", "DEUX", -1)
	line = strings.Replace(line, "THREE", "TROIS", -1)
	line = strings.Replace(line, "ONCE", "1 FOIS", -1)
	line = strings.Replace(line, "TWICE", "2 FOIS", -1)
	line = strings.Replace(line, "WITH MEALS", "AVEC LES REPAS", -1)
	line = strings.Replace(line, "WITH FOOD", "AVEC NOURRITURE", -1)

	line = strings.Replace(line, "MEDICINE", "MÉDICAMENT", -1)

	line = strings.Replace(line, " MLS ", " ML ", -1)
	line = strings.Replace(line, " OR ", " OU ", -1)
	line = strings.Replace(line, " FOR ", " POUR ", -1)
	line = strings.Replace(line, " TO ", " À ", -1)
	line = strings.Replace(line, " PO ", " PAR LA BOUCHE ", -1)

	// Féminin
	line = strings.Replace(line, "UN GOUTTE", "UNE GOUTTE", -1)
	line = strings.Replace(line, "UN BOUFFÉE", "UNE BOUFFÉE", -1)
	line = strings.Replace(line, "UN CAPSULE", "UNE CAPSULE", -1)

	// On enlève les séparateurs de miliers
	re = regexp.MustCompile(`(\d{1,3})(,\d{3})*`)
	line = re.ReplaceAllStringFunc(line, func(s string) string {
		return strings.ReplaceAll(s, ",", "")
	})

	// On remplace les points par des virgules (ex: 12.5 mg → 12,5 mg)
	re = regexp.MustCompile(`([0-9])\.([0-9])`)
	line = re.ReplaceAllString(line, `$1,$2`)

	line = strings.Replace(line, "((", "(", -1)
	line = strings.Replace(line, "))", ")", -1)

	return line
}

func replaceDecimalDoseWithComma(match string) string {
	parts := strings.Split(match, " ")
	number := strings.Replace(parts[0][1:], ".", ",", 1) //remove the opening bracket and replace the first dot with comma
	return "(" + number + " " + parts[1] + ")"
}

func removeCommaInThousand(match string) string {
	parts := strings.Split(match, " ")
	number := strings.Replace(parts[0][1:], ",", "", -1) // Remove the opening bracket and all commas
	return "(" + number + " MG)"
}
