package resource_pool;

import "net/http";

/* generic connection pool */
type InitFunc func() (interface{}, error);

type ConnectionPoolWrapper struct {
  size int;
  conn chan interface{};
}

func (p *ConnectionPoolWrapper) InitPool(size int, initfn InitFunc) (error) {
  p.conn = make(chan interface{}, size);
  for i := 0; i < size; i++ {
    conn, err := initfn();
    if err != nil {
      return err;
    } else {
      p.conn <- conn;
    }
  }
  p.size = size;
  return nil;
}

func (p *ConnectionPoolWrapper) GetConnection() (interface{}) {
  return <- p.conn;
}

func (p *ConnectionPoolWrapper) FreeConnection(conn interface{}) {
  p.conn <- conn;
}

/* index pool */
type IndexPoolWrapper struct {
  size int;
  index chan int;
}

func (p *IndexPoolWrapper) IndexInitPool(size int) {
  p.index = make(chan int, size);
  for i := 0; i < size; i++ {
    p.index <- i;
  }
  p.size = size;
  return;
}

func (p *IndexPoolWrapper) GetIndex() (int) {
  return <- p.index;
}

func (p *IndexPoolWrapper) FreeIndex(index int) {
  p.index <- index;
}

/* http pool */
type HTTPInitFunc func() (*http.Client, error);

type HTTPPoolWrapper struct {
  size int;
  client chan *http.Client;
}

func (p *HTTPPoolWrapper) HTTPInitPool(size int, initfn HTTPInitFunc) (error) {
  p.client = make(chan *http.Client, size);
  for i := 0; i < size; i++ {
    client, err := initfn();
    if err != nil {
      return err;
    } else {
      p.client <- client;
    }
  }
  p.size = size;
  return nil;
}

func (p *HTTPPoolWrapper) GetClient() (*http.Client) {
  return <- p.client;
}

func (p *HTTPPoolWrapper) FreeClient(client *http.Client) {
  p.client <- client;
}

/* transport pool */
type TransportInitFunc func() (*http.Transport, error);

type TransportPoolWrapper struct {
  size int;
  transport chan *http.Transport;
}

func (p *TransportPoolWrapper) TransportInitPool(size int, initfn TransportInitFunc) (error) {
  p.transport = make(chan *http.Transport, size);
  for i := 0; i < size; i++ {
    transport, err := initfn();
    if err != nil {
      return err;
    } else {
      p.transport <- transport;
    }
  }
  p.size = size;
  return nil;
}

func (p *TransportPoolWrapper) GetTransport() (*http.Transport) {
  return <- p.transport;
}

func (p *TransportPoolWrapper) FreeTransport(transport *http.Transport) {
  p.transport <- transport;
}

