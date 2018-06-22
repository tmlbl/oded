#include <stdio.h>
#include <time.h>
#include <stdlib.h>
#include <math.h>
#include <unistd.h>
#include <string.h>

double clock2ms(double t)
{
	return t / (double)CLOCKS_PER_SEC * 1000;
}

double stdev(clock_t *times, int n)
{
	unsigned long total = 0;
	for (int i = 0; i < n; i++) {
		total += times[i];
	}
	unsigned long mean = total / n;
	double diff = 0;
	for (int i = 0; i < n; i++) {
		diff += pow(abs(times[i] - mean), 2);
	}
	return sqrt(diff / n);
}

double bench(int n, int x, int(*f)(int))
{
	clock_t *times = malloc(sizeof(clock_t)*n);
	for (int i = 0; i < n; i++) {
		clock_t t = clock();
		(*f)(x);
		t = clock() - t;
		times[i] = t;
		sleep(1);
	}
	double result = stdev(times, n);
	free(times);
	return result;
}

int fib(int n)
{
	if (n < 2) return n;
	return fib(n - 1) + fib(n - 2);
}

int alloc(int n)
{
	for (int i = 0; i < n; i++) {
		int size = 1024*1024*40;
		void *mem = malloc(size);
		memset(mem, 0, size);
		free(mem);
	}
}

int main()
{
	printf("Starting the CPU benchmark...\n");
	double variance = bench(10, 40, fib);
	printf("Variance: %lfms\n", clock2ms(variance));
	printf("Starting the memory benchmark...\n");
	variance = bench(10, 40, alloc);
	printf("Variance: %lfms\n", clock2ms(variance));
}

