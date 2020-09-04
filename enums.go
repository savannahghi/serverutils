package base

import (
	"fmt"
	"io"
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
	fmt.Fprint(w, strconv.Quote(e.String()))
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
	fmt.Fprint(w, strconv.Quote(e.String()))
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
	fmt.Fprint(w, strconv.Quote(e.String()))
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
	fmt.Fprint(w, strconv.Quote(e.String()))
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
	fmt.Fprint(w, strconv.Quote(e.String()))

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
	fmt.Fprint(w, strconv.Quote(e.String()))

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
	PractitionerSpecialtyOphthalmology                   PractitionerSpecialty = "OPGTHALMOLOGY"
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
	PractitionerSpecialtyOphthalmology,
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
	case PractitionerSpecialtyUnspecified, PractitionerSpecialtyAnaesthesia, PractitionerSpecialtyCardiothoracicSurgery, PractitionerSpecialtyClinicalMedicalGenetics, PractitionerSpecialtyClincicalPathology, PractitionerSpecialtyGeneralPathology, PractitionerSpecialtyAnatomicPathology, PractitionerSpecialtyClinicalOncology, PractitionerSpecialtyDermatology, PractitionerSpecialtyEarNoseAndThroat, PractitionerSpecialtyEmergencyMedicine, PractitionerSpecialtyFamilyMedicine, PractitionerSpecialtyGeneralSurgery, PractitionerSpecialtyGeriatrics, PractitionerSpecialtyImmunology, PractitionerSpecialtyInfectiousDisease, PractitionerSpecialtyInternalMedicine, PractitionerSpecialtyMicrobiology, PractitionerSpecialtyNeurosurgery, PractitionerSpecialtyObstetricsAndGynaecology, PractitionerSpecialtyOccupationalMedicine, PractitionerSpecialtyOphthalmology, PractitionerSpecialtyOrthopaedicSurgery, PractitionerSpecialtyOncology, PractitionerSpecialtyOncologyRadiotherapy, PractitionerSpecialtyPaediatricsAndChildHealth, PractitionerSpecialtyPalliativeMedicine, PractitionerSpecialtyPlasticAndReconstructiveSurgery, PractitionerSpecialtyPsychiatry, PractitionerSpecialtyPublicHealth, PractitionerSpecialtyRadiology, PractitionerSpecialtyUrology:
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
	fmt.Fprint(w, strconv.Quote(e.String()))

}

// IDDocumentType is an internal code system for identification document types.
type IDDocumentType string

// ID type constants
const (
	// IDDocumentTypeNationalID ...
	IDDocumentTypeNationalID IDDocumentType = "national_id"
	// IDDocumentTypePassport ...
	IDDocumentTypePassport IDDocumentType = "passport"
	// IDDocumentTypeAlienID ...
	IDDocumentTypeAlienID IDDocumentType = "alien_id"
)

// AllIDDocumentType is a list of known ID types
var AllIDDocumentType = []IDDocumentType{
	IDDocumentTypeNationalID,
	IDDocumentTypePassport,
	IDDocumentTypeAlienID,
}

// IsValid checks that the ID type is valid
func (e IDDocumentType) IsValid() bool {
	switch e {
	case IDDocumentTypeNationalID, IDDocumentTypePassport, IDDocumentTypeAlienID:
		return true
	}
	return false
}

// String ...
func (e IDDocumentType) String() string {
	return string(e)
}

// UnmarshalGQL translates the input value to an ID type
func (e *IDDocumentType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = IDDocumentType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid IDDocumentType", str)
	}
	return nil
}

// MarshalGQL writes the enum value to the supplied writer
func (e IDDocumentType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))

}

// MaritalStatus is used to code individuals' marital statuses.
//
// See: https://www.hl7.org/fhir/valueset-marital-status.html
type MaritalStatus string

// known marital statuses
const (
	// MaritalStatusA ...
	MaritalStatusA MaritalStatus = "A"
	// MaritalStatusD ...
	MaritalStatusD MaritalStatus = "D"
	// MaritalStatusI ...
	MaritalStatusI MaritalStatus = "I"
	// MaritalStatusL ...
	MaritalStatusL MaritalStatus = "L"
	// MaritalStatusM ...
	MaritalStatusM MaritalStatus = "M"
	// MaritalStatusP ...
	MaritalStatusP MaritalStatus = "P"
	// MaritalStatusS ...
	MaritalStatusS MaritalStatus = "S"
	// MaritalStatusT ...
	MaritalStatusT MaritalStatus = "T"
	// MaritalStatusU ...
	MaritalStatusU MaritalStatus = "U"
	// MaritalStatusW ...
	MaritalStatusW MaritalStatus = "W"
	// MaritalStatusUnk ...
	MaritalStatusUnk MaritalStatus = "UNK"
)

// AllMaritalStatus is a list of known marital statuses
var AllMaritalStatus = []MaritalStatus{
	MaritalStatusA,
	MaritalStatusD,
	MaritalStatusI,
	MaritalStatusL,
	MaritalStatusM,
	MaritalStatusP,
	MaritalStatusS,
	MaritalStatusT,
	MaritalStatusU,
	MaritalStatusW,
	MaritalStatusUnk,
}

