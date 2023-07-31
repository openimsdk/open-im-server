package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"sort"
	"strings"
)

func Get2StringHash(str1, str2 string) uint32 {
	// 将两个字符串进行拼接
	concatenated := strings.Join([]string{str1, str2}, "")
	// 将拼接后的字符串转换为字节数组
	bytes := []byte(concatenated)
	// 对字节数组进行排序
	sort.Slice(bytes, func(i, j int) bool {
		return bytes[i] < bytes[j]
	})
	// 创建一个新的FNV哈希对象
	hash := fnv.New32()
	// 计算排序后的字节数组的哈希值
	hash.Write(bytes)
	hashValue := hash.Sum32()

	fmt.Println("Hash Value:", hashValue)
	return hashValue
}
func GetStringHash(str string) string {
	// 使用MD5哈希算法获取哈希值
	md5Hash := md5.Sum([]byte(str))
	md5HashString := hex.EncodeToString(md5Hash[:])
	fmt.Println("MD5:", md5HashString)

	// 使用SHA1哈希算法获取哈希值
	sha1Hash := sha1.Sum([]byte(str))
	sha1HashString := hex.EncodeToString(sha1Hash[:])
	fmt.Println("SHA1:", sha1HashString)
	return sha1HashString
	//// 使用SHA256哈希算法获取哈希值
	//sha256Hash := sha256.Sum256([]byte(str))
	//sha256HashString := hex.EncodeToString(sha256Hash[:])
	//fmt.Println("SHA256:", sha256HashString)
}
func GetStringHashInt(str string) uint32 {
	// 创建一个新的FNV哈希对象
	hash := fnv.New32()

	// 将字符串转换为字节数组并计算哈希值
	hash.Write([]byte(str))
	hashValue := hash.Sum32()

	fmt.Println("Hash Value:", hashValue)
	return hashValue
}
