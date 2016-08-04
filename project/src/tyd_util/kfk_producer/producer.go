package kfk_producer;

import "errors";
import "fmt";
import "strings";
import "time";

import "github.com/Shopify/sarama";
import "github.com/elodina/siesta";
import producer "github.com/elodina/siesta-producer";

import pool "tyd_util/resource_pool";

type ConnectionPoolWrapper struct {
  size int;
  conn chan sarama.AsyncProducer;
}

type Siesta struct {
  p *producer.KafkaProducer;
  q chan *producer.ProducerRecord;
}

type SiestaPoolWrapper struct {
  size int;
  index *pool.IndexPoolWrapper;
  s []*Siesta;
}

type Wagon struct {
  wagon sarama.AsyncProducer;
  items chan Message;
}

var wagon sarama.AsyncProducer;
var items chan Message;

type Message struct {
  Msg string;
  Topic string;
}

type Producer struct {
  sarama.AsyncProducer;
}

/* default value for siesta */
const MAX_CONN_PER_BROKER = 1;
const STD_RETRY = 3;
const QueueLength = 100;  // need clarify on unit
const EnqueueTimeout = 1 * time.Second;
/* default value for siesta-producer, need further clarify */
const BatchSize = 16384;
const MaxRequests = 10;
const SendRoutines = 10;
const ReceiveRoutines = 10;
const ReadTimeout = 5 * time.Second;
const WriteTimeout = 5 * time.Second;
const RequireAcks = 1;
const AckTimeoutMs = 30000;
const Linger = 1 * time.Second;
const RetryBackoff = 100 * time.Millisecond;

/* errors */
var ERR_TIMEOUT error;
var ERR_UNKNOWN error;

/* KFK Connection Pool */
/* init */
func (p *ConnectionPoolWrapper) InitPool(size int, kafkas string) (error) {
  p.conn = make(chan sarama.AsyncProducer, size);
  for i := 0; i < size; i++ {
    conn, err := initProducer(kafkas);
    if err != nil {
      return err;
    } else {
      //fmt.Printf("p[%d]: %v\n", i, conn);
      p.conn <- conn;
    }
  }
  p.size = size;
  return nil;
}

func (p *SiestaPoolWrapper) InitSiestaPool(size int, kafkas string) (error) {
  p.size = size;
  p.index = &pool.IndexPoolWrapper{};
  p.index.IndexInitPool(size);
  p.s = make([]*Siesta, size);
  ERR_TIMEOUT = errors.New("Enqueue Timeout!");
  ERR_UNKNOWN = errors.New("Unknown Situation!");
  for i := 0; i < size; i++ {
    conn, err := initSiestaProducer(kafkas);
    if err != nil {
      return err;
    } else {
      //fmt.Printf("p[%d]: %v\n", i, conn);
      p.s[i] = &Siesta{p: conn, q: make(chan *producer.ProducerRecord, QueueLength)};
      go p.s[i].siestaSend();
    }
  }
  return nil;
}

func (s *Siesta) siestaSend() {
  var msg *producer.ProducerRecord;
  for msg = range s.q {
    s.p.Send(msg);
  }
}

func (p *SiestaPoolWrapper)WriteSiesta(msg Message) (error) {
  index := p.index.GetIndex();
  select {
  case p.s[index].q <- &producer.ProducerRecord{Topic: msg.Topic, Value: []byte(msg.Msg)}:
    p.index.FreeIndex(index);
    return nil;
  case <-time.After(EnqueueTimeout):
    p.index.FreeIndex(index);
    return ERR_TIMEOUT;
  }
}

func (p *ConnectionPoolWrapper) GetConnection() (sarama.AsyncProducer) {
  return <- p.conn;
}

func (p *ConnectionPoolWrapper) FreeConnection(conn sarama.AsyncProducer) {
  p.conn <- conn;
}

func newAsyncProducer(kfk_list []string) (sarama.AsyncProducer, error) {
  conf := sarama.NewConfig();
  conf.Producer.RequiredAcks = sarama.WaitForLocal;
  conf.Producer.Flush.Frequency = 100 * time.Millisecond;
  conf.Producer.Return.Errors = false;
  return sarama.NewAsyncProducer(kfk_list, conf);
}

