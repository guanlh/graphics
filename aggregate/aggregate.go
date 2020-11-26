package aggregate

import "image"

// ImageContext
type ImageContext struct {
	ImgCarrier *image.RGBA
}

//  处理
type Aggregate interface {
	// 自身的业务
	Do(i *ImageContext) error
	// 设置下一个对象
	SetNext(i Aggregate) Aggregate
	// 执行
	Run(i *ImageContext) error
}

// Next 抽象出来的 可被合成复用的结构体
type Next struct {
	// 下一个对象
	nextStep Aggregate
}

// SetNext 实现好的 可被复用的SetNext方法
// 返回值是下一个对象 方便写成链式代码优雅
// 例如 NullAggregate.SetNext(argumentsAggregate).SetNext(signAggregate).SetNext(frequentAggregate)
func (n *Next) SetNext(a Aggregate) Aggregate {
	n.nextStep = a
	return a
}

// Run 执行
func (n *Next) Run(i *ImageContext) (err error) {
	//  这里无法执行当前aggregate的Do
	// n.Do(i)
	if n.nextStep != nil {
		// 合成复用下的变种
		// 执行下一个的Do
		if err = (n.nextStep).Do(i); err != nil {
			return
		}
		// 执行下一个的Run
		return (n.nextStep).Run(i)
	}
	return
}

type NullAggregate struct {
	// 合成复用Next的`nextStep`成员属性、`SetNext`成员方法、`Run`成员方法
	Next
}

// Do 空的Do
func (w *NullAggregate) Do(i *ImageContext) (err error) {
	return
}
