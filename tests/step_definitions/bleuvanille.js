/* jslint node: true */
'use strict';

var bitterapple = require('bitter-apple');

module.exports = function() {
    // cleanup before every scenario
    this.Before(function(scenario) {
        this.bitterapple = new bitterapple.BitterApple('http', 'localhost:4000');
    });
    
	this.Given(/^I log as test user$/, function(callback) {
		this.bitterapple.setRequestBody("email=testuser@bleuvanille.com;password=xaFqJDeJldIEcdfZS");
		this.bitterapple.addRequestHeader("Content-Type","application/x-www-form-urlencoded; charset=UTF-8");
		var self = this;
		this.bitterapple.post("/users/login", function (error, response) {
			if (error) {
				callback.fail(error);
			}
			
			self.bitterapple.storeValueOfHeaderInGlobalScope("Set-Cookie","cookie");
			
			callback();
		});
	});
	
	this.Given(/^I log as admin user$/, function(callback) {
		this.bitterapple.setRequestBody("email=admin@bleuvanille.com;password=xaFqJDeJldIEcdfZS");
		this.bitterapple.addRequestHeader("Content-Type","application/x-www-form-urlencoded; charset=UTF-8");
		var self = this;
		this.bitterapple.post("/users/login", function (error, response) {
			if (error) {
				callback.fail(error);
			}
			
			self.bitterapple.storeValueOfHeaderInGlobalScope("Set-Cookie","cookie");
			
			callback();
		});
	});
};
