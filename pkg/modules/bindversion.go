package modules

import "github.com/zmap/zdns/pkg/modules/bindversion"

func init() {
	bindversion.RegisterLookup()
}
