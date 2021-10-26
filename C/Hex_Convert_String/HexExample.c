#include <stdio.h>
#include <stdint.h>

int hex_convert_hexstring(uint8_t* data, uint8_t len, uint8_t* result){
	int i=0;
	int idx=0;
	for(i=0;i<len;i++){
		result[idx++]=(*(data+i))>>4 & 0x0f;
		result[idx++]=(*(data+i))& 0x0f;
	}

	for(i=0;i<idx;i++){
		if(result[i]>=10){
			result[i]=result[i]-10+'A';
		}else{
			result[i]=result[i]+'0';
		}
	}
    return idx;
}


int hexstring_convert_hex(uint8_t* data,uint8_t len,uint8_t* result){
	int idx=0;
	int i=0;
	for (i=0;i<len;i++){
		if(*(data+i)>='A'){
			*(data+i)=*(data+i)-'A'+10;
		}else{
			*(data+i)=*(data+i)-'0';
		}
	}
	i=0;
	for(idx=0;idx<len/2;idx++){
		result[idx]=*(data+i++)<<4  | *(data+i++) & 0x0f;
	}
    return idx;
}


int main(){
    int i=0;
    int string_len=0;
    int hex_len=0;
    uint8_t hex[16]={0xA1,0xA2,0xA3,0xA4,0xB5,0xB6,0xB7,0xB8,0xC9,0xC0,0xCA,0xCB,0xDC,0xDD,0xDE,0xDF};
    uint8_t hex_string[64];
    uint8_t hex_result[64];

    string_len=hex_convert_hexstring(hex,sizeof(hex),hex_string);

    for(i=0;i<string_len;i++){
        printf("%C ",hex_string[i]);
    }
    printf("\r\n");
    printf("%s\r\n",hex_string); //A1A2A3A4B5B6B7B8C9C0CACBDCDDDEDF

    hex_len=hexstring_convert_hex(hex_string,string_len,hex_result);
    
    for(i=0;i<hex_len;i++){
        printf("%02X ",hex_result[i]); //A1 A2 A3 A4 B5 B6 B7 B8 C9 C0 CA CB DC DD DE DF
    }
    printf("\r\n");
    return 0;

}