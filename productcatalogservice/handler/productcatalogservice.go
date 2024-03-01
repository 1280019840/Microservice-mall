package handler

import (
	"bytes"
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"

	pb "productcatalogservice/proto"
)

var reloadCatalog bool

// 日志
var (
	buf    bytes.Buffer
	logger = log.New(&buf, "logger: ", log.Lshortfile)
)

// 商品分类结构体
type ProductCatalogService struct {
	sync.Mutex
	products []*pb.Product
}

// 商品列表
func (s *ProductCatalogService) ListProducts(ctx context.Context, in *pb.Empty) (out *pb.ListProductsResponse, e error) {
	out = new(pb.ListProductsResponse)
	out.Products = s.parseCatalog()
	return out, nil
}

// 获得单个商品
func (s *ProductCatalogService) GetProduct(ctx context.Context, in *pb.GetProductRequest) (out *pb.Product, e error) {
	var found *pb.Product
	out = new(pb.Product)
	products := s.parseCatalog()
	for _, p := range products {
		if in.Id == p.Id {
			found = p
		}
	}
	if found == nil {
		return out, status.Errorf(codes.NotFound, "no product with ID %s", in.Id)
	}
	out.Id = found.Id
	out.Name = found.Name
	out.Categories = found.Categories
	out.Description = found.Description
	out.Picture = found.Picture
	out.PriceUsd = found.PriceUsd
	return out, nil
}

// 搜索商品
func (s *ProductCatalogService) SearchProducts(ctx context.Context, in *pb.SearchProductsRequest) (out *pb.SearchProductsResponse, e error) {
	var ps []*pb.Product
	products := s.parseCatalog()
	for _, p := range products {
		if strings.Contains(strings.ToLower(p.Name), strings.ToLower(in.Query)) ||
			strings.Contains(strings.ToLower(p.Description), strings.ToLower(in.Query)) {
			ps = append(ps, p)
		}
	}
	out.Results = ps
	return out, nil
}

// 读配置文件
func (s *ProductCatalogService) readCatalogFile() (*pb.ListProductsResponse, error) {
	s.Lock()
	defer s.Unlock()
	catalogJSON, err := os.ReadFile("data/products.json")
	if err != nil {
		logger.Printf("打开商品 json 文件失败: %v", err)
		return nil, err
	}
	catalog := &pb.ListProductsResponse{}
	if err := protojson.Unmarshal(catalogJSON, catalog); err != nil {
		logger.Printf("解析商品 JSON 文件失败: %v", err)
		return nil, err
	}
	logger.Printf("解析商品 JSON 文件成功")
	return catalog, nil
}

// 解析配置文件
func (s *ProductCatalogService) parseCatalog() []*pb.Product {
	if reloadCatalog || len(s.products) == 0 {
		catalog, err := s.readCatalogFile()
		if err != nil {
			return []*pb.Product{}
		}
		s.products = catalog.Products
	}
	return s.products
}

// 初始化
func init() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGUSR1, syscall.SIGUSR2)
	go func() {
		for {
			sig := <-sigs
			logger.Printf("接收信号: %s", sig)
			if sig == syscall.SIGUSR1 {
				reloadCatalog = true
				logger.Printf("可以加载商品信息")
			} else {
				reloadCatalog = false
				logger.Printf("不能加载商品信息")
			}
		}
	}()
}
