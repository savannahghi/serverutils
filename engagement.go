package base

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/segmentio/ksuid"
	"github.com/xeipuuv/gojsonschema"
)

// defaults
const (
	LogoURL                    = "https://assets.healthcloud.co.ke/bewell_logo.png"
	BlankImageURL              = "https://assets.healthcloud.co.ke/1px.png"
	SampleVideoURL             = "https://www.youtube.com/watch?v=bPiofmZGb8o"
	FallbackSchemaHost         = "https://schema.healthcloud.co.ke"
	SchemaHostEnvVarName       = "SCHEMA_HOST"
	LinkSchemaFile             = "link.schema.json"
	MessageSchemaFile          = "message.schema.json"
	ActionSchemaFile           = "action.schema.json"
	NudgeSchemaFile            = "nudge.schema.json"
	ItemSchemaFile             = "item.schema.json"
	FeedSchemaFile             = "feed.schema.json"
	ContextSchemaFile          = "context.schema.json"
	PayloadSchemaFile          = "payload.schema.json"
	EventSchemaFile            = "event.schema.json"
	StatusSchemaFile           = "status.schema.json"
	VisibilitySchemaFile       = "visibility.schema.json"
	NotificationBodySchemaFile = "notificationbody.schema.json"
)

// Element is a building block of a feed e.g a nudge, action, feed item etc
// An element should know how to validate itself against it's JSON schema
type Element interface {
	ValidateAndUnmarshal(b []byte) error
	ValidateAndMarshal() ([]byte, error)
}

// ValidateAndUnmarshal validates JSON against a named feed schema
// file then unmarshals it into the supplied feed element, which should be a
// pointer.
func ValidateAndUnmarshal(sch string, b []byte, el Element) error {
	err := validateAgainstSchema(sch, b)
	if err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	err = json.Unmarshal(b, el)
	if err != nil {
		return fmt.Errorf("can't unmarshal JSON to struct: %w", err)
	}
	return nil
}

// ValidateAndMarshal marshals a feed element to JSON, checks it against the
// indicated schema file and returns it if it is valid.
func ValidateAndMarshal(sch string, el Element) ([]byte, error) {
	bs, err := json.Marshal(el)
	if err != nil {
		return nil, fmt.Errorf("can't marshal %T to JSON: %w", el, err)
	}
	err = validateAgainstSchema(sch, bs)
	if err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	return bs, nil
}

// Action represents the global and non-global actions that a user can see/do
type Action struct {
	// A unique identifier for each action
	ID string `json:"id" firestore:"id"`

	// A higher sequence number means that it came later
	SequenceNumber int `json:"sequenceNumber" firestore:"sequenceNumber"`

	// A friendly name for the action; rich text with Unicode, can have emoji
	Name string `json:"name" firestore:"name"`

	// A link to a PNG image that would serve as an avatar
	Icon Link `json:"icon" firestore:"icon"`

	// Action types are: primary, secondary, overflow and floating
	// Primary actions get dominant visual treatment;
	// secondary actions less so;
	// overflow actions are hidden;
	// floating actions are material FABs
	ActionType ActionType `json:"actionType" firestore:"actionType"`

	// How the action should be handled e.g inline or full page.
	// This is a hint for frontend logic.
	Handling Handling `json:"handling" firestore:"handling"`

	// indicated whether this action should or can be triggered by na anoymous user
	AllowAnonymous bool `json:"allowAnonymous" firestore:"allowAnonymous"`
}

// ValidateAndUnmarshal checks that the input data is valid as per the
// relevant JSON schema and unmarshals it if it is
func (ac *Action) ValidateAndUnmarshal(b []byte) error {
	err := ValidateAndUnmarshal(ActionSchemaFile, b, ac)
	if err != nil {
		return fmt.Errorf("invalid action JSON: %w", err)
	}
	return nil
}

// ValidateAndMarshal validates against JSON schema then marshals to JSON
func (ac *Action) ValidateAndMarshal() ([]byte, error) {
	return ValidateAndMarshal(ActionSchemaFile, ac)
}

