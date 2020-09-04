package base

import (
	"fmt"
	"io"
	"strconv"
)

// AddressType is used to determine the type of an address.
//
// See: https://www.hl7.org/fhir/valueset-address-type.html
type AddressType string

// known address types
const (
	// AddressTypePostal ...
	AddressTypePostal AddressType = "postal"
	// AddressTypePhysical ...
	AddressTypePhysical AddressType = "physical"
	// AddressTypeBoth ...
	AddressTypeBoth AddressType = "both"
)

// AllAddressType is a list of all known address types
var AllAddressType = []AddressType{
	AddressTypePostal,
	AddressTypePhysical,
	AddressTypeBoth,
}

// IsValid checks that the address type is valid
func (e AddressType) IsValid() bool {
	switch e {
	case AddressTypePostal, AddressTypePhysical, AddressTypeBoth:
		return true
	}
	return false
}

// String renders the address type as a string
func (e AddressType) String() string {
	return string(e)
}

// UnmarshalGQL converts the supplied value to an address type
func (e *AddressType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = AddressType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid AddressType", str)
	}
	return nil
}

// MarshalGQL writes the address type to the supplied writer
func (e AddressType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// AddressUse is used to set address uses.
//
// See: http://hl7.org/fhir/valueset-address-use.html
type AddressUse string

// known address uses
const (
	// AddressUseHome ...
	AddressUseHome AddressUse = "home"
	// AddressUseWork ...
	AddressUseWork AddressUse = "work"
	// AddressUseTemp ...
	AddressUseTemp AddressUse = "temp"
	// AddressUseOld ...
	AddressUseOld AddressUse = "old"
	// AddressUseBilling ...
	AddressUseBilling AddressUse = "billing"
)

// AllAddressUse is a list of all known address uses
var AllAddressUse = []AddressUse{
	AddressUseHome,
	AddressUseWork,
	AddressUseTemp,
	AddressUseOld,
	AddressUseBilling,
}

// IsValid returns true if an address use is valid
func (e AddressUse) IsValid() bool {
	switch e {
	case AddressUseHome, AddressUseWork, AddressUseTemp, AddressUseOld, AddressUseBilling:
		return true
	}
	return false
}

// String ...
func (e AddressUse) String() string {
	return string(e)
}

// UnmarshalGQL converts the supplied value to an address use
func (e *AddressUse) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = AddressUse(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid AddressUse", str)
	}
	return nil
}

// MarshalGQL writes the address to the supplied writer
func (e AddressUse) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// Country codes.
//
// See: https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes and
// https://www.iban.com/country-codes .
type Country string

