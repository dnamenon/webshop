package shopsess

import (
  "net/http"

  "github.com/astaxie/beego/session"
)


func Sessions(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc){
  globalSessions, _ := session.NewManager(
          "cookie", `{"cookieName":"gosessionid","enableSetCookie":false,"gclifetime":3600,"ProviderConfig":"{\"cookieName\":\"gosessionid\",\"securityKey\":\"beegocookiehashkey\"}"}`)
      go globalSessions.GC()

next(rw, r)
}
