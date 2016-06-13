package main;

import "errors";
import "flag";
import "fmt";
import "net/http";
import "strings";

import util "tyd_util";
import logutil "tyd_util/log_util";

const TOPIC = "dn_updater"

const HEAD = "# head of dn_updater";
const TAIL = "# tail of dn_updater";
const HOST = "/etc/hosts";
const DNSMasq = "/var/run/dnsmasq/dnsmasq.pid";

const SIZE_BUF = 65536;

const URI = "/v2/keys/registry/services/endpoints/tyd";

/* flags */
type Options struct {
  build bool;
  debug bool;
  poll int;
  pool int;
  brokers string;
  domain string;
  etcd string;
  host string;
}

/* structure for kube status */
type KubeStatus struct {
  clean bool;
  firstborn bool;
  domain string;
  current *KubeInfo;
  inc *KubeInfo;
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
  hosts map[string]string;
}

var kube_status map[string]KubeStatus;
var opts Options;
var etcds []string;
var domains []string;
var buffer TmpBuffer;

func init() {
  flag.BoolVar(&opts.build, "build", false, "print golang build version");
  flag.BoolVar(&opts.debug, "debug", false, "print debug message");
  flag.IntVar(&opts.poll, "poll", 5, "etcd polling interval");
  flag.IntVar(&opts.pool, "pooll", 3, "pool size for log producers");
  flag.StringVar(&opts.brokers, "brokers", "10.128.112.186:9092", "ip:port for kafka brokers, seperated by comma");
  flag.StringVar(&opts.domain, "domain", "", "domain for etcds, seperated by comma");
  flag.StringVar(&opts.etcd, "etcd", "", "ip:port for etcds, seperated by comma");
  flag.StringVar(&opts.host, "host", "", "host identifier of this machine, usually IP");
}

func initHermes() (error) {
  err := initLog();
  if err != nil {
    return err;
  }
  err = initData();
  if opts.debug {
    util.GenericLogPrinter(opts.host, "DEBUG", fmt.Sprintf("result from initData: %v", err), TOPIC);
  }
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
    kube_status[v] = KubeStatus{firstborn: false, clean: false, domain: domains[k], inc: &KubeInfo{make(map[string]string)}, current: &KubeInfo{make(map[string]string)}};
    buffer.requests[v], err = http.NewRequest("GET", fmt.Sprintf("http://%s%s", v, URI), nil);
    if err != nil {
      return err;
    }
  }
  return nil;
}