func NewAsyncProducer(kfk_list []string) (sarama.AsyncProducer, error) {
  conf := sarama.NewConfig();
  conf.Producer.RequiredAcks = sarama.WaitForLocal;
  conf.Producer.Flush.Frequency = 100 * time.Millisecond;
  conf.Producer.Return.Errors = false;
  return sarama.NewAsyncProducer(kfk_list, conf);
}

func InitWagon(size int, kafkas string) (error) {
  kfk_list := strings.Split(kafkas, ",");
  for i := 0; i < len(kfk_list); i++ {
    if strings.Trim(kfk_list[i], " \t") == "" {
      return errors.New(fmt.Sprintf("invalid kafka brokers: %s\n", kafkas));
    }
  }
  var err error;
  wagon, err = newAsyncProducer(kfk_list);
  if err != nil {
    return err;
  } else {
    items = make(chan Message, size);
    return nil;
  }
}

func LoadItem(msg Message) {
  items <- msg;
}

func getMessage() (Message) {
  return <- items;
}

func UnloadItem() {
//  for {
    metamorphosis(getMessage());
//  }
}

type KFKProducer struct {
  Worker *sarama.AsyncProducer;
}

func InitProducer(kfk_list []string) (*KFKProducer, error) {
  for i := 0; i < len(kfk_list); i++ {
    if strings.Trim(kfk_list[i], " \t") == "" {
      return nil, errors.New(fmt.Sprintf("invalid kafka brokers: %v\n", kfk_list));
    }
  }
  worker, err := newAsyncProducer(kfk_list);
  return &KFKProducer {Worker: &worker}, err;
}

func newSiestaProducer(kfk_list []string) (*producer.KafkaProducer, error) {
  conf := siesta.NewConnectorConfig();
  conf.BrokerList = kfk_list;
  conf.MaxConnectionsPerBroker = MAX_CONN_PER_BROKER;
  conf.MaxConnections = conf.MaxConnectionsPerBroker * len(conf.BrokerList);
  conf.MetadataRetries = STD_RETRY;
  conn, err := siesta.NewDefaultConnector(conf);
  if err != nil {
    return nil, err;
  } else {
    producerConf := producer.NewProducerConfig();
    producerConf.BatchSize = BatchSize;
    producerConf.ClientID = "4710138111115";
    producerConf.RequiredAcks = 0; // Don't wait for ACKs
    producerConf.Retries = STD_RETRY;
//    producerConf.MaxRequests = 1;
//    producerConf.SendRoutines = 1;
    p := producer.NewKafkaProducer(producerConf, producer.ByteSerializer, producer.ByteSerializer, conn);
    return p, nil;
  }
  return nil, ERR_UNKNOWN;
}

func initSiestaProducer(kafkas string) (*producer.KafkaProducer, error) {
  kfk_list := strings.Split(kafkas, ",");
  for i := 0; i < len(kfk_list); i++ {
    if strings.Trim(kfk_list[i], " \t") == "" {
      return nil, errors.New(fmt.Sprintf("[initSiestaProducer] invalid kafka brokers: %s\n", kafkas));
    }
  }
  return newSiestaProducer(kfk_list);
}

func initProducer(kafkas string) (sarama.AsyncProducer, error) {
  kfk_list := strings.Split(kafkas, ",");
  for i := 0; i < len(kfk_list); i++ {
    if strings.Trim(kfk_list[i], " \t") == "" {
      return nil, errors.New(fmt.Sprintf("invalid kafka brokers: %s\n", kafkas));
    }
  }
  return newAsyncProducer(kfk_list);
}

func metamorphosis(m Message) {
  //fmt.Printf("metamorphosis msg: %v\n", m);
  wagon.Input() <- &sarama.ProducerMessage {
    Topic: m.Topic,
    Value: sarama.StringEncoder([]byte(m.Msg)),
  }
  //fmt.Printf("msg sent via %v\n", wagon);
}

func WriteMsg(p sarama.AsyncProducer, m Message) {
  //fmt.Printf("p: %v\n", p);
  p.Input() <- &sarama.ProducerMessage {
    Topic: m.Topic,
    Value: sarama.StringEncoder([]byte(m.Msg)),
  };
}

func Metamorphosis(p sarama.AsyncProducer, msg string, topic string) {
  p.Input() <- &sarama.ProducerMessage {
    Topic: topic,
    Value: sarama.StringEncoder([]byte(msg)),
  }
}

