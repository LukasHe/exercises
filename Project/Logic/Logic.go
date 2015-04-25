package Logic

import (
"time"
"fmt"
"container/list"
"strconv"
"math"
".././NetworkModule"
)

//Initialize Logic and spawn a thread.
func LogicInit(newOrderChan, doneOrderChan, bidChan, sendChan, selfOrderChan, ledOffChan chan string, internalOrderChan chan int) {

	go logic(newOrderChan, doneOrderChan, bidChan, sendChan, selfOrderChan, ledOffChan, internalOrderChan)
}


func logic(newOrderChan, doneOrderChan, bidChan, sendChan, selfOrderChan, ledOffChan chan string, internalOrderChan chan int) {
	selfOrderList := list.New()
	existingBids := make(map[int]int)
	existingIPs := make(map[int]string)
	pendingOrders := make(map[int]string)
	
	for{
		select {

			//Handles incoming newOrders and calculates the cost of handeling it.
			case newOrder := <-newOrderChan:
				timeStamp, origIP, order, _ := NetworkModule.SplitMessage(newOrder)

				//Check if the order is already in the pendingOrder map
				if _, ok := pendingOrders[timeStamp]; ok == false {
				
					pendingOrders[timeStamp] = strconv.Itoa(timeStamp) +  "_" + origIP +  "_" + order +  "_" + "*"
					cost, validCost := calculateCost(newOrder, pendingOrders, selfOrderList)
					
					//Put bid in sendChan to be broadcasted.
					if validCost {
						sendChan <-  "B" + "_" + strconv.Itoa(timeStamp) + "_" + origIP + "_" + strconv.Itoa(cost)+ "_" + NetworkModule.GetOwnIP()
					}					
				}

			//Handles incoming bids and stores the lowest. After a pre-set time
			//the auction is over and no more bids are accepted. 
			case costBid := <-bidChan:
				auction(existingBids, existingIPs, costBid)

			//Handles incoming doneOrders and adjusts the selfOrderList and
			//sends new orders to the hardware.
			case doneOrder := <-doneOrderChan:
				timeStamp, origIP, order, doneIP := NetworkModule.SplitMessage(doneOrder)

				//Check if the doneOrder was from selfOrderList
				if selfOrderList.Len() > 0 && doneIP == NetworkModule.GetOwnIP() {
					ptrFrontElement := selfOrderList.Front()
					openDoor(selfOrderChan, ptrFrontElement, selfOrderList)
					
					//Put new order in selfOrderChan. 
					if selfOrderList.Len() > 1 {
						selfOrderChan <- ptrFrontElement.Next().Value.(string)
					}
					selfOrderList.Remove(ptrFrontElement)
				}

				turnOffLights(origIP, order, ledOffChan)
				delete(pendingOrders, timeStamp)

			//Handles finished auction and sort the order into the correct place in selfOrderList.
			//Rebroadcasts expired Orders and cleans up selfOrderList/pendingOrders
			default:
				sortFinishedAuction(existingBids, existingIPs, pendingOrders, selfOrderList, selfOrderChan)
				rebroadcastPendingOrders(pendingOrders, sendChan)
				cleanupSelfOrderList(selfOrderList)
			
				time.Sleep(10*time.Millisecond)	
		}
	}
}


//Turn of lights correspoding to the order if it was in this elevator.
func turnOffLights(origIP, order string, ledOffChan chan string) {

	if origIP == NetworkModule.GetOwnIP(){

		if string(order[1]) == "U"{
			ledOffChan <- "LIGHT_UP" + string(order[0])

		} else if string(order[1]) == "D"{
			ledOffChan <- "LIGHT_DOWN" + string(order[0])

		} else {
			ledOffChan <- "LIGHT_COMMAND" + string(order[0])
		}
	}
				
}

//Check if the door should open or not. Identical orders will only result in 1 opening.
func openDoor(selfOrderChan chan string, ptrFrontElement *list.Element, ptrselfOrderList *list.List) {

	if  ptrselfOrderList.Len() == 1 {
		selfOrderChan <- "*_*_WAIT_*"

	} else {
		_, _, oldOrder, _ := NetworkModule.SplitMessage(ptrFrontElement.Value.(string))
		_, _, nextOrder, _ := NetworkModule.SplitMessage(ptrFrontElement.Next().Value.(string))

		if oldOrder[0] != nextOrder[0] {
			selfOrderChan <- "*_*_WAIT_*"
		}
	}
}


