#include "fseq.h"

mkfseq::fseq* mkfseq::fseq::from_bytes(const char * b) {
    fseq *f = new fseq;
    memcpy(&(f->_header), b, sizeof(header));
    b += sizeof(header);
    size_t framect = (f->_header.data_format + 1) * 128;
    memcpy(f->_frames, b, sizeof(frame) * framect);
	b += sizeof(frame) * framect;
	f->_end_sysex = *b;
    return f;
}
