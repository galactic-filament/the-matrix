express = require 'express'
raml2html = require 'raml2html'

configWithDefaultTemplates = raml2html.getDefaultConfig()
render = (res, path) ->
  raml2html.render(path, configWithDefaultTemplates).then(
    (result) -> res.send result
    (err) -> res.send err
  )
app = express()
app.get '/', (req, res) -> render res, './schema.raml'
app.get '/:file', (req, res) -> render res, "./#{req.params.file}.raml"

server = app.listen 80, -> console.log 'Listening on 80'

exit = -> server.close ->
  console.log 'Exiting'
  process.exit 0
process.on 'SIGTERM', exit
process.on 'SIGINT', exit
