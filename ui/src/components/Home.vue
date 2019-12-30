<template>
  <div class="hello">
    <exchg-data/>
  </div>
</template>

<script>
import Data from './Data'

export default {
  components: {
    "exchg-data": Data
  },
  name: 'Home',
  data () {
    return {
      msg: 'Exchange Aggregator'
    }
  }
}
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css?family=Zeyada&display=swap');

#box-text {
  padding: 1%;
}

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
