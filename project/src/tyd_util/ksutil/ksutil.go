package ksutil;

import "errors";
import "fmt";
import "strconv";

import util "tyd_util";

/* bit mask for getting directive types */
const SMask = uint8(1);
const IMask = uint8(2);
const BMask = uint8(4);

/* conf type checker structure */
type TypeChecker struct {
  SDirectives map[string]bool;
  IDirectives map[string]bool;
  BDirectives map[string]bool;
}

/* Srv conf structure */
type SrvConf struct {
  S map[string]string;
  I map[string]int;
  B map[string]bool;
  err error;
};

/* init SrvConf */
func initSrvConf() (SrvConf) {
  conf := SrvConf{
    make(map[string]string),  /* string directives */
    make(map[string]int), /* int directives */
    make(map[string]bool), /* bool directives */
    nil,  /* error */
  };
  return conf;
}

/* init type checker */
func initTypeChecker() (TypeChecker) {
  return TypeChecker{
    /* string directives */
    map[string]bool{
      "KEY": true,
      "RSA": true,
      "K_RELEASE": true,
      "R_RELEASE": true,
      "ID_MAP": true,
      "R_MAP": true,
      //"MAL": true,
      "TSDB": true,
      //"EXT_MGO": true,
      "HOST": true, /* keyserver ip */
      "INTERNAL_IP": true, /* keyserver ip */
      "REDIS": true, /* REDIS server ip:port */
    },
    /* int directives */
    map[string]bool{
      "DB_SIZE": true,
      "TIMEOUT": true,
      "ASYNC": true,
      "KL_VERSION": true,
      "MAX_RSA": true,
      "PARTITION": true,
      "MAXGO": true,
      "REKEY": true,
      "REKEY_INTERVAL": true,
      "MAL_INTERVAL": true,
      "NEW_ID_THD": true,
      "WINDOW": true, /* keyserver alive window, in second*/
      "PORT": true, /* keyserver serving port */
      "KS_VERSION": true, /* keyserver version */
      "AGENT_PORT": true, /* port of nginx agent */
      "MC_PORT": true, /* port for memcache */
      "LOG_LEVEL": true, /* log level */
      "REDIS_POOL": true, /* size of connection pool for redis timestamp query */
      "DB": true, /* DB selector, 0: mongo, 1: redis*/
    },
    /* bool directives */
    map[string]bool{
      "DEBUG": true,
      "UTL_DEBUG": true,
      "TEST": true,
    },
  };
}

/* check directive type */
func (tc TypeChecker) getType(key string) (string) {
  /* 
  bit mask:
  1: string;
  2: int;
  3: bool;
  */
  _, sok := tc.SDirectives[key];
  _, iok := tc.IDirectives[key];
  _, bok := tc.BDirectives[key];
  mask := uint8(0)
  if sok {
    mask = mask | uint8(SMask);
  }
  if iok {
    mask = mask | uint8(IMask);
  }
  if bok {
    mask = mask | uint8(BMask);
  }
  switch mask {
  case SMask:
    /* string directive */
    return "S";
  case IMask:
    /* int directive */
    return "I";
  case BMask:
    /* bool directive */
    return "B";
  default:
    /* unknown directive */
    return "U";
  }
  return "U";
}

/* set string directives */
func setSrvConfString(key string, value string, conf *SrvConf) {
  conf.S[key] = value;
}

/* set int directives */
func setSrvConfInt(key string, value string, conf *SrvConf) {
  ivalue, err := strconv.Atoi(value);
  if err != nil {
    util.PrintTime();
    fmt.Printf("Directive %s: expected positive integer, got %s\n", key, value);
  }
  conf.I[key] = ivalue;
}

/* set bool directives */
func setSrvConfBool(key string, value string, conf *SrvConf) {
  if value == "true" {
    conf.B[key] = true;
  } else if value == "false" {
    conf.B[key] = false;
  } else {
    util.PrintTime();
    fmt.Printf("Directive %s: expected boolean, got %s\n", key, value);
  }
}

/* is SrvConf valid? */
func checkSrvConf(conf *SrvConf, tc TypeChecker) (bool) {
  for key, _ := range tc.SDirectives {
    _, ok := conf.S[key];
    if !ok {
      err := fmt.Sprintf("Config error: %s is not set!", key);
      conf.err = errors.New(err);
      return false;
    }
  }
  for key, _ := range tc.IDirectives {
    _, ok := conf.I[key];
    if !ok {
      err := fmt.Sprintf("Config error: %s is not set!", key);
      conf.err = errors.New(err);
      return false;
    }
  }
  for key, _ := range tc.BDirectives {
    _, ok := conf.B[key];
    if !ok {
      err := fmt.Sprintf("Config error: %s is not set!", key);
      conf.err = errors.New(err);
      return false;
    }
  }
  return true;
}

/* SrvConf setting */
func GetSrvConf(p_conf string) (SrvConf, error) {
  hconf := initSrvConf();
  tc := initTypeChecker();
  conf := util.GetConf(p_conf);
  for key, value := range conf.Directives {
    switch tc.getType(value.Key) {
    case "S":
      setSrvConfString(value.Key, value.Value, &hconf);
      break;
    case "I":
      setSrvConfInt(value.Key, value.Value, &hconf);
      break;
    case "B":
      setSrvConfBool(value.Key, value.Value, &hconf);
      break;
    default:
      util.PrintTime();
      fmt.Printf("config error, unknown %d-the directive %s\n", key, value);
      break;
    }
  }
  if checkSrvConf(&hconf, tc) {
    return hconf, nil;
  } else {
    return hconf, hconf.err;
  }
}

