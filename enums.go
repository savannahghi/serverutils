package base

import (
	"fmt"
	"io"
	"log"
	"strconv"
)

// Gender is a code system for administrative gender.
//
// See: https://www.hl7.org/fhir/valueset-administrative-gender.html
type Gender string

// gender constants
const (
	GenderMale    Gender = "male"
	GenderFemale  Gender = "female"
	GenderOther   Gender = "other"
	GenderUnknown Gender = "unknown"
)

// AllGender is a list of known genders
var AllGender = []Gender{
	GenderMale,
	GenderFemale,
	GenderOther,
	GenderUnknown,
}

// IsValid returns True if the enum value is valid
func (e Gender) IsValid() bool {
	switch e {
	case GenderMale, GenderFemale, GenderOther, GenderUnknown:
		return true
	}
	return false
}

func (e Gender) String() string {
	return string(e)
}

// UnmarshalGQL translates from the supplied value to a valid enum value
func (e *Gender) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Gender(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Gender", str)
	}
	return nil
}

// MarshalGQL writes the enum value to the supplied writer
func (e Gender) MarshalGQL(w io.Writer) {
	_, _ = fmt.Fprint(w, strconv.Quote(e.String()))
}

// FieldType is used to represent the GraphQL enum that is used for filter parameters
type FieldType string

const (
	// FieldTypeBoolean represents a boolean filter parameter
	FieldTypeBoolean FieldType = "BOOLEAN"

	// FieldTypeTimestamp represents a timestamp filter parameter
	FieldTypeTimestamp FieldType = "TIMESTAMP"

	// FieldTypeNumber represents a numeric (decimal or float) filter parameter
	FieldTypeNumber FieldType = "NUMBER"

	// FieldTypeInteger represents an integer filter parameter
	FieldTypeInteger FieldType = "INTEGER"

	// FieldTypeString represents a string filter parameter
	FieldTypeString FieldType = "STRING"
)

// AllFieldType is a list of all field types, used to simulate/map to a GraphQL enum
var AllFieldType = []FieldType{
	FieldTypeBoolean,
	FieldTypeTimestamp,
	FieldTypeNumber,
	FieldTypeInteger,
	FieldTypeString,
}

// IsValid returns True if the supplied value is a valid field type
func (e FieldType) IsValid() bool {
	switch e {
	case FieldTypeBoolean, FieldTypeTimestamp, FieldTypeNumber, FieldTypeInteger, FieldTypeString:
		return true
	}
	return false
}

// String represents a GraphQL enum as a string
func (e FieldType) String() string {
	return string(e)
}

// UnmarshalGQL checks whether the supplied value is a valid gqlgen enum
// and returns an error if it is not
func (e *FieldType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = FieldType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid FieldType", str)
	}
	return nil
}

// MarshalGQL serializes the enum value to the supplied writer
func (e FieldType) MarshalGQL(w io.Writer) {
	_, _ = fmt.Fprint(w, strconv.Quote(e.String()))
}

// Operation is used to map to a gqlgen (GraphQL) enum that defines filter/comparison operations
type Operation string

const (
	// OperationLessThan represents < in a GraphQL enum
	OperationLessThan Operation = "LESS_THAN"

	// OperationLessThanOrEqualTo represents <= in a GraphQL enum
	OperationLessThanOrEqualTo Operation = "LESS_THAN_OR_EQUAL_TO"

	// OperationEqual represents = in a GraphQL enum
	OperationEqual Operation = "EQUAL"

	// OperationGreaterThan represents > in a GraphQL enum
	OperationGreaterThan Operation = "GREATER_THAN"

	// OperationGreaterThanOrEqualTo represents >= in a GraphQL enum
	OperationGreaterThanOrEqualTo Operation = "GREATER_THAN_OR_EQUAL_TO"

	// OperationIn represents "in" (for queries that supply a list of parameters)
	// in a GraphQL enum
	OperationIn Operation = "IN"

	// OperationContains represents "contains" (for queries that check that a fragment is contained)
	// in a field(s) in a GraphQL enum
	OperationContains Operation = "CONTAINS"
)

// AllOperation is a list of all valid operations for filter parameters
var AllOperation = []Operation{
	OperationLessThan,
	OperationLessThanOrEqualTo,
	OperationEqual,
	OperationGreaterThan,
	OperationGreaterThanOrEqualTo,
	OperationIn,
	OperationContains,
}

