#ifndef __H2P_FSEQ__
#define __H2P_FSEQ__

#include <cstring>
#include <exception>

namespace mkfseq {
    const size_t SYSEX_HEAD_LEN = 9;
    class fseq_e : public std::exception {};
    
#pragma pack(1)
    class sevenb_short {
    public:
        sevenb_short() { setval(0); }
        explicit sevenb_short(short int val) {
            setval(val);
        }
        
        operator short() const {
            return _bytes[0] + _bytes[1];
        }
        
        sevenb_short & operator=(short int i) throw (fseq_e) {
            setval(i);
            return *this;
        }
                
        const char * getbytes() { return _bytes; }
    private:
        void setval(short int i) throw (fseq_e) {
            if (i >= (1 << 14)) { throw fseq_e(); }
            if (i < 0) { throw fseq_e(); }
            i &= 0x3fff;
            _bytes[0] = i >> 7;
            _bytes[1] = i & 0x7f;
        }
        char _bytes[2];
    };

    template<int Offset>
    class sevenb_char {
    public:
        sevenb_char<Offset>() { setval(0); }
        explicit sevenb_char<Offset>(char val) {
            setval(val);
        }
        
        operator char() const {
            return _val + Offset;
        }
        
        sevenb_char<Offset> & operator=(char i) throw (fseq_e) {
            setval(i);
            return *this;
        }
        
    private:
        void setval(char i) {
            if (i < Offset) { throw fseq_e(); }
            _val= (i - Offset) & 0x7f;
        }
        char _val;
    };

    struct header {
        char sysex_header[SYSEX_HEAD_LEN];
        char name[8];
        char reserved0[8];
        sevenb_short loop_start, loop_end;
        char loop_mode;
        char speed_adjust;
        char velocity_sens;
        char pitch_mode;
        sevenb_char<0> note_assign;
        sevenb_char<-63> pitch_tuning;
        sevenb_char<0> seq_delay;
        char data_format;
        char reserved1[2];
        sevenb_short valid_end;
    };
    
    struct frame {
        sevenb_short fund_pitch;
        sevenb_char<0> voiced_freq_hi[8];
        char reserved[0]; // ???
        sevenb_char<0> voiced_freq_lo[8];
        sevenb_char<0> voiced_level[8];
        sevenb_char<0> unvoiced_freq_hi[8];
        sevenb_char<0> unvoiced_freq_lo[8];
        sevenb_char<0> unvoiced_level[8];
    };

    class fseq {
    public:
        static fseq * from_bytes(const char * b);
        fseq() { memset(this, 0, sizeof(fseq)); }
        const char * to_sysex() const;
    private:
		void setchecksum();
        header _header;
        frame _frames[512];
		char _checksum;
		char _endsysex;
    };
}

#endif