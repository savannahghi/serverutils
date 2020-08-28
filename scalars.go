package base

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/shopspring/decimal"
)

// Base64Binary is a stream of bytes
type Base64Binary string

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (sc *Base64Binary) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("scalar must be serialized as a string")
	}
	*sc = Base64Binary(str)
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (sc Base64Binary) MarshalGQL(w io.Writer) {
	_, _ = w.Write([]byte(strconv.Quote(string(sc))))
}

// Canonical is a URI that is a reference to a canonical URL on a FHIR resource
type Canonical string

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (sc *Canonical) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("scalar must be serialized as a string")
	}
	*sc = Canonical(str)
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (sc Canonical) MarshalGQL(w io.Writer) {
	_, _ = w.Write([]byte(strconv.Quote(string(sc))))
}

// Code is a string which has at least one character and no leading or trailing whitespace and where there is no whitespace other than single spaces in the contents
type Code string

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (sc *Code) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("scalar must be serialized as a string")
	}
	*sc = Code(str)
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (sc Code) MarshalGQL(w io.Writer) {
	_, _ = w.Write([]byte(strconv.Quote(string(sc))))
}

// DateTime is a date, date-time or partial date (e.g. just year or year + month).  If hours and minutes are specified, a time zone SHALL be populated. The format is a union of the schema types gYear, gYearMonth, date and dateTime. Seconds must be provided due to schema type constraints but may be zero-filled and may be ignored.
// Dates SHALL be valid dates.
type DateTime string

// Time converts the DateTime to a time
func (sc DateTime) Time() time.Time {
	d, err := time.Parse(DateTimeFormatLayout, string(sc))
	if err != nil {
		return time.Unix(0, 0)
	}
	return d
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (sc *DateTime) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("scalar must be serialized as a string")
	}
	*sc = DateTime(str)
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (sc DateTime) MarshalGQL(w io.Writer) {
	_, _ = w.Write([]byte(strconv.Quote(string(sc))))
}

// Instant is an instant in time - known at least to the second
type Instant string

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (sc *Instant) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("scalar must be serialized as a string")
	}
	*sc = Instant(str)
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (sc Instant) MarshalGQL(w io.Writer) {
	_, _ = w.Write([]byte(strconv.Quote(string(sc))))
}

// Markdown is a string that may contain Github Flavored Markdown syntax for optional processing by a mark down presentation engine
type Markdown string

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (sc *Markdown) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("scalar must be serialized as a string")
	}
	*sc = Markdown(str)
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (sc Markdown) MarshalGQL(w io.Writer) {
	_, _ = w.Write([]byte(strconv.Quote(string(sc))))
}

// OID is an OID represented as a URI
type OID string

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (sc *OID) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("scalar must be serialized as a string")
	}
	*sc = OID(str)
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (sc OID) MarshalGQL(w io.Writer) {
	_, _ = w.Write([]byte(strconv.Quote(string(sc))))
}

// URI is string of characters used to identify a name or a resource
type URI string

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (sc *URI) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("scalar must be serialized as a string")
	}
	*sc = URI(str)
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (sc URI) MarshalGQL(w io.Writer) {
	_, _ = w.Write([]byte(strconv.Quote(string(sc))))
}

// URL is a URI that is a literal reference
type URL string

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (sc *URL) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("scalar must be serialized as a string")
	}
	valid := govalidator.IsURL(str)
	if !valid {
		errorMessage := fmt.Sprintf("%s is not a valid URL", str)
		return fmt.Errorf(errorMessage)
	}
	*sc = URL(str)
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (sc URL) MarshalGQL(w io.Writer) {
	valid := govalidator.IsURL(string(sc))
	if !valid {
		errorMessage := fmt.Sprintf("%s is not a valid URL", sc)
		_, _ = w.Write([]byte(strconv.Quote(errorMessage)))
		return
	}
	_, _ = w.Write([]byte(strconv.Quote(string(sc))))
}

// UUID is a UUID, represented as a URI
type UUID string

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (sc *UUID) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("scalar must be serialized as a string")
	}
	*sc = UUID(str)
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (sc UUID) MarshalGQL(w io.Writer) {
	_, _ = w.Write([]byte(strconv.Quote(string(sc))))
}

