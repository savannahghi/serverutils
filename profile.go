package base

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"time"
)

// PermissionType defines the type of a permission
type PermissionType string

// bewell admin permissions.
// this is not exhausive. More will be added on a need by need basis after analysis of the application
// and assert what actions need to the admin-permissioned
const (
	PermissionTypeSuperAdmin  PermissionType = "SUPER_ADMIN"
	PermissionTypeAdmin       PermissionType = "ADMIN"
	PermissionTypeCreateAdmin PermissionType = "CREATE_ADMIN"
	PermissionTypeRemoveAdmin PermissionType = "REMOVE_ADMIN"
	PermissionTypeAddSupplier PermissionType = "ADD_SUPPLIER"
	// whether an admin can add a supplier
	PermissionTypeRemoveSupplier PermissionType = "REMOVE_SUPPLIER"
	// whether an admin can suspend a supplier
	PermissionTypeSuspendSupplier PermissionType = "SUSPEND_SUPPLIER"
	// whether an admin can unsuspend a supplier
	PermissionTypeUnSuspendSupplier PermissionType = "UNSUSPEND_SUPPLIER"
	// whether an admin can view and process(approve/reject) kyc requests
	PermissionTypeProcessKYC PermissionType = "PROCESS_KYC"
)

// DefaultSuperAdminPermissions generic permissions for super admins.
// These permissions should be given to the Be.Well dev team.
var DefaultSuperAdminPermissions []PermissionType = []PermissionType{
	PermissionTypeSuperAdmin,
	PermissionTypeCreateAdmin,
	PermissionTypeRemoveAdmin,
	PermissionTypeAddSupplier,
	PermissionTypeRemoveSupplier,
	PermissionTypeSuspendSupplier,
	PermissionTypeUnSuspendSupplier,
	PermissionTypeProcessKYC,
}

// DefaultAdminPermissions generic permissions for admins.
// These permissions should be given to SIL customer happiness and relationship
// management staff.
var DefaultAdminPermissions []PermissionType = []PermissionType{
	PermissionTypeSuperAdmin,
	PermissionTypeAdmin,
	PermissionTypeAddSupplier,
	PermissionTypeSuspendSupplier,
	PermissionTypeUnSuspendSupplier,
	PermissionTypeProcessKYC,
}

// LoginProviderType defines the method of used to login to bewell
type LoginProviderType string

// methods used to login
const (
	LoginProviderTypePhone          LoginProviderType = "PHONE"
	LoginProviderTypeSocialGoogle   LoginProviderType = "SOCIAL_GOOGLE"
	LoginProviderTypeSocialFacebook LoginProviderType = "SOCIAL_FACEBOOK"
	LoginProviderTypeAppleFacebook  LoginProviderType = "SOCIAL_APPLE"
)

// AccountType defines the various supplier account types
type AccountType string

// AccountTypeIndivdual is an example of a suppiler account type
const (
	AccountTypeIndividual   AccountType = "INDIVIDUAL"
	AccountTypeOrganisation AccountType = "ORGANISATION"
)

// AllAccountType is a slice that represents all the account types
var AllAccountType = []AccountType{
	AccountTypeIndividual,
	AccountTypeOrganisation,
}

// IsValid checks if the account type is valid
func (e AccountType) IsValid() bool {
	switch e {
	case AccountTypeIndividual, AccountTypeOrganisation:
		return true
	}
	return false
}

func (e AccountType) String() string {
	return string(e)
}

// UnmarshalGQL converts the input, if valid, into a account type value
func (e *AccountType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = AccountType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid AccountType", str)
	}
	return nil
}