// IsEntity marks this as an Apollo federation GraphQL entity
func (ac Action) IsEntity() {}

// Event An event indicating that this action was triggered
type Event struct {
	// A unique identifier for each action
	ID string `json:"id" firestore:"id"`

	// An event name - two upper case words separated by an underscore
	Name string `json:"name" firestore:"name"`

	// Technical metadata - when/where/why/who/what/how etc
	Context Context `json:"context,omitempty" firestore:"context,omitempty"`

	// The actual 'business data' carried by the event
	Payload Payload `json:"payload,omitempty" firestore:"payload,omitempty"`
}

// ValidateAndUnmarshal checks that the input data is valid as per the
// relevant JSON schema and unmarshals it if it is
func (ev *Event) ValidateAndUnmarshal(b []byte) error {
	err := ValidateAndUnmarshal(EventSchemaFile, b, ev)
	if err != nil {
		return fmt.Errorf("invalid event JSON: %w", err)
	}
	return nil
}

// ValidateAndMarshal validates against JSON schema then marshals to JSON
func (ev *Event) ValidateAndMarshal() ([]byte, error) {
	return ValidateAndMarshal(EventSchemaFile, ev)
}

// IsEntity marks this as an Apollo federation GraphQL entity
func (ev Event) IsEntity() {}

// Context identifies when/where/why/who/what/how an event occurred.
type Context struct {
	// the system or human user that created this event
	UserID string `json:"userID" firestore:"userID"`

	// the flavour of the feed/app that originated this event
	Flavour Flavour `json:"flavour" firestore:"flavour"`

	// the client (organization) that this user belongs to
	OrganizationID string `json:"organizationID" firestore:"organizationID"`

	// the location (e.g branch) from which the event was sent
	LocationID string `json:"locationID" firestore:"locationID"`

	// when this event was sent
	Timestamp time.Time `json:"timestamp" firestore:"timestamp"`
}

// ValidateAndUnmarshal checks that the input data is valid as per the
// relevant JSON schema and unmarshals it if it is
func (ct *Context) ValidateAndUnmarshal(b []byte) error {
	err := ValidateAndUnmarshal(ContextSchemaFile, b, ct)
	if err != nil {
		return fmt.Errorf("invalid context JSON: %w", err)
	}
	return nil
}

// ValidateAndMarshal validates against JSON schema then marshals to JSON
func (ct *Context) ValidateAndMarshal() ([]byte, error) {
	return ValidateAndMarshal(ContextSchemaFile, ct)
}

// Payload carries the actual 'business data' carried by the event.
// It varies from event to event.
type Payload struct {
	Data map[string]interface{} `json:"data" firestore:"data"`
}

// ValidateAndUnmarshal checks that the input data is valid as per the
// relevant JSON schema and unmarshals it if it is
func (pl *Payload) ValidateAndUnmarshal(b []byte) error {
	err := ValidateAndUnmarshal(PayloadSchemaFile, b, pl)
	if err != nil {
		return fmt.Errorf("invalid payload JSON: %w", err)
	}
	return nil
}

// ValidateAndMarshal validates against JSON schema then marshals to JSON
func (pl *Payload) ValidateAndMarshal() ([]byte, error) {
	return ValidateAndMarshal(PayloadSchemaFile, pl)
}

