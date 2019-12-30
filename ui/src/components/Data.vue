<template>
  <div class="container">
    <div class="col">
      <p>{{ sorted_bitmex_bids }}</p>
    </div>
    <div class="col">
      <p>{{ sorted_bitmex_asks }}</p>
    </div>
    <div class="col">
      <!-- <p>{{ data }}</p> -->
    </div>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'Data',
  computed: {
    sorted_binance_asks: function() {
      return
    },

    sorted_binance_bids: function() {
      return
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
    }
  },

  created() {
    try {
      let ws = new WebSocket("ws://localhost:8000/ws")
      ws.onmessage = ((e) => {
        switch (JSON.parse(e.data)["host"]) {
          case "binance":
            this.binance_asks = e.data["asks"]
            this.binance_bids = e.data["bids"]
            break

          case "bitmex":
            switch (e.data["action"]) {
              case "update":
                
                break

              default:
                break
            }

            break

          case "okex":
            this.okex_asks = e.data["data"]["asks"]
            this.okex_bids = e.data["data"]["bids"]
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
