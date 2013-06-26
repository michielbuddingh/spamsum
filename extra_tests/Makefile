# Builds the original spamsum tool, and one additional quick hack,
# required to run the from_git_repo.go script
CC = gcc
CFLAGS = -Wall -W

OBJS = spamsum.o edit_dist.o

all:	spamsum spamsum_compare

spamsum: $(OBJS)
	$(CC) -o spamsum spamsum_main.c $(OBJS)

spamsum_compare: $(OBJS)
	$(CC) -o spamsum_compare spamsum_compare.c $(OBJS)

clean:
	@rm -f $(OBJS) spamsum spamsum_compare *~
