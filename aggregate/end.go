package aggregate

import (
	"fmt"
	"graphics/core"
)

// EndWork 结束，写在最后，把图片合并到一张图上
type EndWork struct {
	// 合成复用Next
	Next
	Output  string
	Quality int
}

// Do 地址逻辑
func (w *EndWork) Do(i *ImageContext) error {
	// 新建文件载体
	merged, err := core.NewOutPutMerged(w.Output)
	if err != nil {
		fmt.Printf("新建最后输出图片的文件错误  err：%v", err)
		return err
	}
	// 合并
	err = core.JpegMerge(merged, i.ImgCarrier, w.Quality)
	if err != nil {
		fmt.Printf("新建最后输出图片错误  err：%v", err)
		return err
	}
	return nil
}
