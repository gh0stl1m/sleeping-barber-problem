package main

import (
	"time"

	"github.com/fatih/color"
)

type BarberShop struct {
  ShopCapacity int
  HairCutDuration time.Duration
  NumberOfBarbers int
  BarbersDoneCh chan bool
  ClientsCh chan string
  Open bool
}

func (bs *BarberShop) addBarber(name string) {

  bs.NumberOfBarbers++

  go func() {

    isSleeping := false
    color.Yellow("%s goes to the waiting room to check for clients", name)

    for {

      if len(bs.ClientsCh) == 0 {

        color.Yellow("There's nothing to do, so %s takes a nap", name)
        isSleeping = true
      }

      client, workToDo := <- bs.ClientsCh

      if workToDo {

        if isSleeping {
  
          color.Yellow("%s wakes %s up", client, name)
          isSleeping = false
        }

        bs.cutHair(name, client)
      } else {

        bs.sendBarberHome(name)
        return
      }
    }
    
  }()
}

func (bs *BarberShop) cutHair(barber, client string) {

  color.Green("%s is cutting the %s's hair", barber, client)

  time.Sleep(bs.HairCutDuration)

  color.Green("%s is finished cutting %s's hair", barber, client)
}

func (bs *BarberShop) sendBarberHome(barber string) {

  color.Cyan("%s is going home", barber)

  bs.BarbersDoneCh <- true
}

func (bs *BarberShop) addClient(client string) {

	color.Green("*** %s arrives!", client)

	if bs.Open {
		select {
		case bs.ClientsCh <- client:
			color.Yellow("%s takes a seat in the waiting room.", client)
		default:
			color.Red("The waiting room is full, so %s leaves.", client)
		}
	} else {

		color.Red("The shop is already closed, so %s leaves!", client)
	}
}

func (bs *BarberShop) close() {

	color.Cyan("Closing shop for the day.")

	close(bs.ClientsCh)
	bs.Open = false

	for a := 1; a <= bs.NumberOfBarbers; a++ {
		<-bs.BarbersDoneCh
	}

	close(bs.BarbersDoneCh)

  color.Green("----------------------------")
  color.Green("The barbershop is now closed")
}
