import Vue from 'vue'

export function showSuccess (msg) {
  Vue.toasted.success(msg)
}

export function showError (e) {
  if (!e.response || !e.response.data) {
    console.log('Error: ' + JSON.stringify(e))
  } else if (!e.response.data.error) {
    console.log('Error: ' + JSON.stringify(e.response.data))
  } else {
    console.log('Error: ' + e.response.data.error)
  }

  if (!e.response || e.response.status !== 401) {
    let msg = 'Something went wrong'
    if (typeof e === 'string') {
      msg = e
    } else if (e.message === 'Network Error') {
      msg = 'Could not reach server'
    }
    Vue.toasted.error(msg)
  }
}
