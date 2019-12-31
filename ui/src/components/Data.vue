<template>
  <div class="container">
    <div class="col">
      <p>{{ sorted_bitmex_bids }} {{ sorted_bitmex_asks }}</p>
    </div>
    <div class="col">
      <p>{{ sorted_binance_bids }} {{ sorted_binance_asks }}</p>
    </div>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'Data',
  computed: {
    sorted_binance_asks: function() {
      let items = this.binance_asks.sort((a, b) => { return a - b })
      return items.slice(0,5)
    },

    sorted_binance_bids: function() {
      let items = this.binance_asks.sort((a, b) => { return b -  a})
      return items.slice(0,5)
    },
  
    sorted_bitmex_asks: function() {
      let items = this.bitmex_asks.sort((a, b) => parseFloat(a.price) - parseFloat(b.price))
      return items.slice(0,5)
    },

    sorted_bitmex_bids: function() {
      let items = this.bitmex_bids.sort((a, b) => parseFloat(b.price) - parseFloat(a.price))
      return items.slice(0,5)
    },

    sorted_okex_asks: function() {
      return
    },

    sorted_okex_bids: function() {
      return
    },

    master: function() {
          
    }
  },

  created() {
    try {
      let ws = new WebSocket("ws://localhost:8000/ws")
      ws.onmessage = ((e) => {
        let dump = JSON.parse(e.data)

        switch (dump["host"]) {
          case "binance":
            this.binance_asks = dump["asks"]
            this.binance_bids = dump["bids"]
            break

          case "bitmex":
            switch (dump["action"]) {
              case "update":
                for (let i = 0; i < dump["data"].length; i++) {
                  if (dump["data"][i]["side"] == "Buy") {
                    this.bitmex_bids[this.bitmex_bids.findIndex(n => n.id === dump["data"][i]["id"])]["size"] = dump["data"][i]["size"]
                  } else if (dump["data"][i]["side"] == "Sell") {
                    this.bitmex_asks[this.bitmex_asks.findIndex(n => n.id === dump["data"][i]["id"])]["size"] = dump["data"][i]["size"]
                  }
                }
                break

              case "delete":
                for (let i = 0; i < dump["data"].length; i++) {
                  if (dump["data"][i]["side"] == "Buy") {
                    this.bitmex_bids.splice(this.bitmex_bids.findIndex(n => n.id === dump["data"][i]["id"]), 1)
                  } else if (dump["data"][i]["side"] == "Sell") {
                    this.bitmex_asks.splice(this.bitmex_asks.findIndex(n => n.id === dump["data"][i]["id"]), 1)
                  }
                }
                break

              case "insert":
                for (let i = 0; i < dump["data"].length; i++) {
                  if (dump["data"][i]["side"] == "Buy") {
                    this.bitmex_bids.push(dump["data"][i])
                  } else if (dump["data"][i]["side"] == "Sell") {
                    this.bitmex_asks.push(dump["data"][i])
                  }
                }
                break

              default:
                break
            }

            break

          case "okex":
            this.okex_asks = dump["data"]["asks"]
            this.okex_bids = dump["data"]["bids"]
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
      binance_asks: [],
      binance_bids: [],
      bitmex_asks: [],
      bitmex_bids: [],
      master_sorted: [],
      okex_asks: [],
      okex_bids: []
    }
  },

  mounted() {
    try {
      axios.get("http://localhost:8000/ids")
      .then(r => {
        this.bitmex_asks = r.data.filter((b) => {
          return b.side == "Sell"
        })

        this.bitmex_bids = r.data.filter((b) => {
          return b.side == "Buy"
        })
      })
    } catch(e) {
      console.log(e)
    }
  }
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
