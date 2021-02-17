#include <stdio.h>
#include <stdlib.h>

int main(int argc, char** argv)
{
	if(argc != 3)
	{
		fprintf(stderr, "usage: %s <oldpath> <newpath>\n", argv[0]);
		exit(EXIT_FAILURE);
	}

	if(rename(argv[1], argv[2]) == -1)
	{
		perror("rename");
		exit(EXIT_FAILURE);
	}

	return 0;
}
