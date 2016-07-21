package http_util;

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

/* HTTP conf structure */
type HTTPConf struct {
  S map[string]string;
  I map[string]int;
  B map[string]bool;
  err error;
};

/* print HTTPConf */
func (hc HTTPConf) PrintHTTPConf() {
  util.PrintPrefix(uint64(0), uint64(0), "info");
  fmt.Printf("--------- Printing HTTP config\n");
  for key, value := range hc.S {
    util.PrintPrefix(uint64(0), uint64(0), "info");
    fmt.Printf("%s: %s\n", key, value);
  }
  for key, value := range hc.I {
    util.PrintPrefix(uint64(0), uint64(0), "info");
    fmt.Printf("%s: %d\n", key, value);
  }
  for key, value := range hc.B {
    util.PrintPrefix(uint64(0), uint64(0), "info");
    fmt.Printf("%s: %v\n", key, value);
  }
}

/* init HTTPConf */
func initHTTPConf() (HTTPConf) {
  conf := HTTPConf{
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
      "URI": true,
      "IP_LIST": true,
      "LOG": true,
      "PORT": true,
      "SRV_LIST": true,
      "PORT_LIST": true,
    },
    /* int directives */
    map[string]bool{
      "MAXGO": true,
      "UPDATE_INTERVAL": true,
      "Q_TIMEOUT": true,
      "MAX_SAMPLE": true,
      "LOG_LEVEL": true,
    },
    /* bool directives */
    map[string]bool{
      "DEBUG": true,
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
func setHTTPConfString(key string, value string, conf *HTTPConf) {
  conf.S[key] = value;
}

/* set int directives */
func setHTTPConfInt(key string, value string, conf *HTTPConf) {
  ivalue, err := strconv.Atoi(value);
  if err != nil {
    util.PrintPrefix(uint64(0), uint64(0), "error");
    fmt.Printf("Directive %s: expected positive integer, got %s\n", key, value);
  }
  conf.I[key] = ivalue;
}

/* set bool directives */
func setHTTPConfBool(key string, value string, conf *HTTPConf) {
  if value == "true" {
    conf.B[key] = true;
  } else if value == "false" {
    conf.B[key] = false;
  } else {
    util.PrintPrefix(uint64(0), uint64(0), "error");
    fmt.Printf("Directive %s: expected boolean, got %s\n", key, value);
  }
}

/* is HTTPConf valid? */
func checkHTTPConf(conf *HTTPConf, tc TypeChecker) (bool) {
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

/* HTTPConf setting */
func GetHTTPConf(p_conf string) (*HTTPConf, error) {
  hconf := initHTTPConf();
  tc := initTypeChecker();
  conf := util.GetConf(p_conf);
  for key, value := range conf.Directives {
    switch tc.getType(value.Key) {
    case "S":
      setHTTPConfString(value.Key, value.Value, &hconf);
      break;
    case "I":
      setHTTPConfInt(value.Key, value.Value, &hconf);
      break;
    case "B":
      setHTTPConfBool(value.Key, value.Value, &hconf);
      break;
    default:
      util.PrintPrefix(uint64(0), uint64(0), "error");
      fmt.Printf("config error, unknown %d-the directive %s\n", key, value);
      break;
    }
  }
  if checkHTTPConf(&hconf, tc) {
    return &hconf, nil;
  } else {
    return nil, hconf.err;
  }
}

