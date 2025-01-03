package util

import (
    "encoding/binary"
    "math/rand"
    "net"
    "strconv"
    "strings"
    "time"
)

func GenerateSpanID(addr string) string {
    // 提取地址中的 IP 部分
    // 例如，addr 可能是 "192.168.1.1:8080"，则 ip 为 "192.168.1.1"
    strAddr := strings.Split(addr, ":")
    if len(strAddr) < 1 {
        return ""
    }
    ip := strAddr[0]
    ipLong, _ := Ip2Long(ip)
    // 纳秒级时间戳
    times := uint64(time.Now().UnixNano())
    rand.Seed(time.Now().UnixNano())
    // 1.当前时间戳与 ipLong 异或操作
    // 2.将结果左移 32 位
    // 3.与随机数进行按位或操作
    spanId := ((times ^ uint64(ipLong)) << 32) | uint64(rand.Int31())
    // 将结果转换为 16 进制字符串
    return strconv.FormatUint(spanId, 16)
}

// Ip2Long 将字符串形式的 IP 地址（如 "192.168.1.1"）转换为 uint32 类型
func Ip2Long(ip string) (uint32, error) {
    ipAddr, err := net.ResolveIPAddr("ip", ip)
    if err != nil {
        return 0, err
    }
    // IPV6 的特殊处理
    if len(ipAddr.IP) == 16 {
        return binary.BigEndian.Uint32(ipAddr.IP[12:16]), nil
    }
    return binary.BigEndian.Uint32(ipAddr.IP), nil
}