// Nudge represents a "prompt" for a user e.g to set a PIN
type Nudge struct {
	// A unique identifier for each nudge
	ID string `json:"id" firestore:"id"`

	// A higher sequence number means that it came later
	SequenceNumber int `json:"sequenceNumber" firestore:"sequenceNumber"`

	// Visibility determines if a nudge should be visible or not
	Visibility Visibility `json:"visibility" firestore:"visibility"`

	// whether the nudge is done (acted on) or pending
	Status Status `json:"status" firestore:"status"`

	// When this nudge should be expired/removed, automatically. RFC3339.
	Expiry time.Time `json:"expiry" firestore:"expiry"`

	// the title (lead line) of the nudge
	Title string `json:"title" firestore:"title"`

	// the text/copy of the nudge
	Text string `json:"text" firestore:"text"`

	// an illustrative image for the nudge
	Links []Link `json:"links" firestore:"links"`

	// actions to include on the nudge
	Actions []Action `json:"actions" firestore:"actions"`

	// Identifiers of all the users that got this message
	Users []string `json:"users,omitempty" firestore:"users,omitempty"`

	// Identifiers of all the groups that got this message
	Groups []string `json:"groups,omitempty" firestore:"groups,omitempty"`

	// How the user should be notified of this new item, if at all
	NotificationChannels []Channel `json:"notificationChannels,omitempty" firestore:"notificationChannels,omitempty"`

	// Text/Message the user will see in their notifications body when an action is performed on a nudge
	NotificationBody NotificationBody `json:"notificationBody,omitempty" firestore:"notificationBody,omitempty"`
}

// ValidateAndUnmarshal checks that the input data is valid as per the
// relevant JSON schema and unmarshals it if it is
func (nu *Nudge) ValidateAndUnmarshal(b []byte) error {
	err := ValidateAndUnmarshal(NudgeSchemaFile, b, nu)
	if err != nil {
		return fmt.Errorf("invalid nudge JSON: %w", err)
	}
	return nil
}

// ValidateAndMarshal verifies against JSON schema then marshals to JSON
func (nu *Nudge) ValidateAndMarshal() ([]byte, error) {
	return ValidateAndMarshal(NudgeSchemaFile, nu)
}

// IsEntity marks this as an Apollo federation GraphQL entity
func (nu Nudge) IsEntity() {}

// Item is a single item in a feed or in an inbox
type Item struct {
	// A unique identifier for each feed item
	ID string `json:"id" firestore:"id"`

	// A higher sequence number means that it came later
	SequenceNumber int `json:"sequenceNumber" firestore:"sequenceNumber"`

	// When this feed item should be expired/removed, automatically. RFC3339.
	Expiry time.Time `json:"expiry" firestore:"expiry"`

	// If a feed item is persistent, it also goes to the inbox
	// AND triggers a push notification.
	// Pinning a feed item makes it persistent.
	Persistent bool `json:"persistent" firestore:"persistent"`

	// Whether the task under a feed item is completed, pending etc
	Status Status `json:"status" firestore:"status"`

	// Whether the feed item is to be shown or hidden
	Visibility Visibility `json:"visibility" firestore:"visibility"`

	// A link to a PNG image that would serve as an avatar
	Icon Link `json:"icon" firestore:"icon"`

	// The person - real or robot - that generated this feed item. Rich text.
	Author string `json:"author" firestore:"author"`

	// An OPTIONAL second title line. Rich text.
	Tagline string `json:"tagline" firestore:"tagline"`

	// A label e.g for the queue that this item belongs to
	Label string `json:"label" firestore:"label"`

	// When this feed item was created. RFC3339.
	// This is used to calculate the feed item's age for display.
	Timestamp time.Time `json:"timestamp" firestore:"timestamp"`

	// An OPTIONAL summary line. Rich text.
	Summary string `json:"summary" firestore:"summary"`

	// Rich text that can include any unicode e.g emoji
	Text string `json:"text" firestore:"text"`

	// TextType determines how the frontend will render the text
	TextType TextType `json:"textType" firestore:"textType"`

	// an illustrative image for the item
	Links []Link `json:"links" firestore:"links"`

	// Actions are the primary, secondary and overflow actions associated
	// with a feed item
	Actions []Action `json:"actions,omitempty" firestore:"actions,omitempty"`

	// Conversations are messages and replies around a feed item
	Conversations []Message `json:"conversations,omitempty" firestore:"conversations,omitempty"`

	// Identifiers of all the users that got this message
	Users []string `json:"users,omitempty" firestore:"users,omitempty"`

	// Identifiers of all the groups that got this message
	Groups []string `json:"groups,omitempty" firestore:"groups,omitempty"`

	// How the user should be notified of this new item, if at all
	NotificationChannels []Channel `json:"notificationChannels,omitempty" firestore:"notificationChannels,omitempty"`
}

