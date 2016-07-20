package log_util;

import "fmt";
import "strings";
import "sync";

import producer "tyd_util/kfk_producer";
//import "github.com/Shopify/sarama";

type LogContent struct {
  Version int;
  Time string;  //micro-second
  Aid string;
  Mod string;
  Lvl string;
  Msg string;
  Opt string;
}

var wagon producer.Wagon;

var kp *(producer.ConnectionPoolWrapper);

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
  out = fmt.Sprintf("%s \"M\": \"%s\"", out, lc.Msg);
  out = fmt.Sprintf("%s}", out);
  return out;
}

func InitKFKProducer(size int, brokers string) (error) {
  kp = &producer.ConnectionPoolWrapper{};
  return kp.InitPool(size, brokers);
}

/*
func InitProducer(size int, brokers string) (error) {
  err := producer.InitWagon(size, brokers);
  return err;
}
//*/
//*
func KFKWriter(lc *LogContent, topic string) {
  msg := producer.Message {
    Msg: logWrapper(lc),
    Topic: topic,
  };
  gregor(msg);
  //go gregor(KFKPool, logWrapper(lc), topic);
}
//*/

/* Log Writer for IPList */
/*
func IPListWriter(lc *LogContent) {
  lc.Aid = "";
  lc.Mod = "ip_list";
  msg := producer.Message {
    Msg: logWrapper(lc),
    Topic: lc.Mod,
  }
  p := kp.GetConnection();
  producer.WriteMsg(p, msg);
  kp.FreeConnection(p);
}
//*/

/* Log Writer for Pooled General Purpose */
func GeneralLogWriter(lc *LogContent, mod string) {
  lc.Aid = "";
  topic := mod;
  fmt.Printf("topic: %s, module: %s\n", mod, lc.Mod);
  if lc.Mod == "" {
    lc.Mod = mod;
  }
  msg := producer.Message {
    Msg: logWrapper(lc),
    Topic: topic,
  }
  p := kp.GetConnection();
  producer.WriteMsg(p, msg);
  kp.FreeConnection(p);
}

/* Log Writer for PodFinder */
/*
func PFWriter(lc *LogContent) {
  lc.Aid = "";
  lc.Mod = "pod_finder";
  msg := producer.Message {
    Msg: logWrapper(lc),
    Topic: lc.Mod,
  }
  p := kp.GetConnection();
  producer.WriteMsg(p, msg);
  kp.FreeConnection(p);
}
//*/

/* Log Writer for KeyServerSynchronizer */
/*
func UpdaterWriter(lc *LogContent) {
  lc.Aid = "";
  lc.Mod = "keyserver_sync";
  msg := producer.Message {
    Msg: logWrapper(lc),
    Topic: lc.Mod,
  }
  p := kp.GetConnection();
  producer.WriteMsg(p, msg);
  kp.FreeConnection(p);
}
//*/

///*
func KSKFKWriter(lc *LogContent) {
  lc.Aid = "";
  lc.Mod = "keyserver";
  msg := producer.Message {
    Msg: logWrapper(lc),
    Topic: lc.Mod,
  };
  gregor(msg);
  go KFKFlush();
}
//*/
/*
func WSKFKWriter(lc *LogContent) {
  msg := producer.Message {
    Msg: logWrapper(lc),
    Topic: "webscraper",
  };
  gregor(msg);
  go KFKFlush();
}
*/

func WSPoolWriter(lc *LogContent) {
  p := kp.GetConnection();
  producer.WriteMsg(p, producer.Message {Msg: logWrapper(lc), Topic: "webscraper"});
  kp.FreeConnection(p);
}

func KFKFlush() {
  producer.UnloadItem();
}

func KFKPush(wg *sync.WaitGroup) {
  producer.UnloadItem();
  //wg.Done();
}

///*
func OneTimeKFKWriter(lc *LogContent, topic string, kafkas string) {
  kfk_list := strings.Split(kafkas, ",");
  samsa(logWrapper(lc), topic, kfk_list);
}
//*/

func gregor(msg producer.Message) {
  producer.LoadItem(msg);
}

/*
func gregor(KFKPool *(producer.ConnectionPoolWrapper), msg string, topic string) {
  p := KFKPool.GetConnection();
  //fmt.Printf("p: %v\n", p);
  producer.Metamorphosis(p, msg, topic);
  KFKPool.FreeConnection(p);
}
//*/

func samsa(msg string, topic string, kfks []string) {
  p, err := producer.NewAsyncProducer(kfks);
  if err != nil {
    ////fmt.Printf("failed to create producer: %s\n", err.Error());
  } else {
    ////fmt.Printf("p: %v\n", p);
    producer.Metamorphosis(p, msg, topic);
    p.Close();
  }
}

/* Unpooled approach, resource pool is implemented via index */
type LogProducer struct {
  Worker *producer.KFKProducer;
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
  producer.WriteMsg(*p.Worker.Worker, msg);
}

