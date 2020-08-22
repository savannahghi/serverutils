package base

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/stretchr/testify/assert"
)

// Dummy is a test node
type Dummy struct{}

func (d Dummy) IsNode() {}

func (d Dummy) GetID() ID {
	return IDValue("dummy id")
}

func (d Dummy) SetID(string) {}

func TestMain(m *testing.M) {
	os.Setenv("ROOT_COLLECTION_SUFFIX", "staging")
	m.Run()
}

func TestSuffixCollection_staging(t *testing.T) {
	col := "otp"
	expect := fmt.Sprintf("%v_%v", col, "staging")
	s := SuffixCollection(col)
	assert.Equal(t, expect, s)
}

func TestServiceSuffixCollection_testing(t *testing.T) {
	col := "otp"
	expect := fmt.Sprintf("%v_%v", col, "staging")
	s := SuffixCollection(col)
	assert.Equal(t, expect, s)
}

func Test_getCollectionName(t *testing.T) {
	n1 := &Dummy{}
	assert.Equal(t, "dummy_staging", GetCollectionName(n1))
}

func Test_validatePaginationParameters(t *testing.T) {
	first := 10
	after := "30"
	last := 10
	before := "20"

	tests := map[string]struct {
		pagination           *PaginationInput
		expectError          bool
		expectedErrorMessage string
	}{
		"first_last_specified": {
			pagination: &PaginationInput{
				First: first,
				Last:  last,
			},
			expectError:          true,
			expectedErrorMessage: "if `first` is specified for pagination, `last` cannot be specified",
		},
		"first_only": {
			pagination: &PaginationInput{
				First: first,
			},
			expectError: false,
		},
		"last_only": {
			pagination: &PaginationInput{
				Last: last,
			},
			expectError: false,
		},
		"first_and_after": {
			pagination: &PaginationInput{
				First: first,
				After: after,
			},
			expectError: false,
		},
		"last_and_before": {
			pagination: &PaginationInput{
				Last:   last,
				Before: before,
			},
			expectError: false,
		},
		"nil_pagination": {
			pagination:  nil,
			expectError: false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := ValidatePaginationParameters(tc.pagination)
			if tc.expectError {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrorMessage)
			}
			if !tc.expectError {
				assert.Nil(t, err)
			}
		})
	}
}

func Test_opstring(t *testing.T) {
	tests := map[string]struct {
		op                   Operation
		expectedOutput       string
		expectError          bool
		expectedErrorMessage string
	}{
		"invalid_operation": {
			op:                   Operation("invalid unknown operation"),
			expectedOutput:       "",
			expectError:          true,
			expectedErrorMessage: "unknown operation; did you forget to update this function after adding new operations in the schema?",
		},
		"less than": {
			op:             OperationLessThan,
			expectedOutput: "<",
			expectError:    false,
		},
		"less than_or_equal_to": {
			op:             OperationLessThanOrEqualTo,
			expectedOutput: "<=",
			expectError:    false,
		},
		"equal_to": {
			op:             OperationEqual,
			expectedOutput: "==",
			expectError:    false,
		},
		"greater_than": {
			op:             OperationGreaterThan,
			expectedOutput: ">",
			expectError:    false,
		},
		"greater_than_or_equal_to": {
			op:             OperationGreaterThanOrEqualTo,
			expectedOutput: ">=",
			expectError:    false,
		},
		"in": {
			op:             OperationIn,
			expectedOutput: "in",
			expectError:    false,
		},
		"contains": {
			op:             OperationContains,
			expectedOutput: "array-contains",
			expectError:    false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			opString, err := OpString(tc.op)
			assert.Equal(t, tc.expectedOutput, opString)
			if tc.expectError {
				assert.Equal(t, tc.expectedErrorMessage, err.Error())
			}
		})
	}
}

func TestGetFirestoreClient(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "good case",
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFirestoreClient(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFirestoreClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotNil(t, got)
		})
	}
}

