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
</style>
