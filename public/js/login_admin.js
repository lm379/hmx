// 自动登陆验证，检测Token是否存在，如果存在，则使用该Token去登录
Cookie.set('Admin', true)
let token = Cookie.get('Token');
// if (token) {
//     HTTP.post('admin_login', { 'Mode': 'user' }).then(json => {
//         if (json.code == 0) { loginDone('admin') }
//     })
// }

// 回车键登陆
function loginEnter(e) {
    let event = e || window.event;
    if (event.keyCode == 13) { login() }
}

// 管理员登录
function login() {
    // 获取输入的账号密码
    let Account = document.getElementById('Account').value,
        Password = document.getElementById('Password').value;

    // 获取记住账号是否勾选
    let Remember = document.getElementById('remember').checked;

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
            // 管理员登录
            HTTP.post('login_admin', options, 'user').then(json => {
                if (json.code === 0) {
                    // 保存Token，以方便日后自动登录
                    Cookie.set('Token', json.token);
                    // 如果勾选记住账号，那么保存账号，否则删除保存的账号
                    if (Remember) { Cookie.set('UserID', Account) }
                    else { Cookie.del('UserID') }
                    // 登录完成
                    loginDone('admin')
                }
                else { alert('登录失败!\n错误：' + json.msg) }
            })
        }
    }
}

// 登录完成后，跳转到后台
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