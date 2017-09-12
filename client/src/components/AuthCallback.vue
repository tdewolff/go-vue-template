<template>
  <div class="authCallback">
    <h1 class="ui header">Logging in...</h1>
  </div>
</template>

<script>
export default {
  name: 'authCallback',
  created () {
    this.$auth.login(this.$route.query.state, this.$route.query.code)
    .then(referer => {
      this.$router.push(referer)
    }, e => {
      this.$router.push('/auth?error=' + encodeURIComponent(e.response.status))
    })
  }
}
</script>
