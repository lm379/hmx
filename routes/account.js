var connection = require('../dataBase/db');
var jwt = require('jsonwebtoken');
var bcrypt = require('bcryptjs');
var express = require('express');
var { getUserInfo } = require('../dataBase/api');
var secretKey = require('../utils/secretKey');
const router = express.Router();
require('dotenv').config();
// var crypto = require('crypto');

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
                    const token = jwt.sign({ User_id: user.User_id, Username: user.Username, Phone: user.Phone, Email: user.Email, Role: user.Role }, secretKey, { expiresIn: 3600 });
                    resolve({ code: 0, token: token, msg: "登录成功" });
                } else {
                    reject({ code: 1, msg: "用户名或密码错误" });
                }
            });
        });
    });
};

// 管理员登录
const adminLogin = (Phone, Username, Email, Password) => {
    return new Promise((resolve, reject) => {
        let query = "SELECT * FROM `users` WHERE Role = 'Administrator' AND ";
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
        connection.query(query, params, (err, data) => {
            if (err) return reject(err);
            if (data.length === 0) return reject(resolve({ code: 1, msg: "用户名或密码错误" }));
            const user = data[0];
            // 验证密码
            bcrypt.compare(Password, user.Password, (err, isMatch) => {
                if (err) return reject(err);
                if (isMatch) {
                    // 生成token
                    const token = jwt.sign({ User_id: user.User_id, Username: user.Username, Phone: user.Phone, Email: user.Email, Role: user.Role }, secretKey, { expiresIn: 3600 });
                    resolve({ code: 0, token: token, msg: "登录成功" });
                } else {
                    reject({ code: 1, msg: "用户名或密码错误" });
                }
            });
        });
    });
};

