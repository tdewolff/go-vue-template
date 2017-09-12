<template>
  <div class="api">
    <h1 class="ui header">API request</h1>
    {{ msg }}
  </div>
</template>

<script>
export default {
  name: 'hello',
  data () {
    return {
      msg: '',
      errors: []
    }
  },
  mounted () {
    this.axios.get('http://localhost:3000/api/')
    .then(response => {
      this.msg = response.data
    })
    .catch(e => {
      this.errors.push(e)
    })
  }
}
</script>
