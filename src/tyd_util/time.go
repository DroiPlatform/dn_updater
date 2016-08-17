package tyd_util;

import "time";
import "fmt";

const LAYOUT = "2006/01/02 15:04:05";

/* print time */
func PrintTime() {
  t := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), time.Now().Nanosecond(), time.Local);
  fmt.Printf("[T:%v] ", t.Format(LAYOUT));
}

/* get time */
func GetTime() (string) {
  t := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), time.Now().Nanosecond(), time.Local);
  return fmt.Sprintf("%v: ", t.Format(LAYOUT));
}

