<template>
  <div>
    <div class="slide-left">
      <h1 id="name">exchange.</h1>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {

    }
  },

  name: 'Header'
}
</script>

<style scoped>
h1, h2 {
  font-weight: normal;
}

ul {
  list-style-type: none;
  padding: 0;
}

li {
  display: inline-block;
  margin: 0 10px;
}

a {
  color: #42b983;
}

div.slide-left {
  width:100%;
  overflow:hidden;
}
div.slide-left h1 {
  animation: slide-left 0.25s;
}

@keyframes slide-left {
  from {
    margin-left: 100%;
    width: 300%;
  }

  to {
    margin-left: 0%;
    width: 100%;
  }
}
</style>
