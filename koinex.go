package main

import(
    "encoding/json"
    "fmt"
    "net/http"
    "io/ioutil"
    "os/exec"
    "time"
    "strconv"
    "log"
    "flag"
)


const (
    VERSION = 1
    KOINEX_JSON = "https://koinex.in/api/ticker"
)

type Prices struct {
    Btc string `json:"BTC"`
    Eth string `json:"ETH"`
    Xrp string `json:"XRP"`
    Bch string `json:"BCH"`
    Ltc string `json:"LTC"`
}

type Response struct {
  Price Prices `json:"prices"`
}

func notify( message string) {
  cmd := exec.Command("notify-send", "-u", "critical",
                    "-i", "notification-message-im", message)
  cmd.Run()
}

func main() {
  min := flag.Float64("min", 45.0, "Lower threshold")
  max := flag.Float64("max", 100.0, "Upper threshold")

  flag.Parse()
  previous := 0.0;
  for {
    url:= "https://koinex.in/api/ticker"
    res,err := http.Get(url)
    if err != nil {
      log.Print(err)
      continue
    }
    body,err := ioutil.ReadAll(res.Body)
    if err != nil {
      log.Print(err)
      continue
    }
    price := Response{}
    err = json.Unmarshal([]byte(body), &price)
    if err != nil {
      log.Print(err)
      continue
    }

    rpc, err:= strconv.ParseFloat(price.Price.Xrp, 32)

    if err != nil {
      log.Print(err)
      continue
    }

    ripple := fmt.Sprintf("XRP -> %s", price.Price.Xrp)
    if rpc > *max || rpc < *min {
      if rpc != previous {
        notify(fmt.Sprintf("%s", ripple))
        previous = rpc
      }
    }
    time.Sleep(time.Second*60)

  }
}
