package Logic

import (
"time"
"fmt"
"container/list"
"strconv"
"strings"
"math"
"net"
)

func LogicInit(newOrderChan, doneOrderChan, bidChan, sendChan, selfOrderChan chan string, internalOrderChan chan int) {
	go logic(newOrderChan, doneOrderChan, bidChan, sendChan, selfOrderChan, internalOrderChan)
}

func logic(newOrderChan, doneOrderChan, bidChan, sendChan, selfOrderChan chan string, internalOrderChan chan int) {
	selfOrderList := list.New()
	existingBids := make(map[int]int)
	existingIPs := make(map[int]string)
	pendingOrders := make(map[int]string)


	allAddrs, _ := net.InterfaceAddrs()
	v4Addr := strings.Split(allAddrs[1].String(), "/")
	completeIP := strings.Split(v4Addr[0],".")
	selfIP := completeIP[3]

	for{
		select {
			case newOrder := <-newOrderChan:
				//fmt.Println(newOrder)
				lengthCounter := 0
				cost := 1000
				//fmt.Println(selfOrderList.Len())
				timeStamp, order, _ := splitMessage(newOrder)

				if _, ok := pendingOrders[timeStamp]; ok == false {
					pendingOrders[timeStamp] = order
					//Calculate Cost improve if possible
					for e := selfOrderList.Front(); e != nil; e = e.Next() {
						esplit := strings.Split(e.Value.(string), "_")
						floorDir := esplit[0]

						if floorDir == order{
							cost = lengthCounter
							break
						}

						lengthCounter = lengthCounter + 1
					}
					if cost < selfOrderList.Len(){
						sendChan <-  "B" + "_" + strconv.Itoa(timeStamp) + "_" + strconv.Itoa(cost)
					} else {
						cost = selfOrderList.Len()
						sendChan <-  "B" + "_" + strconv.Itoa(timeStamp) + "_" + strconv.Itoa(cost)
					}
				}



			case costBid := <-bidChan:
				timeStamp, bid, bidderIP := splitMessage(costBid)
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
				timeStamp, _, _ := splitMessage(doneOrder)
				element := selfOrderList.Front()
				if element.Value == pendingOrders[timeStamp] + "_" + strconv.Itoa(timeStamp){
					selfOrderList.Remove(element)
					
					if selfOrderList.Len() > 0 {
						element := selfOrderList.Front()
						e := element.Value
						fmt.Println("Send to hardware: ", e.(string))
						selfOrderChan <- e.(string)
					}
				}
				delete(pendingOrders, timeStamp)

			case  internalOrder := <-internalOrderChan:
				selfOrderList.PushBack(internalOrder)


			default:
				for timeStamp, _ := range existingBids{
					if (int(time.Now().UnixNano())) - timeStamp > int(math.Pow10(9)){
						if selfIP == existingIPs[timeStamp]{

							fmt.Println("We won! ", timeStamp)

							// Make it more intelligent
							if selfOrderList.Len() == 0{
								selfOrderChan <- pendingOrders[timeStamp] + "_" + strconv.Itoa(timeStamp)
							}
	
						for e := selfOrderList.Front(); e != nil; e = e.Next() {
							esplit := strings.Split(e.Value.(string), "_")
							floorDir := esplit[0]

							if floorDir == pendingOrders[timeStamp]{
								selfOrderList.InsertAfter(pendingOrders[timeStamp] + "_" + strconv.Itoa(timeStamp), e)
								break
								fmt.Println("add middle ", pendingOrders[timeStamp])
							// } else if e == selfOrderList.Back() {
							// 	selfOrderList.PushBack(pendingOrders[timeStamp] + "_" + strconv.Itoa(timeStamp))
							// 	fmt.Println("add to end ", pendingOrders[timeStamp])
							}
						}
							
					}


						delete(existingBids, timeStamp)
						delete(existingIPs, timeStamp)
					}
				}

				for timeStamp, _ := range pendingOrders{
					if (int(time.Now().UnixNano())) - timeStamp > 4*int(math.Pow10(10)){
						sendChan <- "N" + "_" + strconv.Itoa(int(time.Now().UnixNano())) + "_" + pendingOrders[timeStamp]
						delete(pendingOrders, timeStamp)
					}
				}

				
				for e := selfOrderList.Front(); e != nil; e = e.Next() {
					esplit := strings.Split(e.Value.(string), "_")
					timeStamp,_ := strconv.Atoi(esplit[1])
					if  (int(time.Now().UnixNano())) - timeStamp > 4*int(math.Pow10(10)){
						selfOrderList.Remove(e)	
					}
					
				}

				// Check if auction is completed by inspecting timeStamp
				// Remove done auctions and but winner in correct chanel



				time.Sleep(10*time.Millisecond)
		}
	}
}



func splitMessage(message string) (int, string, string) {
	splitMsg := strings.Split(message, "_")
	time, data, remoteIP := splitMsg[0], splitMsg[1], splitMsg[2] 
	splitIP := strings.Split(remoteIP, ":")
	splitIP = strings.Split(splitIP[0], ".")
	IP := splitIP[3]
	timeStamp, _ := strconv.Atoi(time)

	return timeStamp, data, IP
}