// list of known country codes
const (
	// CountryAf ...
	CountryAf Country = "AF"
	// CountryAx ...
	CountryAx Country = "AX"
	// CountryAl ...
	CountryAl Country = "AL"
	// CountryDz ...
	CountryDz Country = "DZ"
	// CountryAs ...
	CountryAs Country = "AS"
	// CountryAd ...
	CountryAd Country = "AD"
	// CountryAo ...
	CountryAo Country = "AO"
	// CountryAi ...
	CountryAi Country = "AI"
	// CountryAq ...
	CountryAq Country = "AQ"
	// CountryAg ...
	CountryAg Country = "AG"
	// CountryAr ...
	CountryAr Country = "AR"
	// CountryAm ...
	CountryAm Country = "AM"
	// CountryAw ...
	CountryAw Country = "AW"
	// CountryAu ...
	CountryAu Country = "AU"
	// CountryAt ...
	CountryAt Country = "AT"
	// CountryAz ...
	CountryAz Country = "AZ"
	// CountryBs ...
	CountryBs Country = "BS"
	// CountryBh ...
	CountryBh Country = "BH"
	// CountryBd ...
	CountryBd Country = "BD"
	// CountryBb ...
	CountryBb Country = "BB"
	// CountryBy ...
	CountryBy Country = "BY"
	// CountryBe ...
	CountryBe Country = "BE"
	// CountryBz ...
	CountryBz Country = "BZ"
	// CountryBj ...
	CountryBj Country = "BJ"
	// CountryBm ...
	CountryBm Country = "BM"
	// CountryBt ...
	CountryBt Country = "BT"
	// CountryBo ...
	CountryBo Country = "BO"
	// CountryBq ...
	CountryBq Country = "BQ"
	// CountryBa ...
	CountryBa Country = "BA"
	// CountryBw ...
	CountryBw Country = "BW"
	// CountryBv ...
	CountryBv Country = "BV"
	// CountryBr ...
	CountryBr Country = "BR"
	// CountryIo ...
	CountryIo Country = "IO"
	// CountryBn ...
	CountryBn Country = "BN"
	// CountryBg ...
	CountryBg Country = "BG"
	// CountryBf ...
	CountryBf Country = "BF"
	// CountryBi ...
	CountryBi Country = "BI"
	// CountryCv ...
	CountryCv Country = "CV"
	// CountryKh ...
	CountryKh Country = "KH"
	// CountryCm ...
	CountryCm Country = "CM"
	// CountryCa ...
	CountryCa Country = "CA"
	// CountryKy ...
	CountryKy Country = "KY"
	// CountryCf ...
	CountryCf Country = "CF"
	// CountryTd ...
	CountryTd Country = "TD"
	// CountryCl ...
	CountryCl Country = "CL"
	// CountryCn ...
	CountryCn Country = "CN"
	// CountryCx ...
	CountryCx Country = "CX"
	// CountryCc ...
	CountryCc Country = "CC"
	// CountryCo ...
	CountryCo Country = "CO"
	// CountryKm ...
	CountryKm Country = "KM"
	// CountryCg ...
	CountryCg Country = "CG"
	// CountryCd ...
	CountryCd Country = "CD"
	// CountryCk ...
	CountryCk Country = "CK"
	// CountryCr ...
	CountryCr Country = "CR"
	// CountryCi ...
	CountryCi Country = "CI"
	// CountryHr ...
	CountryHr Country = "HR"
	// CountryCu ...
	CountryCu Country = "CU"
	// CountryCw ...
	CountryCw Country = "CW"
	// CountryCy ...
	CountryCy Country = "CY"
	// CountryCz ...
	CountryCz Country = "CZ"
	// CountryDk ...
	CountryDk Country = "DK"
	// CountryDj ...
	CountryDj Country = "DJ"
	// CountryDm ...
	CountryDm Country = "DM"
	// CountryDo ...
	CountryDo Country = "DO"
	// CountryEc ...
	CountryEc Country = "EC"
	// CountryEg ...
	CountryEg Country = "EG"
	// CountrySv ...
	CountrySv Country = "SV"
	// CountryGq ...
	CountryGq Country = "GQ"
	// CountryEr ...
	CountryEr Country = "ER"
	// CountryEe ...
	CountryEe Country = "EE"
	// CountrySz ...
	CountrySz Country = "SZ"
	// CountryEt ...
	CountryEt Country = "ET"
	// CountryFk ...
	CountryFk Country = "FK"
	// CountryFo ...
	CountryFo Country = "FO"
	// CountryFj ...
	CountryFj Country = "FJ"
	// CountryFi ...
	CountryFi Country = "FI"
	// CountryFr ...
	CountryFr Country = "FR"
	// CountryGf ...
	CountryGf Country = "GF"
	// CountryPf ...
	CountryPf Country = "PF"
	// CountryTf ...
	CountryTf Country = "TF"
	// CountryGa ...
	CountryGa Country = "GA"
	// CountryGm ...
	CountryGm Country = "GM"
	// CountryGe ...
	CountryGe Country = "GE"
	// CountryDe ...
	CountryDe Country = "DE"
	// CountryGh ...
	CountryGh Country = "GH"
	// CountryGi ...
	CountryGi Country = "GI"
	// CountryGr ...
	CountryGr Country = "GR"
	// CountryGl ...
	CountryGl Country = "GL"
	// CountryGd ...
	CountryGd Country = "GD"
	// CountryGp ...
	CountryGp Country = "GP"
	// CountryGu ...
	CountryGu Country = "GU"
	// CountryGt ...
	CountryGt Country = "GT"
	// CountryGg ...
	CountryGg Country = "GG"
	// CountryGn ...
	CountryGn Country = "GN"
	// CountryGw ...
	CountryGw Country = "GW"
	// CountryGy ...
	CountryGy Country = "GY"
	// CountryHt ...
	CountryHt Country = "HT"
	// CountryHm ...
	CountryHm Country = "HM"
	// CountryVa ...
	CountryVa Country = "VA"
	// CountryHn ...
	CountryHn Country = "HN"
	// CountryHk ...
	CountryHk Country = "HK"
	// CountryHu ...
	CountryHu Country = "HU"
	// CountryIs ...
	CountryIs Country = "IS"
	// CountryIn ...
	CountryIn Country = "IN"
	// CountryID ...
	CountryID Country = "ID"
	// CountryIr ...
	CountryIr Country = "IR"
	// CountryIq ...
	CountryIq Country = "IQ"
	// CountryIe ...
	CountryIe Country = "IE"
	// CountryIm ...
	CountryIm Country = "IM"
	// CountryIl ...
	CountryIl Country = "IL"
	// CountryIt ...
	CountryIt Country = "IT"
	// CountryJm ...
	CountryJm Country = "JM"
	// CountryJp ...
	CountryJp Country = "JP"
	// CountryJe ...
	CountryJe Country = "JE"
	// CountryJo ...
	CountryJo Country = "JO"
	// CountryKz ...
	CountryKz Country = "KZ"
	// CountryKe ...
	CountryKe Country = "KE"
	// CountryKi ...
	CountryKi Country = "KI"
	// CountryKp ...
	CountryKp Country = "KP"
	// CountryKr ...
	CountryKr Country = "KR"
	// CountryKw ...
	CountryKw Country = "KW"
	// CountryKg ...
	CountryKg Country = "KG"
	// CountryLa ...
	CountryLa Country = "LA"
	// CountryLv ...
	CountryLv Country = "LV"
	// CountryLb ...
	CountryLb Country = "LB"
	// CountryLs ...
	CountryLs Country = "LS"
	// CountryLr ...
	CountryLr Country = "LR"
	// CountryLy ...
	CountryLy Country = "LY"
	// CountryLi ...
	CountryLi Country = "LI"
	// CountryLt ...
	CountryLt Country = "LT"
	// CountryLu ...
	CountryLu Country = "LU"
	// CountryMo ...
	CountryMo Country = "MO"
	// CountryMg ...
	CountryMg Country = "MG"
	// CountryMw ...
	CountryMw Country = "MW"
	// CountryMy ...
	CountryMy Country = "MY"
	// CountryMv ...
	CountryMv Country = "MV"
	// CountryMl ...
	CountryMl Country = "ML"
	// CountryMt ...
	CountryMt Country = "MT"
	// CountryMh ...
	CountryMh Country = "MH"
	// CountryMq ...
	CountryMq Country = "MQ"
	// CountryMr ...
	CountryMr Country = "MR"
	// CountryMu ...
	CountryMu Country = "MU"
	// CountryYt ...
	CountryYt Country = "YT"
	// CountryMx ...
	CountryMx Country = "MX"
	// CountryFm ...
	CountryFm Country = "FM"
	// CountryMd ...
	CountryMd Country = "MD"
	// CountryMc ...
	CountryMc Country = "MC"
	// CountryMn ...
	CountryMn Country = "MN"
	// CountryMe ...
	CountryMe Country = "ME"
	// CountryMs ...
	CountryMs Country = "MS"
	// CountryMa ...
	CountryMa Country = "MA"
	// CountryMz ...
	CountryMz Country = "MZ"
	// CountryMm ...
	CountryMm Country = "MM"
	// CountryNa ...
	CountryNa Country = "NA"
	// CountryNr ...
	CountryNr Country = "NR"
	// CountryNp ...
	CountryNp Country = "NP"
	// CountryNl ...
	CountryNl Country = "NL"
	// CountryNc ...
	CountryNc Country = "NC"
	// CountryNz ...
	CountryNz Country = "NZ"
	// CountryNi ...
	CountryNi Country = "NI"
	// CountryNe ...
	CountryNe Country = "NE"
	// CountryNg ...
	CountryNg Country = "NG"
	// CountryNu ...
	CountryNu Country = "NU"
	// CountryNf ...
	CountryNf Country = "NF"
	// CountryMk ...
	CountryMk Country = "MK"
	// CountryMp ...
	CountryMp Country = "MP"
	// CountryNo ...
	CountryNo Country = "NO"
	// CountryOm ...
	CountryOm Country = "OM"
	// CountryPk ...
	CountryPk Country = "PK"
	// CountryPw ...
	CountryPw Country = "PW"
	// CountryPs ...
	CountryPs Country = "PS"
	// CountryPa ...
	CountryPa Country = "PA"
	// CountryPg ...
	CountryPg Country = "PG"
	// CountryPy ...
	CountryPy Country = "PY"
	// CountryPe ...
	CountryPe Country = "PE"
	// CountryPh ...
	CountryPh Country = "PH"
	// CountryPn ...
	CountryPn Country = "PN"
	// CountryPl ...
	CountryPl Country = "PL"
	// CountryPt ...
	CountryPt Country = "PT"
	// CountryPr ...
	CountryPr Country = "PR"
	// CountryQa ...
	CountryQa Country = "QA"
	// CountryRe ...
	CountryRe Country = "RE"
	// CountryRo ...
	CountryRo Country = "RO"
	// CountryRu ...
	CountryRu Country = "RU"
	// CountryRw ...
	CountryRw Country = "RW"
	// CountryBl ...
	CountryBl Country = "BL"
	// CountrySh ...
	CountrySh Country = "SH"
	// CountryKn ...
	CountryKn Country = "KN"
	// CountryLc ...
	CountryLc Country = "LC"
	// CountryMf ...
	CountryMf Country = "MF"
	// CountryPm ...
	CountryPm Country = "PM"
	// CountryVc ...
	CountryVc Country = "VC"
	// CountryWs ...
	CountryWs Country = "WS"
	// CountrySm ...
	CountrySm Country = "SM"
	// CountrySt ...
	CountrySt Country = "ST"
	// CountrySa ...
	CountrySa Country = "SA"
	// CountrySn ...
	CountrySn Country = "SN"
	// CountryRs ...
	CountryRs Country = "RS"
	// CountrySc ...
	CountrySc Country = "SC"
	// CountrySl ...
	CountrySl Country = "SL"
	// CountrySg ...
	CountrySg Country = "SG"
	// CountrySx ...
	CountrySx Country = "SX"
	// CountrySk ...
	CountrySk Country = "SK"
	// CountrySi ...
	CountrySi Country = "SI"
	// CountrySb ...
	CountrySb Country = "SB"
	// CountrySo ...
	CountrySo Country = "SO"
	// CountryZa ...
	CountryZa Country = "ZA"
	// CountryGs ...
	CountryGs Country = "GS"
	// CountrySs ...
	CountrySs Country = "SS"
	// CountryEs ...
	CountryEs Country = "ES"
	// CountryLk ...
	CountryLk Country = "LK"
	// CountrySd ...
	CountrySd Country = "SD"
	// CountrySr ...
	CountrySr Country = "SR"
	// CountrySj ...
	CountrySj Country = "SJ"
	// CountrySe ...
	CountrySe Country = "SE"
	// CountryCh ...
	CountryCh Country = "CH"
	// CountrySy ...
	CountrySy Country = "SY"
	// CountryTw ...
	CountryTw Country = "TW"
	// CountryTj ...
	CountryTj Country = "TJ"
	// CountryTz ...
	CountryTz Country = "TZ"
	// CountryTh ...
	CountryTh Country = "TH"
	// CountryTl ...
	CountryTl Country = "TL"
	// CountryTg ...
	CountryTg Country = "TG"
	// CountryTk ...
	CountryTk Country = "TK"
	// CountryTo ...
	CountryTo Country = "TO"
	// CountryTt ...
	CountryTt Country = "TT"
	// CountryTn ...
	CountryTn Country = "TN"
	// CountryTr ...
	CountryTr Country = "TR"
	// CountryTm ...
	CountryTm Country = "TM"
	// CountryTc ...
	CountryTc Country = "TC"
	// CountryTv ...
	CountryTv Country = "TV"
	// CountryUg ...
	CountryUg Country = "UG"
	// CountryUa ...
	CountryUa Country = "UA"
	// CountryAe ...
	CountryAe Country = "AE"
	// CountryGb ...
	CountryGb Country = "GB"
	// CountryUs ...
	CountryUs Country = "US"
	// CountryUm ...
	CountryUm Country = "UM"
	// CountryUy ...
	CountryUy Country = "UY"
	// CountryUz ...
	CountryUz Country = "UZ"
	// CountryVu ...
	CountryVu Country = "VU"
	// CountryVe ...
	CountryVe Country = "VE"
	// CountryVn ...
	CountryVn Country = "VN"
	// CountryVg ...
	CountryVg Country = "VG"
	// CountryVi ...
	CountryVi Country = "VI"
	// CountryWf ...
	CountryWf Country = "WF"
	// CountryEh ...
	CountryEh Country = "EH"
	// CountryYe ...
	CountryYe Country = "YE"
	// CountryZm ...
	CountryZm Country = "ZM"
	// CountryZw ...
	CountryZw Country = "ZW"
)

