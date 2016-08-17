package conn_pool;

import "net/http";

/* generic pool */
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

