// gcc 4.7.2 +
// gcc -std=gnu99 -Wall -g -o helloworld_c helloworld_c.c -lpthread

#include <pthread.h>
#include <stdio.h>

int i = 0;

// Note the return type: void*
void* thread_1(){
    int j;
    for(j = 0; j < 1000000; j++){
        i++;
    }
    return NULL;
}

void* thread_2(){
    int j;
    for(j = 0; j < 1000000; j++){
        i--;
    }
    return NULL;
}



int main(){
    pthread_t Thread_1;
    pthread_t Thread_2;
    pthread_create(&Thread_1, NULL, thread_1, NULL);
    pthread_create(&Thread_2, NULL, thread_2, NULL);
    // Arguments to a thread would be passed here ---------^
    
    pthread_join(Thread_1, NULL);
    pthread_join(Thread_2, NULL);
    printf("%d\n", i);
    return 0;
    
}