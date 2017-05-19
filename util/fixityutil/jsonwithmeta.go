package fixityutil

import "github.com/leeola/fixity"

// JsonHashWithMetaFields creates JsonHashWithMeta, assigning index fields.
func JsonHashWithMetaFields(fields fixity.Fields) fixity.JsonHashWithMeta {
	return fixity.JsonHashWithMeta{
		JsonWithMeta: fixity.JsonWithMeta{JsonMeta: &fixity.JsonMeta{
			IndexedFields: fields,
		}},
	}
}
