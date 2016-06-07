package main;

import "flag";
import "fmt";
import "net/http";
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
    util.GenericLogPrinter(opts.host, "ERR", fmt.Sprintf("failed to initiate: %s\n"), TOPIC);
    os.Exit(1);
  }
  for _, v := range etcds {
    buffer.requests[v], err = http.NewRequest("GET", fmt.Sprintf("http://%s%s", v, URI), nil);
    if err != nil {
      util.GenericLogPrinter(opts.host, "ERR", fmt.Sprintf("failed to create request: %s\n", err.Error()), TOPIC);
      os.Exit(1);
    }
  }
  Caduceus();
}
