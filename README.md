# Wiktionary IPA Extractor

This Go application processes Wiktionary data dumps to extract the International Phonetic Alphabet (IPA) transcriptions of words. It parses the latest Wiktionary dump, processes the data, and outputs the IPA transcriptions in a JSONL format. The project also includes a nightly job that checks for new Wiktionary dumps and automatically processes them if available.

## Features

- **Automated Data Extraction**: Extracts IPA transcriptions from Wiktionary dumps.
- **JSONL Output**: Outputs IPA transcriptions in a JSONL format with one IPA transcription per line.
- **Nightly Checks**: Automatically checks for new data dumps from Wiktionary and processes them if new data is available.
- **GitHub Release Integration**: Publishes processed data as a release on GitHub, allowing easy access to the latest IPA data.

## Most users do not need to install anything

If you simply want the latest processed IPA data, you can head over to the Releases section of this repository and download the latest JSONL file. No need to install or build the application unless you want to process the data yourself.

## Output Format

The output is in **JSONL format** (JSON Lines), meaning that **each line contains one IPA transcription** for a word. Here's an example of what the output looks like:

```json
{"word":"aardvark","ipa":[{"lang":"en","IPA":["/ˈɑːd.vɑːk/"],"variant":"a=RP"}]}
{"word":"month","ipa":[{"lang":"en","IPA":["/mʌnθ/"]},{"lang":"en","IPA":["/mʌnθ/"]}]}
```

Each entry contains:
- `word`: The word being transcribed.
- `ipa`: A list of IPA transcriptions for different languages or pronunciation variants (e.g., RP for Received Pronunciation, US for American English).
- `lang`: The language code (e.g., "en" for English).
- `IPA`: A list of phonetic transcriptions.
- `variant`: Optional, specifies pronunciation variant (e.g., "RP", "US").

## Installation

1. Install the application using Go:
   ```bash
   go install github.com/terwey/wiktionary-ipa-extract/cmd/wiktionary-ipa-extract@latest
   ```

2. Ensure Go is properly set up and available in your system’s PATH.

## Usage

### Using the Latest Wiktionary Dump

To process the latest Wiktionary dump and extract IPA transcriptions, run the following command:

```bash
curl https://dumps.wikimedia.org/enwiktionary/latest/enwiktionary-latest-pages-articles.xml.bz2 | ./wiktionary-ipa-extract --bz -o ipa-data.jsonl
```

This will download the latest dump, process it, and output the extracted IPA transcriptions into a JSONL file.

### Processing a Different Wiktionary Dump

If you want to process a specific dump (or one downloaded manually), you can provide the file directly using the `--input` flag. Here’s how to do it:

1. Download the Wiktionary dump you wish to process from [Wiktionary Dumps](https://dumps.wikimedia.org/enwiktionary/).
2. Run the command below, replacing `<dump-file>` with the path to your downloaded file and `<output-file>` with the desired output filename.

   ```bash
   ./wiktionary-ipa-extract --input <dump-file> --bz -o <output-file>.jsonl
   ```

   Example:
   ```bash
   ./wiktionary-ipa-extract --input ~/downloads/enwiktionary-20230101-pages-articles.xml.bz2 --bz -o ipa-data.jsonl
   ```

## Automation with Nightly Runs

This project uses GitHub Actions to automate nightly runs that check for new Wiktionary dumps via an [RSS feed](https://dumps.wikimedia.org/enwiktionary/latest/enwiktionary-latest-pages-articles.xml.bz2-rss.xml). If a new dump is found, it is automatically processed, and the extracted IPA data is released on GitHub in JSONL format.

## Contributing

Contributions are welcome! Please submit a pull request or open an issue to discuss improvements or report bugs.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.
