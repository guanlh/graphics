package roundshadow

import (
	"fmt"
	"graphics/aggregate"
	"graphics/core"
	"image"
	"image/png"
	"os"
	"testing"
)

// 生成圆形图片
func TestNewRoundShadow(t *testing.T) {
	file, err := os.Create("../img/circle.png")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	// 处理图片
	// 获取原图，原始尺寸
	//oriImg, oriWidth, oriHeight, err := core.ProcessImagesAccordingToType("../img/ewm.png")
	oriImg, _, _, err := core.ProcessingLocalPictures("../img/ewm.png")
	// 获取压缩图，压缩尺寸
	if err != nil {
		fmt.Println(err)
		return
	}
	// 获取压缩图，压缩尺寸
	zoomImg, ww, hh := core.CalculateRatioFit(oriImg, 0.5)
	// 圆形的直径（压缩）
	d := core.CalculateCircleDiameter(ww, hh)
	//
	maskImg := NewRoundShadow(zoomImg, image.Point{X: 0, Y: 0}, d)
	//
	imgCtx := &aggregate.ImageContext{
		ImgCarrier: core.NewPNG(0, 0, 0, 0),
	}
	bgPoint := image.Point{
		X: 1,
		Y: 2,
	}
	core.MergeImage(imgCtx.ImgCarrier, zoomImg, zoomImg.Bounds().Min.Sub(bgPoint))
	//图片写入
	_ = png.Encode(file, maskImg)
	//if err6 := png.Encode(file, maskImg); err6 != nil {
	//	fmt.Printf("设置质量图片质量错误：[%v]", err6)
	//	return
	//}
}
