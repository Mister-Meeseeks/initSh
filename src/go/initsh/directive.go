
package initsh

type ImportDirector interface {
	importShell (path string, namespace *string) PathIngester
	importNested (path string, namespace *string) PathIngester
	importSubcmd (path string, namespace string) PathIngester
	importNestSubcmd (path string, namespace string) PathIngester
	importData (path string, namespace *string) PathIngester
	importNestData (path string, namespace *string) PathIngester
}

type importDirector struct {
	binPath string
	libPath string
	div string
}

func MakeImporter (binPath string, libPath string, div string) ImportDirector {
	return importDirector{binPath, libPath, div}
}

func (d importDirector) importShell (path string, namespace *string) PathIngester {
	ex := ImportFunnel{ExecFilter{}, d.flatBinTrans(namespace),
		d.linkPathShipper(), path}
	lib := ImportFunnel{ShellLibFilter{}, d.flatLibTrans(namespace),
		d.linkPathShipper(), path}
	return mergeIngesters(ex, lib)
}

func (d importDirector) importNested (path string, namespace *string) PathIngester {
	ex := ImportFunnel{ExecFilter{}, d.nestBinTrans(namespace),
		d.linkPathShipper(), path}
	lib := ImportFunnel{ShellLibFilter{}, d.nestLibTrans(namespace),
		d.linkPathShipper(), path}
	return mergeIngesters(ex, lib)	
}

func (d importDirector) importData (path string, namespace *string) PathIngester {
	ex := ImportFunnel{DataItemFilter{}, d.flatUndropBinTrans(namespace),
		d.dataShipper(), path}
	lib := ImportFunnel{GzDataFilter{}, d.flatBinTrans(namespace),
		d.gzDataShipper(), path}
	return mergeIngesters(ex, lib)
}

func (d importDirector) importNestData (path string, namespace *string) PathIngester {
	ex := ImportFunnel{DataItemFilter{}, d.nestDataTrans(namespace),
		d.dataShipper(), path}
	lib := ImportFunnel{GzDataFilter{}, d.nestDataTrans(namespace),
		d.gzDataShipper(), path}
	return mergeIngesters(ex, lib)
}

func (d importDirector) importSubcmd (path string, namespace string) PathIngester {
	ex := ImportFunnel{ExecFilter{}, d.subcmdFlatTrans(namespace),
		d.subcmdShipper(), path}
	lib := ImportFunnel{ShellLibFilter{}, d.flatLibTrans(&namespace),
		d.linkPathShipper(), path}
	return mergeIngesters(ex, lib)
}

func (d importDirector) importNestSubcmd (path string, namespace string) PathIngester {
	ex := ImportFunnel{ExecFilter{}, d.subcmdBinTrans(namespace),
		d.subcmdShipper(), path}
	lib := ImportFunnel{ShellLibFilter{}, d.nestLibTrans(&namespace),
		d.linkPathShipper(), path}	
	return mergeIngesters(ex, lib)
}

func (d importDirector) subcmdFlatTrans (namespace string) AddressTranslator {
	return stackTrans(FlattenTranslator{"/"}, d.subcmdBinTrans(namespace))
}

func (d importDirector) subcmdBinTrans (namespace string) AddressTranslator {
	return stackTrans(DropExtTranslator{},
		SubcmdTranslator{namespace, d.binPath})
}

func (d importDirector) flatBinTrans (namespace *string) AddressTranslator {
	return stackTrans(FlattenTranslator{"/"}, d.dropBinTrans(namespace))
}

func (d importDirector) nestBinTrans (namespace *string) AddressTranslator {
	return stackTrans(NestedTranslator{"/", d.div}, d.dropBinTrans(namespace))
}

func (d importDirector) nestDataTrans (namespace *string) AddressTranslator {
	return stackTrans(NestedTranslator{"/", d.div}, d.baseBinTrans(namespace))
}

func (d importDirector) flatUndropBinTrans (namespace *string) AddressTranslator {
	return stackTrans(FlattenTranslator{"/"}, d.baseBinTrans(namespace))
}

func (d importDirector) flatLibTrans (namespace *string) AddressTranslator {
	return stackTrans(FlattenTranslator{"/"}, d.baseLibTrans(namespace))
}

func (d importDirector) nestLibTrans (namespace *string) AddressTranslator {
	return stackTrans(NestedTranslator{"/", d.div}, d.baseLibTrans(namespace))
}

func (d importDirector) dropBinTrans (namespace *string) AddressTranslator {
	return stackTrans(DropExtTranslator{}, d.baseBinTrans(namespace))
}

func (d importDirector) baseBinTrans (namespace *string) AddressTranslator {
	return d.baseTrans(d.binPath, namespace)
}

func (d importDirector) baseLibTrans (namespace *string) AddressTranslator {
	return d.baseTrans(d.libPath, namespace)
}

func (d importDirector) baseTrans (dest string, namespace *string) AddressTranslator {
	if (namespace == nil) {
		return MirrorTranslator{dest}
	} else {
		return NamespaceTranslator{dest, *namespace, d.div}
	}
}

func (d importDirector) linkPathShipper() cargoDeliverer {
	return cargoDeliverer{bucketerPathMode{}, symLinkSlotter{}}
}

func (d importDirector) subcmdShipper() cargoDeliverer {
	return cargoDeliverer{bucketerSubcmd{}, subcmdSlotter{}}
}

func (d importDirector) dataShipper() cargoDeliverer {
	return cargoDeliverer{bucketerPathMode{}, dataSlotter{}}
}

func (d importDirector) gzDataShipper() cargoDeliverer {
	return cargoDeliverer{bucketerPathMode{}, gzDataSlotter{}}
}
