package utils

import (
	"fmt"
	"testing"
)

func TestGlobalObj_Reload(t *testing.T) {
	GlobalObject.Reload()
	fmt.Println(GlobalObject)
}
