const { get } = require('../routes');
const connection = require('./db');

// 列出用户列表
const ListUsers = () => {
    return new Promise((resolve, reject) => {
        connection.query("SELECT * FROM `users` WHERE Role = 'User' AND Status = 'Normal'", (err, data) => {
            if (err) return reject(err);
            return resolve({ code: 0, data: data });
        })
    })
}

// 获取用户信息
const getUserInfo = (User_id, Phone, Email, Username) => {
    return new Promise((resolve, reject) => {
        let query = "SELECT * FROM `users` WHERE TRUE ";
        let params = [];
        if (User_id) {
            query += `AND User_id = ?`;
            params.push(User_id);
        }

        if (Phone) {
            query += `AND Phone = ?`;
            params.push(Phone);
        }

        if (Email) {
            query += `AND Email = ?`;
            params.push(Email);
        }

        if (Username) {
            query += 'AND Username = ?';
            params.push(Username);
        }

        if (params.length === 0) {
            return reject(new Error("参数错误"));
        }

        connection.query(query, params, (err, data) => {
            if (err) return reject(err);
            if (data.length === 0) return reject(resolve({ code: 1, msg: "用户不存在" }));
            // console.log(data[0].Status);
            if (data[0].Status === 'Deleted') return reject(resolve({ code: 1, msg: "用户已注销" }));
            return resolve({ code: 0, data: data });
        });
    });
}

// 列出曲目列表
const ListOpera = () => {
    return new Promise((resolve, reject) => {
        connection.query("SELECT * FROM `opera`", (err, data) => {
            if (err) return reject(err);
            if (data.length === 0) return reject(resolve({ msg: "曲目列表为空" }));
            return resolve({ code: 0, data: data });
        })
    })
}

// 获取曲目信息
const getOperaInfo = (Opera_id, Opera_title) => {
    return new Promise((resolve, reject) => {
        let query = "SELECT * FROM `opera` WHERE ";
        let params = [];
        if (Opera_id) {
            query += `Opera_id = ?`;
            params.push(Opera_id);
        } else if (Opera_title) {
            query += `Opera_title = ?`;
            params.push(Opera_title);
        } else {
            return reject(new Error("参数错误"));
        }
        connection.query(query, params, (err, data) => {
            if (err) return reject(err);
            if (data.length === 0 && Opera_id) return reject(resolve({ code: 1, msg: "曲目id不存在" }));
            if (data.length === 0 && Opera_title) return reject(resolve({ code: 1, msg: "曲目标题不存在" }));
            return resolve({ code: 0, data: data });
        });
    });
}

// 获取收藏夹信息
const getFavorites = (User_id, Opera_id) => {
    let query = "SELECT * FROM `favorites` WHERE TRUE ";
    let params = [];
    if (User_id) {
        query += `AND User_id = ? `;
        params.push(User_id);
    }
    if (Opera_id) {
        query += `AND Opera_id = ? `;
        params.push(Opera_id);
    }
    if (params.length === 0) {
        return reject(new Error("参数错误"));
    }
    return new Promise((resolve, reject) => {
        connection.query(query, params, (err, data) => {
            if (err) return reject(err);
            if (data.length === 0) return reject(resolve({ code: 1, msg: "收藏夹为空" }));
            return resolve({ code: 0, data: data });
        });
    });
}

// 获取单用户播放记录
const getPlayHistory = (User_id, Opera_id) => {
    let query = "SELECT * FROM `playhistory` WHERE TRUE ";
    let params = [];
    if (User_id) {
        query += `AND User_id = ?`;
        params.push(User_id);
    }
    if (Opera_id) {
        query += `AND Opera_id = ?`;
        params.push(Opera_id);
    }
    if (params.length === 0) {
        return reject(new Error("参数错误"));
    }
    return new Promise((resolve, reject) => {
        connection.query(query, params, (err, data) => {
            if (err) return reject(err);
            if (data.length === 0) return reject(resolve({ code: 1, msg: "播放记录为空" }));
            return resolve({ code: 0, data: data });
        });
    });
}

// 统计每首曲目的播放次数
const getPlayCount = () => {
    return new Promise((resolve, reject) => {
        connection.query("SELECT Opera_id, Opera_title, (SELECT COUNT(*) FROM playhistory WHERE playhistory.Opera_id = opera.Opera_id) AS PlayNums FROM opera", (err, data) => {
            if (err) return reject(err);
            if (data.length === 0) return reject(resolve({ code: 1, msg: "播放记录为空" }));
            return resolve({ code: 0, data: data });
        });
    });
}

// 根据用户或曲目列出所有评论
const getComments = (Opera_id, User_id, Comment_id) => {
    return new Promise((resolve, reject) => {
        let query = "SELECT * FROM `comments` WHERE TRUE ";
        let params = [];
        if (Opera_id) {
            query += `AND Opera_id = ${Opera_id}`;
            params.push(Opera_id);
        }
        if (User_id) {
            query += `AND User_id = ${User_id}`;
            params.push(User_id);
        }
        if (Comment_id) {
            query += `AND Comment_id = ${Comment_id}`;
            params.push(Comment_id);
        }
        if (params.length === 0) {
            return reject(new Error("参数错误"));
        }
        connection.query(query, params, (err, data) => {
            if (err) return reject(err);
            if (data.length === 0 && Opera_id) return resolve({ code: 1, msg: "该曲目评论为空，快来发表第一条评论吧~" });
            if (data.length === 0 && User_id) return resolve({ code: 1, msg: "该用户评论为空" });
            if (data.length === 0) return resolve({ code: 1, msg: "该用户未在此曲目下发表评论" });
            return resolve({ code: 0, data: data });
        });
    });
}

