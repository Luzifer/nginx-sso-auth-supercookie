# Luzifer / nginx-sso-auth-supercookie

`auth-supercookie` is a plugin for the `nginx-sso` plugin system providing an additional authentication method through a "supercookie": An externally set cookie with a specified name and a specified content. This cookie then serves as a login token for nginx-sso.

The purpose of writing this was on the one hand to have a proof-of-concept for the nginx-sso plugin system and on the other hand to have a persistently signed in device being able to access my own tools being protected by nginx-sso. This way I don't need to re-sign-in after a session timeout on that device which is not always possible.

## Warning

Though I think it's obvious: If you are using this plugin you need to be **very careful** with the device having that super-cookie. In case anyone gets hold of it they are able to impersonate your device. This basically weakens the SSO you set up with `nginx-sso`! Triple-check your cookie settings to ensure that cookie is **never** delivered to any third-party site and only used on secure connections.

Basically: If possible avoid this plugin at all. Because of these security considerations this will never be part of the core-plugins in nginx-sso. You will always need to install it manually after building it yourself.

## Building / installing

To build the plugin you need to [use the same GOPATH as while compiling nginx-sso](https://github.com/golang/go/issues/26759). If you are running nginx-sso from the container I am providing, you can build the plugin just using the `make build` command. Further on in this README I am assuming the Docker version of nginx-sso.

To install the plugin create a directory `plugins` in your data mount (next to your `config.yml`) and copy the `auth-supercookie.so` into that directory. As it is now available to nginx-sso you need to activate plugin loading:

```yaml
plugins:
  directory: /data/plugins
```

And after all you also need to configure the supercookie plugin (it uses the same format as the token plugin with one additional setting):

```yaml
providers:
  supercookie:
    cookie_name: <name of your secret cookie>
    tokens:
      username: verysecretcookiecontent
    groups:
      admins: username
```

After having set this up you should be able to restart `nginx-sso` and it should detect the supercookie when set and like the token provider just skip the login step.

To set the cookie you for example can use the Chrome dev-tools:

```javascript
d = new Date()
d.setTime(d.getTime() + 100 * 365 * 86400)
document.cookie = `secret_cookie_name=verysecretcookiecontent; expires=${d.toUTCString()}; secure=true; domain=yourdomain.local; path=/`
```

This will set a cookie with an approximate 100 years life-span.
