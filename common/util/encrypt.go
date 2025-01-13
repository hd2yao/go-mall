package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

// 存放加密解密相关的工具函数

// AesEncrypt AES加密 | key 长度为 16 字节才能加密成功
func AesEncrypt(origData, key []byte) ([]byte, error) {
	// key 长度必须为 16 字节
	if len(key) != 16 {
		return nil, errors.New("key length must be 16 bytes")
	}

	// 创建新的 AES block
	// block 跟 key 的长度一致
	// 在 aes.NewCipher() 中其实有对 key 长度的判断，必须是 16,24，32
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 获取 block 大小
	blockSize := block.BlockSize()
	// 对原始数据进行 PKCS5 填充，确保数据长度是块大小的倍数
	origData = PKCS5Padding(origData, blockSize)

	// 创建 CBC 模式加密器，IV(初始化向量) 使用 key 的前 blockSize 字节
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	// 用于存储加密后的数据，此时的 origData 已经被填充了
	encrypted := make([]byte, len(origData))
	// 执行 AES 加密操作，将原始数据加密到 encrypted 切片中
	blockMode.CryptBlocks(encrypted, origData)
	return encrypted, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	// 计算需要填充的长度
	padding := blockSize - len(ciphertext)%blockSize
	// 生成填充字节，填充内容为 padding 数量的 byte 值
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

// AesDecrypt AES解密
func AesDecrypt(encrypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(encrypted))
	blockMode.CryptBlocks(origData, encrypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 因为填充的内容是要填充的次数，所以获取最后一位转为 int 即可得到填充次数
	// 去掉最后一个字节 unPadding 次
	unPadding := int(origData[length-1])
	if unPadding < 1 || unPadding > 32 {
		unPadding = 0
	}
	return origData[:(length - unPadding)]
}
