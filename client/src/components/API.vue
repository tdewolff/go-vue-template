<template>
  <div class="api">
    <h1 class="ui header">API request</h1>
    <div class="columns medium-3" v-for="calendar in calendars">
      <div class="card">
        <div class="card-divider">
          {{ calendar.Id }}
        </div>
        <div class="card-section">
          <p>{{ calendar.Name }}.</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'hello',
  data () {
    return {
      msg: 'Welcome to Your Vue.js App',
      calendars: [],
      errors: []
    }
  }
  // created () {
  //   this.axios.get('http://localhost:3000/calendars')
  //   .then(response => {
  //     // JSON responses are automatically parsed.
  //     this.calendars = response.data
  //   })
  //   .catch(e => {
  //     this.errors.push(e)
  //   })
  // }
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
