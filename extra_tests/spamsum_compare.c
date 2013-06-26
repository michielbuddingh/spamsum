/* Copyright 2013, Michiel Buddingh, All rights reserved.  Use of this
   code is governed by version 2.0 or later of the Apache License,
   available at http://www.apache.org/licenses/LICENSE-2.0

   This is a quick-and-dirty hack to compare two spamsums from the
   command line */

#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include "spamsum.h"

int main(int argc, char *argv[]) {
    int result = spamsum_match(argv[1], argv[2]);
    printf("%d", result);
    return 0;
}
