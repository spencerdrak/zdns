package modules

import "github.com/zmap/zdns/pkg/modules/dmarc"

func init() {
	dmarc.RegisterLookup()
}
