const { MAX_SCORE } = require('./constants')

function calculateScore (progress, tasks) {
  const currentTask = progress.current_task
  const lastSeenHintIndex = progress[currentTask]
  const hintCount = tasks[currentTask].hints.length

  if (lastSeenHintIndex === undefined) {
    return MAX_SCORE
  }

  let score = MAX_SCORE

  for (let i = 0; i <= lastSeenHintIndex && i < hintCount; i++) {
    score -= tasks[currentTask].hints[i].penalty
  }

  return score
}

module.exports = {
  calculateScore
}
