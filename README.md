# tkey-led

Example application for TKey written using TinyGo.

## How to run

First compile and flash the TKey with the application:

```shell
make app
```

The LED should start blinking green every half second.

Now you can run the command line client application on your computer:

```shell
go run ./cmd/tkeyled --led 2
```

The LED should now be blinking blue.
