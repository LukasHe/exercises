# Python 3.3.3 and 2.7.6
# python threads_in_python.py

from threading import Thread


i = 0


def thread_1():
	global i
	for j in range(10000):
		i =  i + 1
		print("Inside thread_1 " + str(i))
		print("This is an delay")
def thread_2():
	global i
	for j in range(10000):
		i =  i - 1 
		print("Inside thread_2 " + str(i))
		print("This is an delay")

def main():
    t_1 = Thread(target = thread_1, args = (),)
    t_2 = Thread(target = thread_2, args = (),)
    t_1.start()
    t_2.start()
    t_1.join()
    t_2.join()
    print("Hello from main!")


main()