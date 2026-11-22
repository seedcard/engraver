# SeedCard to-do list

The SeedCard system, end to end:

| Part | What it is | Where |
|---|---|---|
| **Engraver** | Seed Hammer v1.2.1 + SeedCard change: 12 words on SH01, no SeedQR, editable text | this repo, `seedcard` branch |
| **Plate reader** | Reads the 12 words off an engraved plate with the camera | [seedcard/reader](https://github.com/seedcard/reader) (laptop prototype, capture kit, test photos) |
| **SeedSigner fork** | Official SeedSigner + [Haller](https://github.com/haller-project/haller) handwritten-seed module + a Seed Hammer plate module | plan agreed 22 November 2026 |

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

## 3. SeedSigner fork: a Seed Hammer module on Haller

[Haller](https://github.com/haller-project/haller) (released 22 November 2026, MIT) reads
handwritten seed words in a SeedSigner fork
([`seedsigner-haller`, branch `handwritten-seed`](https://github.com/haller-project/seedsigner-haller/tree/handwritten-seed)):
based on official SeedSigner, new files plus three small edits (camera, settings, seed menu),
on-device C library, mandatory word review, and **"verified" only when a written fingerprint
matches**. Its handwriting networks don't read engraved plates reliably (tested on the dummy
photos: 1 of 5 alone), while our font matching reads 4 of 4. So we build a **Seed Hammer plate
module** that plugs into Haller.

Principles (agreed 20 November 2026, still apply):

- **As small as possible.** Official SeedSigner + Haller's module + our plate module; nothing else touched. This device signs for the treasury.
- **Track releases.** Rebase on each SeedSigner release, and on each Haller release.
- **Reproducible build and verified hashes** for the SD card ceremony.
- **Upstream if possible.** Offer plate support to the Haller project once it is built and reviewable; one shared fork is better than two.

Integration point: Haller's screens work from one `SeedReading` (words, each with its crop and
ranked candidates). The plate module returns the same object, so review, trust levels and
fingerprint verification work unchanged.

Build tasks:

- [ ] `seedsigner/helpers/seedhammer/` module: `read_plate(image, progress, on_word) -> SeedReading`
  - [ ] plate finder (corner holes), straightening, orientation, stroke map, line-number alignment (from `seedcard/reader`)
  - [ ] all Seed Hammer seed plates: SH01, SH02, SH03 layouts computed from the engraver code and the plate drawing (see `seedcard/reader/PLATES.md`)
  - [ ] text-grid registration for plates without holes (Seed Hammer v2 machine) or with holes hidden
  - [ ] letter matching against Seed Hammer's font, scored into per-word candidates with probabilities
  - [ ] **engraved fingerprint**: read the 8 hex digits on the left edge (same font, fixed place) and verify the seed with Haller's fingerprint check: plates can reach "verified"
  - [ ] each word's crop for the review screen
- [ ] Hook: in Haller's read view, try the plate path first when Seed Hammer holes or a Seed Hammer text grid are found; otherwise the handwriting path
- [ ] Speed on a Pi Zero: hot loops in C, following Haller's `native/` pattern (no new Python dependency)
- [ ] Layout file per engraver version (`layout.json`), checked by test vectors
- [ ] Tests: dummy-plate photos (`seedcard/reader/testdata`), including Pi camera photos; wrong seed must never be "verified"
- [ ] Device image: depends on Haller's SeedSigner OS build (Buildroot package, not yet published); until then Raspberry Pi OS by hand
- [ ] Short document of what differs from Haller and from official SeedSigner
- [ ] **Then** contact the Haller project (pizza_bitcoiner), once the module is reviewable: a branch on `haller-1` they can run with their own dev setup, dummy-plate test photos (incl. Pi camera) with results per photo, the safety result (no wrong seed ever "verified"), a short design note, Pi Zero timings
- [ ] Meanwhile: build against the `haller-1` tag; watch their repo and posts for any plate/steel work, and get in touch early if it appears

## 4. Key ceremony and btc-dao integration

- [ ] Add the plate check to the btckeys ceremony: after engraving, read the plate back; the fingerprint must match the engraved one and the exported xpub
- [x] Add `docs/drafts/recovery-sheet.md` to btc-dao (drafted 20 November 2026; added in btc-dao PRs #12 and #13): new sheet notice, "I have printed it" confirmation kept on the member's phone, signing gate, spoken check before a migration
- [ ] btc-dao spec and guide: one SeedSigner per council member, no SeedQR card, plates engraved seed-only (SKIP) or stamped, signing at home
- [ ] Mainnet readiness: every council member has confirmed the current recovery sheet; a council change issues a new sheet before the migration

## 5. Not doing (decided)

- Descriptor side on plates: the descriptor is public, changes with every council change, and is kept on the recovery sheet and in the DAO record
- Merging the engraver into the SeedSigner fork: separate SD cards, separate devices
- The "words 4–9" layout (an unsuccessful 2024 test)