// IsValid checks that the marital status is valid
func (e MaritalStatus) IsValid() bool {
	switch e {
	case MaritalStatusA, MaritalStatusD, MaritalStatusI, MaritalStatusL, MaritalStatusM, MaritalStatusP, MaritalStatusS, MaritalStatusT, MaritalStatusU, MaritalStatusW, MaritalStatusUnk:
		return true
	}
	return false
}

// String ...
func (e MaritalStatus) String() string {
	return string(e)
}

// UnmarshalGQL turns the supplied input into a marital status enum value
func (e *MaritalStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = MaritalStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid MaritalStatus", str)
	}
	return nil
}

// MarshalGQL writes the enum value to the supplied writer
func (e MaritalStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// RelationshipType defines relationship types for patients.
//
// See: https://www.hl7.org/fhir/valueset-relatedperson-relationshiptype.html
type RelationshipType string

// list of known relationship types
const (
	// RelationshipTypeC ...
	RelationshipTypeC RelationshipType = "C"
	// RelationshipTypeE ...
	RelationshipTypeE RelationshipType = "E"
	// RelationshipTypeF ...
	RelationshipTypeF RelationshipType = "F"
	// RelationshipTypeI ...
	RelationshipTypeI RelationshipType = "I"
	// RelationshipTypeN ...
	RelationshipTypeN RelationshipType = "N"
	// RelationshipTypeO ...
	RelationshipTypeO RelationshipType = "O"
	// RelationshipTypeS ...
	RelationshipTypeS RelationshipType = "S"
	// RelationshipTypeU ...
	RelationshipTypeU RelationshipType = "U"
)

// AllRelationshipType is a list of all known relationship types
var AllRelationshipType = []RelationshipType{
	RelationshipTypeC,
	RelationshipTypeE,
	RelationshipTypeF,
	RelationshipTypeI,
	RelationshipTypeN,
	RelationshipTypeO,
	RelationshipTypeS,
	RelationshipTypeU,
}

// IsValid ensures that the relationship type is valid
func (e RelationshipType) IsValid() bool {
	switch e {
	case RelationshipTypeC, RelationshipTypeE, RelationshipTypeF, RelationshipTypeI, RelationshipTypeN, RelationshipTypeO, RelationshipTypeS, RelationshipTypeU:
		return true
	}
	return false
}

// String ...
func (e RelationshipType) String() string {
	return string(e)
}

// UnmarshalGQL converts its input (if valid) into a relationship type
func (e *RelationshipType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = RelationshipType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid RelationshipType", str)
	}
	return nil
}

// MarshalGQL writes the relationship type to the supplied writer
func (e RelationshipType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// IdentifierUse is a code system for identifier uses.
//
// See: https://www.hl7.org/fhir/valueset-identifier-use.html
type IdentifierUse string

// known identifier use values
const (
	// IdentifierUseUsual ...
	IdentifierUseUsual IdentifierUse = "usual"
	// IdentifierUseOfficial ...
	IdentifierUseOfficial IdentifierUse = "official"
	// IdentifierUseTemp ...
	IdentifierUseTemp IdentifierUse = "temp"
	// IdentifierUseSecondary ...
	IdentifierUseSecondary IdentifierUse = "secondary"
	// IdentifierUseOld ...
	IdentifierUseOld IdentifierUse = "old"
)

// AllIdentifierUse is a list of all known identifier uses
var AllIdentifierUse = []IdentifierUse{
	IdentifierUseUsual,
	IdentifierUseOfficial,
	IdentifierUseTemp,
	IdentifierUseSecondary,
	IdentifierUseOld,
}

// IsValid returns True if the enum value is valid
func (e IdentifierUse) IsValid() bool {
	switch e {
	case IdentifierUseUsual, IdentifierUseOfficial, IdentifierUseTemp, IdentifierUseSecondary, IdentifierUseOld:
		return true
	}
	return false
}

// String ...
func (e IdentifierUse) String() string {
	return string(e)
}

// UnmarshalGQL translates from the supplied value to a valid enum value
func (e *IdentifierUse) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = IdentifierUse(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid IdentifierUse", str)
	}
	return nil
}

// MarshalGQL writes the enum value to the supplied writer
func (e IdentifierUse) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// ContactPointSystem defines the type of contact it is.
//
// See: http://hl7.org/fhir/valueset-contact-point-system.html
type ContactPointSystem string

// known contact systems
const (
	// ContactPointSystemPhone ...
	ContactPointSystemPhone ContactPointSystem = "phone"
	// ContactPointSystemFax ...
	ContactPointSystemFax ContactPointSystem = "fax"
	// ContactPointSystemEmail ...
	ContactPointSystemEmail ContactPointSystem = "email"
	// ContactPointSystemPager ...
	ContactPointSystemPager ContactPointSystem = "pager"
	// ContactPointSystemURL ...
	ContactPointSystemURL ContactPointSystem = "url"
	// ContactPointSystemSms ...
	ContactPointSystemSms ContactPointSystem = "sms"
	// ContactPointSystemOther ...
	ContactPointSystemOther ContactPointSystem = "other"
)

