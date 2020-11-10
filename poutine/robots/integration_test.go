package robots

import (
	"log"
	"testing"
	"time"

	"github.com/OROHealth/gom-jabbar/poutine/pubsub"
	"github.com/OROHealth/gom-jabbar/poutine/resto"
	"github.com/stretchr/testify/assert"
)

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	assert := assert.New(t)
	done := make(chan bool)
	bus := &pubsub.Local{}

	bus.Subscribe("order-ready", func(msg string) error {
		done <- true
		return nil
	})

	_ = NewOutremona(bus)
	_ = NewMontroyashi(bus)
	_ = NewVerduny(bus)
	_ = NewNordo(bus)
	_ = NewBizar(bus)
	_ = NewOldoporto(bus)

	robot := NewPierre(bus, fullInventory())
	id, err := robot.TakeOrder(&resto.PoutineOrder{
		Size:      "small",
		Potato:    resto.SweetPotato,
		PotatoCut: resto.SmallCut,
		Cheese:    resto.CheeseKindCouicCouic,
		Oil:       resto.OilKindSunflower,
		Gravy:     resto.GravyKindTequila,
	})

	assert.NoError(err)

	select {
	case <-done:
		o, err := robot.Order(id)
		assert.NoError(err)
		assert.Equal(o.Status, resto.OrderStatusDelivered)
		assert.NotNil(o.PoutineDelivered)
		log.Println(o.PoutineDelivered)
		//TODO: validate the poutine
	case <-time.After(resto.OrderTimeout):
		t.Fatal("order completion timed out")
	}
}