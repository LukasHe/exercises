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

func LogicInit(newOrderChan, doneOrderChan, bidChan, sendChan, selfOrderChan chan string) {
	go logic(newOrderChan, doneOrderChan, bidChan, sendChan, selfOrderChan)
}

func logic(newOrderChan, doneOrderChan, bidChan, sendChan, selfOrderChan chan string) {
	selfOrderList := list.New()
	existingBids := make(map[int]int)
	existingIPs := make(map[int]string)
	pendingOrders := make(map[int]string)


	//OBS CHANGE THIS
	selfOrderList.PushFront("2")
	selfOrderList.PushFront("2")
	selfOrderList.PushFront("2")

	allAddrs, _ := net.InterfaceAddrs()
	v4Addr := strings.Split(allAddrs[1].String(), "/")
	selfIP := strings.Split(v4Addr[0],".")
	selfIP = selfIP[3]


	for{
		select {
			case newOrder := <-newOrderChan:
				//fmt.Println(newOrder)
				//fmt.Println(selfOrderList.Len())
				splitOrder := strings.Split(newOrder, "_")
				timeStamp, order := splitOrder[0], splitOrder[1]
				
				if _, ok := pendingOrders[timeStamp]; ok == false {
					pendingOrders[timeStamp] = order
					//Calculate Cost improve if possible
					cost := selfOrderList.Len()
					sendChan <-  "B" + "_" + splitOrder[0] + "_" + strconv.Itoa(cost)
				}



			case costBid := <-bidChan:
				splitBid := strings.Split(costBid, "_")
				timeStamp, bid, bidderIP := splitBid[0], splitBid[1], splitBid[2] 
				splitIP := strings.Split(bidderIP, ":")
				splitIP = strings.Split(splitIP[0], ".")
				bidderIP = splitIP[3]
				//fmt.Println("Time:", timeStamp, "Bid:", bid, "IP:", bidderIP)

				if timeStamp, _ := strconv.Atoi(timeStamp); (int(time.Now().UnixNano())) -  timeStamp > int(math.Pow10(10)){
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
				checkPendingOrders(doneOrder)

			default:
				for timeStamp, _ := range existingBids{
					if (int(time.Now().UnixNano())) - timeStamp > int(math.Pow10(10)){
						if selfIP = existingIPs[timeStamp]{

							// Make it more intelligent
							selfOrderList.PushBack(pendingOrders[timeStamp])
							if selfOrderList.Len() == 0{
								floor := pendingOrders[timeStamp][0]
								selfOrderChan <- floor 
							}
						}


						delete(existingBids, timeStamp)
						delete(existingIPs, timeStamp)
					}
				}

				// Check if auction is completed by inspecting timestamp
				// Remove done auctions and but winner in correct chanel



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