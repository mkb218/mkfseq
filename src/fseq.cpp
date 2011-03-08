#include "fseq.h"
#include <iostream>
mkfseq::fseq* mkfseq::fseq::from_bytes(const char * b) {
    fseq *f = new fseq;
    memcpy(&(f->_header), b, sizeof(header));
    b += sizeof(header);
    size_t framect = (f->_header.data_format + 1) * 128;
    memcpy(f->_frames, b, sizeof(frame) * framect);
	b += sizeof(frame) * framect;
	f->_endsysex = *b;
    return f;
}


void mkfseq::fseq::to_csv(std::ostream & csv) const {
    csv << "header\n";
    for (int i = 0; i < 8; ++i) {
        csv << _header.name[i];
    }
    csv << "\n";
    csv << "loop start," << _header.loop_start << "\n";
    csv << "loop end," << _header.loop_end << "\n";
	csv << "loop_mode," << (int)_header.loop_mode << "\n";
	csv << "speed_adjust," << (int)_header.speed_adjust << "\n";
	csv << "velocity_sens," << (int)_header.velocity_sens << "\n";
	csv << "pitch_mode," << (int)_header.pitch_mode << "\n";
	csv << "note_assign," << _header.note_assign << "\n";
	csv << "pitch_tuning," << (int)_header.pitch_tuning << "\n";
	csv << "seq_delay," << (int)_header.seq_delay << "\n";
	csv << "data_format," << (int)_header.data_format << "\n";
	csv << "valid_end," << _header.valid_end << "\n";
	
	csv << "frame,fund_pitch,voicedop0,voicedlevel0,voicedop1,voicedlevel1,voicedop2,voicedlevel2,voicedop3,voicedlevel3,voicedop4,voicedlevel4,voicedop5,voicedlevel5,voicedop6,voicedlevel6,voicedop7,voicedlevel7,unvoicedop0,unvoicedlevel0,unvoicedop1,unvoicedlevel1,unvoicedop2,unvoicedlevel2,unvoicedop3,unvoicedlevel3,unvoicedop4,unvoicedlevel4,unvoicedop5,unvoicedlevel5,unvoicedop6,unvoicedlevel6,unvoicedop7,unvoicedlevel7\n";
	for (int i = 0; i <= valid_end; ++i) {
		csv << i << ",";
		csv << _frames[i].fund_pitch;
		for (int j = 0; j < 8; ++j) {
			csv << ",";
			int hi = _frames[i].voiced_freq_hi[j];
			int lo = _frames[i].voiced_freq_lo[j];
			int out = ((hi << 7)|lo);
			csv << out << ",";
			int level = _frames[i].voiced_level[j];
			csv << level;
		}
		for (int j = 0; j < 8; ++j) {
			csv << ",";
			int hi = _frames[i].unvoiced_freq_hi[j];
			int lo = _frames[i].unvoiced_freq_lo[j];
			int out = ((hi << 7)|lo);
			csv << out << ",";
			csv << (int)_frames[i].unvoiced_level[j];
		}
		csv << "\n";
	}
		
}