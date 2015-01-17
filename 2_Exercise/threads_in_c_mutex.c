#include <pthread.h>
#include <stdio.h>

// gcc -std=gnu99 -Wall -g -o out threads_in_c.c -lpthread

//Global variable
int i = 100;
pthread_mutex_t lock;



void* thread_1()
{
	pthread_mutex_lock(&lock);
	for (int j = 0; j<1000000; j++)
	{
		i++;
		
	}
	pthread_mutex_unlock(&lock);
	return NULL;
}

void* thread_2()
{

	pthread_mutex_lock(&lock);
	for (int j = 0; j<1000000; j++)
	{
		i--;
		
	}
	pthread_mutex_unlock(&lock);
	return NULL;
}


int main()
{
	pthread_t t_1;
	pthread_t t_2;
	pthread_mutex_init(&lock, NULL);

	//creates an instance of type pthread_t then creates the thread and finally calls it.
	pthread_create(&t_1, NULL, thread_1,NULL); 
	pthread_create(&t_2, NULL, thread_2,NULL); 
	pthread_join(t_1, NULL);
	pthread_join(t_2, NULL);
	pthread_mutex_destroy(&lock);
	printf("Inside the main! This is i:%d\n ",i);

	return 0;
}
