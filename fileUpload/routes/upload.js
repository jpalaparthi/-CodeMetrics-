var express = require('express');
var request = require('request');
var router = express.Router();
var multer = require('multer');

// Need to setup multer module....
var storage =   multer.diskStorage({
  destination: function (req, file, callback) {
    callback(null, './uploads');
  },
  filename: function (req, file, callback) {
    callback(null, file.fieldname + '-' + Date.now());
  }
});
var upload = multer({ storage : storage}).single('fileUpload');

// Post method to upload a file .. 
// Please give any route...
router.post('/json',function(req,res){
    upload(req,res,function(err) {
        if(err) {
            return res.end("Error uploading file.");
        }
        res.end("File is uploaded");
    });
});


module.exports = router;