//Calculates cost for accepting an order by looking at what kind of order it is and the selfOrderList
//and returns the cost and weither or not we should bid.
func calculateCost(newOrder string, pendingOrders map[int]string, selfOrderList *list.List) (int, bool){

	lengthCounter := 0
	_, origIP, order, _ := NetworkModule.SplitMessage(newOrder)
	
		//Check if its an internal order for another elevator
		if string(order[1]) == "I" && origIP == NetworkModule.GetOwnIP(){
			return  0 , true

		} else if string(order[1]) != "I"  {
			
			//Calculate cost by inspecting the selfOrderList
			for elementOfList := selfOrderList.Front(); elementOfList != nil; elementOfList = elementOfList.Next() {
				_, _, floorDir, _ := NetworkModule.SplitMessage(elementOfList.Value.(string))

				if floorDir == order{
					return lengthCounter, true
				}

				lengthCounter = lengthCounter + 1
			}
			return lengthCounter, true
		}
	return -1, false
}

//Check if Auction has timed out and if there is an existing bid on that Auction and compare
//with incoming bid to determin winner.
func auction(existingBids  map[int]int , existingIPs  map[int]string, costBid string) {

	timeStamp, _, bid, bidderIP := NetworkModule.SplitMessage(costBid)

	if  existingBid, ok := existingBids[timeStamp]; ok && (int(time.Now().UnixNano())) -  timeStamp < int(math.Pow10(9)){

		if bid, _ := strconv.Atoi(bid); bid < existingBid{
			existingBids[timeStamp] = bid
			existingIPs[timeStamp] = bidderIP

		} else if bid == existingBid && bidderIP < existingIPs[timeStamp]{
			existingBids[timeStamp] = bid
			existingIPs[timeStamp] = bidderIP
		}	

	} else {
		existingBids[timeStamp],_ = strconv.Atoi(bid)
		existingIPs[timeStamp] = bidderIP
	}
}

//If a order in pendingOrdersList has timed out remove it and rebroadcast it to the whole system
func rebroadcastPendingOrders(pendingOrders map[int]string, sendChan chan string){
	for timeStamp, _ := range pendingOrders{
		if (int(time.Now().UnixNano())) - timeStamp > 2*int(math.Pow10(10)){
			_, origIP, order, _ := NetworkModule.SplitMessage(pendingOrders[timeStamp])
			sendChan <- "N" + "_" + strconv.Itoa(int(time.Now().UnixNano())) + "_" + origIP + "_" + order + "_" + "*"
			delete(pendingOrders, timeStamp)
		}
	}
}

//If an order in selfOrderList has timed out remove it.
func cleanupSelfOrderList(selfOrderList *list.List){
	for elementOfList := selfOrderList.Front(); elementOfList != nil; elementOfList = elementOfList.Next() {
		timeStamp, _, _, _ := NetworkModule.SplitMessage(elementOfList.Value.(string))
		if  (int(time.Now().UnixNano())) - timeStamp > 2*int(math.Pow10(10)){
			selfOrderList.Remove(elementOfList)	
		}	
	}
}

//Sort timed out bids we won into the correct place in the 
//selfOrderList. 
func sortFinishedAuction(existingBids map[int]int, existingIPs, pendingOrders map[int]string, selfOrderList *list.List, selfOrderChan chan string){

	for timeStamp, _ := range existingBids{

		//Checks if the auction has timed out and if we won.
		if (int(time.Now().UnixNano())) - timeStamp > int(math.Pow10(8)){
			if NetworkModule.GetOwnIP()  == existingIPs[timeStamp]{

				fmt.Println("Auction won!: Timestamp: ", timeStamp)
				//If there is no orders to execute put the order in front.							
				if selfOrderList.Len() == 0{
					selfOrderList.PushFront(pendingOrders[timeStamp])
					selfOrderChan <- pendingOrders[timeStamp]


				//If there exists more orders in selfOrderList check for identical orders
				//and sort the new order after it. 
				} else {
					_, _, floorDirInPending, _ := NetworkModule.SplitMessage(pendingOrders[timeStamp])

					for elementOfList := selfOrderList.Front(); elementOfList != nil; elementOfList = elementOfList.Next() {
						_, _, floorDirInList, _ := NetworkModule.SplitMessage(elementOfList.Value.(string))

						if floorDirInList == floorDirInPending || (string(floorDirInPending[1]) == "I" && floorDirInList[0] == floorDirInPending[0]) {
							selfOrderList.InsertAfter(pendingOrders[timeStamp], elementOfList)
							break

						} else if elementOfList == selfOrderList.Back() {
							selfOrderList.PushBack(pendingOrders[timeStamp])
							break
						}
					}
				}
			}

			delete(existingBids, timeStamp)
			delete(existingIPs, timeStamp)
		}
	}
}