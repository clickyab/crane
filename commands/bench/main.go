package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"clickyab.com/crane/openrtb/v2.5"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/safe"
	"github.com/golang/protobuf/jsonpb"
)

var cc = `
{  
   "id":"fd81cb064d4b3384394663762c091f0fddc044246",
   "imp":[  
      {  
         "id":"36f75f0f8647b20ad7290787c61578f832b3b2f2",
         "banner":{  
            "w":1200,
            "h":627,
            "id":"b6f04d3022664f97347182b270c57934b149127c"
         },
         "bidfloor":1
      }
   ],
   "app":{  
      "bundle":"ir.persianfox.messenger"
   },
   "device":{  
      "ua":"Dalvik/2.1.0 (Android)",
      "geo":{  
         "lat":35.693693693693696,
         "lon":51.3597577736758
      },
      "ip":"91.229.215.102",
      "model":"samsung",
      "osv":"24",
      "hwv":"SM-G930F",
      "h":1920,
      "w":1080,
      "ppi":480,
      "language":"en",
      "carrier":"Irancell",
      "mccmnc":"432-35",
      "connectiontype":6
   },
   "user":{  
      "id":"cc56381ef2b4b389c09252d96db89e60",
      "geo":{  
         "lat":35.693693693693696,
         "lon":51.3597577736758
      }
   },
   "at":0
}`

func main() {
	var rq = &openrtb.BidRequest{}
	p := bytes.NewBufferString(cc)
	err := jsonpb.Unmarshal(p, rq)
	if err != nil {
		log.Fatal(err)
	}

	for {
		safe.GoRoutine(context.Background(), func() {
			for i := 0; i < 50; i++ {
				safe.GoRoutine(context.Background(), func() {
					for {
						c := http.Client{}
						p := bytes.NewBufferString(cc)
						rc := ioutil.NopCloser(p)
						ha := http.Header{}
						ha.Add("Content-Type", "application/json")
						url, err := url.Parse("http://demand.cliclyab.ae/api/ortb/f7033f7f55e99da475097798aa611e0b390a8f79")
						assert.Nil(err)
						rq := &http.Request{
							Body:   rc,
							Method: "POST",
							Header: ha,
							URL:    url,
						}
						now := time.Now()
						c.Do(rq)
						fmt.Println("Latency: ", time.Since(now))
					}

				})
			}
		})
		time.Sleep(time.Second * 20)
	}

	time.Sleep(360 * time.Second)
}
