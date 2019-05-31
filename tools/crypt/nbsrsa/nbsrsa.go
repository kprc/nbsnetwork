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

func Save2FileRSAKey(savePath string,privKey *rsa.PrivateKey)  error{
	if savePath == "" {
		return errors.New("Path is none")
	}

	if privKey != nil{
		return errors.New("Private key is none")
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

	//save private key
	pembyte := x509.MarshalPKCS1PrivateKey(privKey)
	block := &pem.Block{Type:"Priv",Bytes:pembyte}
	if f,err:=os.OpenFile(path.Join(savePath,"priv.key"),os.O_CREATE|os.O_TRUNC|os.O_WRONLY,0755);err!=nil{
		return err
	}else{
		if err=pem.Encode(f,block);err!=nil{
			return err
		}
	}

	//save public key
	pubKey := &privKey.PublicKey
	pubbytes:= x509.MarshalPKCS1PublicKey(pubKey)
	block=&pem.Block{Type:"Pub",Bytes:pubbytes}
	if f,err:=os.OpenFile(path.Join(savePath,"pub.key"),os.O_CREATE|os.O_TRUNC|os.O_WRONLY,0755);err!=nil{
		return err
	}else{
		if err=pem.Encode(f,block);err!=nil{
			return err
		}
	}

	return nil
}

func LoadRSAKey(savePath string) (*rsa.PrivateKey,*rsa.PublicKey,error)  {
	x509.ParsePKCS1PrivateKey()
}

