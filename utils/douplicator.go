package utils

import (
	"fmt"
	"time"

)

type ChannelDuplicator struct {
	transfer chan interface{}
	outputs  []chan interface{}
	debug bool
	debugName string
}

func MakeDuplicator() *ChannelDuplicator {
	chDoup := &ChannelDuplicator{
		outputs:  make([]chan interface{}, 0),
		transfer: make(chan interface{}, 100),
		debug: false,
	}

	chDoup.startDuplicator()

	return chDoup
}

func (ch *ChannelDuplicator)EnableDebug(name string){
	ch.debugName = name
	ch.debug = true
}

func (ch *ChannelDuplicator) GetOutput() chan interface{} {
	// make a channel with a 10 buffer size
	if ch.debug{
		fmt.Println("adding output on", ch.debugName)
	}
	newOutput := make(chan interface{}, 10)
	ch.outputs = append(ch.outputs, newOutput)
	return newOutput
}

func (ch *ChannelDuplicator) UnregisterOutput(remove chan interface{}) {
	var removeIndex int
	for i, channel := range ch.outputs {
		if channel == remove {
			removeIndex = i
		}
	}
	//Remove channel by swapping the removed channel to the end and then just trimming the slice
	ch.outputs[len(ch.outputs)-1], ch.outputs[removeIndex] = ch.outputs[removeIndex], ch.outputs[len(ch.outputs)-1]
	ch.outputs = ch.outputs[:len(ch.outputs)-1]
}

func (ch *ChannelDuplicator) RegisterInput(inputChannel <- chan interface{}) {
	go func() {
		if ch.debug{
			fmt.Println("registering input on", ch.debugName)
		}

		for val := range inputChannel {
			if ch.debug{
				fmt.Println("adding to trasfer", ch.debugName,"value", val)
			}
			ch.Offer(val)
			if ch.debug{
				fmt.Println("done transfer on", ch.debugName)
			}

		}
		if ch.debug{
			fmt.Println("closeing input on", ch.debugName)
		}
	}()

}

func (ch *ChannelDuplicator) Offer(value interface{}) {
	if ch.debug{
		fmt.Println("offering to transfer", ch.debugName)
	}
	ch.transfer <- value
}

func (ch *ChannelDuplicator) startDuplicator() {
	go func() {
		for nextValue := range ch.transfer {
			if ch.debug{
				fmt.Println("sending down outputs on", ch.debugName)
			}
			for i, channel := range ch.outputs {
				select {
				case channel <- nextValue:
					if ch.debug{
						fmt.Println("sent to an output of", ch.debugName, "index", i)
					}
					continue
				default:
					if ch.debug{
						fmt.Println("missing messages on", ch.debugName, "index", i)
					}
					continue
				}
			}
		}
	}()

}

func main() {
	input1 := make(chan interface{})
	input2 := make(chan interface{})

	chDoup := MakeDuplicator()
	chDoup.RegisterInput(input1)
	chDoup.RegisterInput(input2)
	for i := 0; i < 3; i++ {
		output := chDoup.GetOutput()
		go func() {
			for value := range output {
				fmt.Println("recieved: ", value)
			}
		}()
	}
	input1 <- "hello"
	input2 <- "world"
	time.Sleep(time.Second * 1)

}
