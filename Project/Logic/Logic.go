package Logic

import (
"time"
"fmt"
"container/list"
"strconv"
"math"
)

func LogicInit(newOrderChan, doneOrderChan, bidChan, sendChan chan string) {
	go logic(newOrderChan, doneOrderChan, bidChan, sendChan)
}

func logic(newOrderChan, doneOrderChan, bidChan, sendChan chan string) {
	selfOrderList := list.New()
	//auctionList := list.New()


	//OBS CHANGE THIS
	selfOrderList.PushFront("2")
	selfOrderList.PushFront("2")
	selfOrderList.PushFront("2")


	for{
		select {
			case newOrder := <-newOrderChan:
				//fmt.Println(newOrder)
				selfOrderList.PushFront(newOrder)
				fmt.Println(selfOrderList.Len())
				cost := selfOrderList.Len()
				sendChan <-  "B" + "_" + newOrder[0:19] + "_" + strconv.Itoa(cost)
				//calculateCost(newOrder, selfOrderList)

			case costBid := <-bidChan:
				fmt.Println(costBid)
				if (int(time.Now().UnixNano()) - strconv.Atoi(string(costBid[0:19]))) > int(math.Pow10(10)){
					add _, to Atoi
				}
					

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