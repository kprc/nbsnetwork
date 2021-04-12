package httphandle

import (
	"github.com/kprc/nbsnetwork/tools/httpfile"
	"log"
	"net/http"
	"path"
	"strings"
)

type insepectorHandler struct {
	afs httpfile.AssetfsFileSystem
	handler http.Handler
	arg interface{}
	fInsepectorHandle func(r *http.Request, arg interface{}) bool
	skipCheckCookie []string
	redirUrl string
}


var New = func(afs httpfile.AssetfsFileSystem,
	           arg interface{}, fInsepectorHandle func( r *http.Request, arg interface{}) bool) *insepectorHandler{
	return &insepectorHandler{
		afs: afs,
		handler: http.FileServer(afs),
		arg: arg,
		fInsepectorHandle: fInsepectorHandle,
		skipCheckCookie: nil,
	}
}

type AssestfsHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

func (ih *insepectorHandler)AddSkipUrls(url ...string)  {
	ih.skipCheckCookie = append(ih.skipCheckCookie, url...)
}

func (ih *insepectorHandler)RedirUrl(url string)  {
	ih.redirUrl = url
}

func (ih *insepectorHandler)ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	if ih.fInsepectorHandle != nil{
		upath:=r.URL.Path
		if !strings.HasPrefix(upath, "/") {
			upath = "/" + upath
			r.URL.Path = upath
		}

		log.Println("insepectorHandler : url path is ",upath)

		log.Println("insepectorHandler : skip cookies ",ih.skipCheckCookie)

		for i:=0;i<len(ih.skipCheckCookie);i++{
			if upath == ih.skipCheckCookie[i] {
				ih.handler.ServeHTTP(w,r)
				return
			}
		}

		log.Println("insepectorHandler : clean path ",path.Clean(upath))
		if err:=ih.afs.TryOpen(path.Clean(upath));err==nil{
			if b:=ih.fInsepectorHandle(r,ih.arg);!b{
				log.Println("redirect to new url",ih.redirUrl)
				http.Redirect(w,r,ih.redirUrl,302)
				return
			}
		}
	}

	log.Println("handle by access file system....")

	ih.handler.ServeHTTP(w,r)
}