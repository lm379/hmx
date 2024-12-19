// 服务器地址
const Server = `${window.location.origin}/api/`;

// 解析JWT Token
function parseJwt(token) {
    try {
        const base64Url = token.split('.')[1];
        const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
        const jsonPayload = decodeURIComponent(atob(base64).split('').map(function (c) {
            return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
        }).join(''));
        return JSON.parse(jsonPayload);
    } catch (e) {
        return null;
    }
}

// 检查用户是否为管理员
function isAdmin(token) {
    const payload = parseJwt(token);
    return payload && payload.role === 'Administrator';
}

// 接口请求函数
const HTTP = {
    // POST请求
    post: async function (endpoint, options, mode) {
        const url = Server + mode + '/' + endpoint;
        const token = Cookie.get('Token')
        const headers = {
            'Authorization': token,
            'Content-Type': 'application/json'
        };

        // 允许空的请求体
        const requestOptions = {
            method: 'POST',
            headers: headers,
        }

        if (options && Object.keys(options).length > 0) {
            requestOptions.body = JSON.stringify(options);
        }

        try {
            const response = await fetch(url, requestOptions);
            return await response.json();
        } catch (error) {
            console.error('Error:', error);
            throw error;
        }
    },

    // GET请求
    get: async function (endpoint, options, mode) {
        let url = new URL(Server + mode + '/' + endpoint);
        const token = Cookie.get('Token');
        const headers = {
            'Authorization': token
        };
        if (options) {
            Object.keys(options).forEach(key => url.searchParams.append(key, options[key]));
        }

        try {
            const response = await fetch(url, {
                method: 'GET',
                headers: headers
            });
            return await response.json();
        } catch (error) {
            console.error('Error:', error);
            throw error;
        }
    },

    // 检查json正确性
    check: function (string) {
        try { if (typeof JSON.parse(string) == "object") { return true; } } catch (e) { }
        return false;
    },

    // 组合Post请求的FormData
    makeBody: function (options) {
        var form = new FormData();
        if (options) {
            for (let i = 0; i < options.length; i++) {
                var key = options[i].Key,
                    val = options[i].Value;
                form.append(key, val)
            }
        }
        return form;
    },
}


// 浏览器Cookie操作
const Cookie = {
    // 设置Cookie
    set: function (name, value) {
        var Days = 3;
        var exp = new Date();
        exp.setTime(exp.getTime() + Days * 24 * 60 * 60 * 1000);
        document.cookie = name + "=" + value + ";expires=" + exp.toGMTString();
    },

    //读取cookies
    get: function (name) {
        var arr, reg = new RegExp("(^| )" + name + "=([^;]*)(;|$)");
        if (arr = document.cookie.match(reg))
            return arr[2];
        else
            return null;
    },

    //删除cookies
    del: function (name) {
        var exp = new Date();
        exp.setTime(exp.getTime() - 1);
        var cval = Cookie.get(name);
        if (cval != null)
            document.cookie = name + "=" + cval + ";expires=" + exp.toGMTString();
    }
}
