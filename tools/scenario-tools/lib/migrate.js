const {loadHintsFile, writeHintsFile, findScenarioHintsFiles} = require('./io')

function groupHintsByTask(hints) {
  return hints.reduce((acc, hint) => {
    // Pull apart the hint text on the first colon
    const [task, hintText] = hint.text.split(': ', 2)

    const sortOrder = Number(task[task.length - 1])

    // Create a new group for the task if it doesn't exist
    if (!acc[task]) acc[task] = { 
      "sort-order": sortOrder, 
      hints: [],
      "starting-point": ""
    }

    // Reconstruct a hint without the task prefix and add it to the group
    acc[task].hints.push({ text: hintText })

    return acc
  }, {})
}

// Takes an array of top level yaml properties represented as a JS object and
// Returns a transformed array of top-level properties in the new format
function transformV0ToV1(hints) {
  const transformed = {}
  // remove hint count
  delete hints[0]['general_overview']["num-hints"]

  transformed.objective = hints[0]['general_overview'].objective
  transformed['starting-point'] = hints[0]['general_overview']['starting-point']
  
  // add version
  transformed.kind = 'cp.simulator/scenario:1.0.0'

  // transform hints
  const hintsList  = Object.entries(hints[1].hints)
    .reduce((acc, [key, val], idx) => { 
      acc.push({ text: val })
      return acc
    }, [])

  transformed.tasks = groupHintsByTask(hintsList)

  return transformed
}

function migrate() {
  hintsFiles = findScenarioHintsFiles('./simulation-scripts/scenario')
  hintsFiles.forEach(hintsFile => {
    const original = loadHintsFile(hintsFile)
    const transformed = transformV0ToV1(original)
    writeHintsFile(transformed, hintsFile)
  })
}

module.exports = {
  groupHintsByTask,
  transformV0ToV1,
  migrate
}
