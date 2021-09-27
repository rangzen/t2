# Traduction Translation aka T2

Check how your test resists to double translation.

I’m French, so when I write in english, I sometimes use double translation to check how the translated english "sounds".
On a translation software, I write directly in English, I translate into French and translate back again into English.  
The name "Traduction Translation" comes from there, "traduction" being "translation" in French.

As a developer, I wanted to try to automate this thing and share it with you.  
Just copy some text, then run `t2 clipboard`, and "voilà".  
I hope it will help you.

Yes, documentation was checked with t2 :)

## Examples

In a terminal, the diff version will appear in color:

![diff screenshot](https://raw.githubusercontent.com/rangzen/t2/main/doc/Screenshot_20210925_diff.png)

Don't forget the `--only-diff` or `-d` flag if you want to display only this part.

### Translate from CLI

```shell
$ t2 "I want speak english."
Using config file: /home/user/.t2.yaml
# Original text
I want speak english.
# Pivot text
Je veux parler anglais.
# Double translated text
I want to speak English.
# Diff version
I want to speak eEnglish.
```

### Translate from the clipboard

Don't bother with copy/paste operations, quoting text, etc. Just copy what you want to check and then `t2 clipboard`.

```shell
$ t2 clipboard
Using config file: /home/user/.t2.yaml
# Original text
Some text this was in clipboard.
# Pivot text
Certains textes étaient dans le presse-papiers.
# Double translated text
Some texts were in the clipboard.
# Diff version
Some text this wasere in the clipboard.
```

### Usage

```shell
$ t2 usage
Using config file: /home/user/.t2.yaml
Usage: 12477/500000
```

## Installation

```shell
git clone https://github.com/rangzen/t2
cd t2
go install -ldflags "-X github.com/rangzen/t2/cmd.Version=`git tag --sort=-version:refname | head -n 1`"
```

## Translation services

#### Configuration

* Create a `.t2.yaml` file configuration with:

```yaml
TranslationServices:
  DeepL:
    Endpoint: https://api-free.deepl.com/v2/translate
    ApiKey: redacted-0123-0123-0123-redacted:fx
  Google:
    Endpoint: https://translation.googleapis.com/language/translate/v2
    ApiKey: redactedredactedredacted
```

See the `t2-example.yaml` file for an example.

### DeepL

The actual default service for translation is [DeepL](https://deepl.com).  
You’ll need a Pro free account because the free account is almost always out of limits.  
I don’t have a Pro paid account, but I think that you just have to change the Endpoint configuration.

### Google Cloud Translation

For using [Google Cloud Translation](https://cloud.google.com/translate/), you need:
* a Google Cloud account,
* a Paid account,
* an API key without restriction.

"The `usage` command doesn't work with Google!"  
I know. If you know the API endpoint for usage, please let me know.
