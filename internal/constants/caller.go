package constants

type Caller string

const (
	CallerZltx Caller = "zltx" // 直连天下
	CallerHyt  Caller = "hyt"  // 货易通
)

func (c Caller) String() string {
	return string(c)
}
