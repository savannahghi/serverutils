package cmd

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	log "github.com/sirupsen/logrus"
)

func Test_pushSchema(t *testing.T) {

	schema := `
	type Query {
		world: String
	  }
	`
	testDir := t.TempDir()
	schemaFile := filepath.Join(testDir, "test.graphql")
	err := ioutil.WriteFile(schemaFile, []byte(schema), 0666)
	if err != nil {
		log.Errorf("error writing to test file: %v", err)
		return
	}

	emptyTestDir := t.TempDir()

	testService := Service{
		Name:    "bewell",
		URL:     "https://bewell-test.com",
		Version: "0.0.1",
	}

	// TODO: Update test url
	testValidationURL := "https://postman-echo.com/post"

	type args struct {
		service   Service
		dir       string
		extension string
		pushURL   string
	}
	tests := []struct {
		name    string
		args    args
		want    *SchemaStatus
		wantErr bool
	}{
		{
			name: "success: push request",
			args: args{
				service:   testService,
				dir:       testDir,
				extension: "graphql",
				pushURL:   testValidationURL,
			},
			wantErr: true,
			want: &SchemaStatus{
				Valid: false,
			},
		},
		{
			name: "fail: no schema files in directory",
			args: args{
				service:   testService,
				dir:       emptyTestDir,
				extension: "graphql",
				pushURL:   testValidationURL,
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "fail: missing url",
			args: args{
				service:   testService,
				dir:       testDir,
				extension: "graphql",
				pushURL:   "",
			},
			wantErr: true,
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := publishSchema(tt.args.service, tt.args.dir, tt.args.extension, tt.args.pushURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("publishSchema() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
