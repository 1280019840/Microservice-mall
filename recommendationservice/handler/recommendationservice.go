package handler

import (
	"bytes"
	"context"
	"log"
	"math/rand"

	pb "recommendationservice/proto"
)

// 日志
var (
	buf    bytes.Buffer
	logger = log.New(&buf, "logger: ", log.Lshortfile)
)

// 推荐服务结构体
type RecommendationService struct {
	ProductCatalogService pb.ProductCatalogServiceClient
}

// 列出推荐
func (s *RecommendationService) ListRecommendations(ctx context.Context, in *pb.ListRecommendationsRequest) (out *pb.ListRecommendationsResponse, e error) {
	maxResponsesCount := 5
	out = new(pb.ListRecommendationsResponse)
	// 查询商品类别
	catalog, err := s.ProductCatalogService.ListProducts(ctx, &pb.Empty{})
	if err != nil {
		return out, err
	}
	filteredProductsIDs := make([]string, 0, len(catalog.Products))
	for _, p := range catalog.Products {
		if contains(p.Id, in.ProductIds) {
			continue
		}
		filteredProductsIDs = append(filteredProductsIDs, p.Id)
	}
	productIDs := sample(filteredProductsIDs, maxResponsesCount)
	logger.Printf("[Recv ListRecommendations] product_ids=%v", productIDs)
	out.ProductIds = productIDs
	return out, nil
}

// 判断是否包含
func contains(target string, source []string) bool {
	for _, s := range source {
		if target == s {
			return true
		}
	}
	return false
}

// 示例
func sample(source []string, c int) []string {
	n := len(source)
	if n <= c {
		return source
	}
	indices := make([]int, n)
	for i := 0; i < n; i++ {
		indices[i] = i
	}
	for i := n - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		indices[i], indices[j] = indices[j], indices[i]
	}
	result := make([]string, 0, c)
	for i := 0; i < c; i++ {
		result = append(result, source[indices[i]])
	}
	return result
}
