package kfk_producer;

import "errors";
import "fmt";
import "strings";
import "time";

import "github.com/elodina/siesta";
import producer "github.com/elodina/siesta-producer";

import pool "tyd_util/resource_pool";

type Siesta struct {
  p *producer.KafkaProducer;
  q chan *producer.ProducerRecord;
}

type SiestaPoolWrapper struct {
  size int;
  index *pool.IndexPoolWrapper;
  s []*Siesta;
}

type Message struct {
  Msg string;
  Topic string;
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

/* Initialization of Siesta Producers */
/* create a siesta producer */
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
    p := producer.NewKafkaProducer(producerConf, producer.ByteSerializer, producer.ByteSerializer, conn);
    return p, nil;
  }
  return nil, ERR_UNKNOWN;
}

/* preprocess the kfk server string */
func initSiestaProducer(kafkas string) (*producer.KafkaProducer, error) {
  kfk_list := strings.Split(kafkas, ",");
  for i := 0; i < len(kfk_list); i++ {
    if strings.Trim(kfk_list[i], " \t") == "" {
      return nil, errors.New(fmt.Sprintf("[initSiestaProducer] invalid kafka brokers: %s\n", kafkas));
    }
  }
  return newSiestaProducer(kfk_list);
}

/* Pooled KFK Producer (producer pool is managed here) */
/* init */
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

/* background go routine for msg queue wipeout */
func (s *Siesta) siestaSend() {
  var msg *producer.ProducerRecord;
  for msg = range s.q {
    s.p.Send(msg);
  }
}

/* enqueue msg to msg queue */
func (p *SiestaPoolWrapper)WriteSiesta(msg Message) (error) {
  index := p.index.GetIndex();
  err := WriteUnpooledSiesta(p.s[index], msg);
  p.index.FreeIndex(index);
  return err;
}

/* Unpooled KFK Producer (producer pool is either managed here or no pool is designed) */
/* init */
func InitProducer(kfk_list []string) (*Siesta, error) {
  for i := 0; i < len(kfk_list); i++ {
    if strings.Trim(kfk_list[i], " \t") == "" {
      return nil, errors.New(fmt.Sprintf("invalid kafka brokers: %v\n", kfk_list));
    }
  }
  worker, err := newSiestaProducer(kfk_list);
  if err != nil {
    return nil, err;
  }
  s := &Siesta {p: worker, q: make(chan *producer.ProducerRecord, QueueLength)};
  s.siestaSend();
  return s, nil;
}

/* enqueue msg to msg queue */
func WriteUnpooledSiesta(p *Siesta, msg Message) (error) {
  select {
  case p.q <- &producer.ProducerRecord{Topic: msg.Topic, Value: []byte(msg.Msg)}:
    return nil;
  case <-time.After(EnqueueTimeout):
    return ERR_TIMEOUT;
  }
}

