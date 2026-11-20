# SeedCard to-do list

The SeedCard system, end to end:

| Part | What it is | Where |
|---|---|---|
| **Engraver** | Seed Hammer v1.2.1 + SeedCard change: 12 words on SH01, no SeedQR, editable text | this repo, `seedcard` branch |
| **Plate reader** | Reads the 12 words off an engraved plate with the camera | [seedcard/reader](https://github.com/seedcard/reader) (laptop prototype, capture kit, test photos) |
| **SeedSigner fork** | Official SeedSigner plus one added module: "Read SeedCard plate" | not started |

Each council member holds their own SeedSigner (fork) and SD card, their plate, their recovery
sheet, their phone and identity card. The engraver is a shared station: the Seed Hammer machine
plus a Pi Zero running the engraver image. Engraving is seed-only (SKIP): no descriptor side.

---

## 1. Engraver

- [x] Clean rebuild from upstream v1.2.1 (`seedcard` branch)
- [x] `seedcard/text.conf`: version mark `SC`, title, `key_number = auto`, `seed_qr = no`
- [x] Image built by GitHub Actions on every push, with SHA-256
- [x] `seedcard` is the default branch; `master` kept as an archive
- [ ] First boot of the image on the engraver's Pi Zero; dry run (long-press middle button on the engraving screen)
- [ ] Engrave a dummy seed with SKIP; compare with the BIFB TREASURY X test plate
- [ ] Tag the first tested version (e.g. `sc-1.0`) and record its image SHA-256 in this file
- [ ] Decide the `title` for DAO plates (recommended: empty, so a plate doesn't name the DAO)

## 2. Plate reader (prototype → module)

- [x] Laptop prototype: corner holes → straighten → stroke map → align on line numbers → letters → BIP39 words → checksum (4 of 4 usable phone photos read; 1 near-tie caught by the checksum)
- [x] Own repo: [seedcard/reader](https://github.com/seedcard/reader), with `run_tests.py` over the dummy-plate photos (5/5 pass)
- [ ] Pi camera photos of a dummy plate with the plate-capture kit (distances, light, angles)
- [ ] Check close-up focus of the OV5647; decide on lens adjustment or a fixed jig/stand
- [ ] Choose the lighting (diffuse works best so far)
- [ ] Speed: half resolution, single-match alignment, hole search inside an on-screen framing guide; target well under 30 s on a Pi Zero
- [ ] Letter templates trained on real plate crops (wider confidence margins)
- [ ] Review screen: show each word's crop with alternatives, as in the SeedSigner demo; never load a seed without the user confirming every word
- [ ] Test set: dummy plates in varied light, every result checked against the known words

## 3. SeedSigner fork

Principles (agreed 20 November 2026):

- **As small as possible.** Official SeedSigner plus one added reader module; nothing else touched. This device signs for the treasury.
- **Track official releases.** Rebase the fork on each SeedSigner release, so security fixes and community review carry over.
- **Reproducible build and verified hashes.** The SD card ceremony checks each member's card exactly as it would an official image.
- **Watch the official feature.** When SeedSigner's handwritten-word reader is released, shrink the fork to a SeedCard plate template on top of it.
- **Upstream if possible.** Contributing plate support to SeedSigner would mean no fork at all.

Steps:

- [ ] Fork the official SeedSigner repo; note the release it is based on
- [ ] Add "Read SeedCard plate" under Load a seed, beside Scan a SeedQR and Enter 12-word seed
- [ ] Port the reader (Python, as SeedSigner is) within the Pi Zero's limits
- [ ] Reproducible image build; publish hashes
- [ ] Keep a short diff-against-upstream document, like `SEEDCARD.md` here
- [ ] Follow SeedSigner's handwritten-word reader (pizza_bitcoiner's post, November 2026); contact the developers about plate support

## 4. Key ceremony and btc-dao integration

- [ ] Add the plate check to the btckeys ceremony: after engraving, read the plate back; the fingerprint must match the engraved one and the exported xpub
- [x] Add `docs/drafts/recovery-sheet.md` to btc-dao (drafted 20 November 2026; added in btc-dao PRs #12 and #13): new sheet notice, "I have printed it" confirmation kept on the member's phone, signing gate, spoken check before a migration
- [ ] btc-dao spec and guide: one SeedSigner per council member, no SeedQR card, plates engraved seed-only (SKIP) or stamped, signing at home
- [ ] Mainnet readiness: every council member has confirmed the current recovery sheet; a council change issues a new sheet before the migration

## 5. Not doing (decided)

- Descriptor side on plates: the descriptor is public, changes with every council change, and is kept on the recovery sheet and in the DAO record
- Merging the engraver into the SeedSigner fork: separate SD cards, separate devices
- The "words 4–9" layout (an unsuccessful 2024 test)
