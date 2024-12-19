var nodemailer = require('nodemailer');

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
        from: process.env.EMAIL_USER,
        to: email,
        subject: subject,
        text: message
    };

    return new Promise((resolve, reject) => {
        transporter.sendMail(mailOptions, (error, info) => {
            if (error) {
                reject({ code: 1, msg: { code: 1, msg: "发送失败" } });
            } else {
                resolve({
                    code: 0,
                    accepted: info.accepted,
                    envelope: info.envelope,
                    msg: {
                        response: info.response,
                        messageId: info.messageId,
                    }
                });
            }
        });
    });
};

module.exports = sendMails;