// 发表评论
const postComment = (User_id, Opera_id, Comment_text) => {
    return new Promise((resolve, reject) => {
        connection.query("INSERT INTO `comments` (User_id, Opera_id, Comment_text) VALUES (?, ?, ?)", [User_id, Opera_id, Comment_text], (err, data) => {
            if (err) return reject({ code: 1, msg: "Failed" });
            return resolve({ code: 0, msg: "Success" });
        });
    });
}

// 删除评论
const deleteComment = (User_id, Comment_id) => {
    return new Promise((resolve, reject) => {
        connection.query("DELETE FROM `comments` WHERE User_id = ? AND Comment_id = ?", [User_id, Comment_id], (err, data) => {
            if (err) return reject({ code: 1, msg: "Failed" });
            return resolve({ code: 0, msg: "Success" });
        });
    });
}

// 新增戏曲
const addOpera = (Opera_title, Artist, Release_date, Duration, Music_path, Video_path, Description) => {
    if (!Opera_title || !Music_path || !Video_path) {
        throw { code: 1, msg: "invalid parameters" };
    }

    // 动态构建SQL语句和参数数组
    let query = "INSERT INTO `opera` (Opera_title, Music_path, Video_path";
    let params = [Opera_title, Music_path, Video_path];
    let placeholders = "?, ?, ?";
    if (Artist) {
        query += ", Artist";
        params.push(Artist);
        placeholders += ", ?";
    }
    if (Release_date) {
        query += ", Release_date";
        params.push(Release_date);
        placeholders += ", ?";
    }
    if (Duration) {
        query += ", Duration";
        params.push(Duration);
        placeholders += ", ?";
    }
    if (Description) {
        query += ", Description";
        params.push(Description);
        placeholders += ", ?";
    }

    query += `) VALUES (${placeholders})`;
    return new Promise((resolve, reject) => {
        connection.query(query, params, (err, data) => {
            if (err) return reject({ code: 1, msg: "Insert Failed" });
            return resolve({ code: 0, msg: "Success" });
        });
    });
}

// 删除戏曲
const deleteOpera = (Opera_id) => {
    if (!Opera_id) {
        throw { code: 1, msg: "invalid parameters" };
    }
    return new Promise((resolve, reject) => {
        connection.query("DELETE FROM `opera` WHERE Opera_id = ?", [Opera_id], (err, data) => {
            if (err) return reject({ code: 1, msg: "Failed" });
            return resolve({ code: 0, msg: "Success" });
        });
    });
}

// 修改戏曲信息
const updateOpera = async (Opera_id, UpdateFields, Value) => {
    if (!Opera_id || !UpdateFields || !Value) {
        throw { code: 1, msg: "invalid parameters" };
    }
    if (UpdateFields !== 'Opera_title' && UpdateFields !== 'Artist' && UpdateFields !== 'Release_date' && UpdateFields !== 'Duration' && UpdateFields !== 'Music_path' && UpdateFields !== 'Video_path' && UpdateFields !== 'Description') {
        throw { code: 1, msg: "invalid field" };
    }

    const data = await getOperaInfo(Opera_id, null);
    if (data.code !== 0) {
        throw { code: 1, msg: "Opera_id not exist" };
    }
    let query = "UPDATE `opera` SET ";
    let params = [Value, Opera_id];
    query += `${UpdateFields} = ? WHERE Opera_id = ?`;

    return new Promise((resolve, reject) => {
        connection.query(query, params, (err, data) => {
            if (err) return reject({ code: 1, msg: "Failed" });
            return resolve({ code: 0, msg: "Success" });
        });
    });
}

// 新增收藏
const addFavorite = (User_id, Opera_id) => {
    if (!User_id || !Opera_id) {
        throw { code: 1, msg: "invalid parameters" };
    }
    return new Promise((resolve, reject) => {
        connection.query("INSERT INTO `favorites` (User_id, Opera_id) VALUES (?, ?)", [User_id, Opera_id], (err, data) => {
            if (err) return reject({ code: 1, msg: "Failed" });
            return resolve({ code: 0, msg: "Success" });
        });
    });
}

// 取消收藏
const deleteFavorite = (User_id, Opera_id) => {
    if (!User_id || !Opera_id) {
        throw { code: 1, msg: "invalid parameters" };
    }
    return new Promise((resolve, reject) => {
        connection.query("DELETE FROM `favorites` WHERE User_id = ? AND Opera_id = ?", [User_id, Opera_id], (err, data) => {
            if (err) return reject({ code: 1, msg: "Failed" });
            return resolve({ code: 0, msg: "Success" });
        });
    });
}

module.exports = {
    ListUsers,
    ListOpera,
    getUserInfo,
    getOperaInfo,
    getFavorites,
    getPlayHistory,
    getPlayCount,
    getComments,
    postComment,
    deleteComment,
    addOpera,
    deleteOpera,
    updateOpera,
    addFavorite,
    deleteFavorite
}