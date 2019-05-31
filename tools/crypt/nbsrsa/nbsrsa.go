package nbsrsa

import (
	"crypto/rsa"
	"crypto/rand"
	"github.com/kprc/nbsnetwork/tools"
	"github.com/pkg/errors"
	"path"
	"os"
	"crypto/x509"
	"encoding/pem"
)

func GenerateKeyPair(bitsCnt int)(*rsa.PrivateKey,*rsa.PublicKey)  {
	priv,err:=rsa.GenerateKey(rand.Reader,bitsCnt)

	if err!=nil{
		return nil,nil
	}

	return priv,&priv.PublicKey

}

func Save2FileRSAKey(savePath string,privKey *rsa.PrivateKey,pubKey *rsa.PublicKey)  error{
	if savePath == "" {
		return errors.New("Path is none")
	}

	keypath:=savePath

	if !path.IsAbs(savePath){
		if homedir,err:=tools.Home();err!=nil{
			return errors.New("Cant get work home directory")
		}else {
			keypath = path.Join(homedir,savePath)
		}
	}
	if !tools.FileExists(keypath) {
		os.MkdirAll(keypath,0755)
	}

	if privKey != nil{
		pembyte := x509.MarshalPKCS1PrivateKey(privKey)
		block := pem.Block{Type:"Priv",Bytes:pembyte}
		tools.Save2File()
	}

	return nil
}

