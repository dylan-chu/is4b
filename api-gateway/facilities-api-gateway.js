var messages = require('./facilities_pb');
var services = require('./facilities_grpc_pb');

var grpc = require('grpc');
var http = require('http');

http.createServer(function(req, res) {
  res.setHeader("Access-Control-Allow-Origin", "*");
  res.setHeader("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept");
  getBuildings(req, res);
}).listen(8888);

function getBuildings(req, res) {
  var client = new services.FacilitiesAdminClient('localhost:53001', grpc.credentials.createInsecure());
  var request = new messages.ListBuildingsEvent();
  // var user;
  // if (process.argv.length >= 3) {
  //   user = process.argv[2];
  // } else {
  //   user = 'world';
  // }
  var results = [];
  request.setUserid(1);
  var call = client.listBuildings(request);
  call.on('data', function(response) {
    var bldg = new Object();
    bldg.id = response.getId();
    bldg.name = response.getName();
    results.push(bldg);
    console.log('Building: ' + response.getName());
  });
  call.on('end', function() {
    res.write(JSON.stringify(results));
    res.end();
    console.log("finished");
  });
  call.on('status', function(status) {
    //console.log("status:" + status);
  });
}
