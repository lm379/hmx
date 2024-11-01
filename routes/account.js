var connection = require('../dataBase/db');
var jwt = require('jsonwebtoken');
var bcrypt = require('bcryptjs');
var express = require('express');
var { getUserInfo } = require('../dataBase/api');
const app = express.Router();
require('dotenv').config();
var crypto = require('crypto');

// 生成随机密钥
const secretKey = crypto.randomBytes(32).toString('hex');

// 手机号输入校验
const checkPhone = (Phone) => {
    if (!Phone) return false;
    var reg = /^1[3456789]\d{9}$/;
    if (!(reg.test(Phone))) return false;
    return true;
}

// 邮箱输入校验
const checkEmail = (Email) => {
    if (!Email) return false;
    var reg = /^(\w-*\.*)+@(\w-?)+(\.\w{2,})+$/;
    if (!(reg.test(Email))) return false;
    return true;
}

// 用户名输入校验
const checkUsername = (Username) => {
    if (!Username) return false;
    var reg = /^[a-zA-Z0-9_-]{4,10}$/;  // 4到10位（字母，数字，下划线，减号）
    if (!(reg.test(Username))) return false;
    return true;
};

// 密码输入校验
const checkPassword = (Password) => {
    if (!Password) return false;
    var reg = /^(?=.*[a-zA-Z])(?=.*\d|.*[\W_]).{8,16}$/;  // 8-16位的大小写字母、数字和符号组成
    if (!(reg.test(Password))) return false;
    return true;
};

// 性别输入校验
const checkSex = (Sex) => {
    if (!Sex) return false;
    if (Sex === 'Male' || Sex === 'Female' || Sex === 'Other') return true;
    else return false;
}

// 用户登录
const userLogin = (Phone, Username, Email, Password) => {
    return new Promise((resolve, reject) => {
        let query = "SELECT * FROM `users` WHERE ";
        if (!Password) return reject(new Error("请输入密码"));
        let params = [];

        // 判断登录方式，支持手机号、用户名、邮箱登录
        if (Phone) {
            query += `Phone = ?`;
            params.push(Phone);
        } else if (Username) {
            query += `Username = ?`;
            params.push(Username);
        } else if (Email) {
            query += `Email = ?`;
            params.push(Email);
        } else {
            return reject(new Error("参数错误"));
        }
        // console.log(query, params);
        connection.query(query, params, (err, data) => {
            if (err) return reject(err);
            if (data.length === 0) return reject(resolve({ code: 1, msg: "用户名或密码错误" }));
            const user = data[0];
            // 验证密码
            bcrypt.compare(Password, user.Password, (err, isMatch) => {
                if (err) return reject(err);
                if (isMatch) {
                    // 生成token
                    const token = jwt.sign({ User_id: user.User_id, username: user.Username, Role: user.Role }, secretKey, { expiresIn: 3600 });
                    resolve({ code: 0, token: token, msg: "登录成功" });
                } else {
                    reject({ code: 1, msg: "用户名或密码错误" });
                }
            });
        });
    });
};

app.post('/login', async (req, res, next) => {
    // const { Phone, Username, Email, Password } = req.body;
    const Phone = req.body.phone;
    const Username = req.body.username;
    const Email = req.body.email;
    const Password = req.body.password;
    if (!Phone && !Username && !Email) {
        return res.status(400).json({ code: 1, message: 'invalid parameters' });
    }
    if (!Password) {
        return res.status(400).json({ code: 1, message: 'invalid password' });
    }
    // 判断密码格式是否正确
    if (!checkPassword(Password)) {
        return res.status(400).json({ code: 1, message: '密码格式错误' });
    }
    // 判断账号是否已注册和输入是否合规
    if (Phone) {
        if (!checkPhone(Phone)) {
            return res.status(400).json({ code: 1, message: '手机号格式错误' });
        }
        const data = await getUserInfo(null, Phone, null, null);
        if (data.code === 1) {
            return res.status(400).json(data);
        }
    }
    else if (Email) {
        if (!checkEmail(Email)) {
            return res.status(400).json({ code: 1, message: '邮箱格式错误' });
        }
        const data = await getUserInfo(null, null, Email, null);
        if (data.code === 1) {
            return res.status(400).json(data);
        }
    }
    else if (Username) {
        if (!checkUsername(Username)) {
            return res.status(400).json({ code: 1, message: '用户名格式错误' });
        }
        const data = await getUserInfo(null, null, null, Username);
        if (data.code === 1) {
            return res.status(400).json(data);
        }
    }

    // 登录
    try {
        const data = await userLogin(Phone, Username, Email, Password);
        res.json(data);
    } catch (err) {
        console.error(err); // 记录错误日志
        res.status(500).json(err);
    }
});

