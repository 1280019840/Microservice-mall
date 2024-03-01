package handler

import (
	"context"
	"math/rand"

	pb "adservice/proto"
)

// 最大展示广告数
const MAX_ADS_TO_SERVE = 2

// 调用后的函数，返回一个map广告数据
var adsMap = createAdsMap()

// 广告服务结构体
type AdService struct{}

// 获得广告方法，传递Context和请求参数，返回响应和错误
func (s *AdService) GetAds(context context.Context, in *pb.AdRequest) (out *pb.AdResponse, err error) {
	// 所有广告切片
	allAds := make([]*pb.Ad, 0)
	// 根据广告种类选择广告
	if len(in.ContextKeys) > 0 {
		// 遍历分类
		for _, category := range in.ContextKeys {
			// 根据分类获得广告
			ads := getAdsByCategory(category)
			// 添加到广告切片
			allAds = append(allAds, ads...)
		}
		// 如果没有随机获取
		if len(allAds) == 0 {
			allAds = getRandomAds()
		}
	} else {
		// 如果没有随机获取
		allAds = getRandomAds()
	}
	// 输出
	out = new(pb.AdResponse)
	// 输出携带广告数据
	out.Ads = allAds
	// 返回
	return out, nil
}

// 根据分类获得广告
func getAdsByCategory(category string) []*pb.Ad {
	return adsMap[category]
}

// 随机获得广告
func getRandomAds() []*pb.Ad {
	ads := make([]*pb.Ad, 0, MAX_ADS_TO_SERVE)
	allAds := make([]*pb.Ad, 0, 7)
	for _, ads := range adsMap {
		allAds = append(allAds, ads...)
	}
	for i := 0; i < MAX_ADS_TO_SERVE; i++ {
		ads = append(ads, allAds[rand.Intn(len(allAds))])
	}
	return ads
}

// 创建广告（也可以查询数据库）
func createAdsMap() map[string][]*pb.Ad {
	hairdryer := &pb.Ad{RedirectUrl: "/product/2ZYFJ3GM2N", Text: "出风机，5折热销"}
	tankTop := &pb.Ad{RedirectUrl: "/product/66VCHSJNUP", Text: "背心8折热销"}
	candleHolder := &pb.Ad{RedirectUrl: "/product/0PUK6V6EV0", Text: "烛台7折热销"}
	bambooGlassJar := &pb.Ad{RedirectUrl: "/product/9SIQT8TOJO", Text: "竹玻璃罐9折"}
	watch := &pb.Ad{RedirectUrl: "/product/1YMWWN1N4O", Text: "手表买一送一"}
	mug := &pb.Ad{RedirectUrl: "/product/6E92ZMYYFZ", Text: "马克杯买二送一"}
	loafers := &pb.Ad{RedirectUrl: "/product/L9ECAV7KIM", Text: "平底鞋，买一送二"}
	return map[string][]*pb.Ad{
		"clothing":    {tankTop},
		"accessories": {watch},
		"footwear":    {loafers},
		"hair":        {hairdryer},
		"decor":       {candleHolder},
		"kitchen":     {bambooGlassJar, mug},
	}
}