// ValidateAndUnmarshal checks that the input data is valid as per the
// relevant JSON schema and unmarshals it if it is
func (it *Item) ValidateAndUnmarshal(b []byte) error {
	err := ValidateAndUnmarshal(ItemSchemaFile, b, it)
	if err != nil {
		return fmt.Errorf("invalid item JSON: %w", err)
	}

	if it.Icon.LinkType != LinkTypePngImage {
		return fmt.Errorf("an icon must be a PNG image")
	}

	return nil
}

// ValidateAndMarshal validates against JSON schema then marshals to JSON
func (it *Item) ValidateAndMarshal() ([]byte, error) {
	if it.Icon.LinkType != LinkTypePngImage {
		return nil, fmt.Errorf("an icon must be a PNG image")
	}

	return ValidateAndMarshal(ItemSchemaFile, it)
}

// IsEntity marks this as an Apollo federation GraphQL entity
func (it Item) IsEntity() {}

// Message is a message in a thread of conversations attached to a feed item
type Message struct {
	// A unique identifier for each message on the thread
	ID string `json:"id" firestore:"id"`

	// A higher sequence number means that it came later
	SequenceNumber int `json:"sequenceNumber" firestore:"sequenceNumber"`

	// Rich text that can include any unicode e.g emoji
	Text string `json:"text" firestore:"text"`

	// The unique ID of any message that this one is replying to - a thread
	ReplyTo string `json:"replyTo" firestore:"replyTo"`

	// The UID of the user that posted the message
	PostedByUID string `json:"postedByUID" firestore:"postedByUID"`

	// The UID of the user that posted the message
	PostedByName string `json:"postedByName" firestore:"postedByName"`

	// when this message was sent
	Timestamp time.Time `json:"timestamp" firestore:"timestamp"`
}

// ValidateAndUnmarshal checks that the input data is valid as per the
// relevant JSON schema and unmarshals it if it is
func (msg *Message) ValidateAndUnmarshal(b []byte) error {
	err := ValidateAndUnmarshal(MessageSchemaFile, b, msg)
	if err != nil {
		return fmt.Errorf("invalid message JSON: %w", err)
	}
	return nil
}

// ValidateAndMarshal validates against JSON schema then marshals to JSON
func (msg *Message) ValidateAndMarshal() ([]byte, error) {
	return ValidateAndMarshal(MessageSchemaFile, msg)
}

// Link holds references to media that is part of the feed.
// The URL should embed authentication details.
// The treatment will depend on the specified asset type.
type Link struct {
	// A unique identifier for each feed item
	ID string `json:"id" firestore:"id"`

	// A URL at which the video can be accessed.
	// For a private video, the URL should include authentication information.
	URL string `json:"url" firestore:"url"`

	// LinkType of link
	LinkType LinkType `json:"linkType" firestore:"linkType"`

	// name or title of the linked item
	Title string `json:"title" firestore:"title"`

	// details about the linked item
	Description string `json:"description" firestore:"description"`

	// A URL to a PNG image that represents a thumbnail for the item
	Thumbnail string `json:"thumbnail" firestore:"thumbnail"`
}

func (l *Link) validateLinkType() error {
	if !govalidator.IsURL(l.URL) {
		return fmt.Errorf("%s is not a valid URL", l.URL)
	}
	switch l.LinkType {
	case LinkTypePdfDocument:
		if !strings.Contains(l.URL, ".png") {
			return fmt.Errorf("%s does not end with .pdf", l.URL)
		}
	case LinkTypePngImage:
		if !strings.Contains(l.URL, ".png") {
			return fmt.Errorf("%s does not end with .png", l.URL)
		}
	case LinkTypeYoutubeVideo:
		if !strings.Contains(l.URL, "youtube.com") {
			return fmt.Errorf("%s is not a youtube.com URL", l.URL)
		}
	}
	return nil
}

