package raw

import (
	"github.com/zmap/dns"
	"github.com/zmap/zdns/pkg/miekg"
	"github.com/zmap/zdns/pkg/zdns"
)

func RegisterLookup() {
	a := new(miekg.GlobalLookupFactory)
	a.SetDNSType(dns.TypeA)
	zdns.RegisterLookup("A", a)

	aaaa := new(miekg.GlobalLookupFactory)
	aaaa.SetDNSType(dns.TypeAAAA)
	zdns.RegisterLookup("AAAA", aaaa)

	afsdb := new(miekg.GlobalLookupFactory)
	afsdb.SetDNSType(dns.TypeAFSDB)
	zdns.RegisterLookup("AFSDB", afsdb)

	atma := new(miekg.GlobalLookupFactory)
	atma.SetDNSType(dns.TypeATMA)
	zdns.RegisterLookup("ATMA", atma)

	avc := new(miekg.GlobalLookupFactory)
	avc.SetDNSType(dns.TypeAVC)
	zdns.RegisterLookup("AVC", avc)

	caa := new(miekg.GlobalLookupFactory)
	caa.SetDNSType(dns.TypeCAA)
	zdns.RegisterLookup("CAA", caa)

	cert := new(miekg.GlobalLookupFactory)
	cert.SetDNSType(dns.TypeCERT)
	zdns.RegisterLookup("CERT", cert)

	cds := new(miekg.GlobalLookupFactory)
	cds.SetDNSType(dns.TypeCDS)
	zdns.RegisterLookup("CDS", cds)

	cdnskey := new(miekg.GlobalLookupFactory)
	cdnskey.SetDNSType(dns.TypeCDNSKEY)
	zdns.RegisterLookup("CDNSKEY", cdnskey)

	cname := new(miekg.GlobalLookupFactory)
	cname.SetDNSType(dns.TypeCNAME)
	zdns.RegisterLookup("CNAME", cname)

	csync := new(miekg.GlobalLookupFactory)
	csync.SetDNSType(dns.TypeCSYNC)
	zdns.RegisterLookup("CSYNC", csync)

	dhcid := new(miekg.GlobalLookupFactory)
	dhcid.SetDNSType(dns.TypeDHCID)
	zdns.RegisterLookup("DHCID", dhcid)

	dnskey := new(miekg.GlobalLookupFactory)
	dnskey.SetDNSType(dns.TypeDNSKEY)
	zdns.RegisterLookup("DNSKEY", dnskey)

	ds := new(miekg.GlobalLookupFactory)
	ds.SetDNSType(dns.TypeDS)
	zdns.RegisterLookup("DS", ds)

	eid := new(miekg.GlobalLookupFactory)
	eid.SetDNSType(dns.TypeEID)
	zdns.RegisterLookup("EID", eid)

	eui48 := new(miekg.GlobalLookupFactory)
	eui48.SetDNSType(dns.TypeEUI48)
	zdns.RegisterLookup("EUI48", eui48)

	eui64 := new(miekg.GlobalLookupFactory)
	eui64.SetDNSType(dns.TypeEUI64)
	zdns.RegisterLookup("EUI64", eui64)

	gid := new(miekg.GlobalLookupFactory)
	gid.SetDNSType(dns.TypeGID)
	zdns.RegisterLookup("GID", gid)

	gpos := new(miekg.GlobalLookupFactory)
	gpos.SetDNSType(dns.TypeGPOS)
	zdns.RegisterLookup("GPOS", gpos)

	hinfo := new(miekg.GlobalLookupFactory)
	hinfo.SetDNSType(dns.TypeHINFO)
	zdns.RegisterLookup("HINFO", hinfo)

	hip := new(miekg.GlobalLookupFactory)
	hip.SetDNSType(dns.TypeHIP)
	zdns.RegisterLookup("HIP", hip)

	https := new(miekg.GlobalLookupFactory)
	https.SetDNSType(dns.TypeHTTPS)
	zdns.RegisterLookup("HTTPS", https)

	isdn := new(miekg.GlobalLookupFactory)
	isdn.SetDNSType(dns.TypeISDN)
	zdns.RegisterLookup("ISDN", isdn)

	key := new(miekg.GlobalLookupFactory)
	key.SetDNSType(dns.TypeKEY)
	zdns.RegisterLookup("KEY", key)

	kx := new(miekg.GlobalLookupFactory)
	kx.SetDNSType(dns.TypeKX)
	zdns.RegisterLookup("KX", kx)

	l32 := new(miekg.GlobalLookupFactory)
	l32.SetDNSType(dns.TypeL32)
	zdns.RegisterLookup("L32", l32)

	l64 := new(miekg.GlobalLookupFactory)
	l64.SetDNSType(dns.TypeL64)
	zdns.RegisterLookup("L64", l64)

	loc := new(miekg.GlobalLookupFactory)
	loc.SetDNSType(dns.TypeLOC)
	zdns.RegisterLookup("LOC", loc)

	lp := new(miekg.GlobalLookupFactory)
	lp.SetDNSType(dns.TypeLP)
	zdns.RegisterLookup("LP", lp)

	md := new(miekg.GlobalLookupFactory)
	md.SetDNSType(dns.TypeMD)
	zdns.RegisterLookup("MD", md)

	mf := new(miekg.GlobalLookupFactory)
	mf.SetDNSType(dns.TypeMF)
	zdns.RegisterLookup("MF", mf)

	mb := new(miekg.GlobalLookupFactory)
	mb.SetDNSType(dns.TypeMB)
	zdns.RegisterLookup("MB", mb)

	mg := new(miekg.GlobalLookupFactory)
	mg.SetDNSType(dns.TypeMG)
	zdns.RegisterLookup("MG", mg)

	mr := new(miekg.GlobalLookupFactory)
	mr.SetDNSType(dns.TypeMR)
	zdns.RegisterLookup("MR", mr)

	mx := new(miekg.GlobalLookupFactory)
	mx.SetDNSType(dns.TypeMX)
	zdns.RegisterLookup("MX", mx)

	naptr := new(miekg.GlobalLookupFactory)
	naptr.SetDNSType(dns.TypeNAPTR)
	zdns.RegisterLookup("NAPTR", naptr)

	nimloc := new(miekg.GlobalLookupFactory)
	nimloc.SetDNSType(dns.TypeNIMLOC)
	zdns.RegisterLookup("NS", nimloc)

	nid := new(miekg.GlobalLookupFactory)
	nid.SetDNSType(dns.TypeNID)
	zdns.RegisterLookup("NID", nid)

	ninfo := new(miekg.GlobalLookupFactory)
	ninfo.SetDNSType(dns.TypeNINFO)
	zdns.RegisterLookup("NINFO", ninfo)

	nsapptr := new(miekg.GlobalLookupFactory)
	nsapptr.SetDNSType(dns.TypeNSAPPTR)
	zdns.RegisterLookup("NSAPPTR", nsapptr)

	ns := new(miekg.GlobalLookupFactory)
	ns.SetDNSType(dns.TypeNS)
	zdns.RegisterLookup("NS", ns)

	nxt := new(miekg.GlobalLookupFactory)
	nxt.SetDNSType(dns.TypeNXT)
	zdns.RegisterLookup("NXT", nxt)

	nsec := new(miekg.GlobalLookupFactory)
	nsec.SetDNSType(dns.TypeNSEC)
	zdns.RegisterLookup("NSEC", nsec)

	nsec3 := new(miekg.GlobalLookupFactory)
	nsec3.SetDNSType(dns.TypeNSEC3)
	zdns.RegisterLookup("NSEC3", nsec3)

	nsec3param := new(miekg.GlobalLookupFactory)
	nsec3param.SetDNSType(dns.TypeNSEC3PARAM)
	zdns.RegisterLookup("NSEC3PARAM", nsec3param)

	null := new(miekg.GlobalLookupFactory)
	null.SetDNSType(dns.TypeNULL)
	zdns.RegisterLookup("NULL", null)

	openpgpkey := new(miekg.GlobalLookupFactory)
	openpgpkey.SetDNSType(dns.TypeOPENPGPKEY)
	zdns.RegisterLookup("OPENPGPKEY", openpgpkey)

	//opt := new(miekg.GlobalLookupFactory)
	//opt.SetDNSType(dns.TypeOPT)
	//zdns.RegisterLookup("OPT", opt)

	ptr := new(miekg.GlobalLookupFactory)
	ptr.SetDNSType(dns.TypePTR)
	zdns.RegisterLookup("PTR", ptr)

	px := new(miekg.GlobalLookupFactory)
	px.SetDNSType(dns.TypePX)
	zdns.RegisterLookup("PX", px)

	rp := new(miekg.GlobalLookupFactory)
	rp.SetDNSType(dns.TypeRP)
	zdns.RegisterLookup("RP", rp)

	rrsig := new(miekg.GlobalLookupFactory)
	rrsig.SetDNSType(dns.TypeRRSIG)
	zdns.RegisterLookup("RRSIG", rrsig)

	rt := new(miekg.GlobalLookupFactory)
	rt.SetDNSType(dns.TypeRT)
	zdns.RegisterLookup("RT", rt)

	smimea := new(miekg.GlobalLookupFactory)
	smimea.SetDNSType(dns.TypeSMIMEA)
	zdns.RegisterLookup("SMIMEA", smimea)

	sshfp := new(miekg.GlobalLookupFactory)
	sshfp.SetDNSType(dns.TypeSSHFP)
	zdns.RegisterLookup("SSHFP", sshfp)

	soa := new(miekg.GlobalLookupFactory)
	soa.SetDNSType(dns.TypeSOA)
	zdns.RegisterLookup("SOA", soa)

	spf := new(miekg.GlobalLookupFactory)
	spf.SetDNSType(dns.TypeSPF)
	zdns.RegisterLookup("SPF", spf)

	srv := new(miekg.GlobalLookupFactory)
	srv.SetDNSType(dns.TypeSRV)
	zdns.RegisterLookup("SRV", srv)

	svcb := new(miekg.GlobalLookupFactory)
	svcb.SetDNSType(dns.TypeSVCB)
	zdns.RegisterLookup("SVCB", svcb)

	talink := new(miekg.GlobalLookupFactory)
	talink.SetDNSType(dns.TypeTALINK)
	zdns.RegisterLookup("TALINK", talink)

	tkey := new(miekg.GlobalLookupFactory)
	tkey.SetDNSType(dns.TypeTKEY)
	zdns.RegisterLookup("TKEY", tkey)

	tlsa := new(miekg.GlobalLookupFactory)
	tlsa.SetDNSType(dns.TypeTLSA)
	zdns.RegisterLookup("TLSA", tlsa)

	txt := new(miekg.GlobalLookupFactory)
	txt.SetDNSType(dns.TypeTXT)
	zdns.RegisterLookup("TXT", txt)

	uid := new(miekg.GlobalLookupFactory)
	uid.SetDNSType(dns.TypeUID)
	zdns.RegisterLookup("UID", uid)

	uinfo := new(miekg.GlobalLookupFactory)
	uinfo.SetDNSType(dns.TypeUINFO)
	zdns.RegisterLookup("UINFO", uinfo)

	unspec := new(miekg.GlobalLookupFactory)
	unspec.SetDNSType(dns.TypeUNSPEC)
	zdns.RegisterLookup("UNSPEC", unspec)

	uri := new(miekg.GlobalLookupFactory)
	uri.SetDNSType(dns.TypeURI)
	zdns.RegisterLookup("URI", uri)

	// Question Only Types
	any := new(miekg.GlobalLookupFactory)
	any.SetDNSType(dns.TypeANY)
	zdns.RegisterLookup("ANY", any)

	// Transfer have their own modules

	//ixfr := new(miekg.GlobalLookupFactory)
	//ixfr.SetDNSType(dns.TypeIXFR)
	//zdns.RegisterLookup("IXFR", ixfr)

	//axfr := new(miekg.GlobalLookupFactory)
	//axfr.SetDNSType(dns.TypeAXFR)
	//zdns.RegisterLookup("AXFR", axfr)

	// TODO(zakir): investigate whether MAILA and MAILB should be supported questions
}
