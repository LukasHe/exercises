package Logic

import (
"time"
"fmt"
)

func LogicInit(newOrderChan, doneOrderChan, bidChan chan string) {
	go logic(newOrderChan, doneOrderChan, bidChan)
}

func logic(newOrderChan, doneOrderChan, bidChan chan string) {
	for{
		select {
			case newOrder := <-newOrderChan:
				calculateCost(newOrder)

			case costBid := <-bidChan:
				auction(costBid)

			case doneOrder := <-doneOrderChan:
				checkPendingOrders(doneOrder)

			default:
				time.Sleep(10*time.Millisecond)
		}
	}
}

func calculateCost(newOrder string) {
	fmt.Println("Cost:", newOrder)
}

func auction(costBid string) {
	fmt.Println("Auction:", costBid)
}

func checkPendingOrders(doneOrder string) {
	fmt.Println("Pending:", doneOrder)
}