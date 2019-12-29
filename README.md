![project_logo](https://user-images.githubusercontent.com/28186014/71556151-f973c000-2a34-11ea-9902-3c1c658b5abe.png)

# rest2sacn
## Translator to send DMX Data over sACN via a simple HTTP Request.
This is a simple translator for devices that are not able to send sACN Data on their own. If you want to send some simple DMX Data over the network you can just send an HTTP Request with a super simple syntax to this API and the server will send out the corresponding DMX Data.

This project currently only supports sACN Unicast. It is not intended for high throughput applications (use a native sACN device for this). If you just want to control a single channel, this is fine. Latency is also not too bad.

## Concept

The original idea that lead to the creation of this project was the following: I needed a way to press some executors on my grandMA2 Lighting Desk via my Win10 Computer. I wanted to press an executor on my MA2, if I press one key on my computer. Because I was heavily invested into the AutoHotKey ecosystem, I just wanted to integrate this thing into one of my existing macros. The Simplest way was to do this via sending an HTTP Request from the AutHotKey Macro. This results in the following workflow:

```
                     +----------------------+  UDP sACN     +-------------------------+
HTTP Request  +----> |      rest2sacn       | +-----------> |      Lighting Desk or   |
                     |  Translating.....    |               |      Lamp, DMX Fixture  |
                     +----------------------+               +-------------------------+
```

## Config File
The server is configured via a YAML config file. An example file is included and looks like the following: 
```
# Simple config to control which universes should be translated/sent.

universe: [69,70]   # Universes that should be available to control.
ip: "127.0.0.1:8080" # On which IP should rest2sacn listen?
destination: "127.0.0.1"    # IP of the device the sACN data should be transmitted.
```
The configuration parameters are pretty self explanatory. Here you should list all universes that you want to control via the REST API. If a universe is not listed in the config file it can't be controlled via the API. 

## API
The REST API of this project is very straight forward. It does not use any sort of authentication as it is intended to be run only on localhost or on protected Networks. 

These are the following endpoints: 
```
/sacn/reset/{universe}
/sacn/send/{universe}/{channel}/{value}
```

### /sacn/send/{universe}/{channel}/{value}
This is the Endpoint that is responsible for sending all the data. It takes the 3 typical DMX Parameter. They should be in the following number range: 
- universe: 1-63999
- channel: 1-512
- value: 0-255

If you send a request with these parameters the program immediately sets the corresponding channel to the sent value. The data doesn't change till the next requests is received. 

### /sacn/reset/{universe}
Takes an Universe Number and sets the output of the universe to `0` on all Channels

## Copyright / Stuff 
This Project is under the MIT License by @localleon 2019. Contributions or Questions are welcome. 
