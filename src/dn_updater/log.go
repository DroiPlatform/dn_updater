package main;

import "fmt";
import "log";
import "time";

import util "tyd_util";
import logutil "tyd_util/log_util";

func getTime() (string) {
  return fmt.Sprintf("%s.%06d", time.Now().Format("2006-01-02 15:04:05"), time.Now().Nanosecond()/1000);
}

func GenericLogPrinter(host, lvl, msg, topic string) {
  if Kafka {
    util.GenericLogPrinter(host, lvl, msg, topic);
  } else {
    lc := logutil.LogContent {
      Version: 1,
      Time: getTime(),
      Mod: mod,
      Lvl: lvl,
      Msg: msg,
    }
    if host != "" {
      lc.Opt = fmt.Sprintf("{\"Nid\": \"%s\"}", host);
    }
    log.Printf("%s\n", logutil.LogWrapper(&lc));
  }
}
