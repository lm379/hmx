// 自动登陆验证，检测Token是否存在，如果存在，则使用该Token去登录
Cookie.set('Admin', false)
let token = Cookie.get('Token');
if (token) {
    HTTP.post('login', { 'Mode': 'user' }).then(json => {
        if (json.code == 0) { loginDone('index') }
    })
}

// 回车键登陆
function loginEnter(e) {
    let event = e || window.event;
    if (event.keyCode == 13) { login() }
}

// 登陆
function login() {
    // 获取输入的账号密码
    let Account = document.getElementById('Account').value,
        Password = document.getElementById('Password').value;

    // 获取记住账号是否勾选
    let Remember = document.getElementById('remember').checked;
    let Mode = 'user';
    // 账号密码不为空的情况下才允许登录
    if (Account != '' && Password != '') {
        // 检查登录方式
        var LoginMethod;
        if (checkPhone(Account)) {
            LoginMethod = 'phone'
        } else if (checkEMail(Account)) {
            LoginMethod = 'email'
        } else if (checkUsername(Account)) {
            LoginMethod = 'username'
        } else {
            alert('账号格式有误')
        }
        if (!checkPassword(Password)) {
            alert('密码格式有误')
        } else {
            let options = { [LoginMethod]: Account, 'password': Password };
            HTTP.post('login', options, Mode).then(json => {
                if (json.code === 0) {
                    // 保存Token，以方便日后自动登录
                    Cookie.set('Token', json.token);
                    // 如果勾选记住账号，那么保存账号，否则删除保存的账号
                    if (Remember) {
                        Cookie.set('UserID', Account);
                    }
                    else { Cookie.del('UserID') }
                    // 登录完成
                    loginDone('index');
                }
                else { alert('登录失败!\n错误：' + json.msg) }
            })
        }
    }
}


// 登录完成后，跳转到首页
function loginDone(page) {
    window.location.href = page + ".html"
}

// 设置记住账号的勾选框是否勾选，以及用户名输入框是否显示记住的账号
function setRemember() {
    let Account = Cookie.get('UserID'),
        Remember = document.getElementById('remember'),
        Account_Input = document.getElementById('Account');

    // 如果存在记录的账号，那么显示账号，并且勾选记住账号的选框
    if (Account) {
        Remember.checked = true;
        Account_Input.value = Account;
    }
    // 否则就取消记住账号的勾选
    else { Remember.checked = false; }
}

// 注册用户
function signup() {
    let Username = document.getElementById('Username').value,
        Password = document.getElementById('Password').value,
        PasswordCheck = document.getElementById('PasswordCheck').value,
        Phone = document.getElementById('Phone').value,
        EMail = document.getElementById('EMail').value;
    // 检查内容是否填写完整
    if (Username == '' || Password == '' || PasswordCheck == '' || Phone == '' || EMail == '') {
        alert('您还有内容尚未填写，请填写完整注册信息')
    }

    // 检查两次输入的密码是否相同
    else if (Password != PasswordCheck) { alert('两次输入的密码不一致') }

    // 检查输入是否合法
    else if (!checkUsername(Username)) { alert('请输入正确的用户名') }

    else if (!checkPhone(Phone)) { alert('请输入正确的手机号码') }

    else if (!checkEMail(EMail)) { alert('请输入正确的邮箱地址') }

    else if (!checkPassword(Password)) { alert('密码格式有误') }

    // 通过以上的验证后，开始调用注册接口
    else {
        let options = { 'username': Username, 'password': Password, 'phone': Phone, 'email': EMail };
        HTTP.post('register', options, 'user').then(json => {
            if (json.code == 0) {
                alert('注册成功!')
                window.location.href = "login.html"
            }
            if (json.code == 1) {
                alert(json.msg)
            }
        })
    }
}

// 发送验证码
function SendCode() {
    let Account = document.getElementById('Account').value;
    if (!Account) {
        alert('请输入手机号码或邮箱地址');
    }
    let options = {};
    let type = '';
    // 自动判断是手机号码还是邮箱地址
    if (checkPhone(Account)) {
        options = { 'phone': Account };
        type = 'send_sms';
    } else if (checkEMail(Account)) {
        options = { 'email': Account };
        type = 'send_email';
    } else {
        alert('请输入正确的手机号码或邮箱地址');
        return; // 阻止发送验证码
    }
    HTTP.post(type, options, 'user').then(json => {
        if (json.code == 0) {
            alert(json.msg)
        } else {
            alert(json.msg)
        }
    });
}

// 忘记密码，修改新密码
function forget() {
    let Account = document.getElementById('Account').value,
        Username = document.getElementById('Username').value,
        VerifyCode = document.getElementById('VerifyCode').value,
        Password = document.getElementById('Password').value;

    // 检查内容是否填写完整
    if (Account == '' || Password == '' || Username == '' || VerifyCode == '') {
        alert('您还有内容尚未填写，请填写完整信息')
    } else {
        if (!checkUsername(Username)) { alert('请输入正确的用户名') }
        if (!checkPassword(Password)) { alert('密码格式有误') }
        let options = {};
        if (checkPhone(Account)) {
            options = { 'phone': Account, 'username': Username, 'password': Password, 'code': VerifyCode };
        } else if (checkEMail(Account)) {
            options = { 'email': Account, 'username': Username, 'password': Password, 'code': VerifyCode };
        } else {
            alert('请输入正确的手机号码或邮箱地址');
            return; // 阻止后续继续执行
        }
        HTTP.post('reset_password', options, 'user').then(json => {
            if (json.code != 0) { alert(json.msg) }
            else {
                alert('重置密码成功')
                window.location.href = 'login.html'
            }
        })
    }
}

// 使用正则表达式，检查手机号码是否正确
function checkPhone(phone) {
    let pattern = /^1[3456789]\d{9}$/;
    if (!pattern.test(phone)) { return false; }
    else { return true; }
}

// 使用正则表达式，检查邮箱是否正确
function checkEMail(email) {
    let pattern = /^([A-Za-z0-9_\-\.\u4e00-\u9fa5])+\@([A-Za-z0-9_\-\.])+\.([A-Za-z]{2,8})$/;
    if (pattern.test(email)) { return true; }
    else { return false; }
}

// 使用正则表达式，检查密码是否符合要求
function checkPassword(password) {
    var pattern = /^(?=.*[a-zA-Z])(?=.*\d|.*[\W_]).{8,16}$/;
    if (pattern.test(password)) { return true; }
    else { return false; }
}

// 使用正则表达式，检查用户名是否符合要求
function checkUsername(account) {
    var pattern = /^[a-zA-Z0-9_-]{4,10}$/;
    if (pattern.test(account)) { return true; }
    else { return false; }
}