package kfk_producer;

import "errors";
import "fmt";
import "strings";
import "time";

import "github.com/Shopify/sarama";

type ConnectionPoolWrapper struct {
  size int;
  conn chan sarama.AsyncProducer;
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
  return sarama.NewAsyncProducer(kfk_list, conf);
}

func NewAsyncProducer(kfk_list []string) (sarama.AsyncProducer, error) {
  conf := sarama.NewConfig();
  conf.Producer.RequiredAcks = sarama.WaitForLocal;
  conf.Producer.Flush.Frequency = 100 * time.Millisecond;
  return sarama.NewAsyncProducer(kfk_list, conf);
}

func InitWagon(size int, kafkas string) (error) {
  kfk_list := strings.Split(kafkas, ",");
  for i := 0; i < len(kfk_list); i++ {
    //if strings.Compare(strings.Trim(" \t"), "") == 0 {
    if strings.Trim(kfk_list[i], " \t") == "" {
      return errors.New(fmt.Sprintf("invalid kafka brokers: %s\n", kafkas));
    }
  }
  var err error;
  wagon, err = newAsyncProducer(kfk_list);
  //fmt.Printf("producer: %v\n", wagon);
  if err != nil {
    return err;
  } else {
    /*
    wagon.wagon = producer;
    wagon.items = make(chan Message, size);
    //*/
    items = make(chan Message, size);
    return nil;
  }
}

func LoadItem(msg Message) {
  items <- msg;
  //fmt.Printf("channel: %v\n", items);
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