// ValidateAndUnmarshal checks that the input data is valid as per the
// relevant JSON schema and unmarshals it if it is
func (l *Link) ValidateAndUnmarshal(b []byte) error {
	err := ValidateAndUnmarshal(LinkSchemaFile, b, l)
	if err != nil {
		return fmt.Errorf("invalid video JSON: %w", err)
	}

	return l.validateLinkType()
}

// ValidateAndMarshal validates against JSON schema then marshals to JSON
func (l *Link) ValidateAndMarshal() ([]byte, error) {
	err := l.validateLinkType()
	if err != nil {
		return nil, fmt.Errorf("can't marshal invalid link: %w", err)
	}
	return ValidateAndMarshal(LinkSchemaFile, l)
}

// ActionType defines the types for global actions
type ActionType string

// the known action types are constants
const (
	ActionTypePrimary   ActionType = "PRIMARY"
	ActionTypeSecondary ActionType = "SECONDARY"
	ActionTypeOverflow  ActionType = "OVERFLOW"
	ActionTypeFloating  ActionType = "FLOATING"
)

// AllActionType has the known set of action types
var AllActionType = []ActionType{
	ActionTypePrimary,
	ActionTypeSecondary,
	ActionTypeOverflow,
	ActionTypeFloating,
}

// IsValid returns true only for valid action types
func (e ActionType) IsValid() bool {
	switch e {
	case ActionTypePrimary,
		ActionTypeSecondary,
		ActionTypeOverflow,
		ActionTypeFloating:
		return true
	}
	return false
}

func (e ActionType) String() string {
	return string(e)
}

// UnmarshalGQL reads an action type from GQL
func (e *ActionType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ActionType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ActionType", str)
	}
	return nil
}

// MarshalGQL writes an action type to the supplied writer
func (e ActionType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// Handling determines whether an action is handled INLINE or
type Handling string

// known action handling strategies
const (
	HandlingInline   Handling = "INLINE"
	HandlingFullPage Handling = "FULL_PAGE"
)

// AllHandling is the set of all valid handling strategies
var AllHandling = []Handling{
	HandlingInline,
	HandlingFullPage,
}

// IsValid returns true only for valid handling strategies
func (e Handling) IsValid() bool {
	switch e {
	case HandlingInline, HandlingFullPage:
		return true
	}
	return false
}

func (e Handling) String() string {
	return string(e)
}

// UnmarshalGQL reads and validates a handling value from the supplied input
func (e *Handling) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Handling(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Handling", str)
	}
	return nil
}

// MarshalGQL writes the Handling value to the supplied writer
func (e Handling) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// Status is the set of known statuses for feed items and nudges
type Status string

// known item and nudge statuses
const (
	StatusPending    Status = "PENDING"
	StatusInProgress Status = "IN_PROGRESS"
	StatusDone       Status = "DONE"
)

// AllStatus is the set of known statuses
var AllStatus = []Status{
	StatusPending,
	StatusInProgress,
	StatusDone,
}

// IsValid returns true if a status is valid
func (e Status) IsValid() bool {
	switch e {
	case StatusPending, StatusInProgress, StatusDone:
		return true
	}
	return false
}

func (e Status) String() string {
	return string(e)
}

// UnmarshalGQL translates the input value given into a status
func (e *Status) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Status(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Status", str)
	}
	return nil
}

// MarshalGQL writes the status to the supplied writer
func (e Status) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// Visibility defines the visibility statuses of feed items
type Visibility string

// known visibility values
const (
	VisibilityShow Visibility = "SHOW"
	VisibilityHide Visibility = "HIDE"
)

// AllVisibility is the set of all known visibility values
var AllVisibility = []Visibility{
	VisibilityShow,
	VisibilityHide,
}

