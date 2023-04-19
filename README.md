# cbot-telegram-whisper
Quick and dirty bot (or more of a PoC) that transcribes Telegram voice messages using OpenAI's whisper models, using CPU inference thanks to [whisper.cpp](https://github.com/ggerganov/whisper.cpp) 

## Built using: 
* [go-telegram-bot-api](https://pkg.go.dev/github.com/go-telegram-bot-api/telegram-bot-api/v5) 
* [whisper.cpp Go bindings](https://github.com/ggerganov/whisper.cpp/tree/master/bindings/go/)
* [u2takey's ffmpeg-go](https://github.com/u2takey/ffmpeg-go)
* [grab](https://github.com/cavaliergopher/grab) 

## TODO
* refactor
* possibly get rid of ffmpeg lib? 
* fix error handling & generating

## How to build
#### You can check out the example [Dockerfile](Dockerfile), but TL;DR:

```bash
git clone git@github.com:chinese-soup/cbot-telegram-whisper.git && \
cd cbot-telegram-whisper && \
git clone https://github.com/ggerganov/whisper.cpp.git && \
cd whisper.cpp/bindings/go && \
make whisper && \
cd ../../.. && \
go get
C_INCLUDE_PATH=/app/whisper.cpp/ LIBRARY_PATH=/app/whisper.cpp/ go build -o whisperbot
```

#### Get a model, e.g.:
```bash
bash whisper.cpp/models/download-ggml-model.sh tiny.en
```
Check out [whisper.cpp's](https://github.com/ggerganov/whisper.cpp) README for more info.


#### Set the environment variables mentioned in `.env.example`, e.g.:

```bash
export TELEGRAM_APITOKEN=<your bot token here>
export MODELPATH=whisper.cpp/models/ggml-tiny.en.bin
```
#### Make sure ffmpeg is in $PATH:
```bash
$ which ffmpeg
/usr/bin/ffmpeg
```
If it isn't install a binary or compile and put it in path.

#### And you're all set, run it:
```bash
./whisperbot
```