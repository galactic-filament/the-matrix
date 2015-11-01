express = require 'express'
raml2html = require 'raml2html'
merge = require 'merge'
_ = require 'underscore'

app = express()

configWithDefaultTemplates = raml2html.getDefaultConfig()
render = (res, path) ->
  raml2html.render(path, configWithDefaultTemplates).then(
    (result) -> res.send result
    (err) -> res.send err
  )
pickNonfalsy = _.partial _.pick, _, _.identity

app.get '/:file?', (req, res) ->
  params = merge { file: 'schema' }, pickNonfalsy(req.params)
  render res, "./#{params.file}.raml"

server = app.listen 80, -> console.log 'Listening on 80'
exit = -> process.exit 0
process.on 'SIGTERM', exit
process.on 'SIGINT', exit
