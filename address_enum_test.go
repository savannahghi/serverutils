package base_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
)

func TestAddressType(t *testing.T) {
	type expects struct {
		isValid      bool
		canUnmarshal bool
	}

	cases := []struct {
		name        string
		args        base.AddressType
		convert     interface{}
		expectation expects
	}{
		{
			name:    "invalid_string",
			args:    "testaddres",
			convert: "testaddress",
			expectation: expects{
				isValid:      false,
				canUnmarshal: false,
			},
		},
		{
			name:    "invalid_int_convert",
			args:    "testaddres",
			convert: 101,
			expectation: expects{
				isValid:      false,
				canUnmarshal: false,
			},
		},
		{
			name:    "valid",
			args:    base.AddressTypePostal,
			convert: base.AddressTypePostal,
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
		{
			name:    "valid_no_convert",
			args:    base.AddressTypePostal,
			convert: "testaddress",
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
		{
			name:    "valid_can_convert",
			args:    base.AddressTypePostal,
			convert: base.AddressTypePostal,
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectation.isValid, tt.args.IsValid())
			assert.NotEmpty(t, tt.args.String())
			err := tt.args.UnmarshalGQL(tt.convert)
			assert.NotNil(t, err)
			tt.args.MarshalGQL(os.Stdout)

		})
	}

}

func TestAddressUse(t *testing.T) {
	type expects struct {
		isValid      bool
		canUnmarshal bool
	}

	cases := []struct {
		name        string
		args        base.AddressUse
		convert     interface{}
		expectation expects
	}{
		{
			name:    "invalid_string",
			args:    "testaddres",
			convert: "testaddress",
			expectation: expects{
				isValid:      false,
				canUnmarshal: false,
			},
		},
		{
			name:    "invalid_int_convert",
			args:    "testaddres",
			convert: 101,
			expectation: expects{
				isValid:      false,
				canUnmarshal: false,
			},
		},
		{
			name:    "valid",
			args:    base.AddressUseWork,
			convert: base.AddressUseWork,
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
		{
			name:    "valid_no_convert",
			args:    base.AddressUseWork,
			convert: "testaddress",
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
		{
			name:    "valid_can_convert",
			args:    base.AddressUseWork,
			convert: base.AddressUseWork,
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectation.isValid, tt.args.IsValid())
			assert.NotEmpty(t, tt.args.String())
			err := tt.args.UnmarshalGQL(tt.convert)
			assert.NotNil(t, err)
			tt.args.MarshalGQL(os.Stdout)

		})
	}

}

func TestCountry(t *testing.T) {
	type expects struct {
		isValid      bool
		canUnmarshal bool
	}

	cases := []struct {
		name        string
		args        base.Country
		convert     interface{}
		expectation expects
	}{
		{
			name:    "invalid_string",
			args:    "testaddres",
			convert: "testaddress",
			expectation: expects{
				isValid:      false,
				canUnmarshal: false,
			},
		},
		{
			name:    "invalid_int_convert",
			args:    "testaddres",
			convert: 101,
			expectation: expects{
				isValid:      false,
				canUnmarshal: false,
			},
		},
		{
			name:    "valid",
			args:    base.CountryBh,
			convert: base.CountryMe,
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
		{
			name:    "valid_no_convert",
			args:    base.CountryHn,
			convert: "testaddress",
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
		{
			name:    "valid_can_convert",
			args:    base.CountryHn,
			convert: base.CountryHn,
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectation.isValid, tt.args.IsValid())
			assert.NotEmpty(t, tt.args.String())
			err := tt.args.UnmarshalGQL(tt.convert)
			assert.NotNil(t, err)
			tt.args.MarshalGQL(os.Stdout)

		})
	}

}
