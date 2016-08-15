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

#define DEBUG 0

int getIPs(char* raw, char* output) {
  cJSON* root = cJSON_Parse(raw);
  int occupied = 0;
  int r;
  if (root == NULL) {
    if (DEBUG) fprintf(stderr, "NF_ROOT\n");
    return NF_ROOT;
  } else {
    cJSON* node = cJSON_GetObjectItem(root, "node");
    if (node == NULL) {
      if (DEBUG) fprintf(stderr, "NF_1ST_NODE\n");
      cJSON_Delete(root);
      return NF_1ST_NODE;
    } else {
      cJSON* nodes = cJSON_GetObjectItem(node, "nodes");
      if (nodes == NULL) {
        if (DEBUG) fprintf(stderr, "NF_2ND_NODES\n");
        cJSON_Delete(root);
        return NF_2ND_NODES;
      } else {
        int num_pod = cJSON_GetArraySize(nodes);
        int i, j, num_ip;
        int pos = 0, index;
        char* key, * raw_key;
        cJSON* current,* value, * subset,* addrs,* ip;
        if (DEBUG) fprintf(stderr, "num_pod: %d\n", num_pod);
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
          if (DEBUG) fprintf(stderr, "%s: occupied: %d, expected subset: %d\n", key, occupied, size_subset);
          if (size_subset == 1) {
            addrs = cJSON_GetObjectItem(cJSON_GetArrayItem(subset, 0), "addresses");
            if (addrs == NULL) {
              r = snprintf(output + occupied, SIZE_BUF - occupied, "-1");
            } else {
              num_ip = cJSON_GetArraySize(addrs);
              if (DEBUG) fprintf(stderr, "%s: occupied: %d, expected ips: %d\n", key, occupied, num_ip);
              for (j = 0; j < num_ip; j++) {
                ip = cJSON_GetArrayItem(addrs, j);
                if (j == 0) {
                  r = snprintf(output + occupied, SIZE_BUF - occupied, "%s", cJSON_GetObjectItem(ip, "ip")->valuestring);
                } else {
                  r = snprintf(output + occupied, SIZE_BUF - occupied, ",%s", cJSON_GetObjectItem(ip, "ip")->valuestring);
                }
                occupied += r;
                if (DEBUG) fprintf(stderr, "%s: size_ip (%d/%d), occupied: %d\n", key, j, num_ip, occupied);
              }
            }
            if (DEBUG) fprintf(stderr, "%s: occupied: %d, ips done: %d\n", key, occupied, j);
          } else if (size_subset == 0) {
            r = snprintf(output + occupied, SIZE_BUF - occupied, "-1");
            occupied += r;
            if (DEBUG) fprintf(stderr, "%s: size_ip (0), occupied: %d\n", key, occupied);
          } else {
            if (DEBUG) fprintf(stderr, "raw: %s\n key: %s\n size_subset: %d\n", raw, key, size_subset);
            cJSON_Delete(root);
            return ERR_SIZE;
          }
        }
        if (DEBUG) fprintf(stderr, "out:\n%s\n", output);
        cJSON_Delete(root);
        return occupied;
      }
      if (DEBUG) fprintf(stderr, "unknown situation for dn_updater (1)\n");
      cJSON_Delete(root);
      return UNKNOWN;
    }
    if (DEBUG) fprintf(stderr, "unknown situation for dn_updater (2)\n");
    cJSON_Delete(root);
    return UNKNOWN;
  }
  if (DEBUG) fprintf(stderr, "unknown situation for dn_updater (3)\n");
  cJSON_Delete(root);
  return UNKNOWN;
}

*/
import "C";

import "bufio";
import "errors";
import "fmt";
import "io/ioutil";
import "os";
import "os/exec";
import "strconv";
import "strings";
import "syscall";
import "time";
import "unicode/utf8";
import "unsafe";

