package base

import (
	"bytes"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestNewDate(t *testing.T) {
	tests := map[string]struct {
		expectError          bool
		expectedErrorMessage string
		expectedMarshalText  string
		expectedJSON         string

		day   int
		month int
		year  int
	}{
		"valid_date": {
			day:                 1,
			month:               1,
			year:                1987,
			expectError:         false,
			expectedMarshalText: "1987-01-01",
			expectedJSON:        "\"1987-01-01\"",
		},
		"feb_31": {
			day:                 31,
			month:               2,
			year:                2024,
			expectError:         false, // will normalize
			expectedMarshalText: "2024-03-02",
			expectedJSON:        "\"2024-03-02\"",
		},
		"invalid_date_invalid_past_year": {
			day:                  1,
			month:                1,
			year:                 1787,
			expectError:          true,
			expectedErrorMessage: "the year must be > 1800",
		},
		"invalid_date_invalid_future_year": {
			day:                  1,
			month:                1,
			year:                 2900,
			expectError:          true,
			expectedErrorMessage: "too far in the future mate, is this software still running?",
		},
		"invalid_date_invalid_low_month": {
			day:                  1,
			month:                0,
			year:                 2021,
			expectError:          true,
			expectedErrorMessage: "the month cannot be < 1",
		},
		"invalid_date_invalid_high_month": {
			day:                  1,
			month:                13,
			year:                 2020,
			expectError:          true,
			expectedErrorMessage: "the month cannot be > 12",
		},
		"invalid_date_invalid_low_day": {
			day:                  0,
			month:                1,
			year:                 2020,
			expectError:          true,
			expectedErrorMessage: "the day cannot be < 1",
		},
		"invalid_date_invalid_high_day": {
			day:                  40,
			month:                1,
			year:                 2020,
			expectError:          true,
			expectedErrorMessage: "the day cannot be > 31",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			d, err := NewDate(tc.day, tc.month, tc.year)
			if tc.expectError {
				assert.NotNil(t, err)
				assert.Nil(t, d)
				assert.Equal(t, tc.expectedErrorMessage, err.Error())
			}
			if !tc.expectError {
				assert.Nil(t, err)
				assert.NotNil(t, d)

				// text marshalling and unmarshalling
				textBs, err := d.MarshalText()
				assert.Nil(t, err)
				assert.NotNil(t, textBs)
				assert.Equal(t, tc.expectedMarshalText, string(textBs))

				newDate := Date{}
				err = newDate.UnmarshalText(textBs)
				assert.Nil(t, err)

				newMarshal, err := newDate.MarshalText()
				assert.Nil(t, err)
				assert.Equal(t, string(textBs), string(newMarshal))

				// JSON Marshalling and unmarshalling
				jsonBs, err := d.MarshalJSON()
				assert.Nil(t, err)
				assert.NotNil(t, jsonBs)
				assert.Equal(t, tc.expectedJSON, string(jsonBs))

				otherDate := Date{}
				err = otherDate.UnmarshalJSON(jsonBs)
				assert.Nil(t, err)

				// GQL marshalling and unmarshalling
				gqlBs := []byte{}
				w := bytes.NewBuffer(gqlBs)
				d.MarshalGQL(w)
				gql := w.String()
				assert.Equal(t, string(jsonBs), gql)

				gqlDate := &Date{}
				err = gqlDate.UnmarshalGQL(gql)
				assert.Nil(t, err)
			}
		})
	}
}

func TestDate_String(t *testing.T) {
	type fields struct {
		Year  int
		Month int
		Day   int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "valid date",
			fields: fields{
				Year:  2020,
				Month: 5,
				Day:   31,
			},
			want: "May 31, 2020",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Date{
				Year:  tt.fields.Year,
				Month: tt.fields.Month,
				Day:   tt.fields.Day,
			}
			if got := d.String(); got != tt.want {
				t.Errorf("Date.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecimal_UnmarshalGQL(t *testing.T) {
	dec, err := decimal.NewFromString("3.14")
	assert.Nil(t, err)

	sc := Decimal(dec)

	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		sc      *Decimal
		args    args
		wantErr bool
	}{
		{
			name:    "good case",
			sc:      &sc,
			args:    args{v: "99.9"},
			wantErr: false,
		},
		{
			name:    "bad case",
			sc:      &sc,
			args:    args{v: dec},
			wantErr: true,
		},
		{
			name:    "invalid decimal",
			sc:      &sc,
			args:    args{v: "not a valid decimal"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.sc.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Decimal.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)

				if err == nil {
					newValString := sc.String()
					inputString, ok := tt.args.v.(string)
					assert.True(t, ok)
					assert.Equal(t, inputString, newValString)
				}
			}
		})
	}
}

