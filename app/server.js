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
app.listen(80);
