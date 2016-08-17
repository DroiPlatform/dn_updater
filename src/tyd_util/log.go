package tyd_util;

import "fmt";
import "time";

import logutil "tyd_util/log_util";

type UID struct {
  Type uint8;
  Upper uint64;
  Lower uint64;
}

/* get formatted time string */
func getTime() (string) {
  return fmt.Sprintf("%s.%06d", time.Now().Format("2006-01-02 15:04:05"), time.Now().Nanosecond()/1000);
}

/* high level log printer for general purpose */
func GenericLogPrinter(host string, lvl string, msg string, topic string) {
  lc := logutil.LogContent {
    Version: 1,
    Time: getTime(),
    Lvl: lvl,
    Msg: msg,
  }
  if host != "" {
    lc.Opt = fmt.Sprintf("{\"Nid\": \"%s\"}", host);
  }
  logutil.GeneralLogWriter(&lc, topic);
}

/* high level log printer for IPList */
func IPListLogPrinter(host string, lvl string, msg string) {
  lc := logutil.LogContent {
    Version: 1,
    Time: getTime(),
    Lvl: lvl,
    Msg: msg,
  }
  if host != "" {
    lc.Opt = fmt.Sprintf("{\"Nid\": \"%s\"}", host);
  }
  logutil.GeneralLogWriter(&lc, "ip_list");
}

/* High Level Log Printer w/ UID */
func UIDLogPrinter(uid UID, host string, lvl string, msg string) {
  lc := logutil.LogContent {
    Version: 1,
    //Time: time.Now().UnixNano() / int64(1000),
    Time: getTime(),
    Lvl: lvl,
    Msg: msg,
  }
  if uid.Upper == uint64(0) && uid.Lower == uint64(0) {
    if host != "" {
      lc.Opt = fmt.Sprintf("{\"Nid\": \"%s\"}", host);
    }
  } else {
    if host != "" {
      lc.Opt = fmt.Sprintf("{\"UID\": \"%d:%d:%d\", \"Nid\": \"%s\"}", uid.Upper, uid.Lower, uid.Type, host);
    } else {
      lc.Opt = fmt.Sprintf("{\"UID\": \"%d:%d:%d\"}", uid.Upper, uid.Lower, uid.Type);
    }
  }
  logutil.GeneralLogWriter(&lc, "keyserver");
}