// IsValid returns true if the operation is valid
func (e Operation) IsValid() bool {
	switch e {
	case OperationLessThan, OperationLessThanOrEqualTo, OperationEqual, OperationGreaterThan, OperationGreaterThanOrEqualTo, OperationIn, OperationContains:
		return true
	}
	return false
}

// String renders an operation enum value as a string
func (e Operation) String() string {
	return string(e)
}

// UnmarshalGQL confirms that an enum value is valid and returns an error if it is not
func (e *Operation) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Operation(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Operation", str)
	}
	return nil
}

// MarshalGQL writes the enum value to the supplied writer
func (e Operation) MarshalGQL(w io.Writer) {
	_, _ = fmt.Fprint(w, strconv.Quote(e.String()))
}

// SortOrder is used to represent map sort directions to a GraphQl enum
type SortOrder string

const (
	// SortOrderAsc is for ascending sorts
	SortOrderAsc SortOrder = "ASC"

	// SortOrderDesc is for descending sorts
	SortOrderDesc SortOrder = "DESC"
)

// AllSortOrder is a list of all valid sort orders
var AllSortOrder = []SortOrder{
	SortOrderAsc,
	SortOrderDesc,
}

// IsValid returns true if the sort order is valid
func (e SortOrder) IsValid() bool {
	switch e {
	case SortOrderAsc, SortOrderDesc:
		return true
	}
	return false
}

// String renders the sort order as a plain string
func (e SortOrder) String() string {
	return string(e)
}

// UnmarshalGQL confirms that the supplied value is a valid sort order
// and returns an error if it is not
func (e *SortOrder) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SortOrder(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SortOrder", str)
	}
	return nil
}

// MarshalGQL writes the sort order to the supplied writer as a quoted string
func (e SortOrder) MarshalGQL(w io.Writer) {
	_, _ = fmt.Fprint(w, strconv.Quote(e.String()))
}

// ContentType defines accepted content types
type ContentType string

// Constants used to map to allowed MIME types
const (
	ContentTypePng ContentType = "PNG"
	ContentTypeJpg ContentType = "JPG"
	ContentTypePdf ContentType = "PDF"
)

// AllContentType is a list of all acceptable content types
var AllContentType = []ContentType{
	ContentTypePng,
	ContentTypeJpg,
	ContentTypePdf,
}

// IsValid ensures that the content type value is valid
func (e ContentType) IsValid() bool {
	switch e {
	case ContentTypePng, ContentTypeJpg, ContentTypePdf:
		return true
	}
	return false
}

func (e ContentType) String() string {
	return string(e)
}

// UnmarshalGQL turns the supplied value into a content type value
func (e *ContentType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ContentType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ContentType", str)
	}
	return nil
}

// MarshalGQL writes the value of this enum to the supplied writer
func (e ContentType) MarshalGQL(w io.Writer) {
	_, err := fmt.Fprint(w, strconv.Quote(e.String()))
	if err != nil {
		log.Printf("%v\n", err)
	}
}

// Language defines allowed languages for uploads
type Language string

// Constants used to map to allowed languages
const (
	LanguageEn Language = "en"
	LanguageSw Language = "sw"
)

// LanguageCodingSystem is the FHIR language coding system
const LanguageCodingSystem = "urn:ietf:bcp:47"

// LanguageCodingVersion is the FHIR language value
const LanguageCodingVersion = ""

// LanguageNames is a map of language codes to language names
var LanguageNames = map[Language]string{
	LanguageEn: "English",
	LanguageSw: "Swahili",
}

// AllLanguage is a list of all allowed languages
var AllLanguage = []Language{
	LanguageEn,
	LanguageSw,
}

// IsValid ensures that the supplied language value is correct
func (e Language) IsValid() bool {
	switch e {
	case LanguageEn, LanguageSw:
		return true
	}
	return false
}

func (e Language) String() string {
	return string(e)
}

// UnmarshalGQL translates the input to a language type value
func (e *Language) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Language(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Language", str)
	}
	return nil
}

