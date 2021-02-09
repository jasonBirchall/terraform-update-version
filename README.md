# terraform-update-version
Go application that runs the terraform version upgrade checker binary and creates a pull request

## Batch run tf 0.13 command
```
$ find . -name '*.tf' | xargs -n1 dirname | uniq | xargs -n1 terraform 0.13upgrade -yes

```

## Git clone, branch and commit
https://github.com/go-git/go-git/tree/master/_examples

## Pull request
https://github.com/google/go-github/blob/master/example/commitpr/main.go