router.post('/login', async (req, res, next) => {
    // const { Phone, Username, Email, Password } = req.body;
    const Phone = req.body.phone;
    const Username = req.body.username;
    const Email = req.body.email;
    const Password = req.body.password;
    if (!Phone && !Username && !Email) {
        return res.status(400).json({ code: 1, msg: 'invalid parameters' });
    }
    if (!Password) {
        return res.status(400).json({ code: 1, msg: 'invalid password' });
    }
    // 判断密码格式是否正确
    if (!checkPassword(Password)) {
        return res.status(400).json({ code: 1, msg: '密码格式错误' });
    }
    // 判断账号是否已注册和输入是否合规
    if (Phone) {
        if (!checkPhone(Phone)) {
            return res.status(400).json({ code: 1, msg: '手机号格式错误' });
        }
        const data = await getUserInfo(null, Phone, null, null);
        if (data.code === 1) {
            return res.status(400).json(data);
        }
    }
    else if (Email) {
        if (!checkEmail(Email)) {
            return res.status(400).json({ code: 1, msg: '邮箱格式错误' });
        }
        const data = await getUserInfo(null, null, Email, null);
        if (data.code === 1) {
            return res.status(400).json(data);
        }
    }
    else if (Username) {
        if (!checkUsername(Username)) {
            return res.status(400).json({ code: 1, msg: '用户名格式错误' });
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

// 管理员登录
router.post('/login_admin', async (req, res, next) => {
    const Phone = req.body.phone;
    const Username = req.body.username;
    const Email = req.body.email;
    const Password = req.body.password;
    if (!Phone && !Username && !Email) {
        return res.status(400).json({ code: 1, msg: 'invalid parameters' });
    }
    if (!Password) {
        return res.status(400).json({ code: 1, msg: 'invalid password' });
    }
    // 判断密码格式是否正确
    if (!checkPassword(Password)) {
        return res.status(400).json({ code: 1, msg: '密码格式错误' });
    }
    // 判断账号是否已注册和输入是否合规
    if (Phone) {
        if (!checkPhone(Phone)) {
            return res.status(400).json({ code: 1, msg: '手机号格式错误' });
        }
        const data = await getUserInfo(null, Phone, null, null);
        if (data.code === 1) {
            return res.status(400).json(data);
        }
        if (data.data[0].Role !== 'Administrator') {
            return res.status(400).json({ code: 1, msg: '权限错误' });
        }
    } else if (Email) {
        if (!checkEmail(Email)) {
            return res.status(400).json({ code: 1, msg: '邮箱格式错误' });
        }
        const data = await getUserInfo(null, null, Email, null);
        if (data.code === 1) {
            return res.status(400).json(data);
        }
        if (data.data[0].Role !== 'Administrator') {
            return res.status(400).json({ code: 1, msg: '权限错误' });
        }
    } else if (Username) {
        if (!checkUsername(Username)) {
            return res.status(400).json({ code: 1, msg: '用户名格式错误' });
        }
        const data = await getUserInfo(null, null, null, Username);
        if (data.code === 1) {
            return res.status(400).json(data);
        }
        if (data.data[0].Role !== 'Administrator') {
            return res.status(400).json({ code: 1, msg: '权限错误' });
        }
    }

    // 管理员登录
    try {
        const data = await adminLogin(Phone, Username, Email, Password);
        if (data.length === 0) {
            res.status(400).json({ code: 1, msg: '用户名或密码错误' });
        }
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

router.post('/register', async (req, res, next) => {
    // const { Username, Phone, Email, Password } = req.body;
    const Username = req.body.username;
    const Phone = req.body.phone;
    const Email = req.body.email;
    const Password = req.body.password;
    if (!Username) {
        return res.status(400).json({ code: 1, msg: 'invalid username' });
    }
    if (!Phone) {
        return res.status(400).json({ code: 1, msg: 'invalid phone' });
    }
    if (!Password) {
        return res.status(400).json({ code: 1, msg: 'invalid password' });
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
    let query = "UPDATE `users` SET Password = ? WHERE Role <> 'Administrator'";

    // 判断输入是否合规
    if (!checkPassword(Password)) throw { code: 1, msg: "密码格式错误" };

    var params;
    if (Phone) {
        if (!checkPhone(Phone)) throw { code: 1, msg: "手机号格式错误" };
        // 查询手机号是否已注册
        let data = await getUserInfo(null, Phone, null, null);
        if (data.code === 1) throw { code: 1, msg: "该手机号未注册" };
        if (data.data[0].Role === 'Administrator') throw { code: 1, msg: "禁止重置管理员密码" };
        query += `AND Phone = ?`;
        params = Phone;
    } else if (Email) {
        if (!checkEmail(Email)) throw { code: 1, msg: "邮箱格式错误" };
        // 查询邮箱是否已注册
        let data = await getUserInfo(null, null, Email, null);
        if (data.code === 1) throw { code: 1, msg: "该邮箱未注册" };
        if (data.data[0].Role === 'Administrator') throw { code: 1, msg: "禁止重置管理员密码" };
        query += `AND Email = ?`;
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
            if (data.length === 0) { return reject({ code: 1, msg: "该手机号或邮箱未注册" }); }
            resolve({ code: 0, msg: "重置密码成功" });
        });
    });
};

router.post('/reset_password', async (req, res, next) => {
    // const { Phone, Email, Password } = req.body;
    const Phone = req.body.phone;
    const Email = req.body.email;
    const Password = req.body.password;
    if (!Phone && !Email) {
        return res.status(400).json({ code: 1, msg: 'invalid parameters' });
    }
    if (!Password) {
        return res.status(400).json({ code: 1, msg: 'invalid password' });
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
const ModifyUserInfo = async (User_id, params) => {
    const user = await getUserInfo(User_id, null, null, null);
    if (user.code === 1) throw { code: 1, msg: "用户不存在" };
    if (params.length === 0) throw { code: 1, msg: "参数错误" };

    // 动态构建SQL语句
    let query = "UPDATE `users` SET ";
    const values = [];

    // 遍历Json数组获得参数
    for (const param of params) {
        for (const [key, value] of Object.entries(param)) {
            if (key === 'Username') {
                if (!checkUsername(value)) throw { code: 1, msg: "用户名格式错误" };
                if (value === user.data[0].Username) throw { code: 1, msg: "用户名未修改" };
                const data = await getUserInfo(null, null, null, value);
                if (data.code === 0) throw { code: 1, msg: "该用户名已注册" };
                query += "Username = ?, ";
                values.push(value);
            }
            if (key === 'Phone') {
                if (!checkPhone(value)) throw { code: 1, msg: "手机号格式错误" };
                if (value === user.data[0].Phone) throw { code: 1, msg: "手机号未修改" };
                const data = await getUserInfo(null, value, null, null);
                if (data.code === 0) throw { code: 1, msg: "该手机号已注册" };
                query += "Phone = ?, ";
                values.push(value);
            }
            if (key === 'Email') {
                if (!checkEmail(value)) throw { code: 1, msg: "邮箱格式错误" };
                if (value === user.data[0].Email) throw { code: 1, msg: "邮箱未修改" };
                const data = await getUserInfo(null, null, value, null);
                if (data.code === 0) throw { code: 1, msg: "该邮箱已注册" };
                query += "Email = ?, ";
                values.push(value);
            }
            if (key === 'Icon') {
                query += "Icon = ?, ";
                values.push(value);
            }
            if (key === 'Sex') {
                if (!checkSex(value)) throw { code: 1, msg: "性别格式错误" };
                if (value === user.data[0].Sex) throw { code: 1, msg: "性别未修改" };
                query += "Sex = ?, ";
                values.push(value);
            }
        }
    }
    
    query = query.slice(0, -2);  // 去掉最后的逗号和空格
    query += " WHERE User_id = ?";
    values.push(User_id);

    return new Promise((resolve, reject) => {
        connection.query(query, values, (err, data) => {
            if (err) return reject({ code: 1, msg: "修改失败" });
            resolve({ code: 0, msg: "修改成功" });
        });
    });
}

// 修改密码
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

router.post('/modify_password', async (req, res, next) => {
    const token = req.headers['authorization'];
    if (!token) {
        return res.status(401).json({ code: 1, msg: 'invalid token' });
    }
    let User_id;
    let user;
    // 验证token
    try {
        const decoded = jwt.verify(token, secretKey);
        User_id = decoded.User_id;
        user = decoded;
    } catch {
        return res.status(401).json({ code: 1, msg: 'invalid token' });
    }
    const OldPassword = req.body.oldpassword;
    const NewPassword = req.body.newpassword;
    // 校验输入的密码
    if (!OldPassword || !NewPassword) {
        return res.status(400).json({ code: 1, msg: 'invalid oldpassword' });
    }
    if (!checkPassword(OldPassword) || !checkPassword(NewPassword)) {
        return res.status(400).json({ code: 1, msg: 'invalid password' });
    }
    try {
        const data = await ModifyPassword(User_id, OldPassword, NewPassword);
        // 生成新的token
        const token = jwt.sign({ User_id: user.User_id, Username: user.Username, Phone: user.Phone, Email: user.Email, Role: user.Role }, secretKey, { expiresIn: 3600 });
        res.json({ code: 0, token: token, msg: "修改成功" });
    } catch (err) {
        console.error(err); // 记录错误日志
        res.status(500).json(err);
    }
});

router.post('/update_user', async (req, res, next) => {
    const token = req.headers['authorization'];
    if (!token) {
        return res.status(401).json({ code: 1, msg: 'invalid token' });
    }
    let User_id;
    // 验证token
    try {
        const decoded = jwt.verify(token, secretKey);
        User_id = decoded.User_id;
    } catch {
        return res.status(401).json({ code: 1, msg: 'invalid token' });
    }

    const Username = req.body.username;
    const Phone = req.body.phone;
    const Email = req.body.email;
    const Icon = req.body.icon;
    const Sex = req.body.Sex;

    if (!User_id) {
        return res.status(400).json({ code: 1, msg: 'invalid user_id' });
    }
    if (Username || Phone || Email || Icon) {
        let params = [];  // 定义json数组用于存储需要变更的信息
        if (Username) {
            if (!checkUsername(Username)) {
                return res.status(401).json({ code: 1, msg: "用户名格式错误" });
            }
            params.push({ "Username": Username });
        }
        if (Phone) {
            if (!checkPhone(Phone)) {
                return res.status(401).json({ code: 1, msg: "手机号格式错误" });
            }
            params.push({ "Phone": Phone });
        }
        if (Email) {
            if (!checkEmail(Email)) {
                return res.status(401).json({ code: 1, msg: "邮箱格式错误" });
            }
            params.push({ "Email": Email });
        }
        if (Icon) {
            params.push({ "Icon": Icon });
        }
        if (Sex) {
            if (!checkSex(Sex)) {
                return res.status(401).json({ code: 1, msg: "性别格式错误" });
            }
            params.push({ "Sex": Sex })
        }

        try {
            const data = await ModifyUserInfo(User_id, params);
            res.json(data);
        } catch (err) {
            console.error(err); // 记录错误日志
            res.status(500).json(err);
        }
    } else {
        return res.status(400).json({ code: 1, msg: 'invalid parameters' });
    }
});

// 注销账号
router.post('/delete_account', async (req, res, next) => {
    const token = req.headers['authorization'];
    if (!token) {
        return res.status(401).json({ code: 1, msg: 'invalid token' });
    }
    let User_id;
    // 验证token
    try {
        const decoded = jwt.verify(token, secretKey);
        User_id = decoded.User_id;
    } catch {
        return res.status(401).json({ code: 1, msg: 'invalid token' });
    }

    if (!User_id) {
        return res.status(400).json({ code: 1, msg: 'invalid user_id' });
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

router.post('/login_confirm', async (req, res, next) => {
    const token = req.headers['authorization'];
    if (!token) {
        return res.status(401).json({ code: 1, msg: 'invalid token' });
    }
    try {
        const decoded = jwt.verify(token, secretKey);
        const User_id = decoded.User_id;
        const Username = decoded.Username;
        const Role = decoded.Role;
        res.json({ code: 0, User_id, Username, Role, msg: "登录成功" });
    } catch {
        return res.status(401).json({ code: 1, msg: 'invalid token' });
    }
});

module.exports = router;