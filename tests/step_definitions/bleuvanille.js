/* jslint node: true */
'use strict';

var bitterapple = require('bitter-apple');

module.exports = function() {
    // cleanup before every scenario
    this.Before(function(scenario) {
        this.bitterapple = new bitterapple.BitterApple('http', 'localhost:4000');
    });
};