// MarshalGQL converts AccountType into a valid JSON string
func (e AccountType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// PartnerType defines the different partner types we have in Be.Well
type PartnerType string

// PartnerTypeRider is an example of a partner type who is involved in delivery of goods
const (
	PartnerTypeRider          PartnerType = "RIDER"
	PartnerTypePractitioner   PartnerType = "PRACTITIONER"
	PartnerTypeProvider       PartnerType = "PROVIDER"
	PartnerTypePharmaceutical PartnerType = "PHARMACEUTICAL"
	PartnerTypeCoach          PartnerType = "COACH"
	PartnerTypeNutrition      PartnerType = "NUTRITION"
	PartnerTypeConsumer       PartnerType = "CONSUMER"
)

// AllPartnerType represents a list of the partner types we offer
var AllPartnerType = []PartnerType{
	PartnerTypeRider,
	PartnerTypePractitioner,
	PartnerTypeProvider,
	PartnerTypePharmaceutical,
	PartnerTypeCoach,
	PartnerTypeNutrition,
	PartnerTypeConsumer,
}

// IsValid checks if a partner type is valid or not
func (e PartnerType) IsValid() bool {
	switch e {
	case PartnerTypeRider, PartnerTypePractitioner, PartnerTypeProvider, PartnerTypePharmaceutical, PartnerTypeCoach, PartnerTypeNutrition, PartnerTypeConsumer:
		return true
	}
	return false
}

func (e PartnerType) String() string {
	return string(e)
}

// UnmarshalGQL converts the input, if valid, into an correct partner type value
func (e *PartnerType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = PartnerType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid PartnerType", str)
	}
	return nil
}

