syncr
=====

hack to run keep dotfiles in sync on multiple hosts (single master)

warning: runs `rsync` with the `--delete` option!

installation
------------

    go get github.com/chlunde/syncr/cmd/syncr

usage
-----

Create a `syncr.yaml`:

```yaml
-
  source: ~/.vim/
  destination: ~/.vim
  hosts:
      - login.foo.net
      - other.com
-
  source: ~/dotfiles/
  destination: ~/dotfiles
  hosts:
      - login.foo.net
      - other.com
```

TODO
-------

1. add support for `rsync` options (exclude etc.)
