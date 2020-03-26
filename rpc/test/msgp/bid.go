//go:generate msgp

package msgp

import (
	"time"
)

// 渠道适配器请求数据
type AdapterRequest struct {
	// 渠道广告位id
	AdxSpotId string `msg:"1"`
	// 渠道id
	AdxId int `msg:"2"`
	// 主渠道id，如果没有，则需要填充AdxId在这
	// MainAdxId int32 `msg:"59"`
	// 渠道用户id
	AdxUID string `msg:"3"`
	// 渠道用户id版本
	// AdxUIDV string `msg:"4"`
	// 曝光序列id
	//SequenceId string `msg:"5"`
	// 内部请求id
	ReqId string `msg:"6"`

	// 竞拍算法:1.明拍;2.暗拍;4.固定CPM
	// Ot int
	// cookieMapping user id
	CmUserId string `msg:"7"`
	// 当前页面URL
	Url string `msg:"8"`
	// 请求 refer
	Referer string `msg:"9"`
	// 访客关键词
	Search string `msg:"64"`

	// 宽
	Width int `msg:"10"`
	// 高
	Height int `msg:"11"`
	// 广告位尺寸 大小为width*10000+height
	SpotSize int `msg:"12"`
	// 最小视频播放时长
	Minduration int `msg:"13"`
	// 最大视频播放时长
	Maxduration int `msg:"14"`
	// 视频广告的播出延时，0及正整数表示前贴，-1表示中贴，-2表示后贴
	StartDelay int `msg:"15"`
	// 是否为线性视频，线性视频是插播视频，非线性是占用部分视频位置加入的视频广告
	Linearity bool `msg:"62"`
	// 视频缩放比例 e.g [4*3] 用英文*号分割
	Zoom []string `msg:"63"`

	// 广告位底价
	Bidfloor int64 `msg:"16"`
	// 广告主底价币种
	// Bidfloorcur string `msg:"17"`
	// 展示类型
	ViewType int `msg:"18"`
	// 展示位置
	Screen int `msg:"19"`
	// 发布商允许的创意类型
	CreativeType []int `msg:"20"`
	// 允许的创意物料类型
	MaterialRenderType int `msg:"21"`
	// 网站分类
	Category int `msg:"22"`

	// 浏览器UA
	Ua string `msg:"23"`
	// 访客IP
	Ip string `msg:"24"`
	// 访客IPv6
	Ipv6 string `msg:"25"`
	// 维度
	GeoLat float64 `msg:"26"`
	// 经度
	GeoLon float64 `msg:"27"`
	// 地图坐标标准
	Standard int `msg:"28"`
	// 运营商或者服务提供商
	Carrier int `msg:"29"`
	// 浏览器语言
	Language int `msg:"30"`
	// 端浏览器
	Browser int `msg:"31"`
	// 设备品牌
	Brand int `msg:"33"`
	// 设备机型
	Model int `msg:"34"`
	// 设备操作系统
	Os int `msg:"36"`
	// 设备操作系统版本
	Osv string `msg:"37"`
	// 设备屏幕宽度
	// DevWidth int `msg:"38"`
	// 设备屏幕高度
	// DevHeight int `msg:"39"`
	// 设备网络类型
	Network int `msg:"40"`
	// 设备类型
	DeviceType int `msg:"41"`

	// 设备id类型
	DeviceIdType int `msg:"42"`
	// 竞价可以获取到的各种设备id
	// DeviceIds          map[string]string
	Imei DeviceFingerprint `msg:"43"`
	Did  DeviceFingerprint `msg:"44"`
	Mac  DeviceFingerprint `msg:"45"`
	GaId string            `msg:"46"`

	// 应用id
	AppBundleId string `msg:"47"`
	// 应用名称
	AppName string `msg:"65"`
	// 应用版本号
	Ver string `msg:"66"`
	// 应用包名或bundle
	Bundle string `msg:"67"`
	// 应用下载地址
	// StoreUlr string `msg:"68"`
	// 应用允许的交互类型
	AppInteractionType []int `msg:"69"`
	// 应用是否第一次启动
	// FirstLaunch bool `msg:"48"`

	// 订单数据
	Deals []Deal `msg:"49"`
	// 不允许的广告主名单，为空则所有均允许
	ExcludedAdvertisers []int `msg:"50"`
	// 排除条件，符合条件的物料均不允许投放
	ExcludedList []ExcludeList `msg:"51"`
	// 图文创意模版id
	TemplateIds []int `msg:"52"`

	// 客户端类型：1.web 2.wap 3.app
	ClientType int `msg:"53"`

	// 国家id
	CountryId int `msg:"54"`
	// 省id
	ProvinceId int `msg:"55"`
	// 市id
	CityId int `msg:"56"`
	// 县id
	// DistrictId int `msg:"57"`

	// 是否要求ssl
	SSL bool `msg:"58"`

	// 允许的创意级别，目前有四个级，0表示创意五审核评级也能投放。1为最高级。
	// CreativeLevel []int `msg:"61"`
	// 过滤规则id
	PublisherFilterids []string `msg:"70"`

	// 竞价规则（现在爱奇艺用）
	BidRule []BidRuleType `msg:"60"`

	// 系统广告位id
	SysSpotId int `msg:"71"`
	// 自有广告位id
	CustomSpotIds []int `msg:"72"`
	// 渠道id和渠道广告位拼接, 拼接符用 “|”
	AdxIdAndAdxSpotId string `msg:"73"`

	// 渠道请求id
	AdxReqId string `msg:"74"`
	// 渠道曝光id
	AdxImpId string `msg:"75"`
	// 内部请求id
	ImpId string `msg:"76"`
	// visitorId
	// PC：cookie id;
	// mobile：did > imei > md5did > md5imei > sha1did > sha1imei
	VisitorId string `msg:"77"`
	// 用户信息
	UserInfo UserInfo `msg:"78"`
	// 来源信息
	SourceInfo *SourceInfo        `msg:"-"`
	Debug      *RequestDebugParam `msg:"90"`
}

