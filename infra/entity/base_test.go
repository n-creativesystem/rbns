package entity

import (
	"sync"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestULIDGenerate(t *testing.T) {
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ul := generate()
			logrus.Println(ul)
		}()
	}
	wg.Wait()
}
