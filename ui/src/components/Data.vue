<template>
  <div class="container">
    <div class="col"
      v-cloak>
      <p>{{ asksBitmex }}</p>
      <p>{{ bidsBitmex }}</p>
    </div>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  computed: {
    
  },

  created() {
    try {
      let ws = new WebSocket("ws://localhost:8000/ws")
      ws.onmessage = ((e) => {
        let dump = JSON.parse(e.data)

        switch (dump["host"]) {
          case "binance":
            this.asksBinance = dump["asks"]
              .map((m) => { return {"exchange": "BNC", "price": m[0], "size": m[1]} })
                .sort((a, b) => { return a.price - b.price })

            this.bidsBinance = dump["bids"]
              .map((m) => { return {"exchange": "BNC", "price": m[0], "size": m[1]} })
                .sort((a, b) => { return b.price - a.price })

            break

          case "bitmex":
            switch (dump["action"]) {
              case "update":
                let bmxBuy = dump["data"]
                  .filter(b => b.side === "Buy")
                    .map((m) => {
                      let n = this.bidsBitmex.findIndex(o => m.id === o.id)
                      this.bidsBitmex[n].size = m.size
                  })

                let bmxSell = dump["data"]
                  .filter(a => a.side === "Sell")
                    .map((m) => {
                      let n = this.asksBitmex.findIndex(o => m.id === o.id)
                      this.asksBitmex[n].size = m.size
                  })

                break

              case "delete":
                
                
                break

              case "insert":
                // for (let i = 0; i < dump["data"].length; i++) {
                  
                // }
                break

              default:
                break
            }

            break

          case "okex":
            this.asksOkex = dump["data"][0]["asks"]
              .map((m) => { return {"exchange": "OKX", "price": m[0], "size": m[1]} })
                .sort((a, b) => { return a.price - b.price })

            this.bidsOkex = dump["data"][0]["bids"]
              .map((m) => { return {"exchange": "OKX", "price": m[0], "size": m[1]} })
                .sort((a, b) => { return b.price - a.price })

            break
          
          default:
            break
        }

        
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
      bidsOkex: []
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

body {
  display: inline-block;
  margin: 0;
  font-weight: 200;
  font-size: 16px;
}

#content {
  padding-top: 3%;
}

@keyframes appear {
  from { opacity: 0; }
  to { opacity: 1; }
}
</style>
