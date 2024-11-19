// config/secretKey.js
const crypto = require('crypto');

// 生成随机密钥
const secretKey = crypto.randomBytes(32).toString('hex');

module.exports = secretKey;