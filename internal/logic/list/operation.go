package list

import (
	"context"
)

func (s *sList) UnionOp(ctx context.Context, A, B, C string) (n int, err error) {
	// 获取 list 数据
	listA := s.GetListData(ctx, A)
	listB := s.GetListData(ctx, B)
	// 数据处理 并集运算
	listC := make(map[string]any)
	for k, v := range listA {
		listC[k] = v
	}
	for k, v := range listB {
		listC[k] = v
	}
	// 保存数据
	return s.AppendListData(ctx, C, listC)
}

func (s *sList) IntersectOp(ctx context.Context, A, B, C string) (n int, err error) {
	// 获取 list 数据
	listA := s.GetListData(ctx, A)
	listB := s.GetListData(ctx, B)
	// 数据处理 交集运算
	listC := make(map[string]any)
	for k, v := range listA {
		if _, ok := listB[k]; ok {
			listC[k] = v
		}
	}
	// 保存数据
	return s.AppendListData(ctx, C, listC)
}

func (s *sList) DifferenceOp(ctx context.Context, A, B, C string) (n int, err error) {
	// 获取 list 数据
	listA := s.GetListData(ctx, A)
	listB := s.GetListData(ctx, B)
	// 数据处理 差集运算
	listC := make(map[string]any)
	for k, v := range listA {
		if _, ok := listB[k]; !ok {
			listC[k] = v
		}
	}
	// 保存数据
	return s.AppendListData(ctx, C, listC)
}
