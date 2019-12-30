<template>
  <div class="container">
    <div class="col">
      <h1>{{ sorted_bitmex_bids }}</h1>
    </div>
    <div class="col">
      <h1>{{ sorted_bitmex_asks }}</h1>
    </div>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'Data',
  computed: {
    sorted_bitmex_asks: function() {
      let items = this.bitmex_asks.sort((a, b) => parseFloat(a.price) - parseFloat(b.price))
      return items.slice(0,5)
    },

    sorted_bitmex_bids: function() {
      let items = this.bitmex_bids.sort((a, b) => parseFloat(b.price) - parseFloat(a.price))
      return items.slice(0,5)
    }
  },

  data () {
    return {
      bitmex_asks: [],
      bitmex_bids: []
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

      let ws = new WebSocket("ws://localhost:8000/ws")
      

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