// 用户注册
const userRegister = async (Username, Phone, Email, Password) => {
    if (!Username) throw { code: 1, msg: "请输入用户名" };
    if (!Phone) throw { code: 1, msg: "请输入手机号" };
    if (!Password) throw { code: 1, msg: "请输入密码" };

    // 判断输入是否合规
    if (!checkUsername(Username)) throw { code: 1, msg: "用户名格式错误" };
    if (!checkPhone(Phone)) throw { code: 1, msg: "手机号格式错误" };
    if (Email && !checkEmail(Email)) throw { code: 1, msg: "邮箱格式错误" };
    if (!checkPassword(Password)) throw { code: 1, msg: "密码格式错误" };

    // 查询手机号是否已存在
    let data = await getUserInfo(null, Phone, null, null);
    if (data.code === 0) throw { code: 1, msg: "该手机号已注册" };

    // 查询邮箱是否已存在
    data = await getUserInfo(null, null, Email, null);
    if (data.code === 0) throw { code: 1, msg: "该邮箱已注册" };

    // 查询用户名是否已存在
    data = await getUserInfo(null, null, null, Username);
    if (data.code === 0) throw { code: 1, msg: "该用户名已注册" };

    // 加密密码
    const salt = await bcrypt.genSalt(10);
    const hash = await bcrypt.hash(Password, salt);

    // 将用户添加到数据库
    return new Promise((resolve, reject) => {
        connection.query("INSERT INTO `users` SET ?", { Username, Phone, Email, Password: hash }, (err, data) => {
            if (err) return reject({ code: 1, msg: "注册失败" });
            resolve({ code: 0, msg: "注册成功" });
        });
    });
};

app.post('/register', async (req, res, next) => {
    // const { Username, Phone, Email, Password } = req.body;
    const Username = req.body.username;
    const Phone = req.body.phone;
    const Email = req.body.email;
    const Password = req.body.password;
    if (!Username) {
        return res.status(400).json({ code: 1, message: 'invalid username' });
    }
    if (!Phone) {
        return res.status(400).json({ code: 1, message: 'invalid phone' });
    }
    if (!Password) {
        return res.status(400).json({ code: 1, message: 'invalid password' });
    }
    // 注册
    try {
        const data = await userRegister(Username, Phone, Email, Password);
        if (data.code === 1) {
            return res.status(400).json(data);
        }
        res.json(data);
    } catch (err) {
        console.error(err); // 记录错误日志
        res.status(500).json(err);
    }
});

// 重置密码
const resetPassword = async (Phone, Email, Password) => {
    // if (!Phone && !Email) throw { code: 1, msg: "请输入手机号或邮箱" };
    if (!Password) throw { code: 1, msg: "请输入密码" };
    let query = "UPDATE `users` SET Password = ? WHERE ";

    // 判断输入是否合规
    if (!checkPassword(Password)) throw { code: 1, msg: "密码格式错误" };

    var params;
    if (Phone) {
        if (!checkPhone(Phone)) throw { code: 1, msg: "手机号格式错误" };
        // 查询手机号是否已注册
        let data = await getUserInfo(null, Phone, null, null);
        if (data.code === 1) throw { code: 1, msg: "该手机号未注册" };
        query += `Phone = ?`;
        params = Phone;
    } else if (Email) {
        if (!checkEmail(Email)) throw { code: 1, msg: "邮箱格式错误" };
        // 查询邮箱是否已注册
        let data = await getUserInfo(null, null, Email, null);
        if (data.code === 1) throw { code: 1, msg: "该邮箱未注册" };
        query += `Email = ?`;
        params = Email;
    } else {
        throw { code: 1, msg: "参数错误" };
    }

    // 加密密码
    const salt = await bcrypt.genSalt(10);
    const hash = await bcrypt.hash(Password, salt);

    // 更新密码
    return new Promise((resolve, reject) => {
        connection.query(query, [hash, params], (err, data) => {
            if (err) return reject({ code: 1, msg: "重置密码失败" });
            resolve({ code: 0, msg: "重置密码成功" });
        });
    });
};

app.post('/reset_password', async (req, res, next) => {
    // const { Phone, Email, Password } = req.body;
    const Phone = req.body.phone;
    const Email = req.body.email;
    const Password = req.body.password;
    if (!Phone && !Email) {
        return res.status(400).json({ code: 1, message: 'invalid parameters' });
    }
    if (!Password) {
        return res.status(400).json({ code: 1, message: 'invalid password' });
    }

    // 重置密码
    try {
        const data = await resetPassword(Phone, Email, Password);
        res.json(data);
    } catch (err) {
        console.error(err); // 记录错误日志
        res.status(500).json(err);
    }
});

