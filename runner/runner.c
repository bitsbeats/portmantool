#include <errno.h>
#include <limits.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <sys/unistd.h>
#include <sys/wait.h>
#include <time.h>

static const char OUTPUT_DIR[] = "./reports";
static const char OUTPUT_FILE[] = "output.xml";

static unsigned int parse_interval_or_die(const char*, char**);

static void make_output_dir(void);

static pid_t run(char* const*);

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

	int num_nmap_args = argc - 2;
	char** const nmap_argv = (char**) malloc(sizeof(char*) * (num_nmap_args+3));
	if(!nmap_argv)
	{
		perror("malloc");
		exit(EXIT_FAILURE);
	}

	int nmap_argc = 0;
	nmap_argv[nmap_argc++] = argv[2];
	nmap_argv[nmap_argc++] = "-oX";
	nmap_argv[nmap_argc++] = OUTPUT_FILE;
	for(int i = 1; i < num_nmap_args; ++i)
	{
		nmap_argv[nmap_argc++] = argv[2+i];
	}
	nmap_argv[nmap_argc++] = NULL;

	make_output_dir();

	char filename[33];
	while(1)
	{
		pid_t pid = run(nmap_argv);
		if(pid == -1)
		{
			sleep(1);
			continue;
		}

		int wstatus;
		if(waitpid(pid, &wstatus, 0) == -1) perror("waitpid");
		if(!WIFEXITED(wstatus) || WEXITSTATUS(wstatus) != 0) continue;

		time_t now = time(NULL);
		if(now == ((time_t) -1))
		{
			perror("time");
			continue;
		}

		struct tm tm;
		if(localtime_r(&now, &tm) == NULL)
		{
			perror("localtime_r");
			continue;
		}

		if(strftime(filename, sizeof(filename), "reports/scan-%Y%m%d-%H%M%S.xml", &tm) == 0) continue;
		if(rename(OUTPUT_FILE, filename) == -1) perror("rename");

		sleep(interval);
	}

	free(nmap_argv);
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

static pid_t run(char* const argv[])
{
	pid_t pid = fork();
	if(pid == -1)
	{
		perror("fork");
		return -1;
	}

	if(pid == 0)
	{
		execvp(argv[0], argv);
		perror("execvp");
		exit(EXIT_FAILURE);
	}

	return pid;
}
