package fixity

import "testing"

func TestBlobTypeConsts(t *testing.T) {
	testCases := []struct {
		Value       BlobType
		ExpectValue BlobType
	}{
		{
			Value:       BlobTypeSchemaless,
			ExpectValue: 0,
		},
		{
			Value:       BlobTypeParts,
			ExpectValue: 1,
		},
		{
			Value:       BlobTypeData,
			ExpectValue: 2,
		},
		{
			Value:       BlobTypeValues,
			ExpectValue: 3,
		},
		{
			Value:       BlobTypeMutation,
			ExpectValue: 4,
		},
	}
	for _, testCase := range testCases {
		if testCase.Value != testCase.ExpectValue {
			t.Errorf("%s want:%d, got:%d", testCase.Value.String(), testCase.ExpectValue, testCase.Value)
		}
	}
}
