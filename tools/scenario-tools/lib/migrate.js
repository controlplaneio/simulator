const { loadYamlFile, writeYamlFile, findScenarioFiles } = require('./io')
const { createLogger } = require('./logger')

const logger = createLogger({})

// A transformation function accepts a tasks manifest as a JS object and
// returns a transformed manifest.  A transformation function *must* include
// a check to validate that the manifest is in the format expected before making
// any changes and warn and short circuit if not

// These transformations are a history of the migrations we have performed over
// time on the tasks.yaml files

function groupHintsByTask (hints) {
  return hints.reduce((acc, hint) => {
    // Pull apart the hint text on the first colon
    const [task, hintText] = hint.text.split(': ', 2)

    const sortOrder = Number(task[task.length - 1])

    // Create a new group for the task if it doesn't exist
    if (!acc[task]) {
      acc[task] = {
        'sort-order': sortOrder,
        hints: [],
        'starting-point': ''
      }
    }

    // Reconstruct a hint without the task prefix and add it to the group
    acc[task].hints.push({ text: hintText })

    return acc
  }, {})
}

// Takes an array of top level yaml properties represented as a JS object and
// Returns a transformed array of top-level properties in the new format
function transformV0ToV1 (hints) {
  if (!Array.isArray(hints)) {
    return logger.warn('skipping manifest because it was not an array')
  }

  const transformed = {}
  // remove hint count
  delete hints[0].general_overview['num-hints']

  transformed.objective = hints[0].general_overview.objective
  transformed['starting-point'] = hints[0].general_overview['starting-point']

  // add version
  transformed.kind = 'cp.simulator/scenario:1.0.0'

  // transform hints
  const hintsList = Object.entries(hints[1].hints)
    .reduce((acc, [key, val], idx) => {
      acc.push({ text: val })
      return acc
    }, [])

  transformed.tasks = groupHintsByTask(hintsList)

  return transformed
}

// Renames tasks from e.g. 'Task 1' -> '1'
function removeTaskPrefix (manifest) {
  const tasks = manifest.tasks

  const unrecognised = Object.keys(tasks).filter(task => !task.startsWith('Task '))

  if (unrecognised.length > 0) {
    return logger.warn('skipping manifest due to unrecognised tasks')
  }

  Object.keys(tasks).forEach(key => {
    // new name is just the digit at the end of the string
    const newname = key[key.length - 1]

    Object.defineProperty(tasks, newname, Object.getOwnPropertyDescriptor(tasks, key))
    delete tasks[key]
  })

  return manifest
}

function migrate (options) {
  // TODO(rem): check options and migrate a single scenario or all

  let manifests
  if (options.type === 'legacy') {
    manifests = findScenarioFiles('./simulation-scripts/scenario', 'hints.yaml')
  } else {
    manifests = findScenarioFiles('./simulation-scripts/scenario', 'tasks.yaml')
  }

  manifests.forEach(manifest => {
    const original = loadYamlFile(manifest)
    let transformed
    if (options.type === 'legacy') {
      transformed = transformV0ToV1(original)
    } else if (options.type === 'remove-task-prefix') {
      transformed = removeTaskPrefix(original)
    } else {
      throw new Error('unrecognised type of migration')
    }

    writeYamlFile(transformed, manifest)
  })
}

module.exports = {
  groupHintsByTask,
  transformV0ToV1,
  removeTaskPrefix,
  migrate
}
