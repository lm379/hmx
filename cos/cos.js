const request = require('request');
const COS = require('cos-nodejs-sdk-v5');
const { response } = require('../app');

const Bucket = process.env.Bucket; // 存储桶名称
const Region = process.env.Region; // 存储桶所在地域
const Prefix = process.env.Prefix; // 存储桶路径

const cos = new COS({
    getAuthorization: function (options, callback) {
        // 初始化时不会调用，只有调用 cos 方法（例如 cos.putObject）时才会进入
        // 异步获取临时密钥
        request({
            url: 'http://127.0.0.1:5000/api/sts', // 替换为自己的获取临时密钥的服务url
            data: {
                Bucket: Bucket,
                Region: Region
            }
        }, function (err, response, body) {
            let data = null;
            let credentials = null;
            try {
                data = JSON.parse(body);
                credentials = data.credentials;
            } catch (e) { }
            if (!data || !credentials) return console.error('credentials invalid');
            callback({
                TmpSecretId: credentials.tmpSecretId,        // 临时密钥的 tmpSecretId
                TmpSecretKey: credentials.tmpSecretKey,      // 临时密钥的 tmpSecretKey
                SecurityToken: credentials.sessionToken, // 临时密钥的 sessionToken
                // 建议返回服务器时间作为签名的开始时间，避免用户浏览器本地时间偏差过大导致签名错误
                StartTime: data.startTime, // 时间戳，单位秒，如：1580000000
                ExpiredTime: data.expiredTime, // 临时密钥失效时间戳，是申请临时密钥时，时间戳加 durationSeconds
                ScopeLimit: true, // 细粒度控制权限需要设为 true，会限制密钥只在相同请求时重复使用
            });
        });
    }
});

// 判断存储桶是否存在
function doesBucketExist() {
    cos.headBucket({
        Bucket: Bucket, // 填入您自己的存储桶，必须字段
        Region: Region,  // 存储桶所在地域，例如 ap-beijing，必须字段
    }, function (err, data) {
        if (err) {
            return new res.JSON({
                code: 1,
                message: 'Bucket not exist'
            });
        }
    });
}

cos.getBucket({
    Bucket: Bucket,
    Region: Region,
    Prefix: Prefix,
    function(err, data) {
        console.log(err || data.Contents);
    }
});