// MarshalGQL writes the value of this enum to the supplied writer
func (e Language) MarshalGQL(w io.Writer) {
	_, err := fmt.Fprint(w, strconv.Quote(e.String()))
	if err != nil {
		log.Printf("%v\n", err)
	}
}

// PractitionerSpecialty is a list of recognised health worker specialties.
//
// See: https://medicalboard.co.ke/resources_page/gazetted-specialties/
type PractitionerSpecialty string

// list of known practitioner specialties
const (
	PractitionerSpecialtyUnspecified                     PractitionerSpecialty = "UNSPECIFIED"
	PractitionerSpecialtyAnaesthesia                     PractitionerSpecialty = "ANAESTHESIA"
	PractitionerSpecialtyCardiothoracicSurgery           PractitionerSpecialty = "CARDIOTHORACIC_SURGERY"
	PractitionerSpecialtyClinicalMedicalGenetics         PractitionerSpecialty = "CLINICAL_MEDICAL_GENETICS"
	PractitionerSpecialtyClincicalPathology              PractitionerSpecialty = "CLINCICAL_PATHOLOGY"
	PractitionerSpecialtyGeneralPathology                PractitionerSpecialty = "GENERAL_PATHOLOGY"
	PractitionerSpecialtyAnatomicPathology               PractitionerSpecialty = "ANATOMIC_PATHOLOGY"
	PractitionerSpecialtyClinicalOncology                PractitionerSpecialty = "CLINICAL_ONCOLOGY"
	PractitionerSpecialtyDermatology                     PractitionerSpecialty = "DERMATOLOGY"
	PractitionerSpecialtyEarNoseAndThroat                PractitionerSpecialty = "EAR_NOSE_AND_THROAT"
	PractitionerSpecialtyEmergencyMedicine               PractitionerSpecialty = "EMERGENCY_MEDICINE"
	PractitionerSpecialtyFamilyMedicine                  PractitionerSpecialty = "FAMILY_MEDICINE"
	PractitionerSpecialtyGeneralSurgery                  PractitionerSpecialty = "GENERAL_SURGERY"
	PractitionerSpecialtyGeriatrics                      PractitionerSpecialty = "GERIATRICS"
	PractitionerSpecialtyImmunology                      PractitionerSpecialty = "IMMUNOLOGY"
	PractitionerSpecialtyInfectiousDisease               PractitionerSpecialty = "INFECTIOUS_DISEASE"
	PractitionerSpecialtyInternalMedicine                PractitionerSpecialty = "INTERNAL_MEDICINE"
	PractitionerSpecialtyMicrobiology                    PractitionerSpecialty = "MICROBIOLOGY"
	PractitionerSpecialtyNeurosurgery                    PractitionerSpecialty = "NEUROSURGERY"
	PractitionerSpecialtyObstetricsAndGynaecology        PractitionerSpecialty = "OBSTETRICS_AND_GYNAECOLOGY"
	PractitionerSpecialtyOccupationalMedicine            PractitionerSpecialty = "OCCUPATIONAL_MEDICINE"
	PractitionerSpecialtyOpthalmology                    PractitionerSpecialty = "OPTHALMOLOGY"
	PractitionerSpecialtyOrthopaedicSurgery              PractitionerSpecialty = "ORTHOPAEDIC_SURGERY"
	PractitionerSpecialtyOncology                        PractitionerSpecialty = "ONCOLOGY"
	PractitionerSpecialtyOncologyRadiotherapy            PractitionerSpecialty = "ONCOLOGY_RADIOTHERAPY"
	PractitionerSpecialtyPaediatricsAndChildHealth       PractitionerSpecialty = "PAEDIATRICS_AND_CHILD_HEALTH"
	PractitionerSpecialtyPalliativeMedicine              PractitionerSpecialty = "PALLIATIVE_MEDICINE"
	PractitionerSpecialtyPlasticAndReconstructiveSurgery PractitionerSpecialty = "PLASTIC_AND_RECONSTRUCTIVE_SURGERY"
	PractitionerSpecialtyPsychiatry                      PractitionerSpecialty = "PSYCHIATRY"
	PractitionerSpecialtyPublicHealth                    PractitionerSpecialty = "PUBLIC_HEALTH"
	PractitionerSpecialtyRadiology                       PractitionerSpecialty = "RADIOLOGY"
	PractitionerSpecialtyUrology                         PractitionerSpecialty = "UROLOGY"
)

