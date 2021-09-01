/*
 * aria.h
 *
 *  Created on: 2021. 8. 6.
 *      Author: piljoo
 */



#ifndef ARIA_H_
#define ARIA_H_


#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
typedef unsigned char Byte;

void ARIA_128test();
int EncKeySetup(const Byte *w0, Byte *e, int keyBits);
int DecKeySetup(const Byte *w0, Byte *d, int keyBits);
void Crypt(const Byte *p, int R, const Byte *e, Byte *c);

uint8_t EnCrypt(uint8_t *p_text, uint8_t *en_text,uint8_t len);
uint8_t DeCrypt(uint8_t *en_text,uint8_t *de_text,uint8_t len);

#endif
