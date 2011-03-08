#include <cstdio>
#include <cstdlib>
#include <iostream>
#include "fseq.h"

static void process(FILE * f) {
    char bytes[sizeof(mkfseq::fseq)];
    size_t bytesread = fread(bytes, 1, sizeof(mkfseq::fseq), f);
    if (bytesread != sizeof(mkfseq::fseq) && ferror(f) ) {
        fprintf(stderr, "couldn't read fseq\n");
        return;
    }
    mkfseq::fseq *fseq = mkfseq::fseq::from_bytes(bytes);
	fseq->to_csv(std::cout);
    return;
}

int main(int argc, char ** argv) {
    if (argc == 1) {
        process(stdin);
    }
    
    for (int i = 1; i < argc; ++i) {
        FILE * f = fopen(argv[i], "r");
        if (!f) { exit(1); }
        process(f);
        fclose(f);
    }
	return 0;
}
