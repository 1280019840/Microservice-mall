package cartstore

import (
	"context"
	"sync"

	pb "cartservice/proto"
)

// 数据保存在内存中结构体，使用嵌套map保存
type memoryCartStore struct {
	// 读写锁
	sync.RWMutex
	carts map[string]map[string]int32
}

// 添加商品
func (s *memoryCartStore) AddItem(ctx context.Context, userID, productID string, quantity int32, out *pb.Empty) (r *pb.Empty, err error) {
	s.Lock()
	defer s.Unlock()

	if cart, ok := s.carts[userID]; ok {
		if currentQuantity, ok := cart[productID]; ok {
			cart[productID] = currentQuantity + quantity
		} else {
			cart[productID] = quantity
		}
		s.carts[userID] = cart
	} else {
		s.carts[userID] = map[string]int32{productID: quantity}
	}
	return out, nil
}

// 清空购物车
func (s *memoryCartStore) EmptyCart(ctx context.Context, userID string) (out *pb.Empty, err error) {
	s.Lock()
	defer s.Unlock()
	out = new(pb.Empty)
	delete(s.carts, userID)
	return out, nil
}

// 获得购物车
func (s *memoryCartStore) GetCart(ctx context.Context, userID string) (*pb.Cart, error) {
	s.RLock()
	defer s.RUnlock()

	if cart, ok := s.carts[userID]; ok {
		items := make([]*pb.CartItem, 0, len(cart))
		for p, q := range cart {
			items = append(items, &pb.CartItem{ProductId: p, Quantity: q})
		}
		return &pb.Cart{UserId: userID, Items: items}, nil
	}
	return &pb.Cart{UserId: userID}, nil
}
