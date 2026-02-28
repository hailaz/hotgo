// Package global
package global

import (
	"context"
	"hotgo/internal/library/addons"
)

func Init(ctx context.Context, sk *addons.Skeleton) {
	skeleton = sk
}

func GetSkeleton() *addons.Skeleton {
	if skeleton == nil {
		panic("addon skeleton not initialized.")
	}
	return skeleton
}
