package Logic

import (
"time"
// "fmt"
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


	// OBS CHANGE THIS
	selfOrderList.PushFront("2")
	selfOrderList.PushFront("2")
	selfOrderList.PushFront("2")

	allAddrs, _ := net.InterfaceAddrs()
	v4Addr := strings.Split(allAddrs[1].String(), "/")
	completeIP := strings.Split(v4Addr[0],".")
	selfIP := completeIP[3]


	for{
		select {
			case newOrder := <-newOrderChan:
				//fmt.Println(newOrder)
				//fmt.Println(selfOrderList.Len())
				timeStamp, order, _ := splitMessage(newOrder)

				
				if _, ok := pendingOrders[timeStamp]; ok == false {
					pendingOrders[timeStamp] = order
					//Calculate Cost improve if possible
					cost := selfOrderList.Len()
					sendChan <-  "B" + "_" + strconv.Itoa(timeStamp) + "_" + strconv.Itoa(cost)
				}



			case costBid := <-bidChan:
				timeStamp, bid, bidderIP := splitMessage(costBid)
				//fmt.Println("Time:", timeStamp, "Bid:", bid, "IP:", bidderIP)

				if (int(time.Now().UnixNano())) -  timeStamp > int(math.Pow10(10)){
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
				delete(pendingOrders, timeStamp)
				if selfOrderList.Front() == 



			default:
				for timeStamp, _ := range existingBids{
					if (int(time.Now().UnixNano())) - timeStamp > int(math.Pow10(10)){
						if selfIP == existingIPs[timeStamp]{

							// Make it more intelligent
							if selfOrderList.Len() == 0{
								floor := string(pendingOrders[timeStamp][0])
								selfOrderChan <- floor
							}
							selfOrderList.PushBack(pendingOrders[timeStamp] + "_" + strconv.Itoa(timeStamp))
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



func splitMessage(message string) (int, string, string) {
	splitMsg := strings.Split(message, "_")
	time, data, remoteIP := splitMsg[0], splitMsg[1], splitMsg[2] 
	splitIP := strings.Split(remoteIP, ":")
	splitIP = strings.Split(splitIP[0], ".")
	IP := splitIP[3]
	timeStamp, _ := strconv.Atoi(time)

	return timeStamp, data, IP
}

