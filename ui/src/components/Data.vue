<template>
  <div class="leveltwo">
    <div 
      class="col"
      v-cloak>
        <div
        v-for="l in levelTwoBids" 
        :key="l">
        <span class="exchange">{{ l.exchange }}</span>
        </div>
    </div>
    <div 
      class="col"
      v-cloak>
        <div
        v-for="l in levelTwoBids" 
        :key="l">
        <span class="price">{{ l.price }}</span>
        </div>
    </div>
    <div 
      class="col"
      v-cloak>
        <div
        v-for="l in levelTwoBids" 
        :key="l">
        <span class="size">{{ parseFloat(l.size).toFixed(3) }}</span>
        </div>
    </div> 
    <div 
      class="col"
      v-cloak>
        <div
        v-for="l in levelTwoAsks" 
        :key="l">
        <span class="exchange">{{ l.exchange }}</span>
        </div>
    </div>
    <div 
      class="col"
      v-cloak>
        <div
        v-for="l in levelTwoAsks" 
        :key="l">
        <span class="price">{{ l.price }}</span>
        </div>
    </div>
    <div 
      class="col"
      v-cloak>
        <div
        v-for="l in levelTwoAsks" 
        :key="l">
        <span class="size">{{ parseFloat(l.size).toFixed(3) }}</span>
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
              .map((m) => { return {"exchange": "BNC", "price": m[0], "size": m[1]} })

            this.bidsBinance = dump["bids"]
              .map((m) => { return {"exchange": "BNC", "price": m[0], "size": m[1]} })

            break

          case "bitmex":
            switch (dump["action"]) {
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
                      this.asksBinance.push({"exchange": "BMX", "id": m.id, "price": m.price, "size": m.size})
                  })

                dump["data"]
                  .filter(b => b.side === "Buy")
                    .map((m) => {
                      this.bidsBinance.push({"exchange": "BMX", "id": m.id, "price": m.price, "size": m.size})
                  })
                  
                break

              default:
                break
            }

            break

          case "okex":
            this.asksOkex = dump["data"][0]["asks"]
              .map((m) => { return {"exchange": "OKX", "price": m[0], "size": m[1]} })

            this.bidsOkex = dump["data"][0]["bids"]
              .map((m) => { return {"exchange": "OKX", "price": m[0], "size": m[1]} })

            break
          
          default:
            break
        }

      let asks = this.asksOkex.concat(this.asksBinance, this.asksBitmex)
      this.levelTwoAsks = asks.sort((a, b) => {
        return a.price - b.price
      }).slice(0, 15)

      let bids = this.bidsOkex.concat(this.bidsBinance, this.bidsBitmex)
      this.levelTwoBids = bids.sort((a, b) => {
        return b.price - a.price
      }).slice(0, 15)

      })
    } catch (e) {
      console.log(e)
    }
  },

  data () {
    return {
      asksBinance: [],
      asksBitmex: [],
      asksOkex: [],
      bidsBinance: [],
      bidsBitmex: [],
      bidsOkex: [],
      levelTwoAsks: [],
      levelTwoBids: [],
    }
  },

  methods: {

  },

  mounted() {
    try {
      axios.get("http://localhost:8000/ids")
      .then(r => {
        this.asksBitmex = r.data
          .filter(b => b.side == "Sell")
            .map((m) => { return {"exchange": "BMX", "id": m.id, "price": m.price, "size": m.size} })
              .sort((a, b) => { return a.price - b.price })

        this.bidsBitmex = r.data
          .filter(b => b.side == "Buy")
            .map((m) => { return {"exchange": "BMX", "id": m.id, "price": m.price, "size": m.size} })
              .sort((a, b) => { return b.price - a.price })
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

.col {
  display: inline-block;
  margin: 0px;
  width: 9%;
}

.exchange {
  font-weight:200;
  font-size:14px;
}

.leveltwo {
  display: inline-block;
  margin: 0px;
  text-align: center;
  padding-top: 3%;
  width: 75%;
}

@keyframes appear {
  from { opacity: 0; }
  to { opacity: 1; }
}

</style>
