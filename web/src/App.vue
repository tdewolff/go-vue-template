<template>
  <div id="app">
    <div class="ui fixed inverted menu">
      <div class="ui container">
        <router-link to="/" class="header item">
          APPNAME
        </router-link>
        <router-link v-if="!$auth.isLoggedIn()" to="/auth" class="item">Login</router-link>
        <router-link v-if="$auth.isLoggedIn()" to="/auth" class="item">Account</router-link>
        <router-link to="/page" class="item">API page</router-link>
      </div>
    </div>
    <div class="ui main text container">
      <router-view></router-view>
    </div>
  </div>
</template>

<script>
export default {
  name: 'app'
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  overflow: hidden;
}

.ui.menu .item img.logo {
  margin-right: 1.5em;
  width: 40px;
}

.main.container {
  margin-top: 7em;
}
</style>
