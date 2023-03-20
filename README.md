# cman

cman is a programming mannual with ChatGPT

## Install
```shell
go install github.com/yah01/cman@latest
```

## Usage
First you need to set the env var `OPENAI_API_KEY` to your OpenAI API key.

Then:
```shell
cman golang waitgroup
```

sometimes we need just a tl;dr:
```shell
cman -s golang waitgroup
```
which provide a short answer.

If you want a short answer but with an example:
```shell
cman -s -e golang waitgroup
```