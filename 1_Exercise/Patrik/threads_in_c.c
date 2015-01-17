#include <pthread.h>
#include <stdio.h>

// gcc -std=gnu99 -Wall -g -o out threads_in_c.c -lpthread

//Global variable
int i = 100;

void* thread_1()
{
	for (int j = 0; j<1000000; j++)
	{
		i++;
		
	}
	printf("Inside the thread_1! This is i:%d\n ",i);
	return NULL;
}

void* thread_2()
{

	for (int j = 0; j<1000000; j++)
	{
		i--;
		
	}
	printf("Inside the thread_2! This is i:%d\n ",i);
	return NULL;
}


int main()
{

	//creates an instance of type pthread_t then creates the thread and finally calls it.
	pthread_t t_1;
	pthread_t t_2;
	pthread_create(&t_1, NULL, thread_1,NULL); 
	pthread_create(&t_2, NULL, thread_2,NULL); 
	pthread_join(t_1, NULL);
	pthread_join(t_2, NULL);

	return 0;
}
