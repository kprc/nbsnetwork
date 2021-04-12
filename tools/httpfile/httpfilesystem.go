package httpfile

import (
	"errors"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"log"
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

	log.Println("AssetfsFileSystem : try open name ",name)

	name = path.Join(ihfs.wfs.Prefix,name)
	if len(name)>0 && name[0] == '/'{
		name = name[1:]
	}
	log.Println("AssetfsFileSystem : try open full name ",name)

	if _,err:=ihfs.wfs.Asset(name);err==nil{
		log.Println("AssetfsFileSystem : try open not a file ",name)
		return nil
	}

	if allfiles, err := ihfs.wfs.AssetDir(name);err != nil{
		log.Println("AssetfsFileSystem : try open not a directory ",name)
		return err
	}else{
		log.Println("AssetfsFileSystem :try open allfiles in directory ",name, allfiles)
		for _,f := range allfiles{
			for i:=0;i<len(homePages);i++{
				if f == homePages[i]{
					log.Println("AssetfsFileSystem : try open is a homepage", f)
					return nil
				}
			}
		}
	}
	log.Println("AssetfsFileSystem : try open not found a  homepage")
	return errors.New("file not found")
}

func (ihfs *insepectorHttpFileSys)Open(name string) (http.File, error)  {
	log.Println("AssetfsFileSystem : open name ",name)
	if !ihfs.bListFile{
		return ihfs.wfs.Open(name)
	}

	name = path.Join(ihfs.wfs.Prefix,name)
	if len(name)>0 && name[0] == '/'{
		name = name[1:]
	}

	log.Println("AssetfsFileSystem : open full name ",name)

	if _,err:=ihfs.wfs.Asset(name);err==nil{
		log.Println("AssetfsFileSystem : open not a file ",name)
		return ihfs.wfs.Open(name)
	}

	if allfiles, err := ihfs.wfs.AssetDir(name);err != nil{
		log.Println("AssetfsFileSystem : open not a directory ",name)
		return ihfs.wfs.Open(name)
	}else{
		log.Println("AssetfsFileSystem :open allfiles in directory ",name, allfiles)
		for _,f := range allfiles{
			for i:=0;i<len(homePages);i++{
				if f == homePages[i]{
					log.Println("AssetfsFileSystem : open is a homepage", f,path.Join(name,f))
					return ihfs.wfs.Open(path.Join(name,f))
				}
			}
		}
	}
	log.Println("AssetfsFileSystem : open not found a  homepage")
	return nil,errors.New("file not found")
}


