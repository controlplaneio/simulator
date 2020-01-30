const { loadYamlFile, getProgress, saveProgress } = require('./io')
const { calculateScore } = require('./scoring')
const { createLogger } = require('./logger')
const { TASKS_FILE_PATH, PROGRESS_FILE_PATH } = require('./constants')
const { prompt } = require('inquirer')

const logger = createLogger({})

async function startTask (task, taskspath = TASKS_FILE_PATH,
  progresspath = PROGRESS_FILE_PATH, log = logger) {
  const { tasks } = loadYamlFile(taskspath)

  if (!tasks[task]) {
    log.warn('Cannot find task')
    return false
  }

  const progress = getProgress(progresspath)
  const currentTask = progress.current_task

  // either the user hasnt started any tasks or they have already been scored
  // for this task so just change the current task
  if (currentTask === undefined || progress[currentTask].score !== undefined) {
    progress.current_task = task
    saveProgress(progress, progresspath)
    log.info(`You are now on ${task}`)

    return true
  }

  const answers = await prompt([{
    type: 'expand',
    messsage: `Would you like to be scored for task ${currentTask}? (If you choose no you cannot be scored on this task in the future)`,
    name: 'scoring',
    choices: [{
      key: 'y',
      name: 'Yes',
      value: 'yes'
    }, {
      key: 'n',
      name: 'No',
      value: 'no'
    }, {
      key: 'c',
      name: 'Cancel',
      value: 'cancel'

    }]
  }])

  if (answers.length !== 1) {
    // should never happen
    throw new Error(
      'No changes made - expected an answer from the scoring prompt')
  }

  if (answers[0] === 'cancel') {
    log.info(`You cancelled... leaving you on ${currentTask}`)
    return false
  }

  if (answers[0] === 'yes') {
    const score = calculateScore(progress, tasks)
    log.info(`Your score for task ${currentTask} was ${score}`)
    progress[currentTask].score = score
  } else if (answers[0] === 'no') {
    log.info(`You chose not to be scored on ${currentTask}`)
    progress[currentTask].score = 'skip'
  } else {
    // should never happen
    throw new Error(
      `No changes made - unrecognised answer from scoring prompt: ${answers[0]}`)
  }

  progress.current_task = task
  saveProgress(progress, progresspath)
  log.info(`You are now on ${task}`)

  return true
}

function getCurrentTask (progresspath = PROGRESS_FILE_PATH) {
  const progress = getProgress(progresspath)
  return progress.current_task
}

module.exports = {
  startTask,
  getCurrentTask
}
