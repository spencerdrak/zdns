package zdns

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

// FactorySet is a map of name (string) -> GlobalLookupFactory, one per module.
type FactorySet map[string]GlobalLookupFactory

// CopyInto copies the modules in s to destination. The sets will be unique, but
// the underlying ScanModule instances will be the same.
func (s FactorySet) CopyInto(destination FactorySet) {
	for name, m := range s {
		if _, ok := destination[strings.ToUpper(name)]; ok {
			log.Warnf("overwriting module %s", name)
		}
		destination[strings.ToUpper(name)] = m
	}
}

// AddModule adds m to the ModuleSet, accessible via the given name. If the name
// is already in the ModuleSet, it is overwritten.
func (s FactorySet) AddModule(name string, m GlobalLookupFactory) {
	if _, ok := s[strings.ToUpper(name)]; ok {
		log.Warnf("overwriting module %s", name)
	}
	s[strings.ToUpper(name)] = m
}

// RemoveModule removes the module at the specified name. If the name is not in
// the module set, nothing happens.
func (s FactorySet) RemoveModule(name string) {
	delete(s, strings.ToUpper(name))
}

// NewFactorySet returns an empty FactorySet.
func NewFactorySet() FactorySet {
	return make(FactorySet)
}
