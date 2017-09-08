<template>
  <div class="authCallback">
    <h1 class="ui header">Logging in...</h1>
  </div>
</template>

<script>
export default {
  name: 'authCallback',
  created () {
    var state = encodeURIComponent(this.$route.query.state)
    var code = encodeURIComponent(this.$route.query.code)
    var sessionId = encodeURIComponent(localStorage.getItem('session_id'))
    this.$http.get('http://localhost:3000/auth/token?state=' + state + '&code=' + code + '&session_id=' + sessionId)
    .then(response => {
      this.$store.dispatch('login', response.data.jwt)
      this.$router.push(response.data.referer)
    })
    .catch(e => {
      console.log(e)
      // this.$router.push('/auth')
    })
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h1, h2 {
  font-weight: normal;
}
</style>
