const { MAX_SCORE } = require('./constants')

function calculateScore (progress, tasks) {
  const currentTask = progress.currentTask
  const taskProgress = progress.tasks.find(t => t.id === currentTask)
  const lastHintIndex = taskProgress.lastHintIndex
  const hintCount = tasks[currentTask].hints.length

  if (lastHintIndex === null) {
    return MAX_SCORE
  }

  let score = MAX_SCORE

  for (let i = 0; i <= lastHintIndex && i < hintCount; i++) {
    score -= tasks[currentTask].hints[i].penalty
  }

  return score
}

module.exports = {
  calculateScore
}
