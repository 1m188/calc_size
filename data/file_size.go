package data

import (
	"math/big"
)

// 文件/目录大小信息
type FileSize struct {
	Size big.Int // 大小（字节）
	Path string  // 路径
}

/*
获取大小（带单位）
  - f: 几位小数
  - prec: 精度
  - 返回值：数字的字符串，单位
*/
func (fs *FileSize) GetSizeWithUnit(f int, prec uint) (string, string) {

	size := big.NewFloat(0).SetPrec(prec).SetInt(&fs.Size)
	x := big.NewFloat(1024).SetPrec(prec)
	unit := "B"

	if size.Cmp(x) > 0 {
		size.Quo(size, x)
		unit = "KB"

		if size.Cmp(x) > 0 {
			size.Quo(size, x)
			unit = "MB"

			if size.Cmp(x) > 0 {
				size.Quo(size, x)
				unit = "GB"

				if size.Cmp(x) > 0 {
					size.Quo(size, x)
					unit = "PB"
				}
			}
		}
	}

	if size.Cmp(big.NewFloat(0).SetPrec(prec)) < 0 {
		unit = ""
	}

	return size.Text('f', f), unit
}
