<div align="center">
<h1>D2-Live</h1>
<h3>Generate live diagram URLs with D2</h3>

[![Publish Status](https://github.com/Watt3r/d2-live/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/Watt3r/d2-live/actions/workflows/docker-publish.yml)
[![Build Status](https://github.com/Watt3r/d2-live/actions/workflows/build.yml/badge.svg)](https://github.com/Watt3r/d2-live/actions/workflows/build.yml)
[![Uptime](https://atlas.lucas.tools/api/badge/2/uptime/240)](https://atlas.lucas.tools/)
[![License](https://img.shields.io/badge/License-MIT-orange.svg)](https://github.com/Watt3r/d2-live/blob/master/LICENSE)
[![Coverage Status](https://coveralls.io/repos/github/Watt3r/d2-live/badge.svg?branch=master)](https://coveralls.io/github/Watt3r/d2-live?branch=master)
</div>

## Demo

Given a diagram:

```
x: That's got to be the best pirate I've ever seen
x -> y
y: So it would seem...
y -> z: Jack Sparrow's theme starts
z: Fade into horizon
```

Paste it in [d2-playground](https://play.d2lang.com/), the URL will have an encoded link with a `script` variable: 

```
LNExDsIwEETR3qf4nSt8ABeUSNCGCzhkRSxIHGUXiHN6ZET9NKORZotcx2ReuRfDCr1go9CLGktekwln_xZaFBWZ3cbhSHU10hWy8Smv59BoCiG42nSPXNLtQbf8FnttnZOgllZTt0dOaRDybIX_ge4bAAD__w%3D%3D
```
! Important: the playground URL adds some extra bits after the `script` variable, don't copy these!

Append the encoded string to the service URL, for example: `https://d2.atlas.lucas.tools/svg/<encoded_string>`, you will get an image from the URL:

![Diagram](https://d2.atlas.lucas.tools/svg/LNExDsIwEETR3qf4nSt8ABeUSNCGCzhkRSxIHGUXiHN6ZET9NKORZotcx2ReuRfDCr1go9CLGktekwln_xZaFBWZ3cbhSHU10hWy8Smv59BoCiG42nSPXNLtQbf8FnttnZOgllZTt0dOaRDybIX_ge4bAAD__w%3D%3D)

It is a normal image svg and can be embedded everywhere you want.
