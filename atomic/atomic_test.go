package atomic

import (
	"fmt"
	"testing"
)

func TestGetFeePerKB(t *testing.T) {
	fmt.Println(GetFeePerKB())
}
