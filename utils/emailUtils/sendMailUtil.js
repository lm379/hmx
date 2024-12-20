var nodemailer = require('nodemailer');
const { getVerificationCode } = require('../redisUtils/redisAccountUtil');

let sendMails = (email, subject, message) => {
    let transporter = nodemailer.createTransport({
        host: process.env.EMAIL_HOST,
        port: process.env.EMAIL_PORT,
        secure: true,
        auth: {
            user: process.env.EMAIL_USER,
            pass: process.env.EMAIL_PASS
        }
    });

    let mailOptions = {
        from: `"黄梅戏文化平台" <${process.env.EMAIL_USER}>`,
        to: email,
        subject: subject,
        text: message
    };

    return new Promise((resolve, reject) => {

        transporter.sendMail(mailOptions, (error, info) => {
            if (error) {
                reject({ code: 1, msg: "邮箱验证码发送失败" });
            } else {
                resolve({ code: 0, msg: "邮箱验证码发送成功，请注意查收" });
            }
        });
    });
};

module.exports = sendMails;
