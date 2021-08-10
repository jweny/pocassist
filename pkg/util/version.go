package util

import (
	"strconv"
	"strings"
)

// SingleVersionCompare 版本对比方法
// <=0没有漏洞    >0有漏洞
func SingleVersionCompare(verCurrent string, verVul string) (int, error) {
	partsCurrent := strings.Split(verCurrent, ".")
	partsVul := strings.Split(verVul, ".")
	// 比较长度 取短的一方
	var parts []string
	if len(partsCurrent) < len(partsVul) {
		//return 0, errors.New("two version have different parts length")
		parts = partsCurrent
	} else {
		parts = partsVul
	}
	for i := range parts {
		partCurrent, err1 := strconv.Atoi(partsCurrent[i])
		partVul, err2 := strconv.Atoi(partsVul[i])
		if err1 != nil {
			return -2, err1
		}
		if err2 != nil {
			return -2, err2
		}
		if partCurrent == partVul {
			continue
		} else {
			return partVul-partCurrent, nil
		}
	}
	return 0, nil
}