package httpfile

import (
	"errors"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/kprc/nbsnetwork/tools/logger"

	"net/http"
	"path"
)

type AssetfsFileSystem interface {
	http.FileSystem
	TryOpen(name string) error
}

type insepectorHttpFileSys struct {
	wfs *assetfs.AssetFS
	bListFile bool
	log logger.Logger
}

var New = func(wfs *assetfs.AssetFS, bList bool, log logger.Logger) AssetfsFileSystem {
	return &insepectorHttpFileSys{
		wfs: wfs,
		bListFile: bList,
		log: log,
	}
}

var homePages = []string{
	"index.html",
	"index.htm",
	"default.html",
	"default.htm",
}

func (ihfs *insepectorHttpFileSys)TryOpen(name string) error  {

	ihfs.log.DebugPrintln("AssetfsFileSystem : try open name ",name)

	name = path.Join(ihfs.wfs.Prefix,name)
	if len(name)>0 && name[0] == '/'{
		name = name[1:]
	}
	ihfs.log.DebugPrintln("AssetfsFileSystem : try open full name ",name)

	if _,err:=ihfs.wfs.Asset(name);err==nil{
		ihfs.log.InfoPrintln("AssetfsFileSystem : try open not a file ",name)
		return nil
	}

	if allfiles, err := ihfs.wfs.AssetDir(name);err != nil{
		ihfs.log.WarnPrintln("AssetfsFileSystem : try open not a directory ",name)
		return err
	}else{
		ihfs.log.DebugPrintln("AssetfsFileSystem :try open allfiles in directory ",name, allfiles)
		for _,f := range allfiles{
			for i:=0;i<len(homePages);i++{
				if f == homePages[i]{
					ihfs.log.InfoPrintln("AssetfsFileSystem : try open is a homepage", f)
					return nil
				}
			}
		}
	}
	ihfs.log.WarnPrintln("AssetfsFileSystem : try open not found a  homepage")
	return errors.New("file not found")
}

func (ihfs *insepectorHttpFileSys)Open(name string) (http.File, error)  {
	ihfs.log.DebugPrintln("AssetfsFileSystem : open name ",name)
	if !ihfs.bListFile{
		return ihfs.wfs.Open(name)
	}

	name = path.Join(ihfs.wfs.Prefix,name)
	if len(name)>0 && name[0] == '/'{
		name = name[1:]
	}

	ihfs.log.DebugPrintln("AssetfsFileSystem : open full name ",name)

	if _,err:=ihfs.wfs.Asset(name);err==nil{
		ihfs.log.InfoPrintln("AssetfsFileSystem : open not a file ",name)
		return ihfs.wfs.Open(name)
	}

	if allfiles, err := ihfs.wfs.AssetDir(name);err != nil{
		ihfs.log.InfoPrintln("AssetfsFileSystem : open not a directory ",name)
		return ihfs.wfs.Open(name)
	}else{
		ihfs.log.DebugPrintln("AssetfsFileSystem : open allfiles in directory ",name, allfiles)
		for _,f := range allfiles{
			for i:=0;i<len(homePages);i++{
				if f == homePages[i]{
					ihfs.log.InfoPrintln("AssetfsFileSystem : open is a homepage", f,path.Join(name,f))
					return ihfs.wfs.Open(path.Join(name,f))
				}
			}
		}
	}
	ihfs.log.WarnPrintln("AssetfsFileSystem : open not found a  homepage")
	return nil,errors.New("file not found")
}


