package main;

import "flag";
import "fmt";
//import "net/http";
import "os";

import logutil "tyd_util/log_util";
import util "tyd_util";

func main() {
  flag.Parse();
  err := logutil.InitKFKProducer(opts.pool, opts.brokers);
  if err != nil {
    fmt.Fprintf(os.Stderr, "failed to initiate log producers: %s", err.Error());
    os.Exit(1);
  }
  err = initHermes();
  if err != nil {
    util.GenericLogPrinter(opts.host, "ERR", fmt.Sprintf("failed to initiate: %s", err.Error()), TOPIC);
    fmt.Fprintf(os.Stderr, "failed to initiate: %s\n", err.Error());
    os.Exit(1);
  }
  /*
  for _, v := range etcds {
    buffer.requests[v], err = http.NewRequest("GET", fmt.Sprintf("http://%s%s", v, URI), nil);
    if err != nil {
      util.GenericLogPrinter(opts.host, "ERR", fmt.Sprintf("failed to create request: %s", err.Error()), TOPIC);
      fmt.Fprintf(os.Stderr, "failed to create request: %s\n", err.Error());
      os.Exit(1);
    }
  }
  */
  util.GenericLogPrinter(opts.host, "INFO", fmt.Sprintf("domain name updater activated, target etcds are %s", opts.etcd), TOPIC);
  Caduceus();
}