func Caduceus() {
  ttl := time.Duration(opts.suicide)
  if opts.test {
    ttl *= time.Second;
  } else {
    ttl *= time.Hour;
  }
  poll := time.Duration(opts.poll) * time.Second;
  death := time.Duration(0);
  for {
    Asclepius();
    rnd_cnt++;
    time.Sleep(poll);
    //GenericLogPrinter(opts.host, "DEBUG", fmt.Sprintf("[Caduceus] ttl: %d", ttl), TOPIC);
    if opts.suicide != 0 {
      ttl -= poll;
      if ttl <= death {
        GenericLogPrinter(opts.host, "WARN", fmt.Sprintf("[Caduceus] time's up, suicide!"), TOPIC);
        os.Remove(HEALTH);
      }
    }
  }
}

func Asclepius() {
  for _, v := range etcds {
    resp, err := buffer.client.Do(buffer.requests[v]);
    if err != nil {
      GenericLogPrinter(opts.host, "ERR", fmt.Sprintf("[Asclepius] failed to do request: %s", err.Error()), TOPIC);
      return;
    } else {
      buffer.raw_json, err = ioutil.ReadAll(resp.Body);
      if err != nil {
        GenericLogPrinter(opts.host, "ERR", fmt.Sprintf("[Asclepius] failed to read from respnose: %s", err.Error()), TOPIC);
        return;
      } else {
        /* clear tmp buffer */
        /*
        for i:= 0; i < SIZE_BUF; i++ {
          buffer.escape_json[i] = byte(0);
        }
        //*/
        if opts.debug {
          GenericLogPrinter(opts.host, "DEBUG", fmt.Sprintf("[Asclepius] Response of %d (%d) bytes from etcd %s received.", len(buffer.raw_json), resp.ContentLength, v), TOPIC);
        }
        /* clear escape */
        replace(string(buffer.raw_json), "\\\"", "\"", buffer.escape_json);
        replace(string(buffer.escape_json), "\\n", "", buffer.escape_json);
        replace(string(buffer.escape_json), "\"{", "{", buffer.escape_json);
        replace(string(buffer.escape_json), "}\"", "}", buffer.escape_json);
        if opts.debug {
//          GenericLogPrinter(opts.host, "DEBUG", fmt.Sprintf("cJSON input: %s", buffer.escape_json), TOPIC);
        }
        //fmt.Fprintf(os.Stderr, "cJSON input: %s\n", buffer.escape_json);
        /* get result from JSON */
        result := int(C.getIPs((*C.char)(unsafe.Pointer(&buffer.escape_json[0])), (*C.char)(unsafe.Pointer(&buffer.raw_json[0]))));
        if result <= 0 {
          GenericLogPrinter(opts.host, "ERR", fmt.Sprintf("[Asclepius] failed to parse JSON from etcd: %d", result), TOPIC);
          return;
        } else {
          if opts.debug {
            GenericLogPrinter(opts.host, "DEBUG", fmt.Sprintf("[Asclepius] JSON parse result: %d", result), TOPIC);
          }
          /* parse and store host/ip pair */
          ///*
          err := fillHosts(string(buffer.raw_json[:result]));
          if err != nil {
            GenericLogPrinter(opts.host, "ERR", fmt.Sprintf("[Asclepius] failed to parse host records from etcd: %s", err.Error()), TOPIC);
            return;
          }
          fillInc();
          throttle := uint8(opts.throttle);
          if throttle == uint8(0) {
            throttle = uint8(1);
          }
          if rnd_cnt % throttle == uint8(0) {
            printHistory();
          }
          //*/
        }
      }
    }
  }
}

func printHistory() {
  for _, v := range etcds {
    observation := kube_status[v].observation;
    info := kube_status[v].current;
    output := fmt.Sprintf("[printHistory] Once offline svcs of %s:", v);
    first := true;
//    fmt.Fprintf(os.Stderr, "etcd %s: observation: %v, current: %v\n", v, kube_status[v].observation, kube_status[v].current);
    for k, cnt := range observation.list {
      if trend, ok := observation.trend[k]; !ok {
        /* no trend */
        if first {
          output = fmt.Sprintf("%s %s[%d/%v]", output, k, cnt, nil);
          first = false;
        } else {
          output = fmt.Sprintf("%s; %s[%d/%v]", output, k, cnt, nil);
        }
      } else {
        /* got trend record */
        if first {
          output = fmt.Sprintf("%s %s[%d/%d]", output, k, cnt, trend);
          first = false;
        } else {
          output = fmt.Sprintf("%s; %s[%d/%d]", output, k, cnt, trend);
        }
      }
    }
    if first {
      output = fmt.Sprintf("All %d services were alive...", len(info.pod_ip));
    }
    GenericLogPrinter(opts.host, "INFO", output, TOPIC);
  }
}

