package zdns

// FactorySet is a map of name (string) -> GlobalLookupFactory, one per module.
type FactorySet map[string]GlobalLookupFactory

// CopyInto copies the modules in s to destination. The sets will be unique, but
// the underlying ScanModule instances will be the same.
func (s FactorySet) CopyInto(destination FactorySet) {
	for name, m := range s {
		destination[name] = m
	}
}

// NewFactorySet returns an empty FactorySet.
func NewFactorySet() FactorySet {
	return make(FactorySet)
}