// AllPractitionerSpecialty is the set of known practitioner specialties
var AllPractitionerSpecialty = []PractitionerSpecialty{
	PractitionerSpecialtyUnspecified,
	PractitionerSpecialtyAnaesthesia,
	PractitionerSpecialtyCardiothoracicSurgery,
	PractitionerSpecialtyClinicalMedicalGenetics,
	PractitionerSpecialtyClincicalPathology,
	PractitionerSpecialtyGeneralPathology,
	PractitionerSpecialtyAnatomicPathology,
	PractitionerSpecialtyClinicalOncology,
	PractitionerSpecialtyDermatology,
	PractitionerSpecialtyEarNoseAndThroat,
	PractitionerSpecialtyEmergencyMedicine,
	PractitionerSpecialtyFamilyMedicine,
	PractitionerSpecialtyGeneralSurgery,
	PractitionerSpecialtyGeriatrics,
	PractitionerSpecialtyImmunology,
	PractitionerSpecialtyInfectiousDisease,
	PractitionerSpecialtyInternalMedicine,
	PractitionerSpecialtyMicrobiology,
	PractitionerSpecialtyNeurosurgery,
	PractitionerSpecialtyObstetricsAndGynaecology,
	PractitionerSpecialtyOccupationalMedicine,
	PractitionerSpecialtyOpthalmology,
	PractitionerSpecialtyOrthopaedicSurgery,
	PractitionerSpecialtyOncology,
	PractitionerSpecialtyOncologyRadiotherapy,
	PractitionerSpecialtyPaediatricsAndChildHealth,
	PractitionerSpecialtyPalliativeMedicine,
	PractitionerSpecialtyPlasticAndReconstructiveSurgery,
	PractitionerSpecialtyPsychiatry,
	PractitionerSpecialtyPublicHealth,
	PractitionerSpecialtyRadiology,
	PractitionerSpecialtyUrology,
}

// IsValid returns True if the practitioner specialty is valid
func (e PractitionerSpecialty) IsValid() bool {
	switch e {
	case PractitionerSpecialtyUnspecified, PractitionerSpecialtyAnaesthesia, PractitionerSpecialtyCardiothoracicSurgery, PractitionerSpecialtyClinicalMedicalGenetics, PractitionerSpecialtyClincicalPathology, PractitionerSpecialtyGeneralPathology, PractitionerSpecialtyAnatomicPathology, PractitionerSpecialtyClinicalOncology, PractitionerSpecialtyDermatology, PractitionerSpecialtyEarNoseAndThroat, PractitionerSpecialtyEmergencyMedicine, PractitionerSpecialtyFamilyMedicine, PractitionerSpecialtyGeneralSurgery, PractitionerSpecialtyGeriatrics, PractitionerSpecialtyImmunology, PractitionerSpecialtyInfectiousDisease, PractitionerSpecialtyInternalMedicine, PractitionerSpecialtyMicrobiology, PractitionerSpecialtyNeurosurgery, PractitionerSpecialtyObstetricsAndGynaecology, PractitionerSpecialtyOccupationalMedicine, PractitionerSpecialtyOpthalmology, PractitionerSpecialtyOrthopaedicSurgery, PractitionerSpecialtyOncology, PractitionerSpecialtyOncologyRadiotherapy, PractitionerSpecialtyPaediatricsAndChildHealth, PractitionerSpecialtyPalliativeMedicine, PractitionerSpecialtyPlasticAndReconstructiveSurgery, PractitionerSpecialtyPsychiatry, PractitionerSpecialtyPublicHealth, PractitionerSpecialtyRadiology, PractitionerSpecialtyUrology:
		return true
	}
	return false
}

func (e PractitionerSpecialty) String() string {
	return string(e)
}

// UnmarshalGQL converts the supplied value to a practitioner specialty
func (e *PractitionerSpecialty) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = PractitionerSpecialty(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid PractitionerSpecialty", str)
	}
	return nil
}

// MarshalGQL writes the practitioner specialty to the supplied writer
func (e PractitionerSpecialty) MarshalGQL(w io.Writer) {
	_, err := fmt.Fprint(w, strconv.Quote(e.String()))
	if err != nil {
		log.Printf("%v\n", err)
	}
}