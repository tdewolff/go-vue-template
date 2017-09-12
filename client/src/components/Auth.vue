<template>
  <div class="auth">
    <h1 class="ui header">Social login</h1>
    <div v-if="$auth.isLoggedIn()">
      <button class="ui button blue" @click="logout()">Logout</button>
      <br>
      {{ profile }}
    </div>
    <div v-if="!$auth.isLoggedIn()">
      <p v-if="$route.query.referer">Please log in first to continue.</p>
      <p v-if="$route.query.error">An error occurred, please try again.</p>
      <div class="ui stacked social-buttons">
        <a v-for="provider in providers" :href="provider.URL" :class="'ui labeled icon button fluid ' + provider.id"><i :class="'icon ' + provider.id"></i>{{ provider.name }}</a>
      </div>
    </div>
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
    this.$watch(() => this.$auth.store.getters.jwt, (isLoggedIn) => {
      this.loadSocialAuths()
    })
  },
  methods: {
    loadSocialAuths () {
      if (!this.$auth.isLoggedIn()) {
        var referer = this.$route.fullPath
        if (this.$route.query.referer) {
          referer = this.$route.query.referer
        }

        this.$auth.getAuthURLs(referer)
        .then(providers => {
          this.providers = providers
        }, e => {
          console.log(e)
        })
      }
    },
    logout () {
      this.$auth.logout()
    }
  },
  computed: {
    loggedIn () {
      return this.$auth.isLoggedIn()
    },
    profile () {
      return this.$auth.getUser()
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
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
