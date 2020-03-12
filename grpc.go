package main

import (
	"log"

	"google.golang.org/grpc"
)
var xgfaceAddr, xgindexAddr string
func init()  {
	xgfaceAddr = Config.XGFace.Base.Address + ":" + Config.XGFace.Base.Port
	xgindexAddr = Config.Index.Base.Address + ":" + Config.Index.Base.Port
}
var xgFaceConn, xgIndexConn *grpc.ClientConn
func newXgFaceConn() (*grpc.ClientConn, error) {
	if xgFaceConn == nil {
		lock.Lock()
		defer lock.Unlock()
		if xgFaceConn == nil {
			var err error
			var opts []grpc.DialOption
			opts = append(opts, grpc.WithInsecure())
			xgFaceConn, err = grpc.Dial(xgfaceAddr, opts...)
			if err != nil {
				log.Println("xgface grpc connect err:", err)
				xgFaceConn = nil
				return xgFaceConn, err
			}
		}
	}

	return xgFaceConn, err
	//只有main.go中才有必要		defer xgFaceConn.Close()
}

//GetXgFaceConn 获取xgface连接
func GetXgFaceConn() (*grpc.ClientConn, error) {
	return newXgFaceConn()
}

func newXgIndexConn() (*grpc.ClientConn, error) {

	if xgIndexConn == nil { //加锁是为了并发，加锁前判断是为了减少操作锁的消耗
		lock.Lock()
		defer lock.Unlock()
		if xgIndexConn == nil {
			var err error
			var opts []grpc.DialOption
			opts = append(opts, grpc.WithInsecure())
			xgIndexConn, err = grpc.Dial(xgindexAddr, opts...)
			if err != nil {
				log.Println("xgindex grpc connect err:", err)
				xgIndexConn = nil
				return xgIndexConn, err
			}
		}
	}
	return xgIndexConn, err
}

//GetXgIndexConn 创建index连接
func GetXgIndexConn() (*grpc.ClientConn, error) {
	return newXgIndexConn()
}
