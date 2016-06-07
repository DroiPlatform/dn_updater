package main;
/*
#cgo LDFLAGS: -lm
#include <errno.h>
#include <stddef.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include "cJSON.h"

#define OK 0
#define NF_ROOT -1
#define NF_1ST_NODE -2
#define NF_2ND_NODES -3
#define ERR_SIZE -4

#define SIZE_BUF 65536

#define UNKNOWN -256

int getIPs(char* raw, char* output) {
  cJSON* root = cJSON_Parse(raw);
  int occupied = 0;
  int r;
  if (root == NULL) {
    return NF_ROOT;
  } else {
    cJSON* node = cJSON_GetObjectItem(root, "node");
    if (node == NULL) {
      return NF_1ST_NODE;
    } else {
      cJSON* nodes = cJSON_GetObjectItem(node, "nodes");
      if (nodes == NULL) {
        return NF_2ND_NODES;
      } else {
        int num_pod = cJSON_GetArraySize(nodes);
        int i, j, num_ip;
        int pos = 0, index;
        char* key, * raw_key;
        cJSON* current,* value, * subset,* addrs,* ip;
        for (i = 0; i < num_pod; i++) {
          current = cJSON_GetArrayItem(nodes, i);
          raw_key = cJSON_GetObjectItem(current, "key")->valuestring;
          for (index = 0; index < strlen(raw_key); index++) {
            if ((char)(raw_key + index)[0] == '/') {
              pos = index;
            }
          }
          key = raw_key + pos + 1;
          if (i == 0) {
            r = snprintf(output + occupied, SIZE_BUF - occupied, "%s?", key);
          } else {
            r = snprintf(output + occupied, SIZE_BUF - occupied, ";%s?", key);
          }
          occupied += r;
          value = cJSON_GetObjectItem(current, "value");
          subset = cJSON_GetObjectItem(value, "subsets");
          int size_subset = cJSON_GetArraySize(subset);
          if (size_subset == 1) {
            addrs = cJSON_GetObjectItem(cJSON_GetArrayItem(subset, 0), "addresses");
            num_ip = cJSON_GetArraySize(addrs);
            for (j = 0; j < num_ip; j++) {
              ip = cJSON_GetArrayItem(addrs, j);
              if (j == 0) {
                r = snprintf(output + occupied, SIZE_BUF - occupied, "%s", cJSON_GetObjectItem(ip, "ip")->valuestring);
              } else {
                r = snprintf(output + occupied, SIZE_BUF - occupied, ",%s", cJSON_GetObjectItem(ip, "ip")->valuestring);
              }
              occupied += r;
            }
          } else if (size_subset == 0) {
            r = snprintf(output + occupied, SIZE_BUF - occupied, "-1");
            occupied += r;
          } else {
            fprintf(stderr, "raw: %s\n key: %s\n size_subset: %d\n", raw, key, size_subset);
            return ERR_SIZE;
          }
        }
        fprintf(stderr, "out:\n%s\n", output);
        return occupied;
      }
      return UNKNOWN;
    }
    return UNKNOWN;
  }
  return UNKNOWN;
}

*/
import "C";

import "fmt";
import "io/ioutil";
import "strings";
import "time";
import "unicode/utf8";
import "unsafe";

import util "tyd_util";

func Caduceus() {
  for {
    Asclepius();
    time.Sleep(time.Duration(opts.poll) * time.Second);
  }
}

func Asclepius() {
  for _, v := range etcds {
    resp, err := buffer.client.Do(buffer.requests[v]);
    if err != nil {
      util.GenericLogPrinter(opts.host, "ERR", fmt.Sprintf("failed to do request: %s", err.Error()), TOPIC);
      return;
    } else {
      buffer.raw_json, err = ioutil.ReadAll(resp.Body);
      if err != nil {
        util.GenericLogPrinter(opts.host, "ERR", fmt.Sprintf("failed to read from respnose: %s", err.Error()), TOPIC);
        return;
      } else {
        replace(string(buffer.raw_json), "\\\"", "\"", buffer.escape_json);
        replace(string(buffer.escape_json), "\\n", "", buffer.escape_json);
        replace(string(buffer.escape_json), "\"{", "{", buffer.escape_json);
        replace(string(buffer.escape_json), "}\"", "}", buffer.escape_json);
        result := int(C.getIPs((*C.char)(unsafe.Pointer(&buffer.escape_json[0])), (*C.char)(unsafe.Pointer(&buffer.raw_json[0]))));
        if result <= 0 {
          util.GenericLogPrinter(opts.host, "ERR", fmt.Sprintf("failed to parse JSON from etcd: %d", result), TOPIC);
          return;
        } else {
          fmt.Printf("---\nresult %s\n===\n", string(buffer.raw_json[:result]));
        }
        return;
      }
    }
  }
}

func replace(s, old, new string, out []byte) {
  // Compute number of replacements.
  var m int;
  if m = strings.Count(s, old); m == 0 {
    return; // avoid allocation
  }

  // Apply replacements to buffer.
  w := 0;
  start := 0;
  for i := 0; i < m; i++ {
    j := start;
    if len(old) == 0 {
      if i > 0 {
        _, wid := utf8.DecodeRuneInString(s[start:]);
        j += wid;
      }
    } else {
      j += strings.Index(s[start:], old);
    }
    w += copy(out[w:], s[start:j]);
    w += copy(out[w:], new);
    start = j + len(old);
  }
  w += copy(out[w:], s[start:]);
}