type RequestDebugParam struct {
	Tags map[string]interface{} `msg:"1"`
	Info struct {
		CampaignId []int `msg:"1"`
		PlanId     []int `msg:"2"`
		CreativeId []int `msg:"3"`
		WhiskyId   []int `msg:"4"`
		BudgetId   []int `msg:"5"`
		UserId     []int `msg:"6"`
	} `msg:"2"`
	Algorithm struct {
		Version     string `msg:"1"`
		GovernorOFF bool   `msg:"2"`
		StatRateOFF bool   `msg:"3"`
	} `msg:"3"`
	Runtime struct {
		PanicON bool `msg:"1"`
	} `msg:"4"`
}

type ResponseDebugParam struct {
	Panic string `msg:"1"`
}

// 竞价核心返回数据
type BidResponse struct {
	// 返回结果状态代码
	// 1 出价，0 不出价
	Code int `msg:"0"`
	// 计划id
	PlanId int `msg:"1"`
	// 活动id
	CampaignId int `msg:"2"`
	// 创意包
	// SweetyId int `msg:"3"`
	// 创意id
	CreativeId int `msg:"4"`
	// 落地页id
	WhiskyId int `msg:"5"`
	// 产品id
	// ProductId int `msg:"6"`
	// dsp用户id
	DspUserId int `msg:"7"`
	// 出价数量
	Cost int64 `msg:"8"`
	// 溢价后的价格
	// PremiumCost int64 `msg:"9"`

	// 资质id
	QualificationId int `msg:"10"`
	// 曝光id
	ImpressionId string `msg:"11"`

	// 订单Id
	// DealId int `msg:"12"`

	// 算法版本
	AlgoVersion int `msg:"13"`
	// 自有广告位
	SpotId int `msg:"14"`
	// 出价类型:CPC/CPM
	OfferType int                 `msg:"15"`
	BudgetId  int                 `msg:"16"`
	Debug     *ResponseDebugParam `msg:"17"`
}

// 竞价规则
type BidRuleType struct {
	BidFloor float64 `msg:"1"`
	Category []int   `msg:"2"`
}

// 排除列表
type ExcludeList struct {
	CategoryIds  []int    `msg:"1"`
	LandingPages []string `msg:"2"`
	Relationship int      `msg:"3"` // 关系 （0 与） （1 或）
	AdxCategory  []string `msg:"4"` // 渠道行业分类
}

type Deal struct {
	PreferredDeal  PreferredDeal  `msg:"1"`
	PrivateAuction PrivateAuction `msg:"2"`
	Type           int            `msg:"3"`
}

type PreferredDeal struct {
	DealId         int   `msg:"1"`
	AdvertiserIds  []int `msg:"2"`
	FixedPrice     int   `msg:"3"`
	FixedPriceUnit int   `msg:"4"`
}

type PrivateAuction struct {
	DealId     int         `msg:"1"`
	BuyerRules []BuyerRule `msg:"2"`
}

type BuyerRule struct {
	AdvertiserIds []int `msg:"1"`
	MinCpmPrice   int   `msg:"2"`
}

// 设备指纹
type DeviceFingerprint struct {
	Value  string `msg:"1"` // 原值
	MValue string `msg:"2"` // md5
	SValue string `msg:"3"` // sha1
}

// 用户信息
type UserInfo struct {
	// 用户年龄
	Age int `msg:"1"`
	// 性别
	// 男: constant.USER_GENDER_MALE；
	// 女: constant.USER_GENDER_FEMALE；
	// 其它: constant.USER_GENDER_OTHER；
	Gender int `msg:"2"`
	// 渠道用户关键字， 兴趣或者意向列表
	KeyWords []string `msg:"3"`
}

// 来源信息
type SourceInfo struct {
	IpInt   uint32    // ip
	Time    time.Time // 请求竞价当前时间实例
	Date    int       // int类型的timeNow 格式：yyyyddmm
	Weekday int       // 星期对应的数字
	Hour    int       // 小时
	Size    int       // 尺寸
	Zone    string    // 地区
}

// 批量请求竞价核心
type MultipleAdapterRequestItem struct {
	SequenceId int
	Request    *AdapterRequest
}
type MultipleAdapterRequest []*MultipleAdapterRequestItem

// 批量返回
type MultiResponse struct {
	SequenceId int
	Response   *BidResponse
}

type MultiBidResponse []*MultiResponse
