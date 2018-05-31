# Ergo [![](https://images.microbadger.com/badges/image/ontrif/ergo.svg)](https://microbadger.com/images/ontrif/ergo) [![](https://images.microbadger.com/badges/version/ontrif/ergo.svg)](https://microbadger.com/images/ontrif/ergo)
<img align="left" src="http://i61.beon.ru/67/83/108367/24/4083124/ergo_proxy10.gif">

Simple http proxy for accessing virtual domains such as `*.test` and `*.local` and exposing microservices as single api during development.

The proxy can be used to organize local projects. You can create your own domain name for each local project and simply add it to the ergo config. If you configure ergo as a web proxy in your browser, this will be enough and the `/etс/hosts` will not need to be edited.

Name and idea are inspired by [cristianoliveira/ergo](https://github.com/cristianoliveira/ergo) project and the "Ergo Proxy" anime.

My project is simpler, but it supports redirects based on regular expressions on the full request URL, which allows you to expose different microservices on single domain as unioned api.

## Features
* can be used as usual reverse proxy (`/etc/hosts` for virtual domains) or http-proxy (no editing is needed for the `/etc/hosts` file).
* websockets support
* simple config based on regexps

## Configuration
The main idea of this proxy is proxying http traffic to different backends based on URL.
It is very similar to plain reverse proxy but in additional to that ergo can be used as usual http-proxy for browser.

Example ergo config for independent frontend (`localhost:8080`) and backend api (`localhost:7777`) which are implemented as separate projects can be as follows:
```
some-site.local/api --> localhost:7777
some-site.local --> localhost:8080 as www.real-name.com
```

For browser configured to use http proxy `localhost:2000` (`2000` is default ergo's listen port) we can access this URL's:
1. `http://some-site.local/some/url` ergo will send request to `http://localhost:8080/some/url` with `Host: www.real-name.com` header.
2. `http://some-site.local/api/some/url` ergo will send request to `http://localhost:7777/api/some/url`

**NOTE:** order of rules in config important! Each regexp will be searched in request URL and first matched rule will be used!

For real-world usage consider to use smart-proxy-selector plugins for your browser (such as FoxyProxy, SwitchyOmega, friGate etc.) Configure ergo proxy usage only for your local dev domains (such as `*.local`, `*.test` etc.)

**NOTE:** don't use `*.dev` domain due to browsers builtin HSTS rules ([read more here](https://stackoverflow.com/a/47768411)).

## Motivation
I want easy-configurable regexp-based proxy which can be configured to merge different microservices at different entry-points
on the same domain. Also I want source code for this proxy to be as simple as possible. And of course writing ergo was fun.

## TODO
1. `CONNECT` method code improvement
2. config auto-reload
3. https support
4. port configuration
5. nginx-style regexp url overwrites via regexp matching groups
