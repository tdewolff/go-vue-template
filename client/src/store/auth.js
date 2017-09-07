import jwtDecode from 'jwt-decode'

export default {
  state: {
    jwt: localStorage.getItem('jwt')
  },
  mutations: {
    LOGIN: (state, jwt) => {
      state.jwt = jwt
    },
    LOGOUT: (state) => {
      state.jwt = null
    }
  },
  actions: {
    login ({commit}, jwt) {
      localStorage.setItem('jwt', jwt)
      commit('LOGIN', jwt)
    },
    logout ({commit}) {
      localStorage.removeItem('jwt')
      commit('LOGOUT')
    }
  },
  getters: {
    isLoggedIn: state => {
      return !!state.jwt
    },
    getUser: state => {
      if (state.jwt) {
        return jwtDecode(state.jwt)
      }
      return null
    }
  }
}
