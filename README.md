shamir
======

A simple command utility that can perform operations using [Shamir's Secret Sharing algorithm](https://en.wikipedia.org/wiki/Shamir's_Secret_Sharing).

## Usage

`shamir` takes and receives base64 encoded key shares.

```
$ shamir -h
Usage: shamir [<options>] (split|combine)
  -in string
    	input file (default "-")
  -out string
    	output file (default "-")
  -parts int
    	number of shares for split operation (default 3)
  -threshold int
    	threshold for split operation (default 2)
```

## Examples


Splitting a secret

```
$ cat secret
I've never been to Spain.
$ shamir -in secret -threshold 5 -parts 8 split
jZYqqJaCxt4kZjh/xCAakXIG0XKdnlxY89/S
1yubXvs8GdSqO8GPX9Xwk4fV22MJ/4F8wxh9
fuVQmRhkVeKQGzpPVT0eNwvLXHxrQmQjFLGH
Yggxmz2bmtTNPAX9ctHffwqsw4Ds8GZHDpa1
Nyns0pAKCfeqLCrzxV+TpoJCvDJKTqJmR2kG
a//rC1IRcRow60BY5SQV+6bZXLFbUxt221zy
uB34M0pUB5vljgIfdItDHetEfKqxDQea6A93
OUQe8CwZBHTbwt0Za/IRLUlLfB39smrEIEpt
```

Combing shares

```
$ cat shares
upOiLkD7A/va03FStBD+NV9feBlKooE8imH7JJ/CXbbV3pQ3EHxZCQ3z
dfJ2ID7zuNNQY+Zqkqvhu2QCBzVLJyuxg4nzT0rSmYrrbR5RBQNPiOVD
+p/3gvuk9sNTjR7UIOVJHi13ZoPvIm/jcBCkJm98tpKgdB1tqWIUK3yt
$ head -n 2 shares | shamir combine
Three can keep a secret if two are dead.

```

## Thanks

Thanks to Hashicorp for implementing the actual algorithm and publishing it as part of the [source code for Vault](https://github.com/hashicorp/vault/tree/master/shamir).
