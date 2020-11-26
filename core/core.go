package core

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/nfnt/resize"
	imgType "github.com/shamsher31/goimgtype"
	"github.com/skip2/go-qrcode"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"math"
	"net/http"
	"os"
)

//新建图片载体
// X1 宽
// Y1 高
func NewPNG(X0 int, Y0 int, X1 int, Y1 int) *image.RGBA {
	return image.NewRGBA(image.Rect(X0, Y0, X1, Y1))
}

//合并图片到载体
func MergeImage(dst draw.Image, image image.Image, sp image.Point) {
	draw.Draw(dst, dst.Bounds(), image, sp, draw.Over)
}

//新建文件载体
func NewOutPutMerged(path string) (*os.File, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return f, nil
}

//合并到图片
func JpegMerge(merged *os.File, dst draw.Image, quality int) error {
	err := jpeg.Encode(merged, dst, &jpeg.Options{Quality: quality})
	if err != nil {
		return err
	}
	return nil
}

//获取二维码图像
func DrawQRImage(url string, level qrcode.RecoveryLevel, size int) (image.Image, error) {
	newQr, err := qrcode.New(url, level)
	if err != nil {
		return nil, err
	}
	qrImage := newQr.Image(size)
	return qrImage, nil
}

///////////////////
//文字切片
type DText struct {
	PNG   draw.Image //合并到的PNG切片,可用image.NewrRGBA设置
	Title string     //文字
	X     int        //横坐标
	Y     int        //纵坐标
	Size  float64
	R     uint8
	G     uint8
	B     uint8
	A     uint8
}

//读取字体类型
func LoadTextType(path string) (*truetype.Font, error) {
	fbyte, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	trueTypeFont, err := freetype.ParseFont(fbyte)
	if err != nil {
		return nil, err
	}
	return trueTypeFont, nil
}

//创建新字体切片
func NewDrawText(png draw.Image) *DText {
	return &DText{
		PNG:  png,
		Size: 18,
		X:    0,
		Y:    0,
		R:    0,
		G:    0,
		B:    0,
		A:    255,
	}
}

//设置字体颜色
func (dtext *DText) SetColor(R uint8, G uint8, B uint8) {
	dtext.R = R
	dtext.G = G
	dtext.B = B
}

//合并字体到载体
func (dtext *DText) MergeText(title string, dpi, size float64, tf *truetype.Font, x int, y int) error {
	fc := freetype.NewContext()
	//设置屏幕每英寸的分辨率
	fc.SetDPI(dpi)
	//设置用于绘制文本的字体
	fc.SetFont(tf)
	//以磅为单位设置字体大小
	fc.SetFontSize(size)
	//设置剪裁矩形以进行绘制
	fc.SetClip(dtext.PNG.Bounds())
	//设置目标图像
	fc.SetDst(dtext.PNG)
	//设置绘制操作的源图像，通常为 image.Uniform
	fc.SetSrc(image.NewUniform(color.RGBA{dtext.R, dtext.G, dtext.B, dtext.A}))

	pt := freetype.Pt(x, y)
	_, err := fc.DrawString(title, pt)
	if err != nil {
		return err
	}
	return nil
}

////////////////////////////////

type processImg struct {
	URL, Path    string
	Scale        float64
	ZoomH, ZoomW int
	initImg      image.Image
}

func NewProcessImgService(url, path string, zoomW, zoomH int, scale float64) *processImg {
	return &processImg{
		URL:   url,
		Path:  path,
		Scale: scale,
		ZoomW: zoomW,
		ZoomH: zoomH,
	}
}

func (p *processImg) ProcessPathImgApp() (image.Image, error) {
	var img image.Image
	var err error
	if p.Path != "" {
		img, _, _, err = ProcessingLocalPictures(p.Path)
	} else if p.URL != "" {
		img, _, _, err = ProcessingRemotePictures(p.URL)
	}
	if err != nil {
		fmt.Printf("获取图片失败！【%v】", err)
		return img, err
	}
	//
	p.initImg = img
	//
	return img, nil
}

func (p *processImg) ProcessZoomImgApp() (image.Image, int, int, error) {
	zoomImg := p.initImg
	var ww, hh int
	if p.Scale != 0 {
		zoomImg, ww, hh = CalculateRatioFit(p.initImg, p.Scale)
	} else if p.ZoomW != 0 || p.ZoomH != 0 {
		zoomImg = ProportionalZoomImage(p.ZoomW, p.ZoomH, p.initImg)
		ww = p.ZoomW
		hh = p.ZoomH
	}
	if zoomImg == nil {
		return zoomImg, ww, hh, errors.New("获取压缩图片失败 aggregate core 174")
	}
	return zoomImg, ww, hh, nil
}

