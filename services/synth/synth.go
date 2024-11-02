package synth

import (
	"clutch/common"
	"fmt"
)

func Synth(synthChan *chan common.Event) {
	fmt.Println("Synthing")
	for event := range *synthChan {
		fmt.Println("Synth event:", event)
	}
}
