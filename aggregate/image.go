package aggregate

import (
	"graphics/core"
	"image"
)

// ImageWork
type ImageWork struct {
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
func (w *ImageWork) Do(i *ImageContext) error {
	//
	processAppSvc := core.NewProcessImgService(w.URL, w.Path, w.ZoomW, w.ZoomH, w.Scale)
	//图像处理
	_, err := processAppSvc.ProcessPathImgApp()
	if err != nil {
		return err
	}
	//压缩
	zoomImg, _, _, err := processAppSvc.ProcessZoomImgApp()
	if err != nil {
		return err
	}
	//坐标
	srcPoint := image.Point{
		X: w.X,
		Y: w.Y,
	}
	core.MergeImage(i.ImgCarrier, zoomImg, zoomImg.Bounds().Min.Sub(srcPoint))
	//
	return nil
}
