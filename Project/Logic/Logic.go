package Logic

import (
"time"
"fmt"
"container/list"
"strconv"
"math"
".././NetworkModule"
)

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
			case newOrder := <-newOrderChan:
				lengthCounter := 0
				cost := 1000
				//fmt.Println(selfOrderList.Len())
				timeStamp, origIP, order, internalOrder := NetworkModule.SplitMessage(newOrder)
				if _, ok := pendingOrders[timeStamp]; ok == false {
					pendingOrders[timeStamp] = strconv.Itoa(timeStamp) +  "_" + origIP +  "_" + order +  "_" + "*"
				if internalOrder == "I" && origIP == NetworkModule.GetOwnIP(){
					fmt.Println("Iternal")
					sendChan <-  "B" + "_" + strconv.Itoa(timeStamp) + "_" + origIP + "_" + "0" + "_" + NetworkModule.GetOwnIP()
					break
				}	
					for e := selfOrderList.Front(); e != nil; e = e.Next() {
						_, _, floorDir, _ := NetworkModule.SplitMessage(e.Value.(string))

						if floorDir == order{
							cost = lengthCounter
							break
						}

						lengthCounter = lengthCounter + 1
					}
					if cost < selfOrderList.Len(){
						sendChan <-  "B" + "_" + strconv.Itoa(timeStamp) + "_" + origIP + "_" + strconv.Itoa(cost) + "_" + NetworkModule.GetOwnIP()
					} else {
						cost = selfOrderList.Len()
						sendChan <-  "B" + "_" + strconv.Itoa(timeStamp) + "_" + origIP + "_" + strconv.Itoa(cost) + "_" + NetworkModule.GetOwnIP()
					}
				}



			case costBid := <-bidChan:
				timeStamp, _, bid, bidderIP := NetworkModule.SplitMessage(costBid)
				//fmt.Println("Time:", timeStamp, "Bid:", bid, "IP:", bidderIP)

				if (int(time.Now().UnixNano())) -  timeStamp > int(math.Pow10(9)){
					break
				} else if  existingBid, ok := existingBids[timeStamp]; ok{
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
						
					

			case doneOrder := <-doneOrderChan:
				timeStamp, origIP, order, _ := NetworkModule.SplitMessage(doneOrder)
				ptrfrontElement := selfOrderList.Front()

				if selfOrderList.Len() > 0 && ptrfrontElement.Value == pendingOrders[timeStamp]{
					//fmt.Println("Done Order: ", doneOrder)
					selfOrderList.Remove(ptrfrontElement)

					if selfOrderList.Len() > 0{
						ptrfrontElement = selfOrderList.Front()
						frontElement := ptrfrontElement.Value
						//fmt.Println("Send to hardware: ", frontElement.(string))
						selfOrderChan <- frontElement.(string)
					}
				}
				if origIP == NetworkModule.GetOwnIP(){
					if string(order[1]) == "U"{
						ledOffChan <- "LIGHT_UP" + string(order[0])
					} else if string(order[1]) == "D"{
						ledOffChan <- "LIGHT_DOWN" + string(order[0])
					} else {
						ledOffChan <- "LIGHT_COMMAND" + string(order[0])
					}
				}
				delete(pendingOrders, timeStamp)


			default:
				for timeStamp, _ := range existingBids{
					if (int(time.Now().UnixNano())) - timeStamp > int(math.Pow10(8)){
						if NetworkModule.GetOwnIP()  == existingIPs[timeStamp]{

							//fmt.Println("We won! ", timeStamp)

							if selfOrderList.Len() == 0{
								selfOrderList.PushBack(pendingOrders[timeStamp])
								selfOrderChan <- pendingOrders[timeStamp]
							} else {
								_, _, floorDirInPending, _ := NetworkModule.SplitMessage(pendingOrders[timeStamp])
								for elementOfList := selfOrderList.Front(); elementOfList != nil; elementOfList = elementOfList.Next() {
									// fmt.Println("self: ", elementOfList.Value.(string))
									_, _, floorDirInList, _ := NetworkModule.SplitMessage(elementOfList.Value.(string))
									// fmt.Println("self: ", floorDirInList)
									if floorDirInList == floorDirInPending || (string(floorDirInPending[1]) == "I" && floorDirInList[0] == floorDirInPending[0]) {
										selfOrderList.InsertAfter(pendingOrders[timeStamp], elementOfList)
										//fmt.Println("add middle ", pendingOrders[timeStamp])
										break
									} else if elementOfList == selfOrderList.Back() {
										selfOrderList.PushBack(pendingOrders[timeStamp])
										//fmt.Println("add to end ", pendingOrders[timeStamp])
										break
									}
								}
							}
						}


						delete(existingBids, timeStamp)
						delete(existingIPs, timeStamp)
					}
				

				for timeStamp, _ := range pendingOrders{
					if (int(time.Now().UnixNano())) - timeStamp > int(math.Pow10(10)){
						_, origIP, order, _ := NetworkModule.SplitMessage(pendingOrders[timeStamp])
						sendChan <- "N" + "_" + strconv.Itoa(int(time.Now().UnixNano())) + "_" + origIP + "_" + order + "_" + "*"
						delete(pendingOrders, timeStamp)
					}
				}

				
				for elementOfList := selfOrderList.Front(); elementOfList != nil; elementOfList = elementOfList.Next() {
					timeStamp, _, _, _ := NetworkModule.SplitMessage(pendingOrders[timeStamp])
					if  (int(time.Now().UnixNano())) - timeStamp > int(math.Pow10(10)){
						selfOrderList.Remove(elementOfList)	
					}	
				}
				time.Sleep(10*time.Millisecond)
			}	
		}
	}
}





