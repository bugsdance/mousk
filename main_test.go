package main

import (
	"fmt"
	"mousk/infra/config"
	"mousk/infra/mousectl"
	"testing"
	"time"
)

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}

func Test_keyboardProcess(t *testing.T) {
	keyboardProcess()
}

func TestScrollMouseFunc(t *testing.T) {
	time.Sleep(2 * time.Second)
	config.Init()
	fmt.Println("hah1")
	mousectl.ScrollMouseCtrl(mousectl.DirectionVerticalDown, 2)
	time.Sleep(time.Second)
	fmt.Println("hah2")
	mousectl.ScrollMouseCtrl(mousectl.DirectionVerticalUp, 2)

}

// -------------------------------------------------------------------------------------------------------------------------------------------------------
// 1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111
// 222222222222222222222222222
// 222222222222222222222222222
// 322222222222222222222222222
// 422222222222222222222222222
// 522222222222222222222222222
// 1
// 213
//12312
//12312
//12312
//12312
//12312
//12312
//12312
//12312
// 1jidoa
