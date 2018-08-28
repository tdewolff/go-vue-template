# Go Vue template <a name="go-vue-template"></a>

Go-Vue-template is a project that serves as a template for future web applications involving a Go backend (API server) and Vue frontend (SPA client). It has the following features:

* Vue for a Single Page Application frontend structure
* OAuth2 authentication using social logins
* Authentication using JWT
* Semantic-UI for interface design
* Webpack to bundle files

This package will be updated as I go. Incomplete todo list:

* Use TLS for everything
* Automatic logout using Vuex and timeout
* Logout blacklists user at the server

## Install
Make sure you have the [Go compiler](https://golang.org/dl/) and [NPM](https://www.npmjs.com/get-npm) installed. Then issue:

``` bash
# Download and install this package
go get github.com/tdewolff/go-vue-template

# Download and install client NPM packages
cd web/
npm install
```

## Run
### Development mode
Server with `DevURL` set in `config.json` will enable CORS headers from localhost:8080. Running the Vue client will automatically open the browser at http://localhost:8080/
``` bash
# Start API server at :3000
go-vue-template

# Start Vue client with hot-reloading at :8080
cd web/
npm run dev
```

### Production mode
Remove `DevURL` from `config.json`, and run:

``` bash
# Build the web files to web/dist/
npm run build

# Start the webserver at :3000
go-vue-template
```

Now navigate to http://localhost:3000/. You can change the port in the `config.json` file.

``` bash
# Analyze filesizes
npm run build --report
```

For detailed explanation on how things work at the client, checkout the [guide](http://vuejs-templates.github.io/webpack/) and [docs for vue-loader](http://vuejs.github.io/vue-loader).

## Ready-up for your application

* Rename config.json.dist to config.json and set name and URLs
* Add social login client IDs and secrets to config.json
* Create your database scheme in scheme.sql
* Set title in web/index.html
* Develop your API server and web interface

## License
Released under the [MIT license](LICENSE.md).

[1]: http://golang.org/ "Go Language"