// MarshalGQL converts partner type into a valid JSON string
func (e PartnerType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// Cover is used to save a user's insurance details.
type Cover struct {
	IdentifierHash *string `json:"identifier_hash" firestore:"identifierHash"`
	PayerName      string  `json:"payer_name,omitempty" firestore:"payerName"`
	PayerSladeCode int     `json:"payer_slade_code,omitempty" firestore:"payerSladeCode"`
	MemberNumber   string  `json:"member_number,omitempty" firestore:"memberNumber"`
	MemberName     string  `json:"member_name,omitempty" firestore:"memberName"`
}

// IsEntity marks this struct as a GraphQL entity
func (c Cover) IsEntity() {}

// BioData structure of bio data information for a user
type BioData struct {
	FirstName   *string `json:"firstName" firestore:"firstName"`
	LastName    *string `json:"lastName" firestore:"lastName"`
	DateOfBirth *Date   `json:"dateOfBirth" firestore:"dateOfBirth"`
	Gender      Gender  `json:"gender" firestore:"gender"`
}

// VerifiedIdentifier metadata of how the user has logged in to bewell
type VerifiedIdentifier struct {
	UID           string            `json:"uid" firestore:"uid"`
	Timestamp     time.Time         `json:"timeStamp" firestore:"timeStamp"`
	LoginProvider LoginProviderType `json:"loginProvider" firestore:"loginProvider"`
}

// UserProfileRepository defines signatures that a repeository that acts on the userprofile should
// implement. Repository heres means a storage unit like firebase or mongodb or pg.
type UserProfileRepository interface {
	UpdateUserName(ctx context.Context, id string, userName string) error
	UpdatePrimaryPhoneNumber(ctx context.Context, id string, phoneNumber string) error
	UpdatePrimaryEmailAddress(ctx context.Context, id string, emailAddress string) error
	UpdateSecondaryPhoneNumbers(ctx context.Context, id string, phoneNumbers []string) error
	UpdateSecondaryEmailAddresses(ctx context.Context, id string, emailAddresses []string) error
	UpdateVerifiedIdentifiers(ctx context.Context, id string, identifiers []VerifiedIdentifier) error
	UpdateVerifiedUIDS(ctx context.Context, id string, uids []string) error
	UpdateSuspended(ctx context.Context, id string, status bool) error
	UpdatePhotoUploadID(ctx context.Context, id string, uploadID string) error
	UpdateCovers(ctx context.Context, id string, covers []Cover) error
	UpdatePushTokens(ctx context.Context, id string, pushToken []string) error
	UpdatePermissions(ctx context.Context, id string, perms []PermissionType) error
	UpdateBioData(ctx context.Context, id string, data BioData) error
	UpdateAddresses(ctx context.Context, id string, address Address, addressType AddressType) error
}

// UserProfile serializes the profile of the logged in user.
type UserProfile struct {
	// globally unique identifier for a profile
	ID string `json:"id" firestore:"id"`

	// unique user name. Synonymous to a handle
	// e.g @juliusowino
	// this will be auto-generated on first login, meaning a user must have a username
	UserName *string `json:"userName" firestore:"userName"`

	// VerifiedIdentifiers represent various ways the user has been able to login
	// and these providers point to the same user
	VerifiedIdentifiers []VerifiedIdentifier `json:"verifiedIdentifiers" firestore:"verifiedIdentifiers"`

	// uids associated with a profile. Theses UIDS should match those in the verfiedIdentifiers.
	// the purpose of having verifiedUIDS is enbale ease querying of the profile using firebase query constructs.
	// when we migrate to postgres, this will be retired
	// the length of verfiedIdentifiers and verifiedUIDS should match
	VerifiedUIDS []string `json:"verifiedUIDS" firestore:"verifiedUIDS"`

	// this is the first class unique attribute of a user profile.  A user profile MUST HAVE A PRIMARY PHONE NUMBER
	PrimaryPhone *string `json:"primaryPhone" firestore:"primaryPhone"`

	// this is the second class unique attribute of a user profile. This can be updated as the user desires
	PrimaryEmailAddress *string `json:"primaryEmailAddress" firestore:"primaryEmailAddress"`

	// these are all phone numbers associated with a user. These phone numbers can be promoted to PRIMARY PHONE NUMBER
	// and/or used for account recovery
	SecondaryPhoneNumbers []string `json:"secondaryPhoneNumbers" firestore:"secondaryPhoneNumbers"`

	SecondaryEmailAddresses []string `json:"secondaryEmailAddresses" firestore:"secondaryEmailAddresses"`

	PushTokens []string `json:"pushTokens,omitempty" firestore:"pushTokens"`

	// what the user is allowed to do. Only valid for admins
	Permissions []PermissionType `json:"permissions,omitempty" firestore:"permissions"`

	// we determine if a user is "live" by examining fields on their profile
	TermsAccepted bool `json:"terms_accepted,omitempty" firestore:"termsAccepted"`

	// determines whether a specific will be visible in query results. If the `true`, means the profile in not
	// in active state and the user should not be allowed to login
	Suspended bool `json:"suspended" firestore:"suspended"`

	// a user's profile photo can be stored as base 64 encoded PNG
	PhotoUploadID string `json:"photoUploadID,omitempty" firestore:"photoUploadID"`

	// a user can have zero or more insurance covers
	Covers []Cover `json:"covers,omitempty" firestore:"covers"`

	// a user's biodata is stored on the profile
	UserBioData BioData `json:"userBioData,omitempty" firestore:"userBioData"`

	// this is the user's home geo location
	HomeAddress *Address `json:"homeAddress,omitempty" firestore:"homeAddress"`

	// this is the user's work geo location
	WorkAddress *Address `json:"workAddress,omitempty" firestore:"workAddress"`
}

// IsEntity marks a profile as a GraphQL entity
func (u UserProfile) IsEntity() {}

// UpdateProfileUserName updates the profiles username attribute
func (u *UserProfile) UpdateProfileUserName(
	ctx context.Context,
	repo UserProfileRepository,
	userName string,
) error {
	return repo.UpdateUserName(ctx, u.ID, userName)
}

// UpdateProfilePrimaryPhoneNumber update the primary phone number for this user profile
func (u *UserProfile) UpdateProfilePrimaryPhoneNumber(
	ctx context.Context,
	repo UserProfileRepository,
	phoneNumber string,
) error {
	return repo.UpdatePrimaryPhoneNumber(ctx, u.ID, phoneNumber)
}

// UpdateProfilePrimaryEmailAddress update the primary phone number for this user profile
func (u *UserProfile) UpdateProfilePrimaryEmailAddress(
	ctx context.Context,
	repo UserProfileRepository,
	email string,
) error {
	return repo.UpdatePrimaryEmailAddress(ctx, u.ID, email)
}

// UpdateProfileSecondaryPhoneNumbers update the primary phone number for this user profile
func (u *UserProfile) UpdateProfileSecondaryPhoneNumbers(
	ctx context.Context,
	repo UserProfileRepository,
	phoneNumbers []string,
) error {
	return repo.UpdateSecondaryPhoneNumbers(ctx, u.ID, phoneNumbers)
}

// UpdateProfileSecondaryEmailAddresses update the primary phone number for this user profile
func (u *UserProfile) UpdateProfileSecondaryEmailAddresses(
	ctx context.Context,
	repo UserProfileRepository,
	emailAddresses []string,
) error {
	return repo.UpdateSecondaryEmailAddresses(ctx, u.ID, emailAddresses)
}

// UpdateProfileVerifiedIdentifiers updatess profile's verified identifiers
func (u *UserProfile) UpdateProfileVerifiedIdentifiers(
	ctx context.Context,
	repo UserProfileRepository,
	identifiers []VerifiedIdentifier,
) error {
	return repo.UpdateVerifiedIdentifiers(ctx, u.ID, identifiers)
}

// UpdateProfileVerifiedUIDS updatess profile's UIDs
func (u *UserProfile) UpdateProfileVerifiedUIDS(
	ctx context.Context,
	repo UserProfileRepository,
	uids []string,
) error {
	return repo.UpdateVerifiedUIDS(ctx, u.ID, uids)
}

// UpdateProfileSuspended update the profiles Suspended attribute
func (u *UserProfile) UpdateProfileSuspended(
	ctx context.Context,
	repo UserProfileRepository,
	status bool,
) error {
	return repo.UpdateSuspended(ctx, u.ID, status)
}

// UpdateProfilePhotoUploadID updates the profiles PhotoUploadID attribute
func (u *UserProfile) UpdateProfilePhotoUploadID(
	ctx context.Context,
	repo UserProfileRepository,
	uploadID string,
) error {
	return repo.UpdatePhotoUploadID(ctx, u.ID, uploadID)
}

// UpdateProfileCovers updates the profile covers attribute
func (u *UserProfile) UpdateProfileCovers(
	ctx context.Context,
	repo UserProfileRepository,
	covers []Cover,
) error {
	return repo.UpdateCovers(ctx, u.ID, covers)
}

// UpdateProfilePushTokens updates the profiles pushTokens
func (u *UserProfile) UpdateProfilePushTokens(
	ctx context.Context,
	repo UserProfileRepository,
	pushToken []string,
) error {
	return repo.UpdatePushTokens(ctx, u.ID, pushToken)
}

// UpdateProfilePermissions updates the profiles persmissions
func (u *UserProfile) UpdateProfilePermissions(
	ctx context.Context,
	repo UserProfileRepository,
	perms []PermissionType,
) error {
	return repo.UpdatePermissions(ctx, u.ID, perms)
}

//UpdateProfileBioData updates the profile biodata
func (u *UserProfile) UpdateProfileBioData(
	ctx context.Context,
	repo UserProfileRepository,
	data BioData,
) error {
	return repo.UpdateBioData(ctx, u.ID, data)
}

//UpdateProfileAddresses updates the profile with a user's addresses
func (u *UserProfile) UpdateProfileAddresses(
	ctx context.Context,
	repo UserProfileRepository,
	address Address,
	addressType AddressType,
) error {
	return repo.UpdateAddresses(
		ctx,
		u.ID,
		address,
		addressType,
	)
}

// UserCommunicationsSetting hold information about the user communication's channels.
// if a channel is true, we will be able to send them marketing or promotional messages
type UserCommunicationsSetting struct {
	ID            string `json:"id" firestore:"id"`
	ProfileID     string `json:"profileID" firestore:"profileID"`
	AllowWhatsApp bool   `json:"allowWhatsApp" firestore:"allowWhatsApp"`
	AllowTextSMS  bool   `json:"allowTextSMS" firestore:"allowTextSMS"`
	AllowPush     bool   `json:"allowPush" firestore:"allowPush"`
	AllowEmail    bool   `json:"allowEmail" firestore:"allowEmail"`
}

// UserResponse returns a user's sign up/in response
type UserResponse struct {
	Profile               *UserProfile               `json:"profile"`
	SupplierProfile       *Supplier                  `json:"supplierProfile"`
	CustomerProfile       *Customer                  `json:"customerProfile"`
	CommunicationSettings *UserCommunicationsSetting `json:"communicationSettings"`
	Auth                  AuthCredentialResponse     `json:"auth"`
}

// AuthCredentialResponse represents a user login response
type AuthCredentialResponse struct {
	CustomToken   *string `json:"customToken"`
	IDToken       *string `json:"id_token"`
	ExpiresIn     string  `json:"expires_in"`
	RefreshToken  string  `json:"refresh_token"`
	UID           string  `json:"uid"`
	IsAdmin       bool    `json:"is_admin"`
	IsAnonymous   bool    `json:"is_anonymous"`
	CanExperiment bool    `json:"can_experiment"`
}

// Customer used to create a customer request payload
type Customer struct {
	ID                 string             `json:"customerID" firestore:"id"`
	ProfileID          *string            `json:"profileID,omitempty" firestore:"profileID"`
	CustomerID         string             `json:"id,omitempty" firestore:"erpCustomerID"`
	ReceivablesAccount ReceivablesAccount `json:"receivables_account" firestore:"receivablesAccount"`
	Active             bool               `json:"active" firestore:"active"`
}

// ReceivablesAccount stores a customer's receivables account info
type ReceivablesAccount struct {
	ID          string `json:"id" firestore:"id"`
	Name        string `json:"name" firestore:"name"`
	IsActive    bool   `json:"is_active" firestore:"isActive"`
	Number      string `json:"number" firestore:"number"`
	Tag         string `json:"tag" firestore:"tag"`
	Description string `json:"description" firestore:"description"`
}

// PayablesAccount stores a supplier's payables account info
type PayablesAccount struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	IsActive    bool   `json:"is_active"`
	Number      string `json:"number"`
	Tag         string `json:"tag"`
	Description string `json:"description"`
}