func printKubeInfo(ki *KubeInfo) {
  if opts.debug {
    for k, v := range ki.pod_ip {
      GenericLogPrinter(opts.host, "DEBUG", fmt.Sprintf("[printKubeInfo] {%s: %s}", k, v), TOPIC);
    }
  }
}

func isIPEqual(l string, r string) (bool) {
  lips := strings.Split(l, ",");
  rips := strings.Split(r, ",");
  if len(lips) != len(rips) {
    return false;
  } else {
    find := false;
    for _, lv := range lips {
      find = false;
      for _, rv := range rips {
        if lv == rv {
          find = true;
          break;
        }
      }
      if !find {
        return find;
      }
    }
    return true;
  }
  return false;
}

func isInfoEqual(l *KubeInfo, r *KubeInfo) (bool) {
  /* dummy code */
  /*
  return true;
  //*/
  if len(l.pod_ip) != len(r.pod_ip) {
    return false;
  } else {
    equal := false;
    for k, _ := range l.pod_ip {
      if _, ok := r.pod_ip[k]; !ok {
        if opts.debug {
          GenericLogPrinter(opts.host, "DEBUG", fmt.Sprintf("[isInfoEqual] %s cannot be found in right set", k), TOPIC);
        }
        return false;
      } else {
        /* testing pod */
        /*
        if opts.debug {
          if k == "hello-curl-http" {
            fmt.Printf("%s: l: %s, r: %s\n", k, l.pod_ip[k], r.pod_ip[k]);
          }
        }
        //*/
        /* check if incoming pod IP set equals to current pod IP set for specific pod */
        equal = isIPEqual(l.pod_ip[k], r.pod_ip[k]);
        if !equal {
          if opts.debug {
            GenericLogPrinter(opts.host, "DEBUG", fmt.Sprintf("[isInfoEqual] %s is clean? %v", k, equal), TOPIC);
          }
          return equal;
        }
      }
    }
    return true;
  }
  return false;
}

func flushInfo(ki *KubeInfo) {
  for k, _ := range ki.pod_ip {
    /* flush pod_ip */
    delete(ki.pod_ip, k);
  }
}

func fillHosts(in string) (error) {
  /* flush old buffer */
  for k, _ := range buffer.hosts {
    delete(buffer.hosts, k);
  }
  /* split return string into records of format "host?ip,ip,..." */
  records := strings.Split(in, ";");
  for _, v := range records {
    hosts := strings.Split(v, "?");
    if len(hosts) != 2 {
      return errors.New(fmt.Sprintf("(ETCD, internal) host format error: %s", in));
    } else {
      buffer.hosts[hosts[0]] = hosts[1];
    }
  }
  return nil;
}

func fillInfo(ki *KubeInfo) (error) {
  for k, v := range buffer.hosts {
    if v == "" {
      continue;
    } else {
      ips := strings.Split(v, ",");
      if len(ips) == 1 && ips[0] == "-1" {
        continue;
      } else {
        for i := 0; i < len(ips); i++ {
          if err := checkIPFmt(ips[i]); err != nil {
            GenericLogPrinter(opts.host, "ERR", fmt.Sprintf("[fillInfo] ips[%s(%d)]: %s of %x\n", k, i, ips[i], []byte(v)), TOPIC);
            flushInfo(ki);
            return err;
          }
        }
        ki.pod_ip[k] = v;
      }
    }
  }
  return nil;
}

func updateKubeStatus(key string, firstborn bool, clean bool, inc *KubeInfo, current *KubeInfo, past *KubeInfo, observation *Observation) {
  /*
  flushInfo(kube_status[key].inc);
  flushInfo(kube_status[key].current);
  */
  domain := kube_status[key].domain;
  delete(kube_status, key);
  kube_status[key] = KubeStatus {firstborn: firstborn, clean: clean, domain: domain, inc: inc, current: current, past: past, observation: observation};
}

