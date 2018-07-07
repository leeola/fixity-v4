package fixity

import "testing"

func TestBlobTypeConsts(t *testing.T) {
	testCases := []struct {
		Value       BlobType
		ExpectValue BlobType
	}{
		{
			Value:       BlobTypeMutation,
			ExpectValue: 1,
		},
		{
			Value:       BlobTypeValues,
			ExpectValue: 2,
		},
		{
			Value:       BlobTypeData,
			ExpectValue: 3,
		},
		{
			Value:       BlobTypeParts,
			ExpectValue: 4,
		},
		{
			Value:       BlobTypeSchemaless,
			ExpectValue: 5,
		},
	}
	for _, testCase := range testCases {
		if testCase.Value != testCase.ExpectValue {
			t.Errorf("%s want:%d, got:%d", testCase.Value.String(), testCase.ExpectValue, testCase.Value)
		}
	}
}
