var express = require('express');
var raml2html = require('raml2html');

var configWithDefaultTemplates = raml2html.getDefaultConfig();

var app = express();
app.get('/', function (req, res) {
  raml2html.render('./schema.raml', configWithDefaultTemplates).then(function (result) {
    res.send(result);
  }, function (err) {
    res.send(err);
  });
});
app.get('/:file', function (req, res) {
  raml2html.render('./' + req.params.file + '.raml', configWithDefaultTemplates).then(function (result) {
    res.send(result);
  }, function (err) {
    res.send(err);
  });
});

var server = app.listen(80, function () {
  console.log('Listening on 80');
});

var exit = function () { server.close(function () {
  console.log('Exiting');
  process.exit(0);
})};

process.on('SIGTERM', exit);
process.on('SIGINT', exit);
