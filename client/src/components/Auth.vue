<template>
  <div class="auth">
    <h1 class="ui header">Login</h1>
    <div class="ui buttons">
      <button class="ui button blue" @click="login()">Login</button>
      <button class="ui button blue" @click="register()">Register</button>
      <button class="ui button blue" @click="logout()">Logout</button>
    </div>
    <div class="ui divider"></div>
    <h1 class="ui header">Social login</h1>
    <div class="ui stacked social-buttons">
      <a v-if="!loggedIn" v-for="provider in providers" :href="provider.URL" :class="'ui labeled icon button fluid ' + provider.id"><i :class="'icon ' + provider.id"></i>{{ provider.name }}</a>
    </div>
    {{ profile }}
  </div>
</template>

<script>
export default {
  name: 'auth',
  data () {
    return {
      providers: []
    }
  },
  mounted () {
    this.loadSocialAuths()
    this.$store.watch(() => this.$store.getters.isLoggedIn, (loggedIn) => {
      this.loadSocialAuths()
    })
  },
  methods: {
    loadSocialAuths () {
      if (!this.$store.getters.isLoggedIn) {
        var referer = encodeURIComponent(this.$route.fullPath)
        this.$http.get('http://localhost:3000/auth/list?referer_uri=' + referer)
        .then(response => {
          localStorage.setItem('session_id', response.data.sessionId)
          this.providers = response.data.providers
        })
        .catch(e => {
          console.log(e)
        })
      } else {
        this.auths = []
      }
    },
    logout () {
      this.$store.dispatch('logout')
    }
  },
  computed: {
    loggedIn () {
      return this.$store.getters.isLoggedIn
    },
    profile () {
      return this.$store.getters.getProfile
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
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

.auth {
  width: 30em;
  margin: 0 auto;
}

.social-buttons {
  width: 10em;
  margin: 0 auto;
}

.ui.button {
  margin-bottom: 0.5em;
  text-align:;
}

.ui.google.button {
  background-color:#4285F4;
  color:#FFFFFF;
}

.ui.github.button {
  background-color:#444444;
  color:#FFFFFF;
}
</style>
