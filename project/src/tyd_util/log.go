package tyd_util;

import "fmt";
import "os";
import "time";

import logutil "tyd_util/log_util";

type UID struct {
  Type uint8;
  Upper uint64;
  Lower uint64;
}

/* print log label */
func PrintPrefix(uid UID, lvl string) {
  fmt.Printf("[S:KeyServer] ");
  //PrintTime();
  if uid.Upper == uint64(0) && uid.Lower == uint64(0) {
    fmt.Printf("[L:%s] MSG: ", lvl);
  } else {
    fmt.Printf("[U:%d,%d,%d] [L:%s] MSG: ", uid.Upper, uid.Lower, uid.Type, lvl);
  }
}

/* kafka log printer */
func KSKFKLogPrinter(uid UID, lvl string, msg string) {
  lc := logutil.LogContent {
    Version: 1,
    //Time: time.Now().UnixNano() / int64(1000),
    Time: fmt.Sprintf("%s.%06d", time.Now().Format("2006-01-02 15:04:05"), time.Now().Nanosecond()/1000),
    Lvl: lvl,
    Msg: msg,
  }
  if uid.Upper == uint64(0) && uid.Lower == uint64(0) {
    lc.Opt = "";
  } else {
    lc.Opt = fmt.Sprintf("{\"UID\": \"%d:%d:%d\"}", uid.Upper, uid.Lower, uid.Type);
  }
  logutil.KSKFKWriter(&lc);
}

func getTime() (string) {
    return fmt.Sprintf("%s.%06d", time.Now().Format("2006-01-02 15:04:05"), time.Now().Nanosecond()/1000);
}

/* high level log printer for general purpose */
func GenericLogPrinter(host string, lvl string, msg string, topic string) {
  hostname, err := os.Hostname();
  if err != nil {
    fmt.Fprintf(os.Stderr, "failed to get hostname: %s\n", err.Error);
    return;
  }
  lc := logutil.LogContent {
    Version: 1,
    Time: getTime(),
    Lvl: lvl,
    Mod: hostname,
    Msg: msg,
  }
  fmt.Printf("hostname: %s\n", lc.Mod);
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
  lc.Opt = fmt.Sprintf("{\"Nid\": \"%s\"}", host);
  logutil.GeneralLogWriter(&lc, "ip_list");
}

/* high level log printer for OAuthFwder */
func OALogPrinter(host string, lvl string, msg string) {
  lc := logutil.LogContent {
    Version: 1,
    Time: getTime(),
    Lvl: lvl,
    Msg: msg,
  }
  lc.Opt = fmt.Sprintf("{\"Nid\": \"%s\"}", host);
  logutil.GeneralLogWriter(&lc, "oauth_fwder");
}

/* high level log printer for PodFinder */
func PFLogPrinter(host string, lvl string, msg string) {
  lc := logutil.LogContent {
    Version: 1,
    Time: getTime(),
    Lvl: lvl,
    Msg: msg,
  }
  lc.Opt = fmt.Sprintf("{\"Nid\": \"%s\"}", host);
  logutil.GeneralLogWriter(&lc, "pod_finder");
}


/* high level log printer for KeyServerSynchronizer */
func UpdaterLogPrinter(host string, lvl string, msg string) {
  lc := logutil.LogContent {
    Version: 1,
    Time: getTime(),
    Lvl: lvl,
    Msg: msg,
  }
  lc.Opt = fmt.Sprintf("{\"Nid\": \"%s\"}", host);
  logutil.GeneralLogWriter(&lc, "keyserver_sync");
}

func KSINDKFKLogPrinter(uid UID, host string, lvl string, msg string) {
  lc := logutil.LogContent {
    Version: 1,
    //Time: time.Now().UnixNano() / int64(1000),
    Time: getTime(),
    Lvl: lvl,
    Msg: msg,
  }
  if uid.Upper == uint64(0) && uid.Lower == uint64(0) {
    lc.Opt = fmt.Sprintf("{\"Nid\": \"%s\"}", host);
  } else {
    lc.Opt = fmt.Sprintf("{\"UID\": \"%d:%d:%d\", \"Nid\": \"%s\"}", uid.Upper, uid.Lower, uid.Type, host);
  }
  logutil.GeneralLogWriter(&lc, "keyserver");
}

