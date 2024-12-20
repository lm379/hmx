(function () {
    let prefix = '';
    let Key = '';

    // 对更多字符编码的 url encode 格式
    const camSafeUrlEncode = function (str) {
        return encodeURIComponent(str)
            .replace(/!/g, '%21')
            .replace(/'/g, '%27')
            .replace(/\(/g, '%28')
            .replace(/\)/g, '%29')
            .replace(/\*/g, '%2A');
    };

    // 获取权限策略
    const getAuthorization = function (opt, callback) {
        // 替换为自己服务端地址 获取post上传签名，demo：https://github.com/tencentyun/cos-demo/blob/main/server/upload-sign/nodejs/app.js
        const url = `${window.location.origin}/api/sts/post-policy?ext=${opt.ext}`;
        const xhr = new XMLHttpRequest();
        xhr.open('GET', url, true);
        xhr.onload = function (e) {
            let credentials;
            try {
                const result = JSON.parse(e.target.responseText);
                credentials = result;
            } catch (e) {
                callback('获取签名出错');
            }
            if (credentials) {
                callback(null, {
                    securityToken: credentials.sessionToken,
                    cosKey: credentials.cosKey,
                    cosHost: credentials.cosHost,
                    policy: credentials.policy,
                    qAk: credentials.qAk,
                    qKeyTime: credentials.qKeyTime,
                    qSignAlgorithm: credentials.qSignAlgorithm,
                    qSignature: credentials.qSignature,
                });
            } else {
                console.error(xhr.responseText);
                callback('获取签名出错');
            }
        };
        xhr.send();
    };

    // 上传文件
    const uploadFile = function (file, callback, onProgress) {
        if (!file) {
            callback('文件不能为空');
            return;
        }
        const fileName = file.name;
        // 获取文件后缀名
        let ext = '';
        const lastDotIndex = fileName.lastIndexOf('.');
        if (lastDotIndex > -1) {
            // 这里获取文件后缀 由服务端生成最终上传的路径
            ext = fileName.substring(lastDotIndex + 1);
        }
        getAuthorization({ ext }, function (err, credentials) {
            if (err) {
                alert(err);
                return;
            }
            const protocol =
                location.protocol === 'https:' ? 'https:' : 'http:';
            prefix = protocol + '//' + credentials.cosHost;
            Key = credentials.cosKey;
            const fd = new FormData();

            // 在当前目录下放一个空的 empty.html 以便让接口上传完成跳转回来
            fd.append('key', Key);

            // 使用 policy 签名保护格式
            credentials.securityToken &&
                fd.append('x-cos-security-token', credentials.securityToken);
            fd.append('q-sign-algorithm', credentials.qSignAlgorithm);
            fd.append('q-ak', credentials.qAk);
            fd.append('q-key-time', credentials.qKeyTime);
            fd.append('q-signature', credentials.qSignature);
            fd.append('policy', credentials.policy);

            // 文件内容，file 字段放在表单最后，避免文件内容过长影响签名判断和鉴权
            fd.append('file', file);

            // xhr
            const url = prefix;
            const xhr = new XMLHttpRequest();
            xhr.open('POST', url, true);
            xhr.upload.onprogress = function (e) {
                if (onProgress && typeof onProgress === 'function') {
                    const percent = Math.round((e.loaded / e.total) * 10000) / 100;
                    onProgress({
                        loaded: e.loaded,
                        total: e.total,
                        percent: percent
                    });
                }
            };
            xhr.onload = function () {
                if (Math.floor(xhr.status / 100) === 2) {
                    const ETag = xhr.getResponseHeader('etag');
                    callback(null, {
                        url:
                            prefix + '/' + camSafeUrlEncode(Key).replace(/%2F/g, '/'),
                        ETag: ETag,
                    });
                } else {
                    callback('文件 ' + Key + ' 上传失败，状态码：' + xhr.status);
                }
            };
            xhr.onerror = function () {
                callback(
                    '文件 ' + Key + ' 上传失败，请检查是否没配置 CORS 跨域规则'
                );
            };
            xhr.send(fd);
        });
    };

    window.uploadFile = uploadFile;
})();