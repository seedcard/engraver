# SeedCard engraver

Controller program for engraving **SeedCard** plates: credit-card sized (SH01) steel plates
carrying the 12 words of a Bitcoin seed, on a [SeedHammer](https://seedhammer.com) engraving
machine.

It's Seed Hammer **v1.2.1**, the last upstream release that supports the credit-card plate
(upstream removed it in v1.3.0), with one small SeedCard change on top:

- the seed side carries the **12 words only**, without the SeedQR
- the version mark reads **SC** and the key number reads **1/1**
- that text (mark, title, key number, SeedQR on/off) is set in one file,
  [`seedcard/text.conf`](seedcard/text.conf)

The seed words, word numbers and master fingerprint always come from the seed. Everything else,
including the descriptor (wallet) side and how seeds are entered, is unchanged from upstream.
[SEEDCARD.md](SEEDCARD.md) has the details.

To see every line that differs from upstream:

```sh
git clone https://github.com/seedcard/engraver && cd engraver
git remote add upstream https://github.com/seedhammer/seedhammer
git fetch upstream tag v1.2.1
git diff v1.2.1 seedcard --
```

## Hardware

The same as a [SeedSigner](https://seedsigner.com/hardware): Raspberry Pi Zero or Zero W, a Waveshare 1.3 inch 240×240 LCD hat, and a Pi Zero camera with an OV5647 sensor,
connected by USB to the SeedHammer engraver.

## How a plate is made

1. On a second SeedSigner, show the seed as a SeedQR. For a multisig wallet, also have the wallet
   descriptor ready as a QR (from Sparrow or another coordinator).
2. Scan the QR codes with the SeedCard engraver and confirm the words on screen.
3. Place an SH01 plate in the engraver and start engraving.

To test without using a plate, long-press the middle button on the engraving screen to turn on
dry-run mode.

## Getting an SD-card image

### Download from GitHub (no tools needed)

Every push to this repository builds the image on GitHub:

1. Open the [Actions tab](https://github.com/seedcard/engraver/actions), choose **Build image**,
   then the newest run of the `seedcard` branch with a green tick.
2. Download the artifact `seedhammer-<version>.img` at the bottom of the page. It contains the
   image and its SHA-256.
3. Unzip it, check the hash (`sha256sum seedhammer-*.img`, or `shasum -a 256` on macOS),
   and write the image to an SD card.

Artifacts are kept for 90 days. To rebuild an older version, open the Build image workflow and
use **Run workflow**.

### Build it yourself

You need [Nix](https://nixos.org/) with flakes enabled:

```sh
nix build github:seedcard/engraver/seedcard
ls result/seedhammer.img
```

For a versioned image, tag the commit (for example `sc-1.0`) and run:

```sh
nix run github:seedcard/engraver/seedcard#mkrelease sc-1.0
```

The tag must be pushed to GitHub first; the version is written into the image.

### Reproducible builds

The build is deterministic: anyone building the same commit gets the same image, bit for bit.
Record the SHA-256 of each image you use for engraving, so you can always tell which version
made a plate.

## Writing the image to an SD card

Linux (replace `/dev/sdX` with your card, checking carefully first):

```sh
dd if=seedhammer-<version>.img of=/dev/sdX bs=1M
```

macOS and Windows: use [balenaEtcher](https://etcher.balena.io/) or Raspberry Pi Imager
("Use custom").

## Changing the engraved text

1. Edit [`seedcard/text.conf`](seedcard/text.conf) and commit.
2. GitHub runs the tests automatically. An invalid value (a character the font can't engrave,
   text too long, a misspelt setting) makes them fail, and the image won't run.
3. Download the new image from the Actions tab as above.

If you change the text, also refresh the reference images that the tests compare against. This
needs Go: run `go test ./backup/ -run 'TestEngrave$' -update`, then look at
`backup/testdata/plate-15-*.png` and `plate-16-*.png` before committing.

## Development

### Tests

```sh
CGO_ENABLED=0 go test ./...
```

### Update through USB

A debug image can replace and restart the controller program on a running device. Build the
debug image from a local clone:

```sh
nix build .#image-debug
```

Write `result/seedhammer-debug.img` to an SD card. Connect the device by USB, using the port
closest to the mini-HDMI port (the one normally used for the engraver). Then upload and run a
new controller program:

```sh
export USBDEV=/dev/cu.usbmodem101 # Or (usually) /dev/ttyUSB0 on Linux.
nix run .#reload $USBDEV
```

In debug mode the controller's log goes to the USB serial device; `cat $USBDEV` shows it
(`nix run .#reload` does this for you).

### Remote control (debug mode)

```sh
echo "input up" > $USBDEV          # buttons: up, down, left, right, center, b1, b2, b3
echo "runes ACCIDENT" > $USBDEV    # types text; each space/newline presses b2
echo "screenshot" > $USBDEV        # saves a screenshot to the SD card
```

## Credits and licence

Built on [Seed Hammer](https://github.com/seedhammer/seedhammer) by the SeedHammer developers.
The files in this repository are in the public domain as described in the [LICENSE](LICENSE)
file, except files in directories with their own LICENSE files.