func TestComposeUnpaginatedQuery(t *testing.T) {
	ctx := context.Background()
	node := &Model{}
	sortAsc := SortInput{
		SortBy: []*SortParam{
			{
				FieldName: "name",
				SortOrder: SortOrderAsc,
			},
		},
	}
	sortDesc := SortInput{
		SortBy: []*SortParam{
			{
				FieldName: "name",
				SortOrder: SortOrderDesc,
			},
		},
	}
	invalidFilter := FilterInput{
		FilterBy: []*FilterParam{
			{
				FieldName:           "name",
				FieldType:           FieldTypeString,
				ComparisonOperation: Operation("not a valid operation"),
				FieldValue:          "val",
			},
		},
	}
	booleanFilter := FilterInput{
		FilterBy: []*FilterParam{
			{
				FieldName:           "deleted",
				FieldType:           FieldTypeBoolean,
				ComparisonOperation: OperationEqual,
				FieldValue:          "false",
			},
		},
	}
	invalidBoolFilterWrongType := FilterInput{
		FilterBy: []*FilterParam{
			{
				FieldName:           "deleted",
				FieldType:           FieldTypeBoolean,
				ComparisonOperation: OperationEqual,
				FieldValue:          false,
			},
		},
	}
	invalidBoolFilterUnparseableString := FilterInput{
		FilterBy: []*FilterParam{
			{
				FieldName:           "deleted",
				FieldType:           FieldTypeBoolean,
				ComparisonOperation: OperationEqual,
				FieldValue:          "bad format",
			},
		},
	}
	intFilter := FilterInput{
		FilterBy: []*FilterParam{
			{
				FieldName:           "count",
				FieldType:           FieldTypeInteger,
				ComparisonOperation: OperationGreaterThan,
				FieldValue:          0,
			},
		},
	}
	invalidIntFilter := FilterInput{
		FilterBy: []*FilterParam{
			{
				FieldName:           "count",
				FieldType:           FieldTypeInteger,
				ComparisonOperation: OperationGreaterThan,
				FieldValue:          "not a valid int",
			},
		},
	}
	timestampFilter := FilterInput{
		FilterBy: []*FilterParam{
			{
				FieldName:           "updated",
				FieldType:           FieldTypeTimestamp,
				ComparisonOperation: OperationGreaterThan,
				FieldValue:          time.Now(),
			},
		},
	}
	numberFilter := FilterInput{
		FilterBy: []*FilterParam{
			{
				FieldName:           "numfield",
				FieldType:           FieldTypeNumber,
				ComparisonOperation: OperationLessThan,
				FieldValue:          1.0,
			},
		},
	}
	stringFilter := FilterInput{
		FilterBy: []*FilterParam{
			{
				FieldName:           "name",
				FieldType:           FieldTypeString,
				ComparisonOperation: OperationEqual,
				FieldValue:          "a string",
			},
		},
	}

	type args struct {
		ctx    context.Context
		filter *FilterInput
		sort   *SortInput
		node   Node
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil sort and filter",
			args: args{
				ctx:    ctx,
				filter: nil,
				sort:   nil,
				node:   node,
			},
			wantErr: false,
		},
		{
			name: "ascending sort",
			args: args{
				ctx:    ctx,
				filter: nil,
				sort:   &sortAsc,
				node:   node,
			},
			wantErr: false,
		},
		{
			name: "descending sort",
			args: args{
				ctx:    ctx,
				filter: nil,
				sort:   &sortDesc,
				node:   node,
			},
			wantErr: false,
		},
		{
			name: "invalid filter",
			args: args{
				ctx:    ctx,
				filter: &invalidFilter,
				sort:   &sortDesc,
				node:   node,
			},
			wantErr: true,
		},
		{
			name: "valid boolean filter",
			args: args{
				ctx:    ctx,
				filter: &booleanFilter,
				sort:   &sortAsc,
				node:   node,
			},
			wantErr: false,
		},
		{
			name: "invalid boolean filter - wrong type",
			args: args{
				ctx:    ctx,
				filter: &invalidBoolFilterWrongType,
				sort:   &sortAsc,
				node:   node,
			},
			wantErr: true,
		},
		{
			name: "invalid boolean filter - unparseable string",
			args: args{
				ctx:    ctx,
				filter: &invalidBoolFilterUnparseableString,
				sort:   &sortAsc,
				node:   node,
			},
			wantErr: true,
		},
		{
			name: "valid integer filter",
			args: args{
				ctx:    ctx,
				filter: &intFilter,
				sort:   &sortAsc,
				node:   node,
			},
			wantErr: false,
		},
		{
			name: "invalid integer filter",
			args: args{
				ctx:    ctx,
				filter: &invalidIntFilter,
				sort:   &sortAsc,
				node:   node,
			},
			wantErr: true,
		},
		{
			name: "valid timestamp filter",
			args: args{
				ctx:    ctx,
				filter: &timestampFilter,
				sort:   &sortAsc,
				node:   node,
			},
			wantErr: false,
		},
		{
			name: "valid number filter",
			args: args{
				ctx:    ctx,
				filter: &numberFilter,
				sort:   &sortAsc,
				node:   node,
			},
			wantErr: false,
		},
		{
			name: "valid string filter",
			args: args{
				ctx:    ctx,
				filter: &stringFilter,
				sort:   &sortAsc,
				node:   node,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ComposeUnpaginatedQuery(tt.args.ctx, tt.args.filter, tt.args.sort, tt.args.node)
			if (err != nil) != tt.wantErr {
				t.Errorf("ComposeUnpaginatedQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})
	}
}

func TestCreateNode(t *testing.T) {
	type args struct {
		ctx  context.Context
		node Node
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "good case",
			args: args{
				ctx:  context.Background(),
				node: &Model{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, timestamp, err := CreateNode(tt.args.ctx, tt.args.node)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotZero(t, id)
			assert.NotZero(t, timestamp)
		})
	}
}

func TestUpdateNode(t *testing.T) {
	ctx := context.Background()
	node := &Model{}
	id, _, err := CreateNode(ctx, node)
	assert.Nil(t, err)
	assert.NotZero(t, id)

	node.Name = "updated" // the update that we are testing

	type args struct {
		ctx  context.Context
		id   string
		node Node
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "good case",
			args: args{
				ctx:  ctx,
				id:   id,
				node: node,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UpdateNode(tt.args.ctx, tt.args.id, tt.args.node)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotZero(t, got)
		})
	}
}

func TestRetrieveNode(t *testing.T) {
	ctx := context.Background()
	node := &Model{}
	id, _, err := CreateNode(ctx, node)
	assert.Nil(t, err)
	assert.NotZero(t, id)

	type args struct {
		ctx  context.Context
		id   string
		node Node
	}
	tests := []struct {
		name    string
		args    args
		want    Node
		wantErr bool
	}{
		{
			name: "good case",
			args: args{
				ctx:  ctx,
				id:   id,
				node: node,
			},
			want:    node,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RetrieveNode(tt.args.ctx, tt.args.id, tt.args.node)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RetrieveNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryNodes(t *testing.T) {
	ctx := context.Background()
	node := &Model{
		Name:        "test model instance",
		Description: "this is a test description",
		Deleted:     false,
	}
	id, _, err := CreateNode(ctx, node)
	assert.Nil(t, err)
	assert.NotZero(t, id)

	sortAsc := SortInput{
		SortBy: []*SortParam{
			{
				FieldName: "name",
				SortOrder: SortOrderAsc,
			},
		},
	}

	type args struct {
		ctx        context.Context
		pagination *PaginationInput
		filter     *FilterInput
		sort       *SortInput
		node       Node
	}
	tests := []struct {
		name    string
		args    args
		want    []*firestore.DocumentSnapshot
		wantErr bool
	}{
		{
			name: "no pagination, filter or sort",
			args: args{
				ctx:        ctx,
				pagination: nil,
				filter:     nil,
				sort:       nil,
				node:       &Model{},
			},
			wantErr: false,
		},
		{
			name: "with pagination, first",
			args: args{
				ctx: ctx,
				pagination: &PaginationInput{
					First: 10,
					After: id,
				},
				filter: nil,
				sort:   &sortAsc,
				node:   &Model{},
			},
			wantErr: false,
		},
		{
			name: "with pagination, last",
			args: args{
				ctx: ctx,
				pagination: &PaginationInput{
					Last:   1,
					Before: id,
				},
				filter: nil,
				sort:   &sortAsc,
				node:   &Model{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			snapshots, pageInfo, err := QueryNodes(tt.args.ctx, tt.args.pagination, tt.args.filter, tt.args.sort, tt.args.node)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryNodes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, pageInfo)
				if pageInfo.StartCursor != nil {
					assert.NotNil(t, snapshots)
				}
			}
		})
	}
}

func TestDeleteNode(t *testing.T) {
	ctx := context.Background()
	node := &Model{
		Name:        "test model instance",
		Description: "this is a test description",
		Deleted:     false,
	}
	id, _, err := CreateNode(ctx, node)
	assert.Nil(t, err)
	assert.NotZero(t, id)

	type args struct {
		ctx  context.Context
		id   string
		node Node
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "existing node",
			args: args{
				ctx:  ctx,
				id:   node.GetID().String(),
				node: &Model{},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "non existent node",
			args: args{
				ctx:  ctx,
				id:   "this should not exist",
				node: &Model{},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeleteNode(tt.args.ctx, tt.args.id, tt.args.node)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DeleteNode() = %v, want %v", got, tt.want)
			}
		})
	}
}
