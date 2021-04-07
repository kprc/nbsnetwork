package httpfile

import (
	"errors"
	"path"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"net/http"
)

type AssetfsFileSystem interface {
	http.FileSystem
	TryOpen(name string) error
}

type insepectorHttpFileSys struct {
	wfs *assetfs.AssetFS
	bListFile bool
}

var New = func(wfs *assetfs.AssetFS, bList bool) AssetfsFileSystem {
	return &insepectorHttpFileSys{
		wfs: wfs,
		bListFile: bList,
	}
}

var homePages = []string{
	"index.html",
	"index.htm",
	"default.html",
	"default.htm",
}

func (ihfs *insepectorHttpFileSys)TryOpen(name string) error  {
	name = path.Join(ihfs.wfs.Prefix,name)
	if len(name)>0 && name[0] == '/'{
		name = name[1:]
	}

	if _,err:=ihfs.wfs.Asset(name);err==nil{
		return nil
	}

	if allfiles, err := ihfs.wfs.AssetDir(name);err != nil{
		return err
	}else{
		for _,f := range allfiles{
			for i:=0;i<len(homePages);i++{
				if f == homePages[i]{
					return nil
				}
			}
		}
	}
	return errors.New("file not found")
}

func (ihfs *insepectorHttpFileSys)Open(name string) (http.File, error)  {
	if !ihfs.bListFile{
		return ihfs.wfs.Open(name)
	}

	name = path.Join(ihfs.wfs.Prefix,name)
	if len(name)>0 && name[0] == '/'{
		name = name[1:]
	}

	if _,err:=ihfs.wfs.Asset(name);err==nil{
		return ihfs.wfs.Open(name)
	}

	if allfiles, err := ihfs.wfs.AssetDir(name);err != nil{
		return ihfs.wfs.Open(name)
	}else{
		for _,f := range allfiles{
			for i:=0;i<len(homePages);i++{
				if f == homePages[i]{
					return ihfs.wfs.Open(path.Join(name,f))
				}
			}
		}
	}
	return nil,errors.New("file not found")
}


