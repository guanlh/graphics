package aggregate

import (
	"fmt"
	"graphics/core"
	"image"
)

//
type BackgroundWork struct {
	// 合成复用Next
	Next
	X     int
	Y     int
	Path  string
	URL   string
	Scale float64 //压缩尺寸
	ZoomW int     //需要压缩的宽
	ZoomH int     //需要压缩的高
}

// Do 地址逻辑
func (w *BackgroundWork) Do(i *ImageContext) error {
	//
	processAppSvc := core.NewProcessImgService(w.URL, w.Path, w.ZoomW, w.ZoomH, w.Scale)
	//图像处理
	_, err := processAppSvc.ProcessPathImgApp()
	if err != nil {
		fmt.Printf("获取背景图失败！【%v】", err)
		return err
	}
	//压缩
	zoomBgImg, _, _, err := processAppSvc.ProcessZoomImgApp()
	if err != nil {
		return err
	}
	//坐标
	bgPoint := image.Point{
		X: w.X,
		Y: w.Y,
	}
	//合并
	core.MergeImage(i.ImgCarrier, zoomBgImg, zoomBgImg.Bounds().Min.Sub(bgPoint))
	//
	return nil
}
