# Go Vue template <a name="go-vue-template"></a>

Go-Vue-template is a project that serves as a template for future web applications involving a Go backend (API server) and Vue frontend (SPA client). It has the following features:

* Vue for a Single Page Application frontend structure
* OAuth2 authentication using social logins
* Semantic-UI for interface design
* Webpack to bundle files

This package will be updated as I go. Incomplete todo list:

* Automatic logout using Vuex and timeout
* Logout blacklists user at the server
* Automatic OAuth token refresh
* Able to add more scopes to OAuth providers
* 401 Unauthorized vs 403 Forbidden in server responses; which is more appropriate? 401 if we want the user to retry logging in

## License
Released under the [MIT license](LICENSE.md).

[1]: http://golang.org/ "Go Language"
