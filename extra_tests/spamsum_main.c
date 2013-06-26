/*
  this is a checksum routine that is specifically designed for spam.
  Copyright Andrew Tridgell <tridge@samba.org> 2002

  This code is released under the GNU General Public License version 2
  or later.  Alteratively, you may also use this code under the terms
  of the Perl Artistic license.

  If you wish to distribute this code under the terms of a different
  free software license then please ask me. If there is a good reason
  then I will probably say yes.
*/
#include <stdio.h>
#include <unistd.h>
#include <stdlib.h>
#include "spamsum.h"

static void show_help(void)
{
	printf("");
}

int main(int argc, char *argv[])
{
	char *sum;
	extern char *optarg;
	extern int optind;
	int c;
	char *dbname = NULL;
	u32 score;
	int i;
	u32 flags = 0;
	u32 block_size = 0;
	u32 threshold = 90;

	while ((c = getopt(argc, argv, "B:WHd:c:C:hT:")) != -1) {
		switch (c) {
		case 'W':
			flags |= FLAG_IGNORE_WHITESPACE;
			break;

		case 'H':
			flags |= FLAG_IGNORE_HEADERS;
			break;

		case 'd':
			dbname = optarg;
			break;

		case 'B':
			block_size = atoi(optarg);
			break;

		case 'T':
			threshold = atoi(optarg);
			break;

		case 'c':
			if (!dbname) {
				show_help();
				exit(1);
			}
			score = spamsum_match_db(dbname, optarg,
						 threshold);
			printf("%u\n", score);
			exit(score >= threshold ? 0 : 2);

		case 'C':
			if (!dbname) {
				show_help();
				exit(1);
			}
			score = spamsum_match_db(dbname,
						 spamsum_file(optarg, flags,
							      block_size),
						 threshold);
			printf("%u\n", score);
			exit(score >= threshold ? 0 : 2);

		case 'h':
		default:
			show_help();
			exit(0);
		}
	}

	argc -= optind;
	argv += optind;

	if (argc == 0) {
		show_help();
		return 0;
	}

	/* compute the spamsum on a list of files */
	for (i=0;i<argc;i++) {
		sum = spamsum_file(argv[i], flags, block_size);
		printf("%s\n", sum);
		free(sum);
	}

	return 0;
}