// AllCountry is a list of all known country codes
var AllCountry = []Country{
	CountryAf,
	CountryAx,
	CountryAl,
	CountryDz,
	CountryAs,
	CountryAd,
	CountryAo,
	CountryAi,
	CountryAq,
	CountryAg,
	CountryAr,
	CountryAm,
	CountryAw,
	CountryAu,
	CountryAt,
	CountryAz,
	CountryBs,
	CountryBh,
	CountryBd,
	CountryBb,
	CountryBy,
	CountryBe,
	CountryBz,
	CountryBj,
	CountryBm,
	CountryBt,
	CountryBo,
	CountryBq,
	CountryBa,
	CountryBw,
	CountryBv,
	CountryBr,
	CountryIo,
	CountryBn,
	CountryBg,
	CountryBf,
	CountryBi,
	CountryCv,
	CountryKh,
	CountryCm,
	CountryCa,
	CountryKy,
	CountryCf,
	CountryTd,
	CountryCl,
	CountryCn,
	CountryCx,
	CountryCc,
	CountryCo,
	CountryKm,
	CountryCg,
	CountryCd,
	CountryCk,
	CountryCr,
	CountryCi,
	CountryHr,
	CountryCu,
	CountryCw,
	CountryCy,
	CountryCz,
	CountryDk,
	CountryDj,
	CountryDm,
	CountryDo,
	CountryEc,
	CountryEg,
	CountrySv,
	CountryGq,
	CountryEr,
	CountryEe,
	CountrySz,
	CountryEt,
	CountryFk,
	CountryFo,
	CountryFj,
	CountryFi,
	CountryFr,
	CountryGf,
	CountryPf,
	CountryTf,
	CountryGa,
	CountryGm,
	CountryGe,
	CountryDe,
	CountryGh,
	CountryGi,
	CountryGr,
	CountryGl,
	CountryGd,
	CountryGp,
	CountryGu,
	CountryGt,
	CountryGg,
	CountryGn,
	CountryGw,
	CountryGy,
	CountryHt,
	CountryHm,
	CountryVa,
	CountryHn,
	CountryHk,
	CountryHu,
	CountryIs,
	CountryIn,
	CountryID,
	CountryIr,
	CountryIq,
	CountryIe,
	CountryIm,
	CountryIl,
	CountryIt,
	CountryJm,
	CountryJp,
	CountryJe,
	CountryJo,
	CountryKz,
	CountryKe,
	CountryKi,
	CountryKp,
	CountryKr,
	CountryKw,
	CountryKg,
	CountryLa,
	CountryLv,
	CountryLb,
	CountryLs,
	CountryLr,
	CountryLy,
	CountryLi,
	CountryLt,
	CountryLu,
	CountryMo,
	CountryMg,
	CountryMw,
	CountryMy,
	CountryMv,
	CountryMl,
	CountryMt,
	CountryMh,
	CountryMq,
	CountryMr,
	CountryMu,
	CountryYt,
	CountryMx,
	CountryFm,
	CountryMd,
	CountryMc,
	CountryMn,
	CountryMe,
	CountryMs,
	CountryMa,
	CountryMz,
	CountryMm,
	CountryNa,
	CountryNr,
	CountryNp,
	CountryNl,
	CountryNc,
	CountryNz,
	CountryNi,
	CountryNe,
	CountryNg,
	CountryNu,
	CountryNf,
	CountryMk,
	CountryMp,
	CountryNo,
	CountryOm,
	CountryPk,
	CountryPw,
	CountryPs,
	CountryPa,
	CountryPg,
	CountryPy,
	CountryPe,
	CountryPh,
	CountryPn,
	CountryPl,
	CountryPt,
	CountryPr,
	CountryQa,
	CountryRe,
	CountryRo,
	CountryRu,
	CountryRw,
	CountryBl,
	CountrySh,
	CountryKn,
	CountryLc,
	CountryMf,
	CountryPm,
	CountryVc,
	CountryWs,
	CountrySm,
	CountrySt,
	CountrySa,
	CountrySn,
	CountryRs,
	CountrySc,
	CountrySl,
	CountrySg,
	CountrySx,
	CountrySk,
	CountrySi,
	CountrySb,
	CountrySo,
	CountryZa,
	CountryGs,
	CountrySs,
	CountryEs,
	CountryLk,
	CountrySd,
	CountrySr,
	CountrySj,
	CountrySe,
	CountryCh,
	CountrySy,
	CountryTw,
	CountryTj,
	CountryTz,
	CountryTh,
	CountryTl,
	CountryTg,
	CountryTk,
	CountryTo,
	CountryTt,
	CountryTn,
	CountryTr,
	CountryTm,
	CountryTc,
	CountryTv,
	CountryUg,
	CountryUa,
	CountryAe,
	CountryGb,
	CountryUs,
	CountryUm,
	CountryUy,
	CountryUz,
	CountryVu,
	CountryVe,
	CountryVn,
	CountryVg,
	CountryVi,
	CountryWf,
	CountryEh,
	CountryYe,
	CountryZm,
	CountryZw,
}

