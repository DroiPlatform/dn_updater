package aes_util;
/*
#cgo LDFLAGS: -lTyd64 -L. -lcrypto -I./engine -Wall -O3 -lm -lz -ldl
#include<errno.h>
#include<stdio.h>
#include<stdlib.h>
#include<string.h>
#include"engine.h"
#define DEC_BUF 1024
int go_aes_dec(unsigned char* in_buf, int in_size, unsigned char* salt, unsigned char* key, int key_len, unsigned char* out_buf, int debug) {
  EVP_CIPHER_CTX e_ctx;
  EVP_CIPHER_CTX d_ctx;
  int out_len = 0;
  if ( aes_init_salt(CBC_128, key, key_len, salt, &e_ctx, &d_ctx ) ) {
  //if ( aes_init(CBC_128, key, key_len, &e_ctx, &d_ctx) ) {
    fprintf(stderr, "context initialization failed. (%s)\n", strerror(errno));
  }
  //out_len =  aes_decrypt(&d_ctx, in_buf, in_size, out_buf, aes_encrypt_mem_size(in_size));
  //out_len =  aes_decrypt(&d_ctx, in_buf, in_size, out_buf, enc_buf_size(in_size));
  out_len =  aes_decrypt(&d_ctx, in_buf, in_size, out_buf, in_size);
  if (out_len < 0 && debug) {
    fprintf(stderr, "decrypted len = %d < 0\n", out_len);
  }
  if (aes_deinit(&e_ctx, &d_ctx)) {
    fprintf(stderr, "context deinit failed. (%s)\n", strerror(errno));
  }
  return out_len;
}
int go_aes_enc(unsigned char* in_buf, int in_size, unsigned char* salt, unsigned char* key, int key_len, unsigned char* out_buf, int debug) {
  EVP_CIPHER_CTX e_ctx;
  EVP_CIPHER_CTX d_ctx;
  int out_len = 0;
  if ( aes_init_salt(CBC_128, key, key_len, salt, &e_ctx, &d_ctx ) ) {
  //if ( aes_init(CBC_128, key, key_len, &e_ctx, &d_ctx) ) {
    fprintf(stderr, "context initialization failed. (%s)\n", strerror(errno));
  }
  //out_len =  aes_encrypt(&e_ctx, in_buf, in_size, out_buf, aes_encrypt_mem_size(in_size));
  out_len =  aes_encrypt(&e_ctx, in_buf, in_size, out_buf, enc_buf_size(in_size));
  if (out_len < 0 && debug) {
    fprintf(stderr, "encrypted len = %d < 0\n", out_len);
  }
  if (aes_deinit(&e_ctx, &d_ctx)) {
    fprintf(stderr, "context deinit failed. (%s)\n", strerror(errno));
  }
  return out_len;
}
int enc_buf_size(int in_size) {
  return aes_encryption_length(in_size);
}
*/
import "C";
import "errors";
import "fmt";
import "unsafe";

import util "tyd_util";

/* App Rev related info */
//const SALT = "#3|3nboe";
//const KEY = "#3|3n";
/* magic prime number */
/******
** rev 1: 239
******/
//const PRIME = uint8(239);
const APP_REV = uint8(1);
const MAGIC = 149;
/******
** Primes for Platform Key
******/
/*
const PRIME_CLOUD = uint64(2286618416893);
const PRIME_MOBILE = uint64(83288996057);
const PRIME_RESTFUL = uint64(1855199054551);
//*/
const PRIME_CLOUD = uint64(997);
const PRIME_MOBILE = uint64(607);
const PRIME_RESTFUL = uint64(877);

/* global variables */
var isInit bool;
var PRIME map[uint8]uint8;
var KEY map[uint8][]byte;
var SALT map[uint8][]byte;

func GetPrime(index uint8) (uint8) {
  if !isInit {
    initAES();
  }
  return PRIME[index];
}

func GetKey(index uint8) ([]byte) {
  if !isInit {
    initAES();
  }
  return KEY[index];
}

