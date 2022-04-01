package cli

import (
	"sort"
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

	sort.Strings(names)
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
	modSet.AddModule("AFSDB", TypedModule{
		zdns.RawModule{},
		dns.TypeAFSDB,
	})
	modSet.AddModule("ATMA", TypedModule{
		zdns.RawModule{},
		dns.TypeATMA,
	})
	modSet.AddModule("AVC", TypedModule{
		zdns.RawModule{},
		dns.TypeAVC,
	})
	modSet.AddModule("CAA", TypedModule{
		zdns.RawModule{},
		dns.TypeCAA,
	})
	modSet.AddModule("CERT", TypedModule{
		zdns.RawModule{},
		dns.TypeCERT,
	})
	modSet.AddModule("CDS", TypedModule{
		zdns.RawModule{},
		dns.TypeCDS,
	})
	modSet.AddModule("CDNSKEY", TypedModule{
		zdns.RawModule{},
		dns.TypeCDNSKEY,
	})
	modSet.AddModule("CNAME", TypedModule{
		zdns.RawModule{},
		dns.TypeCNAME,
	})
	modSet.AddModule("CSYNC", TypedModule{
		zdns.RawModule{},
		dns.TypeCSYNC,
	})
	modSet.AddModule("DHCID", TypedModule{
		zdns.RawModule{},
		dns.TypeDHCID,
	})
	modSet.AddModule("DNSKEY", TypedModule{
		zdns.RawModule{},
		dns.TypeDNSKEY,
	})
	modSet.AddModule("DS", TypedModule{
		zdns.RawModule{},
		dns.TypeDS,
	})
	modSet.AddModule("EID", TypedModule{
		zdns.RawModule{},
		dns.TypeEID,
	})
	modSet.AddModule("EUI48", TypedModule{
		zdns.RawModule{},
		dns.TypeEUI48,
	})
	modSet.AddModule("EUI64", TypedModule{
		zdns.RawModule{},
		dns.TypeEUI64,
	})
	modSet.AddModule("GID", TypedModule{
		zdns.RawModule{},
		dns.TypeGID,
	})
	modSet.AddModule("GPOS", TypedModule{
		zdns.RawModule{},
		dns.TypeGPOS,
	})
	modSet.AddModule("HINFO", TypedModule{
		zdns.RawModule{},
		dns.TypeHINFO,
	})
	modSet.AddModule("HIP", TypedModule{
		zdns.RawModule{},
		dns.TypeHIP,
	})
	modSet.AddModule("HTTPS", TypedModule{
		zdns.RawModule{},
		dns.TypeHTTPS,
	})
	modSet.AddModule("ISDN", TypedModule{
		zdns.RawModule{},
		dns.TypeISDN,
	})
	modSet.AddModule("KEY", TypedModule{
		zdns.RawModule{},
		dns.TypeKEY,
	})
	modSet.AddModule("KX", TypedModule{
		zdns.RawModule{},
		dns.TypeKX,
	})
	modSet.AddModule("L32", TypedModule{
		zdns.RawModule{},
		dns.TypeL32,
	})
	modSet.AddModule("L64", TypedModule{
		zdns.RawModule{},
		dns.TypeL64,
	})
	modSet.AddModule("LOC", TypedModule{
		zdns.RawModule{},
		dns.TypeLOC,
	})
	modSet.AddModule("LP", TypedModule{
		zdns.RawModule{},
		dns.TypeLP,
	})
	modSet.AddModule("MD", TypedModule{
		zdns.RawModule{},
		dns.TypeMD,
	})
	modSet.AddModule("MF", TypedModule{
		zdns.RawModule{},
		dns.TypeMF,
	})
	modSet.AddModule("MB", TypedModule{
		zdns.RawModule{},
		dns.TypeMB,
	})
	modSet.AddModule("MG", TypedModule{
		zdns.RawModule{},
		dns.TypeMG,
	})
	modSet.AddModule("MR", TypedModule{
		zdns.RawModule{},
		dns.TypeMR,
	})
	modSet.AddModule("MX", TypedModule{
		zdns.RawModule{},
		dns.TypeMX,
	})
	modSet.AddModule("NAPTR", TypedModule{
		zdns.RawModule{},
		dns.TypeNAPTR,
	})
	modSet.AddModule("NS", TypedModule{
		zdns.RawModule{},
		dns.TypeNS,
	})
	modSet.AddModule("NID", TypedModule{
		zdns.RawModule{},
		dns.TypeNID,
	})
	modSet.AddModule("NINFO", TypedModule{
		zdns.RawModule{},
		dns.TypeNINFO,
	})
	modSet.AddModule("NSAPPTR", TypedModule{
		zdns.RawModule{},
		dns.TypeNSAPPTR,
	})
	modSet.AddModule("NS", TypedModule{
		zdns.RawModule{},
		dns.TypeNS,
	})
	modSet.AddModule("NXT", TypedModule{
		zdns.RawModule{},
		dns.TypeNXT,
	})
	modSet.AddModule("NSEC", TypedModule{
		zdns.RawModule{},
		dns.TypeNSEC,
	})
	modSet.AddModule("NSEC3", TypedModule{
		zdns.RawModule{},
		dns.TypeNSEC3,
	})
	modSet.AddModule("NSEC3PARAM", TypedModule{
		zdns.RawModule{},
		dns.TypeNSEC3PARAM,
	})
	modSet.AddModule("NULL", TypedModule{
		zdns.RawModule{},
		dns.TypeNULL,
	})
	modSet.AddModule("OPENPGPKEY", TypedModule{
		zdns.RawModule{},
		dns.TypeOPENPGPKEY,
	})
	modSet.AddModule("OPT", TypedModule{
		zdns.RawModule{},
		dns.TypeOPT,
	})
	modSet.AddModule("PTR", TypedModule{
		zdns.RawModule{},
		dns.TypePTR,
	})
	modSet.AddModule("PX", TypedModule{
		zdns.RawModule{},
		dns.TypePX,
	})
	modSet.AddModule("RP", TypedModule{
		zdns.RawModule{},
		dns.TypeRP,
	})
	modSet.AddModule("RRSIG", TypedModule{
		zdns.RawModule{},
		dns.TypeRRSIG,
	})
	modSet.AddModule("RT", TypedModule{
		zdns.RawModule{},
		dns.TypeRT,
	})
	modSet.AddModule("SMIMEA", TypedModule{
		zdns.RawModule{},
		dns.TypeSMIMEA,
	})
	modSet.AddModule("SSHFP", TypedModule{
		zdns.RawModule{},
		dns.TypeSSHFP,
	})
	modSet.AddModule("SOA", TypedModule{
		zdns.RawModule{},
		dns.TypeSOA,
	})
	modSet.AddModule("SPF", TypedModule{
		zdns.RawModule{},
		dns.TypeSPF,
	})
	modSet.AddModule("SRV", TypedModule{
		zdns.RawModule{},
		dns.TypeSRV,
	})
	modSet.AddModule("SVCB", TypedModule{
		zdns.RawModule{},
		dns.TypeSVCB,
	})
	modSet.AddModule("TALINK", TypedModule{
		zdns.RawModule{},
		dns.TypeTALINK,
	})
	modSet.AddModule("TKEY", TypedModule{
		zdns.RawModule{},
		dns.TypeTKEY,
	})
	modSet.AddModule("TLSA", TypedModule{
		zdns.RawModule{},
		dns.TypeTLSA,
	})
	modSet.AddModule("TXT", TypedModule{
		zdns.RawModule{},
		dns.TypeTXT,
	})
	modSet.AddModule("UID", TypedModule{
		zdns.RawModule{},
		dns.TypeUID,
	})
	modSet.AddModule("UINFO", TypedModule{
		zdns.RawModule{},
		dns.TypeUINFO,
	})
	modSet.AddModule("UNSPEC", TypedModule{
		zdns.RawModule{},
		dns.TypeUNSPEC,
	})
	modSet.AddModule("URI", TypedModule{
		zdns.RawModule{},
		dns.TypeURI,
	})
	modSet.AddModule("ANY", TypedModule{
		zdns.RawModule{},
		dns.TypeANY,
	})

	return modSet
}
