# tkey-led

Example application for TKey written using TinyGo for the device application and Go for the client application.

It uses the Tillitis framing protocol for communication (https://dev.tillitis.se/protocol/) between the device and client.

## Device application

To compile and flash the TKey with the device application:

```shell
tinygo flash -size short -target=tkey ./app/blinker
```

The LED should start blinking green every half second.

## Client application

Now you can run the command line client application on your computer:

```shell
go run ./cmd/tkeyled --led 0 --timing 250
```

The LED should now be blinking blue every 250 ms.