func getSalt(index uint8) ([]byte) {
  if !isInit {
    initAES();
  }
  return SALT[index];
}

func initAES() {
  /* prime */
  PRIME = make(map[uint8]uint8);
  PRIME[uint8(1)] = uint8(239);
  PRIME[uint8(2)] = uint8(109);
  PRIME[uint8(3)] = uint8(127);
  PRIME[uint8(4)] = uint8(139);
  PRIME[uint8(5)] = uint8(251);
  PRIME[uint8(6)] = uint8(223);
  /* key */
  KEY = make(map[uint8][]byte);
  KEY[uint8(1)] = util.Ui32ToBL(uint32(2603052521));
  KEY[uint8(2)] = util.Ui32ToBL(uint32(3412296483));
  KEY[uint8(3)] = util.Ui32ToBL(uint32(1412296443));
  KEY[uint8(4)] = util.Ui32ToBL(uint32(2412296719));
  KEY[uint8(5)] = util.Ui32ToBL(uint32(3448256482));
  KEY[uint8(6)] = util.Ui32ToBL(uint32(4122964832));
  /* salt */
  SALT = make(map[uint8][]byte);
  SALT[uint8(1)] = util.Ui64ToBL(uint64(122603052521));
  SALT[uint8(2)] = util.Ui64ToBL(uint64(763412296483));
  SALT[uint8(3)] = util.Ui64ToBL(uint64(941412296443));
  SALT[uint8(4)] = util.Ui64ToBL(uint64(752412296719));
  SALT[uint8(5)] = util.Ui64ToBL(uint64(563448256482));
  SALT[uint8(6)] = util.Ui64ToBL(uint64(284122964864));
}

func Aes_enc(plaintext []byte, key []byte, b_debug bool) ([]byte, error) {
  if !isInit {
    initAES();
  }
  enc_buf := int(C.enc_buf_size(C.int(len(plaintext))));
  debug := 0;
  ciphertext := make([]byte, enc_buf);
  //fmt.Printf("key: %x\n", key);
  enc_len := C.go_aes_enc((*C.uchar)(unsafe.Pointer(&plaintext[0])), C.int(len(plaintext)), (*C.uchar)(unsafe.Pointer(&(getSalt(APP_REV)[0]))), (*C.uchar)(unsafe.Pointer(&key[0])), C.int(len(key)), (*C.uchar)(unsafe.Pointer(&ciphertext[0])), C.int(debug));
  if enc_len < 0 {
    errmsg := fmt.Sprintf("enc_len = %d < 0", enc_len);
    if b_debug {
      util.PrintTime();
      fmt.Printf("%s\n", errmsg);
    }
    e := errors.New(errmsg);
    return nil, e;
  } else {
    return ciphertext[:enc_len], nil;
  }
}

func Aes_dec(app_rev uint8, ciphertext []byte, key []byte, b_debug bool) ([]byte, error) {
  dec_buf := int(C.enc_buf_size(C.int(len(ciphertext))));
  debug := 0;
  /*
  if b_debug {
    debug = 1;
    util.PrintTime();
    fmt.Printf("dec key len: %d\n", len(key));
  }
  */
  decrypted := make([]byte, dec_buf);
  //fmt.Printf("dec key: %x\n", key);
  dec_len := C.go_aes_dec((*C.uchar)(unsafe.Pointer(&ciphertext[0])), C.int(len(ciphertext)), (*C.uchar)(unsafe.Pointer(&(getSalt(app_rev)[0]))), (*C.uchar)(unsafe.Pointer(&(key[0]))), C.int(len(key)), (*C.uchar)(unsafe.Pointer(&decrypted[0])), C.int(debug));
  if dec_len < 0 {
    errmsg := fmt.Sprintf("dec_len = %d < 0", dec_len);
    /*
    if b_debug {
      util.PrintTime();
      fmt.Printf("%s\n", errmsg);
    }
    */
    e := errors.New(errmsg);
    return nil, e;
  } else {
    return decrypted[:dec_len], nil;
  }
}

