package go_utils

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

	// partner management permissions
	PermissionTypeCreatePartner PermissionType = "CREATE_PARTNER"
	PermissionTypeUpdatePartner PermissionType = "UPDATE_PARTNER"
	PermissionTypeDeletePartner PermissionType = "DELETE_PARTNER"

	// consumer management permissions
	PermissionTypeCreateConsumer PermissionType = "CREATE_CONSUMER"
	PermissionTypeUpdateConsumer PermissionType = "UPDATE_CONSUMER"
	PermissionTypeDeleteConsumer PermissionType = "DELETE_CONSUMER"

	// patient management permissions
	PermissionTypeCreatePatient   PermissionType = "CREATE_PATIENT"
	PermissionTypeUpdatePatient   PermissionType = "UPDATE_PATIENT"
	PermissionTypeDeletePatient   PermissionType = "DELETE_PATIENT"
	PermissionTypeIdentifyPatient PermissionType = "IDENTIFY_PATIENT"
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

//DefaultEmployeePermissions generic permissions for field agents
// These permissions should be given to SIL field agents
var DefaultEmployeePermissions []PermissionType = []PermissionType{
	PermissionTypeCreateConsumer,
	PermissionTypeUpdateConsumer,
	PermissionTypeDeleteConsumer,
	PermissionTypeCreatePatient,
	PermissionTypeUpdatePatient,
	PermissionTypeDeletePatient,
	PermissionTypeIdentifyPatient,
}

//DefaultAgentPermissions generic permissions for field agents.
// These permissions should be given to SIL field agents
var DefaultAgentPermissions []PermissionType = []PermissionType{
	PermissionTypeCreatePartner,
	PermissionTypeUpdatePartner,
	PermissionTypeCreateConsumer,
	PermissionTypeUpdateConsumer,
}

// RoleType defines the type of role a subject has
// and the associated permissions
type RoleType string

// Various roles in bewell
const (
	RoleTypeEmployee RoleType = "EMPLOYEE"
	RoleTypeAgent    RoleType = "AGENT"
)

// IsValid checks if the role type is valid
func (r RoleType) IsValid() bool {
	switch r {
	case RoleTypeEmployee, RoleTypeAgent:
		return true
	default:
		return false
	}
}

// Permissions returns permissions for a certain role
func (r RoleType) Permissions() []PermissionType {
	switch r {
	case RoleTypeEmployee:
		return DefaultEmployeePermissions

	case RoleTypeAgent:
		return DefaultAgentPermissions
	default:
		return []PermissionType{}
	}
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
	IdentifierHash        *string   `json:"identifier_hash" firestore:"identifierHash"`
	PayerName             string    `json:"payer_name,omitempty" firestore:"payerName"`
	PayerSladeCode        int       `json:"payer_slade_code,omitempty" firestore:"payerSladeCode"`
	MemberNumber          string    `json:"member_number,omitempty" firestore:"memberNumber"`
	MemberName            string    `json:"member_name,omitempty" firestore:"memberName"`
	BeneficiaryID         int       `json:"beneficiary_id,omitempty" firestore:"beneficiaryID"`
	EffectivePolicyNumber string    `json:"effective_policy_number,omitempty" firestore:"effectivePolicyNumber"`
	ValidFrom             time.Time `json:"valid_from,omitempty" firestore:"validFrom"`
	ValidTo               time.Time `json:"valid_to,omitempty" firestore:"validTo"`
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

// UserInfo is a collection of standard profile information for a user.
type UserInfo struct {
	DisplayName string `json:"displayName,omitempty"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
	PhotoURL    string `json:"photoUrl,omitempty"`
	// In the ProviderUserInfo[] ProviderID can be a short domain name (e.g. google.com),
	// or the identity of an OpenID identity provider.
	// In UserRecord.UserInfo it will return the constant string "firebase".
	ProviderID string `json:"providerId,omitempty"`
	UID        string `json:"rawId,omitempty"`
}

// IsEntity marks a profile as a GraphQL entity
func (u UserProfile) IsEntity() {}

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

// EDIUserProfile is used to (de)serialialize the auth server
// profile of the logged in user.
type EDIUserProfile struct {
	ID              int      `json:"id"`
	GUID            string   `json:"guid"`
	Email           string   `json:"email"`
	FirstName       string   `json:"first_name"`
	LastName        string   `json:"last_name"`
	OtherNames      string   `json:"other_names"`
	IsStaff         bool     `json:"is_staff"`
	IsActive        bool     `json:"is_active"`
	Organisation    int      `json:"organisation"`
	BusinessPartner string   `json:"business_partner"`
	Roles           []string `json:"roles"`
	BPType          string   `json:"bp_type"`
}

//IsEntity ...
func (e EDIUserProfile) IsEntity() {}

// Supplier used to create a supplier request payload
// You can add more or remove fields to suit your organization/project needs
type Supplier struct {
	ID                     string                 `json:"supplierID" firestore:"id"`
	ProfileID              *string                `json:"profileID" firestore:"profileID"`
	SupplierID             string                 `json:"id" firestore:"erpSupplierID"`
	SupplierName           string                 `json:"supplierName" firestore:"supplierName"`
	PayablesAccount        *PayablesAccount       `json:"payables_account" firestore:"payablesAccount"`
	SupplierKYC            map[string]interface{} `json:"supplierKYC" firestore:"supplierKYC"`
	Active                 bool                   `json:"active" firestore:"active"`
	AccountType            *AccountType           `json:"accountType" firestore:"accountType"`
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

// PermissionInput input required to create a permission
type PermissionInput struct {
	Action   string
	Resource string
}

// AuthorizedEmails represent emails to check whether they have access to certain resources
// TODO: make these Env Vars
var AuthorizedEmails = []string{"apa-dev@healthcloud.co.ke", "apa-prod@healthcloud.co.ke"}

// AuthorizedPhones represent phonenumbers to check whether they have access to certain resources
// TODO: make these Env Vars
var AuthorizedPhones = []string{"+254700000000"}

// GetLoggedInUser retrieves logged in user information
func GetLoggedInUser(ctx context.Context) (*UserInfo, error) {
	authToken, err := GetUserTokenFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("user auth token not found in context: %w", err)
	}

	authClient, err := GetFirebaseAuthClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get or create Firebase client: %w", err)
	}

	user, err := authClient.GetUser(ctx, authToken.UID)
	if err != nil {

		return nil, fmt.Errorf("unable to get user: %w", err)
	}
	return &UserInfo{
		UID:         user.UID,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		DisplayName: user.DisplayName,
		ProviderID:  user.ProviderID,
		PhotoURL:    user.PhotoURL,
	}, nil
}
