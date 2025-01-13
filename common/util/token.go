package util

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

const aesKEY = "fc49607d05e1a1ba" // hd2yao md5 必须是 16 个字节
const md5Len = 4                  // MD5 的部分保留的字节数
const aesLen = 16                 // AES 加密后的字节数，12-->16

// 将 userId 和 MD5 揉到一起
// 类似于md5(userId+time)(4字节)+aes(userId+time)(16字节)，最终40个字符
func genAccessToken(uid int64) (string, error) {
	// 12 字节，前 8 字节是 userId，后 4 字节是当前纳秒级的时间戳
	byteInfo := make([]byte, 12)
	binary.BigEndian.PutUint64(byteInfo, uint64(uid))
	binary.BigEndian.PutUint32(byteInfo[8:], uint32(time.Now().UnixNano()))

	// 使用密钥 aesKEY 进行 AES 加密
	// 基于 byteInfo 的结构，加密 12 字节，得到 16 字节
	// 因为 AES-128 在 CBC 模式下总是输出 16 字节的数据，AES 的块大小为 16 字节
	// aes(userId+time)
	encodeByte, err := AesEncrypt(byteInfo, []byte(aesKEY))
	if err != nil {
		return "", err
	}

	// 计算 byteInfo 的 MD5 值
	// md5(userId+time)
	md5Byte := md5.Sum(byteInfo)

	// 将 MD5 值和加密后的数据拼接起来，并以十六进制编码的字符串返回
	// md5(userId+time)(4字节) + aes(userId+time)(16字节)
	data := append(md5Byte[:md5Len], encodeByte...)

	// 二进制数据（例如 AES 加密字节或 MD5 哈希）编码为十六进制时，每个字节被表示为 2 个十六进制字符
	// 因此，20 字节的数据将编码为 40 个十六进制字符
	return hex.EncodeToString(data), nil
}

func genRefreshToken(uid int64) (string, error) {
	return genAccessToken(uid)
}

func GenUserAuthToken(uid int64) (accessToken, refreshToken string, err error) {
	accessToken, err = genAccessToken(uid)
	if err != nil {
		return "", "", err
	}
	refreshToken, err = genRefreshToken(uid)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func GenSessionId(userId int64) string {
	return fmt.Sprintf("%d-%d-%s", userId, time.Now().Unix(), RandNumStr(6))
}

// ParseUserIdFromToken 从 Token 中解析出 userId
// 后端服务 redis 不可用也无法立即恢复时可以使用这个方法保持产品最基本功能的使用，不至于直接白屏
func ParseUserIdFromToken(accessToken string) (userId int64, err error) {
	if len(accessToken) != 2*(md5Len+aesLen) {
		// Token 格式不对
		return 0, errors.New("invalid token")
	}
	encodeStr := accessToken[md5Len*2:]
	data, err := hex.DecodeString(encodeStr)
	if err != nil {
		return 0, err
	}
	decodeByte, err := AesDecrypt(data, []byte(aesKEY))
	if err != nil {
		return 0, err
	}
	uid := binary.BigEndian.Uint64(decodeByte)
	if uid == 0 {
		return 0, errors.New("invalid token")
	}
	return int64(uid), nil
}
