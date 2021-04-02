const fs = require('fs');
const path = require('path');
const yaml = require('js-yaml');


var destination = "";
var base = "";

try {
    var filename = path.join(__dirname, '../.', 'blox.yaml'),
        contents = fs.readFileSync(filename, 'utf8'),
        config = yaml.load(contents);


    destination = config.destination;
    base = config.base;
} catch (err) {
    console.log(err.stack || String(err));
}



const datafile = path.join(__dirname, '../.', base, destination, "data.json");


var data = {}

const file = JSON.parse(fs.readFileSync(datafile, 'utf8'));
console.log(file)
data = file;


module.exports = data;
