package cli

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/zmap/dns"
	"github.com/zmap/zdns/pkg/zdns"
)

type TypedModule struct {
	Module zdns.Module
	Type   uint16
}

type ModuleSet map[string]TypedModule

func (m ModuleSet) AddModule(name string, mod TypedModule) {
	logger := log.WithFields(log.Fields{
		"Module": "cli",
	})

	name = strings.ToUpper(name)

	if _, ok := m[name]; ok {
		logger.Warnf("module %s already exists in cli moduleset and will be overwritten.")
	}

	m[name] = mod
}

func (m ModuleSet) HasModule(name string) bool {
	if _, ok := m[name]; ok {
		return true
	}
	return false
}

func (m ModuleSet) ValidModulesString() string {
	names := make([]string, 0, len(m))
	for key := range m {
		names = append(names, key)
	}

	return strings.Join(names, ",")
}

func GenerateModSet() ModuleSet {
	modSet := ModuleSet{}

	// In addition to writing the new module, devs are expected to add new modules to the CLI below.
	// The CLI should support all modules that are tracked within ZDNS itself.
	modSet.AddModule("A", TypedModule{
		zdns.RawModule{},
		dns.TypeA,
	})
	modSet.AddModule("AAAA", TypedModule{
		zdns.RawModule{},
		dns.TypeAAAA,
	})

	return modSet
}
