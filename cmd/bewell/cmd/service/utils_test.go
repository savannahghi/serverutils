package cmd

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"

	log "github.com/sirupsen/logrus"
)

func Test_readSchemaFile(t *testing.T) {
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

	type args struct {
		schemaFile string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "success:existing test file",
			args: args{
				schemaFile: schemaFile,
			},
			want:    schema,
			wantErr: false,
		},
		{
			name: "fail:missing file",
			args: args{
				schemaFile: "doesn't exist",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readSchemaFile(tt.args.schemaFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("readSchemaFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("readSchemaFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readSchemaFilesInDirectory(t *testing.T) {
	mutation := `
	type Mutation {
		world: String
	  }
	`
	query := `
	type Query {
		world: String
	  }
	`

	combinedSchema := mutation + "\n" + query + "\n"

	testDir1 := t.TempDir()

	schemaFile1 := filepath.Join(testDir1, "mutation.graphql")
	schemaFile2 := filepath.Join(testDir1, "query.graphql")

	err := ioutil.WriteFile(schemaFile1, []byte(mutation), 0666)
	if err != nil {
		log.Errorf("error writing to test file: %v", err)
		return
	}
	err = ioutil.WriteFile(schemaFile2, []byte(query), 0666)
	if err != nil {
		log.Errorf("error writing to test file: %v", err)
		return
	}

	testDir2 := t.TempDir()

	type args struct {
		dir       string
		extension string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "success:directory with schema files",
			args: args{
				dir:       testDir1,
				extension: "graphql",
			},
			want:    combinedSchema,
			wantErr: false,
		},
		{
			name: "fail:directory without schema files",
			args: args{
				dir:       testDir2,
				extension: "graphql",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "success:single schema file",
			args: args{
				dir:       schemaFile1,
				extension: "graphql",
			},
			want:    mutation + "\n",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readSchemaFilesInDirectory(tt.args.dir, tt.args.extension)
			if (err != nil) != tt.wantErr {
				t.Errorf("readSchemaFilesInDirectory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("readSchemaFilesInDirectory() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO: Change placeholder request for testing
func Test_makeRequest(t *testing.T) {
	schema := `
	type Query {
		world: String
	  }
	`
	// TODO: Update the test URL
	url := "https://postman-echo.com/post"

	type args struct {
		url  string
		body GraphqlSchemaPayload
	}
	tests := []struct {
		name    string
		args    args
		want    *http.Response
		wantErr bool
	}{
		{
			name: "success: make request",
			args: args{
				url: url,
				body: GraphqlSchemaPayload{
					Name:     "test",
					Version:  "0.0.1",
					TypeDefs: schema,
				},
			},
			wantErr: false,
		},
		{
			name: "fail: empty url in request",
			args: args{
				url: "",
				body: GraphqlSchemaPayload{
					Name:     "test",
					Version:  "0.0.1",
					TypeDefs: schema,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := makeRequest(tt.args.url, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("makeRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
