package background

import (
	"fmt"
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
)

func TestBackgroundTaks(t *testing.T) {
	ticker := NewBackgroundWorker()
	ticker.SetInterval(1 * time.Second)
	ticker.AddTask(NewTask("1223", func() {
		fmt.Println("Task 1")
	}))

	ticker.RunUntil(4 * time.Second)

	time.Sleep(5 * time.Second)

	assert.Equal(t, ticker.PastInterval(), 3, fmt.Sprintf("Cycles does not match %v", ticker.PastInterval()))

}