func printHosts(dest *os.File, key string) {
  info := kube_status[key].current;
  for k, v := range info.pod_ip {
    ips := strings.Split(v, ",");
    for _, ip := range ips {
      fmt.Fprintf(dest, "%s %s.%s %s\n", ip, k, kube_status[key].domain, k);
    }
  }
  observation := kube_status[key].observation;
  if opts.debug {
    GenericLogPrinter(opts.host, "DEBUG", fmt.Sprintf("[printHosts] IPs of %d pods are written, prepared to write %d inactive pods", len(info.pod_ip), len(observation.trend)), TOPIC);
  }
  /* for those unavailable pods */
  for k, _ := range observation.trend {
    fmt.Fprintf(dest, "%s %s.%s %s\n", opts.redirect, k, kube_status[key].domain, k);
  }
}

func writeKubeInfo(key string) {
  src, err := os.Open(HOST);
  defer src.Close();
  if err != nil {
    GenericLogPrinter(opts.host, "ERR", fmt.Sprintf("[writeKubeInfo] failed to open host file (%s): %s", HOST, err.Error()), TOPIC);
  } else {
    dest, err := os.Create(HOST + ".tmp");
    defer dest.Close();
    if err != nil {
      GenericLogPrinter(opts.host, "ERR", fmt.Sprintf("[writeKubeInfo] failed to open tmp host file (%s): %s", HOST + ".tmp", err.Error()), TOPIC);
    } else {
      scanner := bufio.NewScanner(src);
      head := false;
      tail := false;
      save := false;
      target_hd := fmt.Sprintf("%s-%s", HEAD, key);
      target_tl := fmt.Sprintf("%s-%s", TAIL, key);
      for scanner.Scan() {
        raw := scanner.Text();
        current_line := strings.TrimSpace(raw);
        if current_line == target_hd {
          /* head of target found */
          head = true;
          fmt.Fprintf(dest, "%s\n", current_line);
        } else if head && !tail && current_line == target_tl {
          /* tail of target found, write out updated records */
          printHosts(dest, key);
          save = true;
          fmt.Fprintf(dest, "%s\n", current_line);
        } else if head && !tail {
          /* inside target block, skip all and do nothing */
        } else if head && tail {
          /* after target block, write out */
          fmt.Fprintf(dest, "%s\n", current_line);
        } else if !head {
          /* b4 target block, write out */
          fmt.Fprintf(dest, "%s\n", current_line);
        }
      }
      if !head {
        /* not deployed, write new block */
        fmt.Fprintf(dest, "%s\n", target_hd);
        printHosts(dest, key);
        fmt.Fprintf(dest, "%s\n", target_tl);
        save = true;
      }
      if save {
        timestamp := time.Now().UnixNano()/int64(1000000);
        //err := os.Rename(HOST, fmt.Sprintf("%s.%d", HOST, timestamp));
        err := exec.Command("cp", HOST, fmt.Sprintf("%s.%d", HOST, timestamp)).Run();
        if err != nil {
          GenericLogPrinter(opts.host, "ERR", fmt.Sprintf("[writeKubeInfo] failed to move host file (%s) to backup (%s): %s", HOST, fmt.Sprintf("%s.%d", HOST, timestamp), err.Error()), TOPIC);
        }
        //err = os.Rename(HOST + ".tmp", HOST);
        err = exec.Command("cp", HOST + ".tmp", HOST).Run();
        if err != nil {
          GenericLogPrinter(opts.host, "ERR", fmt.Sprintf("[writeKubeInfo] failed to move tmp host file (%s) to host file (%s): %s", HOST + ".tmp", HOST, err.Error()), TOPIC);
        }
        reloadDNSMasq();
      }
    }
  }
}

func getDNSMasqPID() (int, error) {
  fp, err := os.Open(DNSMasq);
  if err != nil {
    return -1, err;
  } else {
    scanner := bufio.NewScanner(fp);
    for scanner.Scan() {
      raw := scanner.Text();
      current_line := strings.TrimSpace(raw);
      pid, err := strconv.Atoi(current_line);
      if err != nil {
        return -1, err;
      } else {
        return pid, nil;
      }
    }
  }
  return -1, nil;
}