// Supplier used to create a supplier request payload
type Supplier struct {
	ID                     string                 `json:"supplierID" firestore:"id"`
	ProfileID              *string                `json:"profileID" firestore:"profileID"`
	SupplierID             string                 `json:"id" firestore:"erpSupplierID"`
	SupplierName           string                 `json:"supplierName" firestore:"supplierName"`
	PayablesAccount        *PayablesAccount       `json:"payables_account" firestore:"payablesAccount"`
	SupplierKYC            map[string]interface{} `json:"supplierKYC" firestore:"supplierKYC"`
	Active                 bool                   `json:"active" firestore:"active"`
	AccountType            AccountType            `json:"accountType" firestore:"accountType"`
	UnderOrganization      bool                   `json:"underOrganization" firestore:"underOrganization"`
	IsOrganizationVerified bool                   `json:"isOrganizationVerified" firestore:"isOrganizationVerified"`
	SladeCode              string                 `json:"sladeCode" firestore:"sladeCode"`
	ParentOrganizationID   string                 `json:"parentOrganizationID" firestore:"parentOrganizationID"`
	OrganizationName       string                 `json:"organizationName" firestore:"organizationName"`
	HasBranches            bool                   `json:"hasBranches,omitempty" firestore:"hasBranches"`
	Location               *Location              `json:"location,omitempty" firestore:"location"`
	PartnerType            PartnerType            `json:"partnerType" firestore:"partnerType"`
	EDIUserProfile         *EDIUserProfile        `json:"ediuserprofile" firestore:"ediserprofile"`
	PartnerSetupComplete   bool                   `json:"partnerSetupComplete" firestore:"partnerSetupComplete"`
	KYCSubmitted           bool                   `json:"kycSubmitted" firestore:"kycSubmitted"`
}

// Location is used to store a user's branch or organisation
type Location struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	BranchSladeCode *string `json:"branchSladeCode"`
}

// OtpResponse returns an otp
type OtpResponse struct {
	OTP string `json:"otp"`
}

// Address holds Google Map's location data
type Address struct {
	Latitude         string  `json:"latitude"`
	Longitude        string  `json:"longitude"`
	Locality         *string `json:"locality"`
	Name             *string `json:"name"`
	PlaceID          *string `json:"placeID"`
	FormattedAddress *string `json:"formattedAddress"`
}
