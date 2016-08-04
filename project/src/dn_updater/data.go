package main;

import "errors";
import "flag";
import "fmt";
import "net/http";
import "os";
import "strings";

const TOPIC = "dn_updater"

const HEAD = "# head of dn_updater";
const TAIL = "# tail of dn_updater";
const HOST = "/etc/hosts";
const HEALTH = "/opt/healthz";
const DNSMasq = "/var/run/dnsmasq/dnsmasq.pid";

//const SIZE_BUF = 65536;
const SIZE_BUF = 1048576;

const URI = "/v2/keys/registry/services/endpoints/tyd";

/* flags */
type Options struct {
  build bool;
  debug bool;
  test bool;
  json int;
  poll int;
  pool int;
  suicide int;
  throttle int;
  brokers string;
  domain string;
  etcd string;
  host string;
  log string;
}

/* structure for kube status */
type KubeStatus struct {
  clean bool;
  firstborn bool;
  domain string;
  past *KubeInfo;
  current *KubeInfo;
  inc *KubeInfo;
  observation *Observation;
}

/* structure for kube info */
type KubeInfo struct {
  /* mapping from service name to pod ip */
  pod_ip map[string]string;
}

/* observation list */
type Observation struct {
  list map[string]int;
  trend map[string]int;
}

/* buffer */
type TmpBuffer struct {
  client *http.Client;
  requests map[string]*http.Request;
  raw_json []byte;  // raw json from etcd, allocate 64KB = 65536
  escape_json []byte;  // json for espaced, allocate 64KB = 65536
  hosts map[string]string;
}

var kube_status map[string]KubeStatus;
var opts Options;
var etcds []string;
var domains []string;
var buffer TmpBuffer;
var rnd_cnt uint8;
var host string;
var mod string;
var Kafka bool;

func init() {
  flag.BoolVar(&opts.build, "build", false, "print golang build version");
  flag.BoolVar(&opts.debug, "debug", false, "print debug message");
  flag.BoolVar(&opts.test, "test", false, "trigger test mode");
  flag.IntVar(&opts.json, "json", 134217728, "buffer size of json from etcd (byte)");
  flag.IntVar(&opts.poll, "poll", 5, "etcd polling interval");
  flag.IntVar(&opts.pool, "pool", 3, "pool size for log producers");
  flag.IntVar(&opts.suicide, "suicide", 0, "suicide after <suicide> hours, 0 for never");
  flag.IntVar(&opts.throttle, "throttle", 100, "throttle for debug msg frequency, 1-255");
  flag.StringVar(&opts.brokers, "brokers", "10.128.112.186:9092", "ip:port for kafka brokers, seperated by comma");
  flag.StringVar(&opts.domain, "domain", "", "domain for etcds, seperated by comma");
  flag.StringVar(&opts.etcd, "etcd", "", "ip:port for etcds, seperated by comma");
  flag.StringVar(&opts.host, "host", "", "host identifier of this machine, usually IP");
  flag.StringVar(&opts.log, "log", "dn_updater.log", "path to local log file.");
}

func initHermes() (error) {
  err := initData();
  if opts.debug {
    GenericLogPrinter(opts.host, "DEBUG", fmt.Sprintf("result from initData: %v", err), TOPIC);
  }
  return err;
}

func initData() (error){
  /* initial global variables */
  rnd_cnt = uint8(0);
  mod, _ = os.Hostname();
  /* get pod ip addr */
  host = getLocalIP();
  fmt.Printf("[initData] host: %s\n", host);
  /* initiate tmp buffers and shared resources */
  buffer.client = &http.Client{};
  buffer.requests = make(map[string]*http.Request);
//  buffer.raw_json = make([]byte, opts.json);
  buffer.escape_json = make([]byte, opts.json);
  buffer.hosts = make(map[string]string);

  /* initiate etcd related resources */
  etcds = strings.Split(opts.etcd, ",");
  err := checkETCD();
  if err != nil {
    return err;
  }
  domains = strings.Split(opts.domain, ",");
  if len(etcds) != len(domains) {
    return errors.New(fmt.Sprintf("# of etcds and # of domains not match: %s %s", opts.etcd, opts.domain));
  }
  kube_status = make(map[string]KubeStatus);
  for k, v := range etcds {
    kube_status[v] = KubeStatus{firstborn: false, clean: false, domain: domains[k], inc: &KubeInfo{make(map[string]string)}, past: &KubeInfo{make(map[string]string)}, current: &KubeInfo{make(map[string]string)}, observation: &Observation{list: make(map[string]int), trend: make(map[string]int)}};
//    fmt.Fprintf(os.Stderr, "observation of %s: %v\n", v, kube_status[v].observation);
    buffer.requests[v], err = http.NewRequest("GET", fmt.Sprintf("http://%s%s", v, URI), nil);
    if err != nil {
      return err;
    }
  }
  return nil;
}