// AllContactPointSystem is a list of known contact systems
var AllContactPointSystem = []ContactPointSystem{
	ContactPointSystemPhone,
	ContactPointSystemFax,
	ContactPointSystemEmail,
	ContactPointSystemPager,
	ContactPointSystemURL,
	ContactPointSystemSms,
	ContactPointSystemOther,
}

// IsValid checks that the contact system is valid
func (e ContactPointSystem) IsValid() bool {
	switch e {
	case ContactPointSystemPhone, ContactPointSystemFax, ContactPointSystemEmail, ContactPointSystemPager, ContactPointSystemURL, ContactPointSystemSms, ContactPointSystemOther:
		return true
	}
	return false
}

// String ...
func (e ContactPointSystem) String() string {
	return string(e)
}

// UnmarshalGQL converts the supplied value to a contact point system
func (e *ContactPointSystem) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ContactPointSystem(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ContactPointSystem", str)
	}
	return nil
}

// MarshalGQL writes the enum value to the supplied writer
func (e ContactPointSystem) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// ContactPointUse defines the uses of a contact.
//
// See: https://www.hl7.org/fhir/valueset-contact-point-use.html
type ContactPointUse string

// contact point uses
const (
	// ContactPointUseHome ...
	ContactPointUseHome ContactPointUse = "home"
	// ContactPointUseWork ...
	ContactPointUseWork ContactPointUse = "work"
	// ContactPointUseTemp ...
	ContactPointUseTemp ContactPointUse = "temp"
	// ContactPointUseOld ...
	ContactPointUseOld ContactPointUse = "old"
	// ContactPointUseMobile ...
	ContactPointUseMobile ContactPointUse = "mobile"
)

// AllContactPointUse is a list of known contact point uses
var AllContactPointUse = []ContactPointUse{
	ContactPointUseHome,
	ContactPointUseWork,
	ContactPointUseTemp,
	ContactPointUseOld,
	ContactPointUseMobile,
}

// IsValid checks that the enum value is valid
func (e ContactPointUse) IsValid() bool {
	switch e {
	case ContactPointUseHome, ContactPointUseWork, ContactPointUseTemp, ContactPointUseOld, ContactPointUseMobile:
		return true
	}
	return false
}

// String ...
func (e ContactPointUse) String() string {
	return string(e)
}

// UnmarshalGQL converts the supplied interface to a contact point use value
func (e *ContactPointUse) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ContactPointUse(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ContactPointUse", str)
	}
	return nil
}

// MarshalGQL writes the enum value to the supplied writer
func (e ContactPointUse) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// NameUse is used to define the uses of a human name.
//
// See: https://www.hl7.org/fhir/valueset-name-use.html
type NameUse string

// known name uses
const (
	// NameUseUsual ...
	NameUseUsual NameUse = "usual"
	// NameUseOfficial ...
	NameUseOfficial NameUse = "official"
	// NameUseTemp ...
	NameUseTemp NameUse = "temp"
	// NameUseNickname ...
	NameUseNickname NameUse = "nickname"
	// NameUseAnonymous ...
	NameUseAnonymous NameUse = "anonymous"
	// NameUseOld ...
	NameUseOld NameUse = "old"
	// NameUseMaiden ...
	NameUseMaiden NameUse = "maiden"
)

// AllNameUse is a list of known name uses
var AllNameUse = []NameUse{
	NameUseUsual,
	NameUseOfficial,
	NameUseTemp,
	NameUseNickname,
	NameUseAnonymous,
	NameUseOld,
	NameUseMaiden,
}

// IsValid checks that the name use is valid
func (e NameUse) IsValid() bool {
	switch e {
	case NameUseUsual, NameUseOfficial, NameUseTemp, NameUseNickname, NameUseAnonymous, NameUseOld, NameUseMaiden:
		return true
	}
	return false
}

// String ...
func (e NameUse) String() string {
	return string(e)
}

// UnmarshalGQL turns the supplied value into a name use enum
func (e *NameUse) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = NameUse(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid NameUse", str)
	}
	return nil
}

// MarshalGQL writes the name use enum value to the supplied writer
func (e NameUse) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// CalendarView is used to determine what view of a calendar to render
type CalendarView string

// calendar view constants
const (
	// CalendarViewDay ...
	CalendarViewDay CalendarView = "DAY"
	// CalendarViewWeek ...
	CalendarViewWeek CalendarView = "WEEK"
)

// AllCalendarView is a list of calendar views
var AllCalendarView = []CalendarView{
	CalendarViewDay,
	CalendarViewWeek,
}

// IsValid returns true if a calendar view is valid
func (e CalendarView) IsValid() bool {
	switch e {
	case CalendarViewDay, CalendarViewWeek:
		return true
	}
	return false
}

// String ...
func (e CalendarView) String() string {
	return string(e)
}

// UnmarshalGQL converts the input value into a calendar view
func (e *CalendarView) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = CalendarView(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid CalendarView", str)
	}
	return nil
}

// MarshalGQL writes the calendar view value to the supplied writer
func (e CalendarView) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
