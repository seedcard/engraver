# SeedCard engraver

This branch is Seed Hammer **v1.2.1**, the last upstream release that supports the
credit-card plate (SH01; upstream removed it in v1.3.0), plus one small SeedCard change.
To see everything that differs from upstream:

    git diff v1.2.1 seedcard --

## What SeedCard changes

| | Seed Hammer v1.2.1 | SeedCard |
|---|---|---|
| Seed side, SH01 | 12 words + SeedQR | 12 words, **no SeedQR** |
| Version mark (SH01: up the left edge) | `V1` | **`SC`** (editable) |
| Title (SH01: up the right edge) | wallet name from the descriptor | the same, or a fixed title (editable) |
| Seed words, word numbers, fingerprint, key number (1/2) | from seed and wallet | **unchanged** |
| Descriptor side (wallet backup with QR) | upstream | **unchanged** |
| Input: scanning SeedQR / descriptor, entering words | upstream | **unchanged** |
| SH02 / SH03 plates | upstream | same SeedCard text settings apply |

## Editing the engraved text

All editable text is in **`seedcard/text.conf`**:

```
version_mark = SC     # 1–4 characters, replaces "V1"
title =               # empty = wallet name; otherwise up to 18 characters
seed_qr = no          # yes = engrave the SeedQR again
```

Only capitals A–Z, digits, spaces and `: , / ' . ( ) [ ] { } # * @ -` can be engraved; lower
case is turned into capitals. The file is built into the program, so:

1. Edit `seedcard/text.conf`.
2. Run `go test ./seedcard/ ./backup/`. An invalid value (unsupported character, too long,
   misspelt setting) fails here and never reaches a plate. If you changed the text, also run
   `go test ./backup/ -run 'TestEngrave$' -update` to refresh the reference images, and look
   at the new `backup/testdata/plate-15-*.png` and `plate-16-*.png` before committing.
3. Build the SD-card image as described in the README (`nix build`) and flash it.

Because the text is built in, the image's SHA-256 identifies exactly what it engraves.
Record the hash with each release.

## Plate layout reference

`backup/testdata/plate-15-side-1-1-of-1-words-12.png` (single key) and
`plate-16-side-1-1-of-2-words-12.png` (key 1 of 2) show the SH01 seed side. The layout
reproduces, pixel for pixel, the BIFB TREASURY X test plate engraved in 2024, and is the layout
the plate-reader prototype is calibrated for.

## Not changed on purpose

- The descriptor side still carries the full descriptor as text **and** QR. For a multisig wallet
  it's the only engraved copy of what's needed to rebuild the wallet, and the QR alone is
  often hard to scan off steel.
- Upstream's check that a descriptor fits the plate (`ErrDescriptorTooLarge`) still works.
