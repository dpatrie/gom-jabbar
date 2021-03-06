package robots

import (
	"testing"
	"time"

	"github.com/dpatrie/gom-jabbar/poutine/pkg/pubsub"
	"github.com/dpatrie/gom-jabbar/poutine/pkg/resto"
	"github.com/stretchr/testify/assert"
)

func TestOutremonaSuccess(t *testing.T) {
	assert := assert.New(t)
	done := make(chan bool)
	screamed := false
	curds := resto.SqueezedCheeseCurdsReady{}

	bus := &pubsub.Local{}
	bus.Subscribe("squeezed-cheese-ready", func(msg string) error {
		fromJSON(&curds, msg)
		done <- true
		return nil
	})
	bus.Subscribe("cheese-screams", func(msg string) error {
		screamed = true
		return nil
	})

	_ = NewOutremona(bus)
	o := &resto.PoutineOrder{
		ID:     "test",
		Size:   "small",
		Cheese: resto.CheeseKindCouicCouic,
	}
	bus.Publish("order-received", toJSON(o))

	select {
	case <-done:
		assert.Equal(curds.OrderID, o.ID)
		assert.Equal(curds.Kind, o.Cheese)
		assert.Equal(curds.Quantity, o.Size.Template().CurdsCount)
		assert.True(screamed)
	case <-time.After(20 * time.Second):
		t.Fatal("failed to pick and squeeze curds quickly enough")
	}

}
