# Fast CGI

For example purposes (and for fun), kram can be deployed with FastCGI behind for example nginx. 

## Requirements

OS X: `brew install fcgi spawn-fcgi`
Debian / Ubuntu: `apt-get install libfcgi-dev spawn-fcgi`

## Building

```bash
make fcgi
```

## Running

```bash
spawn-fcgi -p 8000 -n bin/fcgi
```

## Configuring nginx

You could then configure nginx to look something like this.

```
server {
    listen 8888 default_server;
    root /var/www;

    server_name example.com;

    location / {
        fastcgi_pass 127.0.0.1:8000;
        fastcgi_index index.kr;
        fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
        include fastcgi_params;
        fastcgi_split_path_info ^(.+\.php)(.*)$;
        try_files /$uri /index.kr?$args;
    }
}
```

Aaaaand, hopefully that would work. :)
