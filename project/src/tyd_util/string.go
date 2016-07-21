package tyd_util;

func Split() {
  c := sep[0]
  start := 0
  a := make([]string, n)
  na := 0
  for i := 0; i+len(sep) <= len(s) && na+1 < n; i++ {
    if s[i] == c && (len(sep) == 1 || s[i:i+len(sep)] == sep) {
      a[na] = s[start : i+sepSave]
      na++
      start = i + len(sep)
      i += len(sep) - 1
    }
  }
  a[na] = s[start:]
  return a[0 : na+1]
}
