作者：陈鼎
ConfigCreator需要的Go依赖，内容被我重新修改过。

1、captcha.go：
const (
	// 生成数字的长度
	DefaultLen = 4
)

2、font.go:
const (
	// 字体宽度
	fontWidth  = 11
	// 字体高度
	fontHeight = 18
)


3、image.go:
const (
	// 图片长宽
	StdWidth  = 160
	StdHeight = 80
	// 每个数字的最大倾斜因子
	maxSkew = 0.7
	// 背景原点个数
	circleCount = 20
)