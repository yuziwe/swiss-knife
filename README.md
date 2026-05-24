# swiss-knife

All kinds of tools improve your coding experience.

## Introduction

~~~shell
.
├── LICENSE
├── Makefile
├── README.md
├── openwrt-ruleset-update # Update ruleset for openwrt route system
└── terminal-translator # Translator driven by LLM API, modify `SYSTEM_PROMPT` to fit your flavor:)
~~~

## Prequisties

- Go >= 1.25.0
- GNU make
- fish >= 4.0.0

## Usage

### terminal-translator

~~~sh
# build
make
# Install, default in /usr/local/bin
# You can specify the location by passing `prefix` variable e.g `make prefix=~/opt/ install`
make install
# help
make help
~~~

### openwrt-ruleset-update

> [!NOTE]
>
> DO NOT FORGET TO GENERATE SSH KEY PAIR FOR COPY FILES SMOOTHLY

~~~sh
make openwrt-ruleset-update
~~~

Enjoy it!