// IsValid returns true if a visibility value is valid
func (e Visibility) IsValid() bool {
	switch e {
	case VisibilityShow, VisibilityHide:
		return true
	}
	return false
}

func (e Visibility) String() string {
	return string(e)
}

// UnmarshalGQL reads and validates a visibility value
func (e *Visibility) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Visibility(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Visibility", str)
	}
	return nil
}

// MarshalGQL writes a visibility value into the supplied writer
func (e Visibility) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// Channel represents a notification challen
type Channel string

// known notification channels
const (
	ChannelFcm      Channel = "FCM"
	ChannelEmail    Channel = "EMAIL"
	ChannelSms      Channel = "SMS"
	ChannelWhatsapp Channel = "WHATSAPP"
)

// AllChannel is the set of all supported notification channels
var AllChannel = []Channel{
	ChannelFcm,
	ChannelEmail,
	ChannelSms,
	ChannelWhatsapp,
}

// IsValid returns True only for a valid channel
func (e Channel) IsValid() bool {
	switch e {
	case ChannelFcm, ChannelEmail, ChannelSms, ChannelWhatsapp:
		return true
	}
	return false
}

func (e Channel) String() string {
	return string(e)
}

// UnmarshalGQL converts the supplied input into a channel value
func (e *Channel) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Channel(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Channel", str)
	}
	return nil
}

// MarshalGQL writes the channel to the supplied writer
func (e Channel) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// LinkType determines how a linked asset is handled on the feed
type LinkType string

// known link types
const (
	LinkTypeYoutubeVideo LinkType = "YOUTUBE_VIDEO"
	LinkTypePngImage     LinkType = "PNG_IMAGE"
	LinkTypePdfDocument  LinkType = "PDF_DOCUMENT"
	LinkTypeSvgImage     LinkType = "SVG_IMAGE"
	LinkTypeDefault      LinkType = "DEFAULT"
)

// AllLinkType is the set of all known link types
var AllLinkType = []LinkType{
	LinkTypeYoutubeVideo,
	LinkTypePngImage,
	LinkTypePdfDocument,
	LinkTypeSvgImage,
	LinkTypeDefault,
}

// IsValid is true only when a link type is avalid
func (e LinkType) IsValid() bool {
	switch e {
	case LinkTypeYoutubeVideo, LinkTypePngImage, LinkTypePdfDocument, LinkTypeSvgImage, LinkTypeDefault:
		return true
	}
	return false
}

func (e LinkType) String() string {
	return string(e)
}

// UnmarshalGQL reads a link type from the supplied input
func (e *LinkType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = LinkType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid LinkType", str)
	}
	return nil
}

// MarshalGQL writes a link type to the supplied writer
func (e LinkType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// TextType determines how clients render the text
type TextType string

// known text types
const (
	TextTypeHTML     TextType = "HTML"
	TextTypeMarkdown TextType = "MARKDOWN"
	TextTypePlain    TextType = "PLAIN"
)

// AllTextType is the set of all known text types
var AllTextType = []TextType{
	TextTypeHTML,
	TextTypeMarkdown,
	TextTypePlain,
}

// IsValid returns true only for valid text types
func (e TextType) IsValid() bool {
	switch e {
	case TextTypeHTML, TextTypeMarkdown, TextTypePlain:
		return true
	}
	return false
}

func (e TextType) String() string {
	return string(e)
}

// UnmarshalGQL translates the supplied interface into a text type
func (e *TextType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = TextType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid TextType", str)
	}
	return nil
}

// MarshalGQL writes the text type to the supplied writer
func (e TextType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// Flavour is the flavour of a feed i.e consumer or pro
type Flavour string

// known flavours
const (
	FlavourPro      Flavour = "PRO"
	FlavourConsumer Flavour = "CONSUMER"
)

// AllFlavour is a set of all valid flavours
var AllFlavour = []Flavour{
	FlavourPro,
	FlavourConsumer,
}

// IsValid returns True if a feed is valid
func (e Flavour) IsValid() bool {
	switch e {
	case FlavourPro, FlavourConsumer:
		return true
	}
	return false
}

func (e Flavour) String() string {
	return string(e)
}

// UnmarshalGQL translates and validates the input flavour
func (e *Flavour) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Flavour(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Flavour", str)
	}
	return nil
}