// 处理本地图片
func ProcessingLocalPictures(imgPath string) (image.Image, int, int, error) {
	var img image.Image
	var oriWidth, oriHeight int
	//
	imgFile, err := os.Open(imgPath)
	if err != nil {
		fmt.Printf("打开图片出错：[%v]", err)
		return img, oriWidth, oriHeight, err
	}
	defer imgFile.Close()
	// 获取图片的类型
	datatype, err := imgType.Get(imgFile.Name())
	if err != nil {
		fmt.Printf("不是图片文件：[%v]", err)
		return img, oriWidth, oriHeight, err
	}
	// 根据文件类型执行响应的操作
	switch datatype {
	case `image/jpeg`:
		img, err = jpeg.Decode(imgFile)
	case `image/png`:
		img, err = png.Decode(imgFile)
	}
	if err != nil {
		fmt.Printf("把图片解码为结构体时出错：[%v]", err)
		return img, oriWidth, oriHeight, err
	}
	//尺寸，宽 高
	oriWidth = img.Bounds().Max.X - img.Bounds().Min.X
	oriHeight = img.Bounds().Max.Y - img.Bounds().Min.Y
	//
	return img, oriWidth, oriHeight, nil
}

//处理远程图片
func ProcessingRemotePictures(imgURL string) (image.Image, int, int, error) {
	var oriWidth, oriHeight int
	srcReader, err := GetResourceReader(imgURL)
	if err != nil {
		return nil, oriWidth, oriHeight, err
	}
	srcImage, _, err := image.Decode(srcReader)
	if err != nil {
		fmt.Printf("远程图片处理错误 err：%v", err)
		return srcImage, oriWidth, oriHeight, err
	}
	//尺寸，宽 高
	oriWidth = srcImage.Bounds().Max.X - srcImage.Bounds().Min.X
	oriHeight = srcImage.Bounds().Max.Y - srcImage.Bounds().Min.Y
	//
	return srcImage, oriWidth, oriHeight, nil
}

// 等比例缩放计算图片缩放后的尺寸---以大的一边为基准
// 占原图的25%
// 压缩图，压缩尺寸
func CalculateRatioFit(img image.Image, scale float64) (image.Image, int, int) {
	//原始尺寸，宽 高
	oriWidth := img.Bounds().Max.X - img.Bounds().Min.X
	oriHeight := img.Bounds().Max.Y - img.Bounds().Min.Y
	//获取等比例压缩尺寸
	baseSize := oriWidth
	if oriWidth < oriHeight {
		baseSize = oriHeight
	}
	var resWidth, resHeight float64
	if baseSize == oriWidth {
		resWidth = float64(baseSize) * scale //现在的图宽
		ratio := resWidth / float64(oriWidth)
		resHeight = float64(oriHeight) * ratio
	} else if baseSize == oriHeight {
		resHeight = float64(baseSize) * scale //现在的图高
		ratio := resHeight / float64(oriHeight)
		resWidth = float64(oriWidth) * ratio
	}
	//压缩图片
	zoomImg := ProportionalZoomImage(int(math.Ceil(resWidth)), int(math.Ceil(resHeight)), img)
	//
	return zoomImg, int(math.Ceil(resWidth)), int(math.Ceil(resHeight))
}

// 圆形返回合适的直径
// 圆的直径以短边为准
func CalculateCircleDiameter(oriWidth, oriHeight int) int {
	//
	diameter := oriWidth
	if oriWidth > oriHeight {
		diameter = oriHeight
	}
	return diameter
}

// 等比例缩放图片
func ProportionalZoomImage(width, height int, oriImage image.Image) image.Image {
	return resize.Resize(uint(width), uint(height), oriImage, resize.Lanczos3)
}

// 根据http 连接获取图片处理
func GetResourceReader(url string) (r *bytes.Reader, err error) {
	if url[0:4] == "http" {
		resp, err := http.Get(url)
		if err != nil {
			return r, err
		}
		defer resp.Body.Close()
		fileBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return r, err
		}
		r = bytes.NewReader(fileBytes)
	} else {
		fileBytes, err := ioutil.ReadFile(url)
		if err != nil {
			return nil, err
		}
		r = bytes.NewReader(fileBytes)
	}
	return r, nil
}
