package nbscrypt

import (
	"crypto/rsa"
	"crypto/rand"
	"github.com/kprc/nbsnetwork/tools"
	"github.com/pkg/errors"
	"path"
	"os"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
)

func GenerateKeyPair(bitsCnt int)(*rsa.PrivateKey,*rsa.PublicKey)  {
	priv,err:=rsa.GenerateKey(rand.Reader,bitsCnt)

	if err!=nil{
		return nil,nil
	}

	return priv,&priv.PublicKey

}

func RsaKeyIsExists(savaPath string) bool  {
	if tools.FileExists(path.Join(savaPath,"priv.key")){
		return true
	}

	return false
}

func Save2FileRSAKey(savePath string,privKey *rsa.PrivateKey)  error{
	if savePath == "" {
		return errors.New("Path is none")
	}

	if privKey == nil{
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

	if f,err:=os.OpenFile(path.Join(keypath,"priv.key"),os.O_CREATE|os.O_TRUNC|os.O_WRONLY,0755);err!=nil{
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
	if f,err:=os.OpenFile(path.Join(keypath,"pub.key"),os.O_CREATE|os.O_TRUNC|os.O_WRONLY,0755);err!=nil{
		return err
	}else{
		if err=pem.Encode(f,block);err!=nil{
			return err
		}
	}

	return nil
}

func LoadRSAKey(savePath string) (priv *rsa.PrivateKey,pub *rsa.PublicKey,err error)  {
	//load privKey
	var block *pem.Block
	if privKey,err:=ioutil.ReadFile(path.Join(savePath,"priv.key"));err!=nil{
		return nil,nil,errors.New("read priv.key error")
	}else {
		block, _ = pem.Decode(privKey)
		if block == nil{
			return nil,nil,errors.New("recover privKey error")
		}
	}
	if priv,err = x509.ParsePKCS1PrivateKey(block.Bytes);err!=nil{
		return nil,nil,errors.New("Parse privKey error")
	}

	pub = &priv.PublicKey

	//fmt.Println("privkey:",len(block.Bytes),base58.FastBase58Encoding(block.Bytes))
	//pubbytes:=x509.MarshalPKCS1PublicKey(pub)
	//fmt.Println("pubkey:",len(pubbytes),base64.StdEncoding.EncodeToString(pubbytes))

	if pubkey,err:=ioutil.ReadFile(path.Join(savePath,"pub.key"));err!=nil{
		return nil,nil,errors.New("read pub.key error")
	}else {
		block, _ = pem.Decode(pubkey)
		if block == nil{
			return nil,nil,errors.New("recover pubKey error")
		}
		//fmt.Println("pubkey:",base64.StdEncoding.EncodeToString(block.Bytes))

	}

	if pub,err = x509.ParsePKCS1PublicKey(block.Bytes);err!=nil{
		return nil,nil,errors.New("Parse PubKey error")
	}

	return

}

func EncryptRSA(data []byte,pub *rsa.PublicKey) (encData []byte ,err error)  {
	return rsa.EncryptPKCS1v15(rand.Reader,pub,data)
}

func DecryptRsa(encData []byte,priv *rsa.PrivateKey) (data []byte ,err error)   {
	return rsa.DecryptPKCS1v15(rand.Reader,priv,encData)
}