func reloadDNSMasq() {
  pid, err := getDNSMasqPID();
  if err != nil {
    GenericLogPrinter(opts.host, "ERR", fmt.Sprintf("[reloadDNSMasq] failed to get pid of dnsmasq: %s", err.Error()), TOPIC);
  } else {
    if pid == -1 {
      GenericLogPrinter(opts.host, "ERR", "[reloadDNSMasq] impossible condition while getting pid of dnsmasq", TOPIC);
    }
    pc, err := os.FindProcess(pid);
    if err != nil {
      GenericLogPrinter(opts.host, "ERR", fmt.Sprintf("[reloadDNSMasq] failed to find proccess via pid %d", pid), TOPIC);
    } else {
      err =  pc.Signal(syscall.SIGHUP);
      if err != nil {
        GenericLogPrinter(opts.host, "ERR", fmt.Sprintf("[reloadDNSMasq] failed to send signal to %d: %s", pid, err.Error()), TOPIC);
      }
    }
  }
}

func copyInfo(dst, src *KubeInfo) {
  flushInfo(dst);
  for k, v := range src.pod_ip {
    dst.pod_ip[k] = v;
  }
}

func monitorInfo(key string) {
  observation := kube_status[key].observation;
  current := kube_status[key].current;
  past := kube_status[key].past;
  for k, v := range observation.list {
    /* is anyone back online? */
    if _, ok := current.pod_ip[k]; !ok {
      /* nope, still offline */
      observation.list[k] = v + 1;
      /* record the trend */
      if cnt, ok := observation.trend[k]; !ok {
        /* no trend, it's not recently offline */
        observation.trend[k] = 1;
      } else {
        /* got trend, it's recently offline */
        observation.trend[k] = cnt + 1;
        GenericLogPrinter(opts.host, "WARN", fmt.Sprintf("[monitorInfo] %v is still offline for %d seconds...", k, observation.trend[k] * opts.poll), TOPIC);
      }
    } else {
      /* k is back, remove k from observation list? nope, decrease the cnt for now */
      //delete(observation.list, k);
      if v == 0 {
        delete(observation.list, k);
      } else {
        observation.list[k] = v - 1;
      }
      /* deal with trend record */
      if _, ok := observation.trend[k]; !ok {
        /* no trend, do nothing */
      } else {
        /* got trend, erase the trend record */
        delete(observation.trend, k);
      }
    }
  }
  for k, _ := range past.pod_ip {
    if _, ok := current.pod_ip[k]; !ok {
      cnt, ok := observation.list[k];
      if !ok {
        observation.list[k] = 1;
        if trend, ok := observation.trend[k]; !ok {
          observation.trend[k] = 1;
        } else {
          observation.trend[k] = trend + 1;
        }
      } else {
        observation.list[k] = cnt + 1;
        if trend, ok := observation.trend[k]; !ok {
          observation.trend[k] = 1;
        } else {
          observation.trend[k] = trend + 1;
        }
      }
    }
  }
}

