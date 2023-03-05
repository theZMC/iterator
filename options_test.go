package iterator

import (
	"sync"
	"testing"
)

// need to enable the race detector for this test to really be valuable
func Test_ThreadSafe(t *testing.T) {
	iter := From(
		[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		ThreadSafe(true),
	)
	wg := sync.WaitGroup{}
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			for {
				if _, ok := iter.Next(); !ok {
					break
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func Test_CopySource(t *testing.T) {
	source := []int{1, 2, 3, 4, 5}
	iter := From(source, CopySource(true))
	source[0] = 100
	if num, ok := iter.Next(); !ok || num != 1 {
		t.Errorf("Expected 1, got %d", num)
	}
}

func Test_BufferLen(t *testing.T) {
	fromOpts := new(fromOptions)
	fromOpts.bufferLen = 64
	BufferLen(32)(fromOpts)
	if fromOpts.bufferLen != 32 {
		t.Errorf("Expected 32, got %d", fromOpts.bufferLen)
	}
}

func Test_CloseChannel(t *testing.T) {
	ico := new(intoChannelOptions)
	ico.closeChannel = false
	CloseChannel(true)(ico)
	if !ico.closeChannel {
		t.Error("Expected true, got false")
	}
}
