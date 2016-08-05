package tyd_util;

import "bufio";
import "fmt";
import "os";
import "strings";

type Conf struct {
  Key string;
  Value string;
}

type GlobalConf struct {
  Directives []Conf;
};

/* print config */
func (gc GlobalConf) PrintConfig() {
  for _, value := range gc.Directives {
    UIDLogPrinter(UID {uint8(0), uint64(0), uint64(0)}, "", "INFO", fmt.Sprintf("Key: %s, Value: %s\n", value.Key, value.Value));
  }
}

/* get Conf structure */
func GetConf(p_conf string) (GlobalConf) {
  conf := GlobalConf{};
  pcfp, err := os.Open(p_conf);
  if err != nil {
    panic(err);
  }
  defer pcfp.Close();
  prescanner := bufio.NewScanner(pcfp);
  cf_line := 0;
  for prescanner.Scan() {
    if strings.Contains(prescanner.Text(), "=") {
      cf_line++;
    }
  }
  cfp, err := os.Open(p_conf);
  if err != nil {
    panic(err);
  }
  defer cfp.Close();
  key := make([]string, cf_line);
  value := make([]string, cf_line);
  scanner := bufio.NewScanner(cfp);
  cptr := 0;
  lptr := 1;
  for scanner.Scan() {
    config := strings.Split(scanner.Text(), "#");
    rec := strings.Split(config[0], "=");
    if len(rec) == 2 {
      key[cptr] = strings.TrimSpace(rec[0]);
      value[cptr] = strings.TrimSpace(rec[1]);
      /*
      PrintTime();
      fmt.Printf("%d-th key: %s, value: %s\n", cptr, key[cptr], value[cptr]);
      */
      cptr++;
    } else if len(rec) == 1 && rec[0] == "" {
    } else {
      UIDLogPrinter(UID {uint8(0), uint64(0), uint64(0)}, "", "ERR", fmt.Sprintf("config error (%s) @ line %d\n", rec, lptr));
      os.Exit(1);
    }
    lptr++;
  }
  if cf_line < cptr {
    UIDLogPrinter(UID {uint8(0), uint64(0), uint64(0)}, "", "ERR", fmt.Sprintf("config error, cf_line (%d) < cptr (%d)\n", cf_line, cptr));
  }
  for i := 0; i < cptr; i++ {
    conf.Directives = append(conf.Directives, Conf{key[i], value[i]});
  }
  return conf;
}