// XHTML is xhtml - escaped html (see specfication)
type XHTML string

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (sc *XHTML) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("scalar must be serialized as a string")
	}
	*sc = XHTML(str)
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (sc XHTML) MarshalGQL(w io.Writer) {
	_, _ = w.Write([]byte(strconv.Quote(string(sc))))
}

// Decimal is a rational number with implicit precision
type Decimal decimal.Decimal

// String renders the underlying decimal value as a string
func (sc *Decimal) String() string {
	dec := decimal.Decimal(*sc)
	return dec.String()
}

// Decimal returns the underlying decimal
func (sc *Decimal) Decimal() decimal.Decimal {
	return decimal.Decimal(*sc)
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (sc *Decimal) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("scalar must be serialized as a string")
	}
	dec, err := decimal.NewFromString(str)
	if err != nil {
		return fmt.Errorf(
			"can't parse '%s' into decimal, error: %s", str, err)
	}
	deci := Decimal(dec)
	*sc = deci
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (sc Decimal) MarshalGQL(w io.Writer) {
	strRepr := strconv.Quote(sc.String())
	_, _ = w.Write([]byte(strRepr))
}

// NewDate initializes and validates a date
func NewDate(day, month, year int) (*Date, error) {
	date := Date{
		Year:  year,
		Month: month,
		Day:   day,
	}
	err := date.Validate()
	if err != nil {
		return nil, err
	}
	return &date, err
}

// Date is a custom date type that maintains only date level precision
type Date struct {
	Year  int
	Month int
	Day   int
}

// Validate checks that the date makes sense
func (d *Date) Validate() error {
	if d.Year < 1800 {
		return fmt.Errorf("the year must be > 1800")
	}

	if d.Year > 2100 {
		return fmt.Errorf(
			"too far in the future mate, is this software still running?")
	}

	if d.Month < 1 {
		return fmt.Errorf("the month cannot be < 1")
	}

	if d.Month > 12 {
		return fmt.Errorf("the month cannot be > 12")
	}

	if d.Day < 1 {
		return fmt.Errorf("the day cannot be < 1")
	}

	if d.Day > 31 {
		return fmt.Errorf("the day cannot be > 31")
	}

	return nil
}

// AsTime returns a Go stdlib time that corresponds to this date
func (d Date) AsTime() time.Time {
	return time.Date(
		int(d.Year),
		time.Month(d.Month),
		int(d.Day),
		0,
		0,
		0,
		0,
		time.UTC,
	)
}

// MarshalText translates the date into text
func (d Date) MarshalText() (text []byte, err error) {
	err = d.Validate()
	if err != nil {
		return nil, err
	}
	t := d.AsTime()
	text = []byte(t.Format(DateLayout))
	return // implicit return of text and err, to match encoding.TextMarshaler precisely
}

func (d Date) String() string {
	t := time.Date(
		int(d.Year),
		time.Month(d.Month),
		int(d.Day),
		0,
		0,
		0,
		0,
		time.UTC,
	)
	return t.Format("Jan 02, 2006")
}

// UnmarshalText parses the value from text
func (d *Date) UnmarshalText(text []byte) error {
	inp := string(text)
	t, err := time.Parse(DateLayout, inp)
	if err != nil {
		return fmt.Errorf("can't parse '%s' into date: %v", inp, err)
	}
	d.Year = t.Year()
	d.Month = int(t.Month())
	d.Day = t.Day()
	if err := d.Validate(); err != nil {
		return fmt.Errorf("invalid date (%v): %v", d, err)
	}
	return nil
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (d *Date) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("dates must be serialized as strings")
	}
	bs := []byte(str)
	return d.UnmarshalJSON(bs)
}

// MarshalGQL implements the graphql.Marshaler interface
func (d Date) MarshalGQL(w io.Writer) {
	bs, err := d.MarshalJSON()
	if err != nil {
		errMsg := fmt.Sprintf(
			"can't marshal date, error: %v", err)
		_, _ = w.Write([]byte(errMsg))
	}
	_, _ = w.Write(bs)
}

// MarshalJSON implements the json.Marshaler interface.
func (d Date) MarshalJSON() ([]byte, error) {
	text, err := d.MarshalText()
	if err != nil {
		return nil, err
	}
	return []byte(strconv.Quote(string(text))), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *Date) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	text := string(data)
	if text == "null" {
		return nil
	}
	stripped := strings.ReplaceAll(text, "\"", "")
	return d.UnmarshalText([]byte(stripped))
}
