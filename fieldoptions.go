package kala

const (
	FOKeyFullTextSearch = "FullTextSearch"
)

// FieldOptions optionally store index configuration for the specific field.
//
// It's up to the Index implementor to use these options. Because FieldOptions
// are *optional*, the implementor should gracefully ignore the majority of
// these options if they are not used or supported properly.
//
// NOTE: These options should remain as generic as possible so that they are
// not too bound to a specific Index implementor. Being too bound to the index
// implementor effectively ties the data store to only a single Index.
//
// This may or may not be problematic to specific use cases, it's just worth
// trying to avoid.
//
// To aid this, common field options are defined below.
type FieldOptions map[string]interface{}

// FullTextSearch enable full text search for this specific field.
//
// Note that many Indexes do not implement this, and should ignore it gracefully
// if not implemented.
func (f FieldOptions) FullTextSearch() FieldOptions {
	f[FOKeyFullTextSearch] = true
	return f
}
