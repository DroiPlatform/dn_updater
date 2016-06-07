package main;

import "flag";
import "fmt";
import "net/http";
import "strings";

import logutil "tyd_util/log_util";

const TOPIC = "dns_updater"

const SIZE_BUF = 65536;

const URI = "/v2/keys/registry/services/endpoints/tyd";

/* flags */
type Options struct {
  build bool;
  debug bool;
  poll int;
  pool int;
  brokers string;
  etcd string;
  host string;
}

/* structure for kube status */
type KubeStatus struct {
  clean bool;
  inc *KubeInfo;
  current *KubeInfo;
}

/* structure for kube info */
type KubeInfo struct {
  /* mapping from service name to pod ip */
  pod_ip map[string]string;
}

/* buffer */
type TmpBuffer struct {
  client *http.Client;
  requests map[string]*http.Request;
  raw_json []byte;  // raw json from etcd, allocate 64KB = 65536
  escape_json []byte;  // json for espaced, allocate 64KB = 65536
}

var kube_status map[string]KubeStatus;
var opts Options;
var etcds []string;
var buffer TmpBuffer;

func init() {
  flag.BoolVar(&opts.build, "build", false, "print golang build version");
  flag.IntVar(&opts.poll, "poll", 5, "etcd polling interval");
  flag.IntVar(&opts.pool, "pooll", 3, "pool size for log producers");
  flag.StringVar(&opts.brokers, "brokers", "10.128.112.78:9092", "ip:port for kafka brokers, seperated by comma");
  flag.StringVar(&opts.etcd, "etcd", "", "ip:port for etcds, seperated by comma");
}

func initHermes() (error) {
  err := initLog();
  if err != nil {
    return err;
  }
  err = initData();
  return err;
}

func initLog() (error) {
  return logutil.InitKFKProducer(opts.pool, opts.brokers);
}

func initData() (error){
  /* initiate tmp buffers and shared resources */
  buffer.client = &http.Client{};
  buffer.requests = make(map[string]*http.Request);
  buffer.raw_json = make([]byte, SIZE_BUF);
  buffer.escape_json = make([]byte, SIZE_BUF);

  /* initiate etcd related resources */
  etcds = strings.Split(opts.etcd, ",");
  kube_status = make(map[string]KubeStatus);
  var err error
  for _, v := range etcds {
    kube_status[v] = KubeStatus{clean: false, inc: &KubeInfo{make(map[string]string)}, current: &KubeInfo{make(map[string]string)}};
    buffer.requests[v], err = http.NewRequest("GET", fmt.Sprintf("http://%s%s", v, URI), nil);
    if err != nil {
      return err;
    }
  }
  return nil;
}