func TestDecimal_MarshalGQL(t *testing.T) {
	dec, err := decimal.NewFromString("3.14")
	assert.Nil(t, err)

	sc := Decimal(dec)
	tests := []struct {
		name  string
		sc    Decimal
		wantW string
	}{
		{
			name:  "good case",
			sc:    sc,
			wantW: `"` + sc.String() + `"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.sc.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Decimal.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestURL_UnmarshalGQL(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    URL
		wantErr bool
	}{
		{
			name:    "good case-valid url",
			args:    args{v: "http://www.example.com/index.html"},
			want:    URL("http://www.example.com/index.html"),
			wantErr: false,
		},
		{
			name:    "sad case-unmarshal non-string input",
			args:    args{v: 119.12},
			want:    URL(""),
			wantErr: true,
		},
		{
			name:    "sad case-invalid URL",
			args:    args{v: "not a link"},
			want:    URL(""),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var link URL
			err := link.UnmarshalGQL(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("URL.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr == false {
				expected, ok := tt.args.v.(string)
				if ok {
					if string(link) != expected {
						t.Errorf("URL error, expected %s, got %s", expected, string(link))
					}
				}

			}
			if link != tt.want {
				t.Errorf("expected unmarshaled URL to be %v, got %v", tt.want, link)
			}
		})
	}
}

func TestURL_MarshalGQL(t *testing.T) {
	goodLink := "www.example.com/index.html"
	badLink := "this is not a link"
	errorMessage := "this is not a link is not a valid URL"

	tests := []struct {
		name  string
		sc    URL
		wantW string
	}{
		{
			name:  "good case",
			sc:    URL(goodLink),
			wantW: `"` + goodLink + `"`,
		},
		{
			name:  "invalid URL",
			sc:    URL(badLink),
			wantW: `"` + errorMessage + `"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.sc.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("URL.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestDateTime_Time(t *testing.T) {
	timeStr := "2020-04-03T12:10:34+03:00"
	validTime, err := time.Parse(dateTimeFormatLayout, timeStr)
	assert.Nil(t, err)

	tests := []struct {
		name string
		sc   DateTime
		want time.Time
	}{
		{
			name: "good case",
			sc:   DateTime(timeStr),
			want: validTime,
		},
		{
			name: "bad case",
			sc:   DateTime("this is not a valid date string"),
			want: time.Unix(0, 0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sc.Time(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DateTime.Time() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_MarshalText(t *testing.T) {
	type fields struct {
		Year  int
		Month int
		Day   int
	}
	tests := []struct {
		name     string
		fields   fields
		wantText []byte
		wantErr  bool
	}{
		{
			name: "invalid date",
			fields: fields{
				Year:  1600,
				Month: 13,
				Day:   32,
			},
			wantText: nil,
			wantErr:  true,
		},
		{
			name: "valid date",
			fields: fields{
				Year:  1920,
				Month: 11,
				Day:   22,
			},
			wantText: []byte("1920-11-22"),
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Date{
				Year:  tt.fields.Year,
				Month: tt.fields.Month,
				Day:   tt.fields.Day,
			}
			gotText, err := d.MarshalText()
			if (err != nil) != tt.wantErr {
				t.Errorf("Date.MarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotText, tt.wantText) {
				t.Errorf("Date.MarshalText() = %v, want %v", gotText, tt.wantText)
			}
		})
	}
}

func TestDate_UnmarshalText(t *testing.T) {
	type fields struct {
		Year  int
		Month int
		Day   int
	}
	type args struct {
		text []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid date",
			fields: fields{
				Year:  1905,
				Month: 1,
				Day:   1,
			},
			args: args{
				text: []byte("1939-04-09"),
			},
			wantErr: false,
		},
		{
			name: "invalid layout",
			fields: fields{
				Year:  1905,
				Month: 1,
				Day:   1,
			},
			args: args{
				text: []byte("this is not a valid date layout"),
			},
			wantErr: true,
		},
		{
			name: "valid layout with invalid values",
			fields: fields{
				Year:  1905,
				Month: 1,
				Day:   1,
			},
			args: args{
				text: []byte("1639-04-09"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Date{
				Year:  tt.fields.Year,
				Month: tt.fields.Month,
				Day:   tt.fields.Day,
			}
			if err := d.UnmarshalText(tt.args.text); (err != nil) != tt.wantErr {
				t.Errorf("Date.UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDate_UnmarshalJSON(t *testing.T) {
	type fields struct {
		Year  int
		Month int
		Day   int
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid date",
			fields: fields{
				Year:  1905,
				Month: 1,
				Day:   1,
			},
			args: args{
				data: []byte("1939-04-09"),
			},
			wantErr: false,
		},
		{
			name: "invalid layout",
			fields: fields{
				Year:  1905,
				Month: 1,
				Day:   1,
			},
			args: args{
				data: []byte("this is not a valid date layout"),
			},
			wantErr: true,
		},
		{
			name: "valid layout with invalid values",
			fields: fields{
				Year:  1905,
				Month: 1,
				Day:   1,
			},
			args: args{
				data: []byte("1639-04-09"),
			},
			wantErr: true,
		},
		{
			name: "null special case",
			fields: fields{
				Year:  1905,
				Month: 1,
				Day:   1,
			},
			args: args{
				data: []byte("null"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Date{
				Year:  tt.fields.Year,
				Month: tt.fields.Month,
				Day:   tt.fields.Day,
			}
			if err := d.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Date.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDate_MarshalJSON(t *testing.T) {
	type fields struct {
		Year  int
		Month int
		Day   int
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "valid date",
			fields: fields{
				Year:  1905,
				Month: 1,
				Day:   1,
			},
			want:    []byte(strconv.Quote("1905-01-01")),
			wantErr: false,
		},
		{
			name: "invalid values",
			fields: fields{
				Year:  1665,
				Month: 1,
				Day:   1,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Date{
				Year:  tt.fields.Year,
				Month: tt.fields.Month,
				Day:   tt.fields.Day,
			}
			got, err := d.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Date.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Date.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_MarshalGQL(t *testing.T) {
	type fields struct {
		Year  int
		Month int
		Day   int
	}
	tests := []struct {
		name   string
		fields fields
		wantW  string
	}{
		{
			name: "valid date",
			fields: fields{
				Year:  1905,
				Month: 1,
				Day:   1,
			},
			wantW: strconv.Quote("1905-01-01"),
		},
		{
			name: "invalid values",
			fields: fields{
				Year:  1665,
				Month: 1,
				Day:   1,
			},
			wantW: "can't marshal date, error: the year must be > 1800",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Date{
				Year:  tt.fields.Year,
				Month: tt.fields.Month,
				Day:   tt.fields.Day,
			}
			w := &bytes.Buffer{}
			d.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Date.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestDate_UnmarshalGQL(t *testing.T) {
	type fields struct {
		Year  int
		Month int
		Day   int
	}
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid date",
			fields: fields{
				Year:  1905,
				Month: 1,
				Day:   1,
			},
			args: args{
				v: "1939-04-09",
			},
			wantErr: false,
		},
		{
			name: "invalid layout",
			fields: fields{
				Year:  1905,
				Month: 1,
				Day:   1,
			},
			args: args{
				v: "this is not a valid date layout",
			},
			wantErr: true,
		},
		{
			name: "valid layout with invalid values",
			fields: fields{
				Year:  1905,
				Month: 1,
				Day:   1,
			},
			args: args{
				v: "1639-04-09",
			},
			wantErr: true,
		},
		{
			name: "invalid non string value",
			fields: fields{
				Year:  1905,
				Month: 1,
				Day:   1,
			},
			args: args{
				v: 6978,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Date{
				Year:  tt.fields.Year,
				Month: tt.fields.Month,
				Day:   tt.fields.Day,
			}
			if err := d.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Date.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBase64Binary_UnmarshalGQL(t *testing.T) {
	var sc Base64Binary
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		sc      *Base64Binary
		args    args
		wantErr bool
	}{
		{
			name: "good case",
			sc:   &sc,
			args: args{
				v: "Y2FydG9vbg==",
			},
			wantErr: false,
		},
		{
			name: "non string input",
			sc:   &sc,
			args: args{
				v: 6452,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.sc.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Base64Binary.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBase64Binary_MarshalGQL(t *testing.T) {
	sc := Base64Binary("Y2FydG9vbg==")
	tests := []struct {
		name  string
		sc    Base64Binary
		wantW string
	}{
		{
			name:  "good case",
			sc:    sc,
			wantW: strconv.Quote("Y2FydG9vbg=="),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.sc.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Base64Binary.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestCanonical_UnmarshalGQL(t *testing.T) {
	var sc Canonical
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		sc      *Canonical
		args    args
		wantErr bool
	}{
		{
			name: "good case",
			sc:   &sc,
			args: args{
				v: "http://hl7.org/fhir/ValueSet/my-valueset|0.8",
			},
			wantErr: false,
		},
		{
			name: "non string input",
			sc:   &sc,
			args: args{
				v: 897987,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.sc.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Canonical.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCanonical_MarshalGQL(t *testing.T) {
	sc := Canonical("http://hl7.org/fhir/ValueSet/my-valueset|0.8")
	tests := []struct {
		name  string
		sc    Canonical
		wantW string
	}{
		{
			name:  "good case",
			sc:    sc,
			wantW: strconv.Quote("http://hl7.org/fhir/ValueSet/my-valueset|0.8"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.sc.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Canonical.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestCode_UnmarshalGQL(t *testing.T) {
	var sc Code

	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		sc      *Code
		args    args
		wantErr bool
	}{
		{
			name: "good case",
			sc:   &sc,
			args: args{
				v: "J10",
			},
			wantErr: false,
		},
		{
			name: "non string input",
			sc:   &sc,
			args: args{
				v: 897987,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.sc.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Code.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCode_MarshalGQL(t *testing.T) {
	sc := Code("J10")
	tests := []struct {
		name  string
		sc    Code
		wantW string
	}{
		{
			name:  "good case",
			sc:    sc,
			wantW: strconv.Quote("J10"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.sc.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Code.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestDateTime_UnmarshalGQL(t *testing.T) {
	var sc DateTime
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		sc      *DateTime
		args    args
		wantErr bool
	}{
		{
			name: "good case",
			sc:   &sc,
			args: args{
				v: "2020-01-01",
			},
			wantErr: false,
		},
		{
			name: "wrong input type",
			sc:   &sc,
			args: args{
				v: 879798,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.sc.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("DateTime.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDateTime_MarshalGQL(t *testing.T) {
	sc := DateTime("2020-01-01")
	tests := []struct {
		name  string
		sc    DateTime
		wantW string
	}{
		{
			name:  "good case",
			sc:    sc,
			wantW: strconv.Quote("2020-01-01"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.sc.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("DateTime.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestInstant_UnmarshalGQL(t *testing.T) {
	var sc Instant
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		sc      *Instant
		args    args
		wantErr bool
	}{
		{
			name: "good case",
			sc:   &sc,
			args: args{
				v: "2020-01-01",
			},
			wantErr: false,
		},
		{
			name: "wrong input type",
			sc:   &sc,
			args: args{
				v: 879798,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.sc.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Instant.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInstant_MarshalGQL(t *testing.T) {
	sc := Instant("2020-01-01")
	tests := []struct {
		name  string
		sc    Instant
		wantW string
	}{
		{
			name:  "good case",
			sc:    sc,
			wantW: strconv.Quote("2020-01-01"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.sc.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Instant.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestMarkdown_UnmarshalGQL(t *testing.T) {
	var sc Markdown
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		sc      *Markdown
		args    args
		wantErr bool
	}{
		{
			name: "good case",
			sc:   &sc,
			args: args{
				v: "this is valid Markdown",
			},
			wantErr: false,
		},
		{
			name: "wrong input type",
			sc:   &sc,
			args: args{
				v: 879798,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.sc.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Markdown.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMarkdown_MarshalGQL(t *testing.T) {
	sc := Markdown("this is Markdown")
	tests := []struct {
		name  string
		sc    Markdown
		wantW string
	}{
		{
			name:  "good case",
			sc:    sc,
			wantW: strconv.Quote("this is Markdown"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.sc.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Markdown.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestOID_UnmarshalGQL(t *testing.T) {
	var sc OID
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		sc      *OID
		args    args
		wantErr bool
	}{
		{
			name: "good case",
			sc:   &sc,
			args: args{
				v: "oid:an-oid",
			},
			wantErr: false,
		},
		{
			name: "wrong input type",
			sc:   &sc,
			args: args{
				v: 879798,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.sc.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("OID.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOID_MarshalGQL(t *testing.T) {
	sc := OID("oid:an-oid")
	tests := []struct {
		name  string
		sc    OID
		wantW string
	}{
		{
			name:  "good case",
			sc:    sc,
			wantW: strconv.Quote("oid:an-oid"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.sc.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("OID.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestURI_UnmarshalGQL(t *testing.T) {
	var sc URI
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		sc      *URI
		args    args
		wantErr bool
	}{
		{
			name: "good case",
			sc:   &sc,
			args: args{
				v: "ftp://a.b.c",
			},
			wantErr: false,
		},
		{
			name: "wrong input type",
			sc:   &sc,
			args: args{
				v: 879798,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.sc.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("URI.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestURI_MarshalGQL(t *testing.T) {
	sc := URI("ftp://a.b.c")
	tests := []struct {
		name  string
		sc    URI
		wantW string
	}{
		{
			name:  "good case",
			sc:    sc,
			wantW: strconv.Quote("ftp://a.b.c"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.sc.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("URI.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestUUID_UnmarshalGQL(t *testing.T) {
	var sc UUID
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		sc      *UUID
		args    args
		wantErr bool
	}{
		{
			name: "good case",
			sc:   &sc,
			args: args{
				v: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "wrong input type",
			sc:   &sc,
			args: args{
				v: 879798,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.sc.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("UUID.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUUID_MarshalGQL(t *testing.T) {
	randomUUID := uuid.New().String()
	sc := UUID(randomUUID)
	tests := []struct {
		name  string
		sc    UUID
		wantW string
	}{
		{
			name:  "good case",
			sc:    sc,
			wantW: strconv.Quote(randomUUID),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.sc.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("UUID.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestXHTML_UnmarshalGQL(t *testing.T) {
	var sc XHTML
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		sc      *XHTML
		args    args
		wantErr bool
	}{
		{
			name: "good case",
			sc:   &sc,
			args: args{
				v: "<p>some fragment</p>",
			},
			wantErr: false,
		},
		{
			name: "wrong input type",
			sc:   &sc,
			args: args{
				v: 879798,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.sc.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("XHTML.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestXHTML_MarshalGQL(t *testing.T) {
	fragment := "<p>a fragment</p>"
	sc := XHTML(fragment)
	tests := []struct {
		name  string
		sc    XHTML
		wantW string
	}{
		{
			name:  "good case",
			sc:    sc,
			wantW: strconv.Quote(fragment),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.sc.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("XHTML.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestDecimal_String(t *testing.T) {
	dec := decimal.NewFromFloat(1.23)
	sc := Decimal(dec)
	tests := []struct {
		name string
		sc   *Decimal
		want string
	}{
		{
			name: "happy case",
			sc:   &sc,
			want: "1.23",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sc.String(); got != tt.want {
				t.Errorf("Decimal.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecimal_Decimal(t *testing.T) {
	dec := decimal.NewFromFloat(1.23)
	sc := Decimal(dec)
	tests := []struct {
		name string
		sc   *Decimal
		want decimal.Decimal
	}{
		{
			name: "happy case",
			sc:   &sc,
			want: dec,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sc.Decimal(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decimal.Decimal() = %v, want %v", got, tt.want)
			}
		})
	}
}
