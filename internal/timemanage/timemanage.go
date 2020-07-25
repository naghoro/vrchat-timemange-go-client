package timemanage

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"vrchat-timemanage/internal/oscclient"
)

// TimeManage ... main struct
type TimeManage struct {
	SecondOfHour time.Duration
	osc          oscclient.OscIface
	doneCh       chan string
}

// New ... return TimeManage
func New(opts ...Option) *TimeManage {
	// default value
	t := &TimeManage{
		SecondOfHour: 3600 * time.Second,
	}

	for _, opt := range opts {
		opt(t)
	}

	return t
}

// Option ... argument of New
type Option func(*TimeManage)

// SetSecondOfHour ... set second of hour
func SetSecondOfHour(second time.Duration) Option {
	return func(t *TimeManage) {
		t.SecondOfHour = second
	}
}

// periodic loop
func (tmanage *TimeManage) periodicLoop(ctx context.Context) {
	defer close(tmanage.doneCh)

	// process osc message 24 times per days
	// day's second is 86400 default. We can change Timemanage struct
	fmt.Println("time", tmanage.SecondOfHour)
	ticker := time.NewTicker(tmanage.SecondOfHour)
	defer ticker.Stop()

	now := time.Now()
	hour := now.Hour()

	// exec firstly
	fmt.Printf("send hour:%d\n", hour)
	tmanage.osc.Sendhour(hour, 1)

	for {
		select {
		case <-ctx.Done():
			return
		case t := <-ticker.C:
			_ = t
			// send osc
			oneHourAgo := hour

			if err := tmanage.osc.Sendhour(oneHourAgo, 0); err != nil {
				fmt.Printf("Send OSC error: %d, %+v\n", oneHourAgo, err)
			}

			hour = (hour + 1) % 24
			fmt.Printf("send hour:%d\n", hour)

			if err := tmanage.osc.Sendhour(hour, 1); err != nil {
				fmt.Printf("Send OSC error: %d, %+v\n", hour, err)
			}

		}
	}
}

// ManageStart ... start
func (tmanage *TimeManage) ManageStart() int {

	termCh := make(chan os.Signal, 1)
	signal.Notify(termCh, syscall.SIGTERM, syscall.SIGINT)

	// context
	ctx := context.Background()
	ctxC, cancel := context.WithCancel(ctx)

	// gracefull stop
	tmanage.doneCh = make(chan string)

	// osc
	tmanage.osc = oscclient.New()

	go tmanage.periodicLoop(ctxC)

	select {
	case <-termCh:
		cancel()
		// finish collect
		<-tmanage.doneCh
	}

	return 0
}
