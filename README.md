# Traduction Translation aka T2

Check how your test resists to double translation.

I’m French, so when I write in english, I sometimes use double translation to check how the translated english "sounds".
On a translation software, I write directly in English, I translate into French and translate back again into English.  
The name "Traduction Translation" comes from there, "traduction" being "translation" in French.

As a developer, I wanted to try to automate this thing and share it with you.  
Just copy some text, then run `t2 clipboard`, and "voilà".  
I hope it will help you.

Yes, this paragraph was checked with t2 :)

## Installation

```shell
git clone https://github.com/rangzen/t2
cd t2
go install
```

## Translation services

### DeepL

The actual default, and only, service is for translation is [DeepL](https://deepl.com).  
You’ll need a Pro free account cause the free account is almost out of limits.  
I don’t have a Pro paid account but I think that you just have to change the Endpoint configuration.

#### Configuration

* Create a `.t2.yaml` file configuration with:

```yaml
Endpoint: https://api-free.deepl.com/v2/translate
ApiKey: redacted-0123-0123-0123-redacted:fx
```

You can also see the `t2-example.yaml` file for an example.
