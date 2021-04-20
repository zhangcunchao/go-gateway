package inc

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"fmt"
	"os"
	"runtime"

	"go-gateway/debug"

	"github.com/Unknwon/goconfig"
	"github.com/gin-gonic/gin"
)

var (
	Cfg *goconfig.ConfigFile
	key = []byte("cdgbdfef")
)

//主动初始化
func init() {

	defer func() {
		if r := recover(); r != nil {
			debug.DebugPrint("捕获到inc层错误：%s", r)
			os.Exit(2)
		}
	}()

	//初始化conf.ini
	var err error
	Cfg, err = goconfig.LoadConfigFile("conf.ini")
	if err != nil {
		panic("conf.ini不存在")
	}

}

func InitGatway() {
	//多核配置
	numCPU := Cfg.MustInt("http", "numCPU", 0)
	if numCPU == 0 {
		numCPU = runtime.NumCPU()
	}
	runtime.GOMAXPROCS(numCPU)
	debug.DebugPrint("run cpuNum %d", numCPU)
	//gin运行模式
	DebugModel := Cfg.MustValue("http", "DebugModel", "debug")
	gin.SetMode(DebugModel)
}

func MD5(s string) string {
	hash := md5.New()
	_, err := hash.Write([]byte(s))
	if err != nil {
		panic(err)
	}
	sum := hash.Sum(nil)
	return fmt.Sprintf("%x\n", sum)
}

/*des加密解密*/

//des的cbc模式加密函数
//src ：明文
//key :密钥，大小为8byte
func DesEncrypt_CBC(src []byte) []byte { //传递两个参数,src为需要加密的明文，返回[]byte类型数据
	//1：创建并返回一个DES算法的cipher.block接口
	block, err := des.NewCipher(key)
	//2：判断是否创建成功
	if err != nil {
		panic(err)
	}
	//3：对最后一个明文分组进行数据填充
	src = PKCS5Padding(src, block.BlockSize())
	//4：创建一个密码分组为链接模式的，底层使用DES加密的BLockMode接口
	//    参数iv的长度，必须等于b的块尺寸
	tmp := []byte("helloAAA") //初始化向量
	blockmode := cipher.NewCBCEncrypter(block, tmp)
	//5:加密连续的数据块
	dst := make([]byte, len(src))
	blockmode.CryptBlocks(dst, src)
	//fmt.Println("加密之后的数据:",dst)
	//6:将加密后的数据返回
	return dst
}

//des解密函数
//src :要解密的密文
//key :密钥，和加密密钥相同，大小为8byte
func DesDecrypt_CBC(src []byte) []byte {
	//1:创建并返回一个使用DES算法的cipher.block接口
	block, err := des.NewCipher(key)
	//2:判断是否创建成功
	if err != nil {
		panic(err)
	}
	//创建一个密码分组为链接模式的，底层使用DES解密的BlockMode接口
	tmp := []byte("helloAAA")
	blockMode := cipher.NewCBCDecrypter(block, tmp)
	//解密数据
	dst := make([]byte, len(src))
	blockMode.CryptBlocks(dst, src)
	//5:去掉最后一组填充数据
	dst = PKCS5UnPadding(dst)

	//返回结果
	return dst
}

//使用pkcs5的方式填充

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	//1:计算最后一个分组缺多少字节
	padding := blockSize - (len(ciphertext) % blockSize)
	//2:创建一个大小为padding的切片,每个字节的值为padding
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	//3:将padText添加到原始数据的后边，将最后一个分组缺少的字节数补齐
	newText := append(ciphertext, padText...)

	return newText
}

//删除pkcs5填充的尾部数据
func PKCS5UnPadding(origData []byte) []byte {
	//1:计算数据的总长度
	length := len(origData)
	//2:根据填充的字节值得到填充的次数
	number := int(origData[length-1])
	//3:将尾部填充的number个字节去掉
	return origData[:(length - number)]
}
