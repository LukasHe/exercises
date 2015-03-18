package Logic

import (
"time"
"fmt"
"container/list"
"strconv"
)

func LogicInit(newOrderChan, doneOrderChan, bidChan, sendChan chan string) {
	go logic(newOrderChan, doneOrderChan, bidChan, sendChan)
}

func logic(newOrderChan, doneOrderChan, bidChan, sendChan chan string) {
	selfOrderList := list.New()

	for{
		select {
			case newOrder := <-newOrderChan:
				//fmt.Println(newOrder)
				selfOrderList.PushFront(newOrder)
				fmt.Println(selfOrderList.Len())
				cost := selfOrderList.Len()
				sendChan <-  "B" + newOrder[0:19] + strconv.Itoa(cost)
				//calculateCost(newOrder, selfOrderList)

			case costBid := <-bidChan:
				auction(costBid)

			case doneOrder := <-doneOrderChan:
				checkPendingOrders(doneOrder)

			default:
				time.Sleep(10*time.Millisecond)
		}
	}
}


// func calculateCost(newOrder string, selfOrderList list.List) {
// 	fmt.Println("Cost:", newOrder)
// }

func auction(costBid string) {
	fmt.Println("Auction:", costBid)
}

func checkPendingOrders(doneOrder string) {
	fmt.Println("Pending:", doneOrder)
}