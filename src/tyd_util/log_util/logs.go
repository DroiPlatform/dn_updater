package log_util;

import "fmt";
import "os";

import producer "tyd_util/kfk_producer";

type LogContent struct {
  Version int;
  Time string;  //micro-second
  Aid string;
  Mod string;
  Lvl string;
  Msg string;
  Opt string;
  /* required fields for webscraper */
  RequestSpent string;    //transfer to int
  RequestLength string;   //transfer to int
  ResponseLength string;  //transfer to int
  Identifier string;      //transfer to int
  AccessTime string;      //transfer to int
}

var sp *(producer.SiestaPoolWrapper);

/* transform log structure to log JSON string, public function */
func LogWrapper(lc *LogContent) (string) {
  out := fmt.Sprintf("{");
  out = fmt.Sprintf("%s \"V\": \"%d\",", out, lc.Version);
  out = fmt.Sprintf("%s \"T\": \"%s\",", out, lc.Time);
  out = fmt.Sprintf("%s \"L\": \"%s\",", out, lc.Lvl);
  if lc.Aid != "" {
    out = fmt.Sprintf("%s \"Aid\": \"%s\",", out, lc.Aid);
  }
  out = fmt.Sprintf("%s \"Pd\": \"%s\",", out, lc.Mod);
  if lc.Opt != "" {
    out = fmt.Sprintf("%s \"Op\": %s,", out, lc.Opt);
  }
  /* webscraper */
  if lc.RequestSpent != "" {
    out = fmt.Sprintf("%s \"RequestSpent\": %s,", out, lc.RequestSpent);
  }
  if lc.RequestLength != "" {
    out = fmt.Sprintf("%s \"RequestLength\": %s,", out, lc.RequestLength);
  }
  if lc.ResponseLength != "" {
    out = fmt.Sprintf("%s \"ResponseLength\": %s,", out, lc.ResponseLength);
  }
  if lc.Identifier != "" {
    out = fmt.Sprintf("%s \"Identifier\": %s,", out, lc.Identifier);
  }
  if lc.AccessTime != "" {
    out = fmt.Sprintf("%s \"AccessTime\": \"%s\",", out, lc.AccessTime);
  }
  out = fmt.Sprintf("%s \"M\": \"%s\"", out, lc.Msg);
  out = fmt.Sprintf("%s}", out);
  return out;
}

/* transform log structure to log JSON string, private function, deprecated */
func logWrapper(lc *LogContent) (string) {
  out := fmt.Sprintf("{");
  out = fmt.Sprintf("%s \"V\": \"%d\",", out, lc.Version);
  out = fmt.Sprintf("%s \"T\": \"%s\",", out, lc.Time);
  out = fmt.Sprintf("%s \"L\": \"%s\",", out, lc.Lvl);
  if lc.Aid != "" {
    out = fmt.Sprintf("%s \"Aid\": \"%s\",", out, lc.Aid);
  }
  out = fmt.Sprintf("%s \"Pd\": \"%s\",", out, lc.Mod);
  if lc.Opt != "" {
    out = fmt.Sprintf("%s \"Op\": %s,", out, lc.Opt);
  }
  /* webscraper */
  if lc.RequestSpent != "" {
    out = fmt.Sprintf("%s \"RequestSpent\": %s,", out, lc.RequestSpent);
  }
  if lc.RequestLength != "" {
    out = fmt.Sprintf("%s \"RequestLength\": %s,", out, lc.RequestLength);
  }
  if lc.ResponseLength != "" {
    out = fmt.Sprintf("%s \"ResponseLength\": %s,", out, lc.ResponseLength);
  }
  if lc.Identifier != "" {
    out = fmt.Sprintf("%s \"Identifier\": %s,", out, lc.Identifier);
  }
  if lc.AccessTime != "" {
    out = fmt.Sprintf("%s \"AccessTime\": \"%s\",", out, lc.AccessTime);
  }
  out = fmt.Sprintf("%s \"M\": \"%s\"", out, lc.Msg);
  out = fmt.Sprintf("%s}", out);
  return out;
}

/* initilize producer pool, producer pool is managed by this pkg */
func InitKFKProducer(size int, brokers string) (error) {
  sp = &producer.SiestaPoolWrapper{};
  return sp.InitSiestaPool(size, brokers);
}

/* send log structure to producer */
func GeneralLogWriter(lc *LogContent, topic string) {
  lc.Aid = "";
  mod, err := os.Hostname();
  if err != nil {
    lc.Mod = "";
  } else {
    lc.Mod = mod;
  }
  msg := producer.Message {
    Msg: logWrapper(lc),
    Topic: topic,
  }
  sp.WriteSiesta(msg);
}

/* Unpooled approach, resource pool is implemented via index */
type LogProducer struct {
  Worker *producer.Siesta;
}

/* Producer initialization */
func InitProducer(kfk_list []string) (*LogProducer, error) {
  worker, err := producer.InitProducer(kfk_list);
  return &LogProducer {Worker: worker}, err;
}

/* Log Writer for Un-pooled General Purpose (pool is implemented by index pool) */
func UnpooledGeneralLogWriter(lc *LogContent, mod string, p *LogProducer) {
  lc.Aid = "";
  lc.Mod = mod;
  msg := producer.Message {
    Msg: logWrapper(lc),
    Topic: lc.Mod,
  }
  producer.WriteUnpooledSiesta(p.Worker, msg);
}

