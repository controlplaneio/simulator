#!/usr/bin/env node

// Usage: run `make devtools` and then `migrate-hints` in root of the simulator dir

const {loadHintsFile, writeHintsFile, findScenarioHintsFiles} = require('../lib/io')
const {transformV0ToV1} = require('../lib/transforms')

hintsFiles = findScenarioHintsFiles('./simulation-scripts/scenario')
hintsFiles.forEach(hintsFile => {
  const original = loadHintsFile(hintsFile)
  const transformed = transformV0ToV1(original)
  writeHintsFile(transformed, hintsFile)
})

