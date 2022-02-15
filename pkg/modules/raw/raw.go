package raw

import (
	"github.com/zmap/dns"
	"github.com/zmap/zdns/pkg/miekg"
	"github.com/zmap/zdns/pkg/zdns"
)

func AddRawLookupsTo(s zdns.FactorySet) {
	a := new(miekg.GlobalLookupFactory)
	a.SetDNSType(dns.TypeA)
	s.AddModule("A", a)

	aaaa := new(miekg.GlobalLookupFactory)
	aaaa.SetDNSType(dns.TypeAAAA)
	s.AddModule("AAAA", aaaa)

	afsdb := new(miekg.GlobalLookupFactory)
	afsdb.SetDNSType(dns.TypeAFSDB)
	s.AddModule("AFSDB", afsdb)

	atma := new(miekg.GlobalLookupFactory)
	atma.SetDNSType(dns.TypeATMA)
	s.AddModule("ATMA", atma)

	avc := new(miekg.GlobalLookupFactory)
	avc.SetDNSType(dns.TypeAVC)
	s.AddModule("AVC", avc)

	caa := new(miekg.GlobalLookupFactory)
	caa.SetDNSType(dns.TypeCAA)
	s.AddModule("CAA", caa)

	cert := new(miekg.GlobalLookupFactory)
	cert.SetDNSType(dns.TypeCERT)
	s.AddModule("CERT", cert)

	cds := new(miekg.GlobalLookupFactory)
	cds.SetDNSType(dns.TypeCDS)
	s.AddModule("CDS", cds)

	cdnskey := new(miekg.GlobalLookupFactory)
	cdnskey.SetDNSType(dns.TypeCDNSKEY)
	s.AddModule("CDNSKEY", cdnskey)

	cname := new(miekg.GlobalLookupFactory)
	cname.SetDNSType(dns.TypeCNAME)
	s.AddModule("CNAME", cname)

	csync := new(miekg.GlobalLookupFactory)
	csync.SetDNSType(dns.TypeCSYNC)
	s.AddModule("CSYNC", csync)

	dhcid := new(miekg.GlobalLookupFactory)
	dhcid.SetDNSType(dns.TypeDHCID)
	s.AddModule("DHCID", dhcid)

	dnskey := new(miekg.GlobalLookupFactory)
	dnskey.SetDNSType(dns.TypeDNSKEY)
	s.AddModule("DNSKEY", dnskey)

	ds := new(miekg.GlobalLookupFactory)
	ds.SetDNSType(dns.TypeDS)
	s.AddModule("DS", ds)

	eid := new(miekg.GlobalLookupFactory)
	eid.SetDNSType(dns.TypeEID)
	s.AddModule("EID", eid)

	eui48 := new(miekg.GlobalLookupFactory)
	eui48.SetDNSType(dns.TypeEUI48)
	s.AddModule("EUI48", eui48)

	eui64 := new(miekg.GlobalLookupFactory)
	eui64.SetDNSType(dns.TypeEUI64)
	s.AddModule("EUI64", eui64)

	gid := new(miekg.GlobalLookupFactory)
	gid.SetDNSType(dns.TypeGID)
	s.AddModule("GID", gid)

	gpos := new(miekg.GlobalLookupFactory)
	gpos.SetDNSType(dns.TypeGPOS)
	s.AddModule("GPOS", gpos)

	hinfo := new(miekg.GlobalLookupFactory)
	hinfo.SetDNSType(dns.TypeHINFO)
	s.AddModule("HINFO", hinfo)

	hip := new(miekg.GlobalLookupFactory)
	hip.SetDNSType(dns.TypeHIP)
	s.AddModule("HIP", hip)

	https := new(miekg.GlobalLookupFactory)
	https.SetDNSType(dns.TypeHTTPS)
	s.AddModule("HTTPS", https)

	isdn := new(miekg.GlobalLookupFactory)
	isdn.SetDNSType(dns.TypeISDN)
	s.AddModule("ISDN", isdn)

	key := new(miekg.GlobalLookupFactory)
	key.SetDNSType(dns.TypeKEY)
	s.AddModule("KEY", key)

	kx := new(miekg.GlobalLookupFactory)
	kx.SetDNSType(dns.TypeKX)
	s.AddModule("KX", kx)

	l32 := new(miekg.GlobalLookupFactory)
	l32.SetDNSType(dns.TypeL32)
	s.AddModule("L32", l32)

	l64 := new(miekg.GlobalLookupFactory)
	l64.SetDNSType(dns.TypeL64)
	s.AddModule("L64", l64)

	loc := new(miekg.GlobalLookupFactory)
	loc.SetDNSType(dns.TypeLOC)
	s.AddModule("LOC", loc)

	lp := new(miekg.GlobalLookupFactory)
	lp.SetDNSType(dns.TypeLP)
	s.AddModule("LP", lp)

	md := new(miekg.GlobalLookupFactory)
	md.SetDNSType(dns.TypeMD)
	s.AddModule("MD", md)

	mf := new(miekg.GlobalLookupFactory)
	mf.SetDNSType(dns.TypeMF)
	s.AddModule("MF", mf)

	mb := new(miekg.GlobalLookupFactory)
	mb.SetDNSType(dns.TypeMB)
	s.AddModule("MB", mb)

	mg := new(miekg.GlobalLookupFactory)
	mg.SetDNSType(dns.TypeMG)
	s.AddModule("MG", mg)

	mr := new(miekg.GlobalLookupFactory)
	mr.SetDNSType(dns.TypeMR)
	s.AddModule("MR", mr)

	mx := new(miekg.GlobalLookupFactory)
	mx.SetDNSType(dns.TypeMX)
	s.AddModule("MX", mx)

	naptr := new(miekg.GlobalLookupFactory)
	naptr.SetDNSType(dns.TypeNAPTR)
	s.AddModule("NAPTR", naptr)

	nimloc := new(miekg.GlobalLookupFactory)
	nimloc.SetDNSType(dns.TypeNIMLOC)
	s.AddModule("NS", nimloc)

	nid := new(miekg.GlobalLookupFactory)
	nid.SetDNSType(dns.TypeNID)
	s.AddModule("NID", nid)

	ninfo := new(miekg.GlobalLookupFactory)
	ninfo.SetDNSType(dns.TypeNINFO)
	s.AddModule("NINFO", ninfo)

	nsapptr := new(miekg.GlobalLookupFactory)
	nsapptr.SetDNSType(dns.TypeNSAPPTR)
	s.AddModule("NSAPPTR", nsapptr)

	ns := new(miekg.GlobalLookupFactory)
	ns.SetDNSType(dns.TypeNS)
	s.AddModule("NS", ns)

	nxt := new(miekg.GlobalLookupFactory)
	nxt.SetDNSType(dns.TypeNXT)
	s.AddModule("NXT", nxt)

	nsec := new(miekg.GlobalLookupFactory)
	nsec.SetDNSType(dns.TypeNSEC)
	s.AddModule("NSEC", nsec)

	nsec3 := new(miekg.GlobalLookupFactory)
	nsec3.SetDNSType(dns.TypeNSEC3)
	s.AddModule("NSEC3", nsec3)

	nsec3param := new(miekg.GlobalLookupFactory)
	nsec3param.SetDNSType(dns.TypeNSEC3PARAM)
	s.AddModule("NSEC3PARAM", nsec3param)

	null := new(miekg.GlobalLookupFactory)
	null.SetDNSType(dns.TypeNULL)
	s.AddModule("NULL", null)

	openpgpkey := new(miekg.GlobalLookupFactory)
	openpgpkey.SetDNSType(dns.TypeOPENPGPKEY)
	s.AddModule("OPENPGPKEY", openpgpkey)

	//opt := new(miekg.GlobalLookupFactory)
	//opt.SetDNSType(dns.TypeOPT)
	//s.AddModule("OPT", opt)

	ptr := new(miekg.GlobalLookupFactory)
	ptr.SetDNSType(dns.TypePTR)
	s.AddModule("PTR", ptr)

	px := new(miekg.GlobalLookupFactory)
	px.SetDNSType(dns.TypePX)
	s.AddModule("PX", px)

	rp := new(miekg.GlobalLookupFactory)
	rp.SetDNSType(dns.TypeRP)
	s.AddModule("RP", rp)

	rrsig := new(miekg.GlobalLookupFactory)
	rrsig.SetDNSType(dns.TypeRRSIG)
	s.AddModule("RRSIG", rrsig)

	rt := new(miekg.GlobalLookupFactory)
	rt.SetDNSType(dns.TypeRT)
	s.AddModule("RT", rt)

	smimea := new(miekg.GlobalLookupFactory)
	smimea.SetDNSType(dns.TypeSMIMEA)
	s.AddModule("SMIMEA", smimea)

	sshfp := new(miekg.GlobalLookupFactory)
	sshfp.SetDNSType(dns.TypeSSHFP)
	s.AddModule("SSHFP", sshfp)

	soa := new(miekg.GlobalLookupFactory)
	soa.SetDNSType(dns.TypeSOA)
	s.AddModule("SOA", soa)

	spf := new(miekg.GlobalLookupFactory)
	spf.SetDNSType(dns.TypeSPF)
	s.AddModule("SPF", spf)

	srv := new(miekg.GlobalLookupFactory)
	srv.SetDNSType(dns.TypeSRV)
	s.AddModule("SRV", srv)

	svcb := new(miekg.GlobalLookupFactory)
	svcb.SetDNSType(dns.TypeSVCB)
	s.AddModule("SVCB", svcb)

	talink := new(miekg.GlobalLookupFactory)
	talink.SetDNSType(dns.TypeTALINK)
	s.AddModule("TALINK", talink)

	tkey := new(miekg.GlobalLookupFactory)
	tkey.SetDNSType(dns.TypeTKEY)
	s.AddModule("TKEY", tkey)

	tlsa := new(miekg.GlobalLookupFactory)
	tlsa.SetDNSType(dns.TypeTLSA)
	s.AddModule("TLSA", tlsa)

	txt := new(miekg.GlobalLookupFactory)
	txt.SetDNSType(dns.TypeTXT)
	s.AddModule("TXT", txt)

	uid := new(miekg.GlobalLookupFactory)
	uid.SetDNSType(dns.TypeUID)
	s.AddModule("UID", uid)

	uinfo := new(miekg.GlobalLookupFactory)
	uinfo.SetDNSType(dns.TypeUINFO)
	s.AddModule("UINFO", uinfo)

	unspec := new(miekg.GlobalLookupFactory)
	unspec.SetDNSType(dns.TypeUNSPEC)
	s.AddModule("UNSPEC", unspec)

	uri := new(miekg.GlobalLookupFactory)
	uri.SetDNSType(dns.TypeURI)
	s.AddModule("URI", uri)

	// Question Only Types
	any := new(miekg.GlobalLookupFactory)
	any.SetDNSType(dns.TypeANY)
	s.AddModule("ANY", any)

	// Transfer have their own modules

	//ixfr := new(miekg.GlobalLookupFactory)
	//ixfr.SetDNSType(dns.TypeIXFR)
	//s.AddModule("IXFR", ixfr)

	//axfr := new(miekg.GlobalLookupFactory)
	//axfr.SetDNSType(dns.TypeAXFR)
	//s.AddModule("AXFR", axfr)

	// TODO(zakir): investigate whether MAILA and MAILB should be supported questions
}
