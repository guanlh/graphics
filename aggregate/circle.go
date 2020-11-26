package aggregate

import (
	"graphics/core"
	"graphics/roundshadow"
	"image"
)

// ImageCircleWork 根据Path路径设置圆形图片
type ImageCircleWork struct {
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
func (w *ImageCircleWork) Do(i *ImageContext) error {
	//
	processAppSvc := core.NewProcessImgService(w.URL, w.Path, w.ZoomW, w.ZoomH, w.Scale)
	//图像处理
	_, err := processAppSvc.ProcessPathImgApp()
	if err != nil {
		return err
	}
	//压缩
	zoomCirCleImg, ww, hh, err := processAppSvc.ProcessZoomImgApp()
	if err != nil {
		return err
	}
	// 圆形的直径（压缩）
	diameter := core.CalculateCircleDiameter(ww, hh)
	//把头像转成Png,否则会有白底
	srcPng := core.NewPNG(0, 0, ww, hh)
	//
	core.MergeImage(srcPng, zoomCirCleImg, zoomCirCleImg.Bounds().Min)
	// 遮罩
	srcMask := roundshadow.NewRoundShadow(srcPng, image.Point{0, 0}, diameter)
	//坐标
	srcPoint := image.Point{
		X: w.X,
		Y: w.Y,
	}
	core.MergeImage(i.ImgCarrier, srcMask, zoomCirCleImg.Bounds().Min.Sub(srcPoint))
	//
	return nil
}
