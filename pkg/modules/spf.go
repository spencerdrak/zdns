package modules

import "github.com/zmap/zdns/pkg/modules/spf"

func init() {
	spf.RegisterLookup()
}
