expect = require('chai').expect
exec = require('child_process').exec
async = require 'async'
assert = require 'assert'

repos = [
  'omega-jazz'
  'pho-sho'
  'go-home'
  'py-lyfe'
  'es-bueno'
]

# utility commands
runCmd = (cmd, cb) ->
  exec cmd, (err, stdout, stderr) ->
    cb err
getFromCmd = (cmd, cb) ->
  exec cmd, cb
repoCmd = (repoName, cmd, cb) ->
  runCmd "cd ./repos/#{repoName} && #{cmd}", (err) ->
    cb err
getFromRepoCmd = (repoName, cmd, cb) ->
  getFromCmd "cd ./repos/#{repoName} && #{cmd}", cb
withRepos = (repos, done, iterator) ->
  async.each repos, iterator, (err) ->
    done err

# derived commands
cloneRepo = (repoName, cb) ->
  runCmd(
    "git clone https://github.com/ihsw/#{repoName}.git ./repos/#{repoName}"
    cb
  )
deleteRepo = (repoName, cb) ->
  runCmd "rm -rf ./repos/#{repoName}", cb

# derived repo commands
buildImage = (repoName, cb) ->
  repoCmd repoName, "./bin/build-images", cb
upWeb = (repoName, cb) ->
  repoCmd repoName, 'docker-compose up -d web-test', cb
getContainerId = (repoName, cb) ->
  getFromRepoCmd(
    repoName
    "docker-compose ps -q web-test"
    (err, stdout, stderr) ->
      return cb err, null if err
      cb null, stdout.trim()
  )
stopWeb = (repoName, cb) ->
  repoCmd repoName, 'docker-compose stop web-test', cb

# derived container commands
isRepoUp = (repoName, cb) ->
  getContainerId repoName, (err, containerId) ->
    return cb err, null if err
    getFromCmd "docker inspect #{containerId}", (err, stdout, stderr) ->
      return cb err, null if err
      info = JSON.parse(stdout)
      cb null, info[0].State.Running
removeRepoContainer = (repoName, cb) ->
  getContainerId repoName, (err, containerId) ->
    return cb err if err
    runCmd "docker rm -v #{containerId}", cb
removeTestContainer = (repoName, cb) ->
  runCmd "docker rm -v #{repoName}_test", cb
runTests = (repoName, cb) ->
  getContainerId repoName, (err, containerId) ->
    return cb err if err

    link = "#{containerId}:ApiServer"
    name = "#{repoName}_test"
    getFromCmd(
      "docker run -t --link #{link} --name #{name} ihsw/the-matrix-tests"
      cb
    )

describe 'Api Servers', ->
  before (done) ->
    withRepos repos, done, (repoName, eachNext) ->
      tasks = [
        (seriesNext) -> cloneRepo repoName, (err) -> seriesNext err
        (seriesNext) -> buildImage repoName, (err) -> seriesNext err
      ]
      async.series tasks, (err) -> eachNext err
  after (done) ->
    withRepos repos, done, (repoName, eachNext) ->
      tasks = [
        (seriesNext) ->
          isRepoUp repoName, (err, isUp) ->
            return seriesNext null if !isUp
            stopWeb repoName, (err) -> seriesNext err
        (seriesNext) -> removeRepoContainer repoName, (err) -> seriesNext err
        (seriesNext) -> removeTestContainer repoName, (err) -> seriesNext err
        (seriesNext) -> deleteRepo repoName, (err) -> seriesNext err
      ]
      async.series tasks, (err) -> eachNext err
  it 'should run the test suite', (done) ->
    withRepos repos, done, (repoName, eachNext) ->
      tasks = [
        (seriesNext) -> upWeb repoName, (err) -> seriesNext err
        (seriesNext) ->
          isRepoUp repoName, (err, isUp) ->
            return seriesNext new Error "repo #{repoName} is not up" if !isUp
            runTests repoName, (err, stdout, stderr) ->
              expect(err).to.equal null, "Could not run tests for #{repoName}"
              seriesNext null
        (seriesNext) -> stopWeb repoName, (err) -> seriesNext err
      ]
      async.series tasks, (err) -> eachNext err
