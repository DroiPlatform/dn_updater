package main;

import "flag";
import "fmt";
//import "net/http";
import "log";
import "os";
import "time";

import logutil "tyd_util/log_util";

func main() {
  flag.Parse();
  fp, err := os.OpenFile(opts.log, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666);
  if err != nil {
    fmt.Fprintf(os.Stderr, "[restfulInit] failed to open log file %s: %s\n", opts.log, err.Error());
    os.Exit(1);
  }
  defer fp.Close();
  log.SetOutput(fp);
  Kafka = true;
  err = logutil.InitKFKProducer(opts.pool, opts.brokers);
  if err != nil {
    Kafka = false;
    GenericLogPrinter(opts.host, "ERR", "[main] failed to initiate log producers: %s", err.Error());
  }
  err = initHermes();
  if err != nil {
    GenericLogPrinter(opts.host, "ERR", fmt.Sprintf("[main] failed to initiate: %s", err.Error()), TOPIC);
    time.Sleep(time.Duration(1) * time.Second);
    os.Exit(1);
  }
  GenericLogPrinter(opts.host, "INFO", fmt.Sprintf("[main] domain name updater activated, target etcds are %s", opts.etcd), TOPIC);
  Caduceus();
}