// MarshalGQL writes the flavour to the supplied writer
func (e Flavour) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// Keys are the top level keys in a feed
type Keys string

// known feed keys
const (
	KeysActions Keys = "actions"
	KeysNudges  Keys = "nudges"
	KeysItems   Keys = "items"
)

// AllKeys is the set of all valid feed keys
var AllKeys = []Keys{
	KeysActions,
	KeysNudges,
	KeysItems,
}

// IsValid returns true if a feed key is valid
func (e Keys) IsValid() bool {
	switch e {
	case KeysActions, KeysNudges, KeysItems:
		return true
	}
	return false
}

func (e Keys) String() string {
	return string(e)
}

// UnmarshalGQL translates a feed key from a string
func (e *Keys) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Keys(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid FeedKeys", str)
	}
	return nil
}

// MarshalGQL writes the feed key to the supplied writer
func (e Keys) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// BooleanFilter defines true/false/both for filtering against bools
type BooleanFilter string

// known boolean filter value
const (
	BooleanFilterTrue  BooleanFilter = "TRUE"
	BooleanFilterFalse BooleanFilter = "FALSE"
	BooleanFilterBoth  BooleanFilter = "BOTH"
)

// IsValid is a set of known boolean filters
var IsValid = []BooleanFilter{
	BooleanFilterTrue,
	BooleanFilterFalse,
	BooleanFilterBoth,
}

// IsValid returns True if the boolean filter value is valid
func (e BooleanFilter) IsValid() bool {
	switch e {
	case BooleanFilterTrue, BooleanFilterFalse, BooleanFilterBoth:
		return true
	}
	return false
}

func (e BooleanFilter) String() string {
	return string(e)
}

// UnmarshalGQL reads the bool value in from input
func (e *BooleanFilter) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = BooleanFilter(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid BooleanFilter", str)
	}
	return nil
}

// MarshalGQL writes the bool value to the supplied writer
func (e BooleanFilter) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// GetPNGImageLink returns an initialized PNG image link.
//
// It is used in testing and default data generation.
func GetPNGImageLink(url string, title string, description string, thumbnailURL string) Link {
	return Link{
		ID:          ksuid.New().String(),
		URL:         url,
		LinkType:    LinkTypePngImage,
		Title:       title,
		Description: description,
		Thumbnail:   thumbnailURL,
	}
}

// GetSVGImageLink returns an initialized PNG image link.
//
// It is used in testing and default data generation.
func GetSVGImageLink(url string, title string, description string, thumbnailURL string) Link {
	return Link{
		ID:          ksuid.New().String(),
		URL:         url,
		LinkType:    LinkTypeSvgImage,
		Title:       title,
		Description: description,
		Thumbnail:   thumbnailURL,
	}
}

// GetYoutubeVideoLink returns an initialized YouTube video link.
//
// It is used in testing and default data generation.
func GetYoutubeVideoLink(url string, title string, description string, thumbnailURL string) Link {
	return Link{
		ID:          ksuid.New().String(),
		URL:         url,
		LinkType:    LinkTypeYoutubeVideo,
		Title:       title,
		Description: description,
		Thumbnail:   thumbnailURL,
	}
}

// GetPDFDocumentLink returns an initialized PDF document link.
//
// It is used in testing and default data generation.
func GetPDFDocumentLink(url string, title string, description string, thumbnailURL string) Link {
	return Link{
		ID:          ksuid.New().String(),
		URL:         url,
		LinkType:    LinkTypePdfDocument,
		Title:       title,
		Description: description,
		Thumbnail:   thumbnailURL,
	}
}

