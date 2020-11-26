package aggregate

import (
	"errors"
	"fmt"
	"github.com/skip2/go-qrcode"
	"graphics/core"
	"image"
)

// QRCodeWork
type QRCodeWork struct {
	// 合成复用Next
	Next
	X     int
	Y     int
	URL   string
	Level int //Error detection/recovery capacity. (0-7%,1-15%,2-25%,3-30%)
	Size  int
}

var QrCodeLevel map[int]qrcode.RecoveryLevel

func init() {
	QrCodeLevel = map[int]qrcode.RecoveryLevel{
		0: qrcode.Low,
		1: qrcode.Medium,
		2: qrcode.High,
		3: qrcode.Highest,
	}
}

// Do 地址逻辑
func (w *QRCodeWork) Do(i *ImageContext) error {
	//生成二维码
	value, ok := QrCodeLevel[w.Level]
	if !ok {
		return errors.New("Error detection/recovery capacity ")
	}
	//
	qrImage, err := core.DrawQRImage(w.URL, value, w.Size)
	if err != nil {
		fmt.Printf(" 生成二维码错误 err：%v", err)
		return err
	}
	// 把二维码合并到pngCarrier
	qrPoint := image.Point{X: w.X, Y: w.Y}
	core.MergeImage(i.ImgCarrier, qrImage, qrImage.Bounds().Min.Sub(qrPoint))
	//
	return nil
}
