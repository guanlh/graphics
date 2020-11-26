package aggregate

import (
	"fmt"
	"graphics/core"
)

// TextWork
type TextWork struct {
	// 合成复用Next
	Next
	X        int
	Y        int
	Size     float64
	Dpi      float64
	R        uint8
	G        uint8
	B        uint8
	Text     string
	FontPath string
}

// Do 地址逻辑
func (w *TextWork) Do(i *ImageContext) error {
	//设置字体切片
	if w.Size == 0 {
		w.Size = 30
	}
	trueTypeFont, err := core.LoadTextType(w.FontPath)
	if err != nil {
		fmt.Printf("加载字体失败 err：%v", err)
		return err
	}
	//
	dText := core.NewDrawText(i.ImgCarrier)
	//设置颜色
	dText.SetColor(w.R, w.G, w.B)
	err = dText.MergeText(w.Text, w.Dpi, w.Size, trueTypeFont, w.X, w.Y)
	if err != nil {
		fmt.Printf("设置字失败 err：%v", err)
		return err
	}
	return nil
}
