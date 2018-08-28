<template>
  <div class="api">
    <h1 class="ui header">API request</h1>
    {{ msg }}
  </div>
</template>

<script>
import * as util from '../util'

export default {
  name: 'hello',
  data () {
    return {
      msg: ''
    }
  },
  mounted () {
    this.axios.get('http://localhost:3000/api/user/')
    .then(response => {
      this.msg = response.data
    }, e => {
      util.showError(e)
    })
  }
}
</script>
