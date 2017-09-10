# Go Vue template <a name="go-vue-template"></a>

Go-Vue-template is a project that serves as a template for future web applications involving a Go backend (API server) and Vue frontend (SPA client). It has the following features:

* Vue for a Single Page Application frontend structure
* OAuth2 authentication using social logins
* Semantic-UI for interface design
* Webpack to bundle files

This package will be updated as I go. Incomplete todo list:

* Use TLS for everything
* Automatic logout using Vuex and timeout
* Logout blacklists user at the server
* Automatic OAuth token refresh
* Able to add more scopes to OAuth providers
* 401 Unauthorized vs 403 Forbidden in server responses; which is more appropriate? 401 if we want the user to retry logging in
* Move client authorization code into a plugin
* Use github.com/markbates/goth for OAuth providers

## Install
Make sure you have the [Go compiler](https://golang.org/dl/) and [NPM](https://www.npmjs.com/get-npm) installed.

``` bash
# Download and install this package
go get github.com/tdewolff/go-vue-template

# Download and install client NPM packages
cd client
npm install
```

## Run
### Development mode
Server in dev mode will enable CORS headers from localhost:8080. Running the Vue client will automatically open the browser at http://localhost:8080/
``` bash
# Start API server at :3000
go-vue-template --dev

# Start Vue client with hot-reloading at :8080
cd client
npm run dev
```

### Production mode
``` bash
# Build the client files to client/dist/
npm run build

# Start the webserver at :3000
go-vue-template
```

Now navigate to http://localhost:3000/. You can change the port in the config.json file.

``` bash
# Analyze filesizes
npm run build --report
```

For detailed explanation on how things work at the client, checkout the [guide](http://vuejs-templates.github.io/webpack/) and [docs for vue-loader](http://vuejs.github.io/vue-loader).

## Ready-up for your application

* Set name and host in config.json
* Set title in client/index.html
* Set JWT secret in main.go
* Add social login client IDs and secrets to config.json
* Create your database scheme in scheme.sql
* Develop your API server and client interface

## License
Released under the [MIT license](LICENSE.md).

[1]: http://golang.org/ "Go Language"
