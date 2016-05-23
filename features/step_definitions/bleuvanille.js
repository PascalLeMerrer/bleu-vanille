/* jslint node: true */
'use strict';

var bitterapple = require('bitter-apple');
var http = require('http');
var querystring = require('querystring');
var _ = require('lodash');
var curry = require('lodash.curry');

module.exports = function() {
	// cleanup before every scenario
	this.Before(function(scenario) {
		var server = process.env.BleuVanilleName
		var port = process.env.BleuVanillePort
		this.bitterapple = new bitterapple.BitterApple('http', server + ':' + port);
	});

	this.Given(/^I log as test user$/, function(callback) {
		authenticate("testuser@bleuvanille.com", "xaFqJDeJldIEcdfZS", this, callback)
	});

	this.Given(/^I log as admin user$/, function(callback) {
		authenticate("admin@bleuvanille.com", process.env.AdminPassword, this, callback)
	});

	function authenticate(login, password, self, callback) {
		self.bitterapple.setRequestBody("email=" + login + ";password=" + password);
		self.bitterapple.addRequestHeader("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8");
		self.bitterapple.post("/users/login", function(error, response) {
			if (error) {
				callback.fail(error);
			}
			else {
				self.bitterapple.sendCookie("token")
				callback();
			}
		});
	}

	this.Given(/^the value of body path (\S+) should be (\S+) \+ (\d)$/, function(path, data1, data2, callback) {
		data1 = parseInt(dereferenceData(this.bitterapple, data1));
		data2 = parseInt(dereferenceData(this.bitterapple, data2));
		this.bitterapple.realValue = this.bitterapple.getResponseObject().body;
		var valueAtPath = this.bitterapple.evaluatePath(path, this.bitterapple.realValue)
		if(data1 + data2 !== parseInt(valueAtPath)) {
			callback(new Error('Body is ' + this.bitterapple.realValue + '\nExpected count to be ' + (data1+data2))) ;
		} else {
			callback();
		}
	});

	// if the given data is a variable name, returns its value; otherwise returns the data itself
	function dereferenceData(bitterapple, data) {
		var value
		if (typeof(data) === 'string') {
			value = bitterapple.getGlobalVariable(data)
			if (value === undefined) {
				value = bitterapple.scenarioVariables[data]
			}
		}
		if (value === undefined) {
			return data
		}
		return value
	}

	function is_int(value){
	  if((parseFloat(value) == parseInt(value)) && !isNaN(value)){
	      return true;
	  } else {
	      return false;
	  }
	}
	this.Given(/^I truncate the database collection (.*)$/, function(name, callback) {
		var postData = JSON.stringify({
			"collection": name,
			"example": {}
		});

		var protocol = (process.env.DatabaseProtocol || 'http') + ':';
		var hostname = process.env.DatabaseHost || 'localhost';
		var port = process.env.DatabasePort || 8529;
		var database = process.env.DatabaseName || 'bleuvanille';

		var options = {
			'protocol': protocol,
			'hostname': hostname,
			'port': port,
			'path': '/_db/' + database +'/_api/simple/remove-by-example',
			'method': 'PUT',
			'headers': {
				'Content-Type': 'application/x-www-form-urlencoded',
				'Content-Length': postData.length
			}
		}

		var curried = _.curry(truncated);
		var req = http.request(options, curried(callback));
		req.write(postData);
		req.end();
	});

	function truncated(callback, response) {
		console.log("response status %s message %s", response.statusCode, response.statusMessage )
		callback()
	}

};
