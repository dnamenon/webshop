package shopsess

import (
  "net/http"

  "github.com/astaxie/beego/session"
)

type Shopsess struct {
    tf bool
    Nsess session.Manager()
}

// Middleware is a struct that has a ServeHTTP method
func NewShopsess() *Shopsess {
    return &Shopsess{true}
}





func (l *Shopsess) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc){
        Nsess, _ := session.NewManager(
          "cookie", `{"cookieName":"gosessionid","enableSetCookie":false,"gclifetime":3600,"ProviderConfig":"{\"cookieName\":\"gosessionid\",\"securityKey\":\"beegocookiehashkey\"}"}`)
      go Nsess.GC()

next(rw, r)

return Nsess()
}
