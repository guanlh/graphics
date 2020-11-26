package graphics

import (
	"fmt"
	"graphics/aggregate"
	"graphics/core"
	"testing"
	"time"
)

func TestEntrance(t *testing.T) {
	//
	t1 := time.Now().UnixNano() / 1e6
	nullAggregate := aggregate.NullAggregate{}
	//新建图片载体
	imgCtx := &aggregate.ImageContext{
		ImgCarrier: core.NewPNG(0, 0, 500, 1000),
	}
	//
	//绘制背景图----的zoomW和zoomY应该和载体的X1和Y1对应
	backgroundWork := &aggregate.BackgroundWork{
		X: 0,
		Y: 0,
		//Path:  "./img/dt.jpg",
		URL:   "https://cf-markting-test.oss-cn-hangzhou.aliyuncs.com/W00000000109/material/poster/185f1850-d5ee-4efe-b471-88b8b5752c65/jjud061dxyf1606244674137.jpg",
		Scale: 0,
		ZoomW: 500,
		ZoomH: 1000,
	}
	//绘制圆形图像
	imageCircleWork := &aggregate.ImageCircleWork{
		X: 30,
		Y: 50,
		//Path:  "./img/0.jpeg",
		URL:   "http://wework.qpic.cn/bizmail/JmYIOjuHGC1BJFHNgqyeQnwy7xrnb0jnjZ56fxNqicrCpuHKbncdAvQ/0",
		Scale: 0,
		ZoomW: 100,
		ZoomH: 100,
	}
	//绘制图像
	imageWork := &aggregate.ImageWork{
		X: 30,
		Y: 800,
		//Path:  "./img/ewm.png",
		URL:   "http://cf-qw-test.oss-cn-hangzhou.aliyuncs.com/wwpic/293512_-wGAgpLtTCK1jPG_1606244686/0",
		Scale: 0,
		ZoomW: 100,
		ZoomH: 100,
	}
	//绘制文字
	textWork1 := &aggregate.TextWork{
		X:        130,
		Y:        120,
		Size:     25,
		Dpi:      80,
		R:        255,
		G:        255,
		B:        255,
		Text:     "朱启铭",
		FontPath: "./img/dustess.ttf",
	}
	//结束绘制，把前面的内容合并成一张图片
	endWork := &aggregate.EndWork{
		Output:  "./img/poster.png",
		Quality: 100,
	}

	// 链式调用绘制过程
	nullAggregate.
		SetNext(backgroundWork).
		SetNext(imageCircleWork).
		SetNext(imageWork).
		SetNext(textWork1).
		SetNext(endWork)

	// 开始执行业务
	if err := nullAggregate.Run(imgCtx); err != nil {
		// 异常
		panic("Fail | Error:" + err.Error())
		return
	}
	t2 := time.Now().UnixNano() / 1e6
	fmt.Println(t2 - t1)
	fmt.Println("合成成功！")
	return
}
