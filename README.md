# tinygo-tkey

Develop applications for the Tillitis TKey-1 using TinyGo. 

It includes an implementation of the Tillitis framing protocol for communication (https://dev.tillitis.se/protocol/) between the device and client that can be used for applications that run on the TKey written using TinyGo. 

## Examples

![tkey led](./images/tkey-led.gif)

Example application for TKey written using TinyGo for the device application and Go for the client application.

## Device application

To compile and flash the TKey with the device application:

```shell
tinygo flash -size short -target=tkey ./examples/blinker/app
```

The LED should start blinking green every half second.

## Client application

Now you can run the command line client application on your computer:

```shell
go run ./examples/blinker/cmd --led 0 --timing 250
```

The LED should now be blinking blue every 250 ms.