func fillInc() {
  for _, v := range etcds {
    inc := kube_status[v].inc;
    current := kube_status[v].current;
    past := kube_status[v].past;
    observation := kube_status[v].observation;
    if kube_status[v].firstborn {
      /* not first time, flush incoming buffer for incoming info  */
      flushInfo(inc);
      fillInfo(inc);
      clean := isInfoEqual(inc, current);
      if !clean {
        flushInfo(current);
        fillInfo(current);
//        fmt.Fprintf(os.Stderr, "current b4: %v\n", kube_status[v].current);
        updateKubeStatus(v, true, clean, inc, current, past, observation);
//        fmt.Fprintf(os.Stderr, "current after: %v\n", kube_status[v].current);
        if opts.debug {
          GenericLogPrinter(opts.host, "DEBUG", fmt.Sprintf("[fillInc] firstborn: %v, clean: %v", kube_status[v].firstborn, kube_status[v].clean), TOPIC);
          GenericLogPrinter(opts.host, "DEBUG", fmt.Sprintf("[fillInc] inc: "), TOPIC);
          printKubeInfo(kube_status[v].inc);
          GenericLogPrinter(opts.host, "DEBUG", fmt.Sprintf("[fillInc] current:"), TOPIC);
          printKubeInfo(kube_status[v].current);
        }
      } else {
        if opts.debug {
          GenericLogPrinter(opts.host, "DEBUG", fmt.Sprintf("[fillInc] %s is clean? %v!", v, clean), TOPIC);
        }
      }
      monitorInfo(v);
      if !clean {
        writeKubeInfo(v);
      }
      copyInfo(past, current);
    } else {
      /* first time, make space for incoming info and current info */
      err := fillInfo(inc);
      if err != nil {
        GenericLogPrinter(opts.host, "ERR", fmt.Sprintf("[fillInc] failed to fill inc: %s", err.Error()), TOPIC);
      }
      err = fillInfo(current);
      if err != nil {
        GenericLogPrinter(opts.host, "ERR", fmt.Sprintf("[fillInc] failed to fill current: %s", err.Error()), TOPIC);
      }
      updateKubeStatus(v, true, true, inc, current, past, observation);
      if opts.debug {
        GenericLogPrinter(opts.host, "DEBUG", fmt.Sprintf("[fillInc] first! firstborn: %v, clean: %v", kube_status[v].firstborn, kube_status[v].clean), TOPIC);
      }
      writeKubeInfo(v);
      if opts.debug {
        GenericLogPrinter(opts.host, "DEBUG", fmt.Sprintf("[fillInc] inc: "), TOPIC);
        printKubeInfo(kube_status[v].inc);
        GenericLogPrinter(opts.host, "DEBUG", fmt.Sprintf("[fillInc] current:"), TOPIC);
        printKubeInfo(kube_status[v].current);
      }
    }
  }
}

func replace(s, old_str, new_str string, out []byte) {
  // Compute number of replacements.
  var m int;
  if m = strings.Count(s, old_str); m == 0 {
    return; // avoid allocation
  }

  // Apply replacements to buffer.
  w := 0;
  start := 0;
  for i := 0; i < m; i++ {
    j := start;
    if len(old_str) == 0 {
      if i > 0 {
        _, wid := utf8.DecodeRuneInString(s[start:]);
        j += wid;
      }
    } else {
      j += strings.Index(s[start:], old_str);
    }
    w += copy(out[w:], s[start:j]);
    w += copy(out[w:], new_str);
    start = j + len(old_str);
  }
  w += copy(out[w:], s[start:]);
}

func checkIPFmt(ip string) (error) {
  ips := strings.Split(ip, ".");
  if len(ips) != 4 {
    return errors.New(fmt.Sprintf("[checkIPFmt] Invalid IP format, expected <num>.<num>.<num>.<num>, got %s", ip));
  } else {
    for _, v := range ips {
      num, err := strconv.Atoi(v);
      if err != nil {
        return errors.New(fmt.Sprintf("[checkIPFmt] Failed to parse IP: %s", err.Error()));
      } else {
        if num > 255 || num < 0 {
          return errors.New(fmt.Sprintf("[checkIPFmt] Invalid IP: %s", ip));
        }
      }
    }
    return nil;
  }
  return errors.New("(IP) Unknown situation");
}

func checkETCDIPFmt(ip string) (error) {
  tokens := strings.Split(ip, ":");
  if len(tokens) != 2 {
    GenericLogPrinter(opts.host, "ERR", fmt.Sprintf("[checkETCDIPFmt] (ETCD) IP format error, expected <IP>:<Port>, got %s", ip), TOPIC);
    return errors.New(fmt.Sprintf("(ETCD) IP format error, expected <IP>:<Port>, got %s", ip));
  } else {
    port, err := strconv.Atoi(tokens[1]);
    if err != nil {
      return errors.New(fmt.Sprintf("(ETCD) Failed to parse port of ETCD %s: %s", ip, err.Error()));
    } else {
      if port > 65535 || port < 1000 {
        return errors.New(fmt.Sprintf("(ETCD) Invalid port range, expected 1000~65535, got %d", port));
      } else {
        return checkIPFmt(tokens[0]);
      }
      return errors.New("(ETCD) Unknown situation");
    }
    return errors.New("(ETCD) Unknown situation");
  }
  return errors.New("(ETCD) Unknown situation");
}

func checkETCD() (error) {
  var err error;
  for i := 0; i < len(etcds); i++ {
    err = checkETCDIPFmt(etcds[i]);
    if err != nil {
      return err;
    }
  }
  return nil;
}

