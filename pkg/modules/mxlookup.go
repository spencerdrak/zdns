package modules

import "github.com/zmap/zdns/pkg/modules/mxlookup"

func init() {
	mxlookup.RegisterLookup()
}
