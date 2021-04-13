package httphandle

import (
	"github.com/kprc/nbsnetwork/tools/httpfile"
	"github.com/kprc/nbsnetwork/tools/logger"
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
	log logger.Logger
}

var homePages = []string{
	"index.html",
	"index.htm",
	"default.html",
	"default.htm",
}


var New = func(afs httpfile.AssetfsFileSystem,
	           arg interface{}, fInsepectorHandle func( r *http.Request, arg interface{}) bool,log logger.Logger) *insepectorHandler{
	return &insepectorHandler{
		afs: afs,
		handler: http.FileServer(afs),
		arg: arg,
		fInsepectorHandle: fInsepectorHandle,
		skipCheckCookie: nil,
		log: log,
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

		ih.log.DebugPrintln("insepectorHandler : url path is ",upath)

		ih.log.DebugPrintln("insepectorHandler : skip cookies ",ih.skipCheckCookie)

		for i:=0;i<len(ih.skipCheckCookie);i++{
			if upath == ih.skipCheckCookie[i] {
				ih.handler.ServeHTTP(w,r)
				ih.log.InfoPrintln("insepectorHandler : skip by ",ih.skipCheckCookie[i])
				return
			}
		}

		for i:=0;i<len(homePages);i++{
			p:=path.Join(upath,homePages[i])
			for j:=0;j<len(ih.skipCheckCookie);j++{
				if p == ih.skipCheckCookie[j]{
					ih.log.InfoPrintln("insepectorHandler : skip by homepage ",p)
					ih.handler.ServeHTTP(w,r)
					return
				}
			}
		}


		ih.log.DebugPrintln("insepectorHandler : clean path ",path.Clean(upath))
		if err:=ih.afs.TryOpen(path.Clean(upath));err==nil{
			if b:=ih.fInsepectorHandle(r,ih.arg);!b{
				ih.log.InfoPrintln("redirect to new url",ih.redirUrl)
				http.Redirect(w,r,ih.redirUrl,302)
				return
			}
		}
	}

	ih.log.DebugPrintln("insepectorHandler : handle by access file system....")

	ih.handler.ServeHTTP(w,r)
}