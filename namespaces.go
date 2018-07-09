package fixity

// Reserved namespaces for fixity implementation and metadata.
const (
	ReservedNamespaceSigners = "_signers"
	ReservedNamespaceAuthors = "_authors"
)

func IsReservedNamespace(s string) bool {
	switch s {
	case ReservedNamespaceSigners:
		return true
	case ReservedNamespaceAuthors:
		return true
	default:
		return false
	}
}
