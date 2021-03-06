const swaggermerge = require('swagger-merge');
const fs = require('fs');
const path = require('path');

const swaggerFilesPath = path.resolve(__dirname, '../../protobufs/gen/swagger');
const swaggerJsonPath = path.resolve(__dirname, '../../protobufs/gen/swagger.json');
const swaggerConfig = require(path.resolve(__dirname, '../swagger_config'));

/* # build_swagger.js
 *
 * This script recursively searches the swaggerFilesPath for any file ending in .swagger.json
 * loading all matching files and merging them into one for easier portability. In the second
 * step it then creates an html version of the documentation using spectacles.
 *
 * The defaults are defined in ../swagger_config.js
 */


// Find matching swagger files in path
// From: https://gist.github.com/kethinov/6658166
let getSwaggerFiles = function(dir, filelist) {
    let files = fs.readdirSync(dir);
    filelist = filelist || [];
    files.forEach(function(file) {
        if (fs.statSync(dir + '/' + file).isDirectory()) {
            getSwaggerFiles(dir + '/' + file, filelist);
        }
        else {
            if (file.indexOf(".swagger.json") > 0) {
                filelist.push(path.join(dir, file));
            }
        }
    });
    return filelist;
};

let files = getSwaggerFiles(swaggerFilesPath);
// There is a default swagger definition in swaggerConfig.defaultSwagger which we add first
swaggerModules = [swaggerConfig.defaultSwagger,];
files.forEach(function (f) {
    swaggerModules.push(require(f))
});

swaggermerge.on('warn', function (msg) {
    console.log(msg)
});

swaggermerge.on('err', function (msg) {
    console.error(msg);
    process.exit(1)
});

let merged = swaggermerge.merge(swaggerModules, swaggerConfig.info, swaggerConfig.pathPrefix, swaggerConfig.host, swaggerConfig.schemes);
let json = JSON.stringify(merged);
console.log("Merged swagger.json, writing to:", swaggerJsonPath);
fs.writeFileSync(swaggerJsonPath, json);