// IsValid returns True if a country code is valid
func (e Country) IsValid() bool {
	switch e {
	case CountryAf, CountryAx, CountryAl, CountryDz, CountryAs, CountryAd, CountryAo, CountryAi, CountryAq, CountryAg, CountryAr, CountryAm, CountryAw, CountryAu, CountryAt, CountryAz, CountryBs, CountryBh, CountryBd, CountryBb, CountryBy, CountryBe, CountryBz, CountryBj, CountryBm, CountryBt, CountryBo, CountryBq, CountryBa, CountryBw, CountryBv, CountryBr, CountryIo, CountryBn, CountryBg, CountryBf, CountryBi, CountryCv, CountryKh, CountryCm, CountryCa, CountryKy, CountryCf, CountryTd, CountryCl, CountryCn, CountryCx, CountryCc, CountryCo, CountryKm, CountryCg, CountryCd, CountryCk, CountryCr, CountryCi, CountryHr, CountryCu, CountryCw, CountryCy, CountryCz, CountryDk, CountryDj, CountryDm, CountryDo, CountryEc, CountryEg, CountrySv, CountryGq, CountryEr, CountryEe, CountrySz, CountryEt, CountryFk, CountryFo, CountryFj, CountryFi, CountryFr, CountryGf, CountryPf, CountryTf, CountryGa, CountryGm, CountryGe, CountryDe, CountryGh, CountryGi, CountryGr, CountryGl, CountryGd, CountryGp, CountryGu, CountryGt, CountryGg, CountryGn, CountryGw, CountryGy, CountryHt, CountryHm, CountryVa, CountryHn, CountryHk, CountryHu, CountryIs, CountryIn, CountryID, CountryIr, CountryIq, CountryIe, CountryIm, CountryIl, CountryIt, CountryJm, CountryJp, CountryJe, CountryJo, CountryKz, CountryKe, CountryKi, CountryKp, CountryKr, CountryKw, CountryKg, CountryLa, CountryLv, CountryLb, CountryLs, CountryLr, CountryLy, CountryLi, CountryLt, CountryLu, CountryMo, CountryMg, CountryMw, CountryMy, CountryMv, CountryMl, CountryMt, CountryMh, CountryMq, CountryMr, CountryMu, CountryYt, CountryMx, CountryFm, CountryMd, CountryMc, CountryMn, CountryMe, CountryMs, CountryMa, CountryMz, CountryMm, CountryNa, CountryNr, CountryNp, CountryNl, CountryNc, CountryNz, CountryNi, CountryNe, CountryNg, CountryNu, CountryNf, CountryMk, CountryMp, CountryNo, CountryOm, CountryPk, CountryPw, CountryPs, CountryPa, CountryPg, CountryPy, CountryPe, CountryPh, CountryPn, CountryPl, CountryPt, CountryPr, CountryQa, CountryRe, CountryRo, CountryRu, CountryRw, CountryBl, CountrySh, CountryKn, CountryLc, CountryMf, CountryPm, CountryVc, CountryWs, CountrySm, CountrySt, CountrySa, CountrySn, CountryRs, CountrySc, CountrySl, CountrySg, CountrySx, CountrySk, CountrySi, CountrySb, CountrySo, CountryZa, CountryGs, CountrySs, CountryEs, CountryLk, CountrySd, CountrySr, CountrySj, CountrySe, CountryCh, CountrySy, CountryTw, CountryTj, CountryTz, CountryTh, CountryTl, CountryTg, CountryTk, CountryTo, CountryTt, CountryTn, CountryTr, CountryTm, CountryTc, CountryTv, CountryUg, CountryUa, CountryAe, CountryGb, CountryUs, CountryUm, CountryUy, CountryUz, CountryVu, CountryVe, CountryVn, CountryVg, CountryVi, CountryWf, CountryEh, CountryYe, CountryZm, CountryZw:
		return true
	}
	return false
}

// String ...
func (e Country) String() string {
	return string(e)
}

// UnmarshalGQL turns the value into a country
func (e *Country) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Country(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Country", str)
	}
	return nil
}

// MarshalGQL writes the enum value to the supplied writer
func (e Country) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
