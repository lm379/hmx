var createError = require('http-errors');
var express = require('express');
var path = require('path');
var cookieParser = require('cookie-parser');
var logger = require('morgan');
var bodyParser = require("body-parser");
var dotenv = require('dotenv').config();
// var jwt = require('jsonwebtoken');

var indexRouter = require('./routes/index');
var usersRouter = require('./routes/data');
var stsRouter = require('./routes/sts-server');
var accountRouter = require('./routes/account');
var app = express();

// 后端HTTP端口
const port = process.env.port;

// view engine setup
app.set('views', path.join(__dirname, 'views'));
app.set('view engine', 'jade');

app.use(logger('dev'));
app.use(express.json());
app.use(express.urlencoded({ extended: false }));
app.use(cookieParser());
app.use(express.static(path.join(__dirname, 'public')));

app.use('/', indexRouter);
app.use('/api', usersRouter);
app.use('/api', stsRouter);
app.use('/', accountRouter);

app.use(bodyParser.urlencoded({ extended: false }));
app.use(bodyParser.json());

//解决跨域
app.all("*", function (req, res, next) {
  res.header("Access-Control-Allow-Origin", "*")
  res.header("Access-Control-Allow-Headers", "Content-Type")
  res.header("Access-Control-Allow-Methods", "*")
  res.header("Content-Type", "application/json;charset=utf-8")
  next()
});

// catch 404 and forward to error handler
app.use(function (req, res, next) {
  next(createError(404));
});

// error handler
app.use(function (err, req, res, next) {
  // set locals, only providing error in development
  res.locals.message = err.message;
  res.locals.error = req.app.get('env') === 'development' ? err : {};

  // render the error page
  res.status(err.status || 500);
  res.render('error');
});

// app.listen(port, () => {
//   console.log(`Server is running on host: http://127.0.0.1:${port}`);
// });

module.exports = app;
