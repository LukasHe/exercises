# Python 3.3.3 and 2.7.6
# python threads_in_python.py

import threading 
from threading import Thread

i = 0
lock = threading.Lock()

def thread_1():
	lock.acquire()
	global i
	for j in range(10000):
		i =  i + 1
	lock.release()

def thread_2():
	lock.acquire()
	global i
	for j in range(10000):
		i =  i - 1 
	lock.release()


def main():
	
	t_1 = Thread(target = thread_1, args = (),)
	t_2 = Thread(target = thread_2, args = (),)
	t_1.start()
	t_2.start()
	t_1.join()
	t_2.join()
	print("Hello from main! i: " + str(i))


main()