<template>
  <div id="app">
    <hdr id="hdr-id"/>
    <transition
      name="fade"
      mode="out-in">
      <router-view />
    </transition>
  </div>
</template>

<script>
import Header from './components/Header'

export default {
  components: {
    'hdr': Header
  },

  name: 'App'
}
</script>

<style>
@import url('https://fonts.googleapis.com/css?family=Montserrat&display=swap');

#app {
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  font-family: 'Montserrat', sans-serif;
  text-align: center;
}

.fade-enter-active,
.fade-leave-active {
  transition-duration: 0.8s;
  transition-property: opacity;
  transition-timing-function: ease;
}

.fade-enter,
.fade-leave-active {
  opacity: 0
}

html {
  margin-left: calc(100vw - 100%);
}

#hdr-id {
  align-content: center;
  height: 100%;
  margin-top: 5%;
  padding: 5px;
  width: 100%;
}
</style>
