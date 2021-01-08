package base

import (
	"context"
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

// Cover is used to save a user's insurance details.
type Cover struct {
	PayerName      string `json:"payer_name,omitempty" firestore:"payerName"`
	PayerSladeCode int    `json:"payer_slade_code,omitempty" firestore:"payerSladeCode"`
	MemberNumber   string `json:"member_number,omitempty" firestore:"memberNumber"`
	MemberName     string `json:"member_name,omitempty" firestore:"memberName"`
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
	UpdateBioData(ctx context.Context, id string, data BioData) error
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

	SecondaryEmailAddresses []string `json:"secondaryEmailAddresses " firestore:"secondaryEmailAddresses"`

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
}

// IsEntity marks a profile as a GraphQL entity
func (u UserProfile) IsEntity() {}

// UpdateProfileUserName updates the profiles username attribute
func (u *UserProfile) UpdateProfileUserName(ctx context.Context, repo UserProfileRepository, userName string) error {
	return repo.UpdateUserName(ctx, u.ID, userName)
}

// UpdateProfilePrimaryPhoneNumber update the primary phone number for this user profile
func (u *UserProfile) UpdateProfilePrimaryPhoneNumber(ctx context.Context, repo UserProfileRepository, phoneNumber string) error {
	return repo.UpdatePrimaryPhoneNumber(ctx, u.ID, phoneNumber)
}

// UpdateProfilePrimaryEmailAddress update the primary phone number for this user profile
func (u *UserProfile) UpdateProfilePrimaryEmailAddress(ctx context.Context, repo UserProfileRepository, email string) error {
	return repo.UpdatePrimaryEmailAddress(ctx, u.ID, email)
}

// UpdateProfileSecondaryPhoneNumbers update the primary phone number for this user profile
func (u *UserProfile) UpdateProfileSecondaryPhoneNumbers(ctx context.Context, repo UserProfileRepository, phoneNumbers []string) error {
	return repo.UpdateSecondaryPhoneNumbers(ctx, u.ID, phoneNumbers)
}

// UpdateProfileSecondaryEmailAddresses update the primary phone number for this user profile
func (u *UserProfile) UpdateProfileSecondaryEmailAddresses(ctx context.Context, repo UserProfileRepository, emailAddresses []string) error {
	return repo.UpdateSecondaryEmailAddresses(ctx, u.ID, emailAddresses)
}

// UpdateProfileVerifiedIdentifiers updatess profile's verified identifiers
func (u *UserProfile) UpdateProfileVerifiedIdentifiers(ctx context.Context, repo UserProfileRepository, identifiers []VerifiedIdentifier) error {
	return repo.UpdateVerifiedIdentifiers(ctx, u.ID, identifiers)
}

// UpdateProfileVerifiedUIDS updatess profile's UIDs
func (u *UserProfile) UpdateProfileVerifiedUIDS(ctx context.Context, repo UserProfileRepository, uids []string) error {
	return repo.UpdateVerifiedUIDS(ctx, u.ID, uids)
}

// UpdateProfileSuspended update the profiles Suspended attribute
func (u *UserProfile) UpdateProfileSuspended(ctx context.Context, repo UserProfileRepository, status bool) error {
	return repo.UpdateSuspended(ctx, u.ID, status)
}

// UpdateProfilePhotoUploadID updates the profiles PhotoUploadID attribute
func (u *UserProfile) UpdateProfilePhotoUploadID(ctx context.Context, repo UserProfileRepository, uploadID string) error {
	return repo.UpdatePhotoUploadID(ctx, u.ID, uploadID)
}

// UpdateProfileCovers updates the profile covers attribute
func (u *UserProfile) UpdateProfileCovers(ctx context.Context, repo UserProfileRepository, covers []Cover) error {
	return repo.UpdateCovers(ctx, u.ID, covers)
}

// UpdateProfilePushTokens updates the profiles pushTokens
func (u *UserProfile) UpdateProfilePushTokens(ctx context.Context, repo UserProfileRepository, pushToken []string) error {
	return repo.UpdatePushTokens(ctx, u.ID, pushToken)
}

//UpdateProfileBioData updates the profile biodata
func (u *UserProfile) UpdateProfileBioData(ctx context.Context, repo UserProfileRepository, data BioData) error {
	return repo.UpdateBioData(ctx, u.ID, data)
}
