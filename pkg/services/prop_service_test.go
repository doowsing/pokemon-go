package services

import (
	"fmt"
	"testing"
)

func TestPropService_UsePropById(t *testing.T) {
	opt := NewOptService()
	result := opt.PropSrv.UsePropById(1, 10001)
	fmt.Printf("use prop result:%s\n", result)
}
