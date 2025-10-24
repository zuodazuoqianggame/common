package errorCode

const (
	Sucess     = 0   //请求成功
	ServiceErr = 500 //微服异常

	ParameterInvalid = 1801 //参数异常
	AuthFails        = 1802 //授权失败
	AuthInvalid      = 1803 //鉴权失败
	NoUser           = 1804 //没有此用户
	DbErr            = 1805 //数据库操作失败
	ConfErr          = 1806 //配置错误
	ThridErr         = 1807 //第三方服务器异常
	RequestInvalid   = 1808 //非法请求
	SystemError      = 1809 //服务器的错误
	NoFoundData      = 1810 //没有发现数据
)
