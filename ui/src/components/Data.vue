<template>
  <div class="levelTwo">
    <div
    class="bids wrap">
      <div 
        class="col"
        v-cloak>
          <div
          v-for="(l, index) in levelTwoBids" 
          :key="index">
          <span class="exchange bids">{{ l.exchange }}</span>
          </div>
      </div>
      <div 
        class="col"
        v-cloak>
          <div
          v-for="(l, index) in levelTwoBids" 
          :key="index">
          <span class="price bids">{{ l.price }}</span>
          </div>
      </div>
      <div 
        class="col"
        v-cloak>
          <div
          v-for="(l, index) in levelTwoBids" 
          :key="index">
          <span class="size bids" v-if="l.exchange === 'BMX'">{{ l.size }}</span>
          <span class="size bids" v-else>{{ l.size }}</span>
          </div>
      </div>

    </div>
    <div
    class="asks wrap">
      <div 
        class="col"
        v-cloak>
          <div
          v-for="(l, index) in levelTwoAsks" 
          :key="index">
          <span class="exchange asks">{{ l.exchange }}</span>
          </div>
      </div>
      <div 
        class="col"
        v-cloak>
          <div
          v-for="(l, index) in levelTwoAsks" 
          :key="index">
          <span class="price asks">{{ l.price }}</span>
          </div>
      </div>
      <div 
        class="col"
        v-cloak>
          <div
          v-for="(l, index) in levelTwoAsks" 
          :key="index">
          <span class="size asks" v-if="l.exchange === 'BMX'">{{ l.size }}</span>
          <span class="size asks" v-else>{{ l.size }}</span>
          </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  created() {
    try {
      let ws = new WebSocket("ws://localhost:8000/ws")
      ws.onmessage = ((e) => {
        let dump = JSON.parse(e.data)

        switch (dump["host"]) {
          case "binance":
            this.asksBinance = dump["asks"]
              .map((m) => { return {"exchange": "BNC", "price": parseFloat(m[0]).toFixed(1), "size": m[1]} })

            this.bidsBinance = dump["bids"]
              .map((m) => { return {"exchange": "BNC", "price": parseFloat(m[0]).toFixed(1), "size": m[1]} })

            break

          case "bitmex":
            switch (dump["action"]) {
              case "delete":
                dump["data"]
                  .filter(a => a.side === "Sell")
                    .map((m) => {
                      let n = this.asksBitmex.findIndex(o => m.id === o.id)
                      this.asksBitmex.splice(n, 1)
                    })

                dump["data"]
                  .filter(b => b.side === "Bids")
                    .map((m) => {
                      let n = this.bidsBitmex.findIndex(o => m.id === o.id)
                      this.bidsBitmex.splice(n, 1)
                    })

                break

              case "insert":
                dump["data"]
                  .filter(a => a.side === "Sell")
                    .map((m) => {
                      this.asksBitmex.push({"exchange": "BMX", "id": m.id, "price": m.price, "size": m.size})
                  })

                dump["data"]
                  .filter(b => b.side === "Buy")
                    .map((m) => {
                      this.bidsBitmex.push({"exchange": "BMX", "id": m.id, "price": m.price, "size": m.size})
                  })
                  
                break

              case "partial":
                this.asksBitmex = []
                dump["data"]
                  .filter(a => a.side === "Sell")
                    .map((m) => {
                      this.asksBitmex.push({"exchange": "BMX", "id": m.id, "price": m.price, "size": m.size})
                  })
                
                this.bidsBitmex = []
                dump["data"]
                  .filter(b => b.side === "Buy")
                    .map((m) => {
                      this.bidsBitmex.push({"exchange": "BMX", "id": m.id, "price": m.price, "size": m.size})
                  })

                break

              case "update":
                let bmxSell = dump["data"]
                  .filter(a => a.side === "Sell")
                    .map((m) => {
                      let n = this.asksBitmex.findIndex(o => m.id === o.id)
                      this.asksBitmex[n].size = m.size
                  })

                let bmxBuy = dump["data"]
                  .filter(b => b.side === "Buy")
                    .map((m) => {
                      let n = this.bidsBitmex.findIndex(o => m.id === o.id)
                      this.bidsBitmex[n].size = m.size
                  })
                
                break

              default:
                break
            }

            break

          case "bitstamp":
            this.asksBitstamp = dump["data"]["asks"]
              .map((m) => { return {"exchange": "BTS", "price": parseFloat(m[0]).toFixed(1), "size": parseInt(parseFloat(m[1]) * parseFloat(m[0]))}})

            this.bidsBitstamp = dump["data"]["bids"]
              .map((m) => { return {"exchange": "BTS", "price": parseFloat(m[0]).toFixed(1), "size": parseInt(parseFloat(m[1]) * parseFloat(m[0]))}})

            break

          case "huobi":
            this.asksHuobi = dump["tick"]["asks"]
              .map((m) => { return {"exchange": "HBI", "price": parseFloat(m[0]).toFixed(1), "size": parseInt(parseFloat(m[1]) * parseFloat(m[0]))}})

            this.bidsHuobi = dump["tick"]["bids"]
              .map((m) => { return {"exchange": "HBI", "price": parseFloat(m[0]).toFixed(1), "size": parseInt(parseFloat(m[1]) * parseFloat(m[0]))}})

          case "okex":
            this.asksOkex = dump["data"][0]["asks"]
              .map((m) => { return {"exchange": "OKX", "price": parseFloat(m[0]).toFixed(1), "size": parseInt(m[1] * 100)}})

            this.bidsOkex = dump["data"][0]["bids"]
              .map((m) => { return {"exchange": "OKX", "price": parseFloat(m[0]).toFixed(1), "size": parseInt(m[1] * 100)}})

            break
          
          default:
            break
        }

      let asks = this.asksOkex.concat(this.asksBinance, this.asksBitmex, this.asksBitstamp, this.asksHuobi)
      this.levelTwoAsks = asks.sort((a, b) => {
        return a.price - b.price
      }).slice(0, 20)

      let bids = this.bidsOkex.concat(this.bidsBinance, this.bidsBitmex, this.bidsBitstamp, this.bidsHuobi)
      this.levelTwoBids = bids.sort((a, b) => {
        return b.price - a.price
      }).slice(0, 20)

      })
    } catch (e) {
      console.log(e)
    }
  },

  data () {
    return {
      asksBinance: [],
      asksBitmex: [],
      asksBitstamp: [],
      asksHuobi: [],
      asksOkex: [],
      bidsBinance: [],
      bidsBitmex: [],
      bidsBitstamp: [],
      bidsHuobi: [],
      bidsOkex: [],
      levelTwoAsks: [],
      levelTwoBids: [],
      symbol: null,
    }
  },

  mounted() {
    try {
      axios.get("http://localhost:8000/ids")
      .then(r => {
        this.asksBitmex = JSON.parse(r.data)
          .filter(b => b.side == "Sell")
            .map((m) => { return {"exchange": "BMX", "id": m.id, "price": parseFloat(m.price).toFixed(1), "size": m.size}})
              .sort((a, b) => { return a.price - b.price })

        this.bidsBitmex = JSON.parse(r.data)
          .filter(b => b.side == "Buy")
            .map((m) => { return {"exchange": "BMX", "id": m.id, "price": parseFloat(m.price).toFixed(1), "size": m.size} })
              .sort((a, b) => { return b.price - a.price })

        this.symbol = r.data[0].symbol
      })
    } catch(e) {
      console.log(e)
    }
  },

  name: 'Data'
}
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css?family=Zeyada&display=swap');

.asks {
  color: red;
  display: inline-block;
  padding: 2%;
}

.bids {
  color: green;
  display: inline-block;
  padding: 2%
}

.col {
  display: inline-block;
  margin: 0 25px;
  text-align: left;
  width: 35px;
}

.exchange {
  font-weight:200;
  font-size:14px;
  text-align: left;
}

.levelTwo {
  display: inline-block;
  margin: 0px;
  text-align: center;
  padding-top: 1%;
  width: 75%;
}

.price {
  text-align: left;
}

.size {
  text-align: left;
}

.wrap {
  border: dashed 1px;
}

@keyframes appear {
  from { opacity: 0; }
  to { opacity: 1; }
}
</style>
