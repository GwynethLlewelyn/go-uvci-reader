# go-uvci-reader
A simple tool to validate EU Digital COVID-19 Certificates

## Minimalistic instructions

(because this tool is, well, minimalistic...)

1. Clone
2. Compile with `go build`/`go install`
3. Scan the QR code on the EU Digital COVID-19 Certificate, save the scan as ASCII
4. `echo 'the long, random string of garbage from the QR scan' > ./go-uvci-reader`
5. Confirm that the result contains your personal data (as printed in plaintext on the EU Digital COVID-19 Certificate) and that the signature _is_ valid.

That's it. No options, no CLI commands, no config files. It just does one thing. Not necessarily well.

## Bugs

Too many to list.
**Most notably:** the signature checker fails and panics. In other words: you _can_ get the encoded data from the QR code, but it doesn't do the vital step of _validating_ it. As such, the original purpose — _validation!_ — is defeated.
**Reason for the above:** probably you need to know in advance which authority emitted the signature and get its public key. This is non-trivial. Most EU organisations seem to distribute their public keys only to a very limited set of entities; also, each EU member state is free to generate as many signatures as they wish, and somehow (internally) 'decide' which entities are allowed to emit valid signatures _and_ freely exchange their public keys. It's true that this works across borders, but _how_ it works is beyond my understanding of the (very long!) implementation details. They're written in the most dense Eurocratese.
**Why such an obscure, opaque method?** That, unfortunately, is not for me to answer. I was genuinely naïve enough to think that you'd read the QR code, see if the signature was valid, show the entity signing it, and that would be it. Apparently there are further steps that need to be in place before you can actually do that.
**But other non-Go tools seem to work!** Well, that's very likely because they were written by knowledgeable, professional programmers, not utterly clueless amateurs like me. _Or_ it's because I've interpreted the relevant documentation wrongly.

[![Codacy Badge](https://app.codacy.com/project/badge/Grade/e42ead2d0ade40f0a30dab46ac2c0625)](https://www.codacy.com/gh/GwynethLlewelyn/go-uvci-reader/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=GwynethLlewelyn/go-uvci-reader&amp;utm_campaign=Badge_Grade)
![CodeQL](https://github.com/GwynethLlewelyn/go-uvci-reader/workflows/CodeQL/badge.svg)
