package modules

import "github.com/zmap/zdns/pkg/modules/nslookup"

func init() {
	nslookup.RegisterLookup()
}