// 修改信息
const ModifyUserInfo = async (User_id, ModifyFields, Value) => {
    const user = await getUserInfo(User_id, null, null, null);
    if (user.code === 1) throw { code: 1, msg: "用户不存在" };

    if (ModifyFields !== 'Username' && ModifyFields !== 'Phone' && ModifyFields !== 'Email' && ModifyFields !== 'Icon' && ModifyFields !== 'Sex') {
        throw { code: 1, msg: "参数错误" };
    }
    if (ModifyFields === 'Phone') {
        if (!checkPhone(Value)) throw { code: 1, msg: "手机号格式错误" };
        if (Value === user.data[0].Phone) throw { code: 1, msg: "手机号未修改" };
        const data = await getUserInfo(null, Value, null, null);
        if (data.code === 0) throw { code: 1, msg: "该手机号已注册" };
    }
    if (ModifyFields === 'Email') {
        if (!checkEmail(Value)) throw { code: 1, msg: "邮箱格式错误" };
        if (Value === user.data[0].Email) throw { code: 1, msg: "邮箱未修改" };
        const data = await getUserInfo(null, null, Value, null);
        if (data.code === 0) throw { code: 1, msg: "该邮箱已注册" };
    }
    if (ModifyFields === 'Username') {
        if (!checkUsername(Value)) throw { code: 1, msg: "用户名格式错误" };
        if (Value === user.data[0].Username) throw { code: 1, msg: "用户名未修改" };
        const data = await getUserInfo(null, null, null, Value);
        if (data.code === 0) throw { code: 1, msg: "该用户名已注册" };
    }
    if (ModifyFields === 'Sex') {
        if (!checkSex(Value)) throw { code: 1, msg: "性别格式错误" };
    }

    let query = `UPDATE \`users\` SET \`${ModifyFields}\` = ? WHERE User_id = ?`;
    return new Promise((resolve, reject) => {
        connection.query(query, [Value, User_id], (err, data) => {
            if (err) return reject({ code: 1, msg: "修改失败" });
            resolve({ code: 0, msg: "修改成功" });
        });
    });
}

const ModifyPassword = async (User_id, OldPassword, NewPassword) => {
    try {
        const user = await getUserInfo(User_id, null, null, null);
        if (user.code === 1) throw { code: 1, msg: "用户不存在" };
        const currentPassword = user.data[0].Password; // 当前用户密码
        if (!OldPassword) throw { code: 1, msg: "请输入旧密码" };
        if (!NewPassword) throw { code: 1, msg: "请输入新密码" };

        // 判断密码格式是否正确
        if (!checkPassword(OldPassword) || !checkPassword(NewPassword)) throw { code: 1, msg: "密码格式错误" };

        // 确保用户信息中的密码字段有效
        if (!currentPassword) throw { code: 1, msg: "用户密码不存在" };

        // 验证密码
        const isMatch = bcrypt.compare(OldPassword, currentPassword);
        if (!isMatch) throw { code: 1, msg: "旧密码错误" };

        // 加密新密码
        const salt = await bcrypt.genSalt(10);
        const hash = await bcrypt.hash(NewPassword, salt);

        // 更新密码
        return new Promise((resolve, reject) => {
            connection.query("UPDATE `users` SET Password = ? WHERE User_id = ?", [hash, User_id], (err, data) => {
                if (err) return reject({ code: 1, msg: "修改失败" });
                resolve({ code: 0, msg: "修改成功" });
            });
        });
    } catch (err) {
        throw err;
    };
}

app.post('/update_user', async (req, res, next) => {
    const token = req.headers['authorization'];
    if (!token) {
        return res.status(401).json({ code: 1, message: 'invalid token' });
    }
    let User_id;
    // 验证token
    try {
        const decoded = jwt.verify(token, secretKey);
        // console.log(decoded);
        User_id = decoded.User_id;
    } catch {
        return res.status(401).json({ code: 1, message: 'invalid token' });
    }

    // const { ModifyFields, Value, OldPassword, NewPassword } = req.body;
    const ModifyFields = req.body.modifyfields;
    const Value = req.body.value;
    const OldPassword = req.body.oldpassword;
    const NewPassword = req.body.newpassword;

    if (!User_id) {
        return res.status(400).json({ code: 1, message: 'invalid user_id' });
    }
    if (ModifyFields) {
        try {
            const data = await ModifyUserInfo(User_id, ModifyFields, Value);
            res.json(data);
        } catch (err) {
            console.error(err); // 记录错误日志
            res.status(500).json(err);
        }
    } else if (OldPassword && NewPassword) {
        try {
            const data = await ModifyPassword(User_id, OldPassword, NewPassword);
            res.json(data);
        } catch (err) {
            console.error(err); // 记录错误日志
            res.status(500).json(err);
        }
    } else {
        return res.status(400).json({ code: 1, message: 'invalid parameters' });
    }
});

// 注销账号
app.post('/delete_account', async (req, res, next) => {
    const token = req.headers['authorization'];
    if (!token) {
        return res.status(401).json({ code: 1, message: 'invalid token' });
    }
    let User_id;
    // 验证token
    try {
        const decoded = jwt.verify(token, secretKey);
        User_id = decoded.User_id;
    } catch {
        return res.status(401).json({ code: 1, message: 'invalid token' });
    }

    if (!User_id) {
        return res.status(400).json({ code: 1, message: 'invalid user_id' });
    }

    return new Promise((resolve, reject) => {
        connection.query("UPDATE `users` SET Status = 'Deleted' WHERE User_id = ?", User_id, (err, data) => {
            if (err) return reject({ code: 1, msg: "注销失败" });
            resolve({ code: 0, msg: "注销成功" });
        });
    }).then(data => {
        res.json(data);
    }).catch(err => {
        console.error(err); // 记录错误日志
        res.status(500).json(err);
    });
});

module.exports = app;
module.exports.secretKey = secretKey;