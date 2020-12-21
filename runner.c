#include <errno.h>
#include <limits.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <sys/unistd.h>

static unsigned int parse_interval_or_die(const char*, char**);

static void make_output_dir(void);

int main(int argc, char* argv[])
{
	if(argc < 3)
	{
		fprintf(stderr, "usage: %s <interval>['s'|'m'|'h'|'d'] <nmap executable> [nmap arg...]\n", argv[0]);
		exit(EXIT_FAILURE);
	}

	char* endptr;
	unsigned int interval = parse_interval_or_die(argv[1], &endptr);
	while(*endptr != 0)
	{
		interval += parse_interval_or_die(endptr, &endptr);
	}

	make_output_dir();

	fprintf(stderr, "call %s every %u seconds\n", argv[2], interval);
	return 0;
}

static unsigned int parse_interval_or_die(const char* nptr, char** endptr)
{
	errno = 0;
	long l = strtol(nptr, endptr, 10);
	if(errno)
	{
		perror("strtol");
		exit(EXIT_FAILURE);
	}

	if(*nptr != '\0' && *endptr == nptr)
	{
		fprintf(stderr, "interval must be a number\n");
		exit(EXIT_FAILURE);
	}

	if(l < 0 || l > UINT_MAX)
	{
		fprintf(stderr, "number out of range: %ld\n", l);
		exit(EXIT_FAILURE);
	}

	switch(**endptr)
	{
	case '\0':
	case 's':
		break;
	case 'm':
		l *= 60;
		break;
	case 'h':
		l *= 3600;
		break;
	case 'd':
		l *= 86400;
		break;
	default:
		fprintf(stderr, "invalid suffix '%c'\n", **endptr);
		exit(EXIT_FAILURE);
	}
	if(**endptr != '\0') ++*endptr;

	return (unsigned int) l;
}

static void make_output_dir(void)
{
	const char OUTPUT_DIR[] = "./reports";

	struct stat statbuf;
	if(stat(OUTPUT_DIR, &statbuf) == 0)
	{
		if(S_ISDIR(statbuf.st_mode)) return;

		fprintf(stderr, "%s exists, but is not (or does not point to) a directory\n", OUTPUT_DIR);
		exit(EXIT_FAILURE);
	}

	switch(errno)
	{
	case ENOENT:
		break;
	default:
		perror("stat");
		exit(EXIT_FAILURE);
	}

	if(mkdir(OUTPUT_DIR, 0755) == -1)
	{
		perror("mkdir");
		exit(EXIT_FAILURE);
	}
}