// getSchemaURL serves JSON schema from this server and only falls back to a
// remote schema host when the local server cannot serve the JSON schema files.
// This has been done so as to reduce the impact of the network and DNS on the
// schema validation process - a critical path activity.
func getSchemaURL() string {
	schemaHost, err := GetEnvVar(SchemaHostEnvVarName)
	if err != nil {
		log.Printf("can't get env var `%s`: %s", SchemaHostEnvVarName, err)
	}

	client := http.Client{
		Timeout: time.Second * 1, // aggressive timeout
	}
	req, err := http.NewRequest(http.MethodGet, schemaHost, nil)
	if err != nil {
		log.Printf("can't create request to local schema URL: %s", err)
	}
	if err == nil {
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("error accessing schema URL: %s", err)
		}
		if err == nil {
			if resp.StatusCode != http.StatusOK {
				log.Printf("schema URL error status code: %s", resp.Status)
			}
			if resp.StatusCode == http.StatusOK {
				return schemaHost // we want this case to be the most common
			}
		}
	}

	// fall back to an externally hosted schema
	return FallbackSchemaHost
}

func validateAgainstSchema(sch string, b []byte) error {
	schemaURL := fmt.Sprintf("%s/%s", getSchemaURL(), sch)
	schemaLoader := gojsonschema.NewReferenceLoader(schemaURL)
	documentLoader := gojsonschema.NewStringLoader(string(b))
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf(
			"failed to validate `%s` against %s, got %#v: %w",
			string(b),
			sch,
			result,
			err,
		)
	}
	if !result.Valid() {
		errMsgs := []string{}
		for _, vErr := range result.Errors() {
			errType := vErr.Type()
			val := vErr.Value()
			context := vErr.Context().String()
			field := vErr.Field()
			desc := vErr.Description()
			descFormat := vErr.DescriptionFormat()
			details := vErr.Details()
			errMsg := fmt.Sprintf(
				"errType: %s\nval: %s\ncontext: %s\nfield: %s\ndesc: %s\ndescFormat: %s\ndetails: %s\n",
				errType,
				val,
				context,
				field,
				desc,
				descFormat,
				details,
			)
			errMsgs = append(errMsgs, errMsg)
		}
		return fmt.Errorf(
			"the result of validating `%s` against %s is not valid: %#v",
			string(b),
			sch,
			errMsgs,
		)
	}
	return nil
}

// NotificationBody represents human readable messages sent in notifications
type NotificationBody struct {
	// Human readable rich text sent when an item/nudge is published to a user's Feed
	PublishMessage string `json:"publishMessage" firestore:"publishMessage"`

	// Human readable rich text sent when item/nudge is deleted to a user's Feed
	DeleteMessage string `json:"deleteMessage" firestore:"deleteMessage"`

	// Human readable rich text sent when a user does a RESOLVE action
	ResolveMessage string `json:"resolveMessage" firestore:"resolveMessage"`

	// Human readable rich text sent when a user does an UNRESOLVE action
	UnresolveMessage string `json:"unresolveMessage" firestore:"unresolveMessage"`

	// Human readable rich text sent when a user does a SHOW action
	ShowMessage string `json:"showMessage" firestore:"showMessage"`

	// Human readable rich text sent when a user does a HIDE action
	HideMessage string `json:"hideMessage" firestore:"hideMessage"`
}

// ValidateAndUnmarshal checks that the input data is valid as per the
// relevant JSON schema and unmarshals it if it is
func (nb *NotificationBody) ValidateAndUnmarshal(b []byte) error {
	err := ValidateAndUnmarshal(NotificationBodySchemaFile, b, nb)
	if err != nil {
		return fmt.Errorf("invalid notification body JSON: %w", err)
	}
	return nil
}

// ValidateAndMarshal validates against JSON schema then marshals to JSON
func (nb *NotificationBody) ValidateAndMarshal() ([]byte, error) {
	return ValidateAndMarshal(NotificationBodySchemaFile, nb)
}

// IsEntity marks this as an Apollo federation GraphQL entity
func (nb *NotificationBody) IsEntity() {}
