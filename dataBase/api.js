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
        let query = "SELECT * FROM `users` WHERE TRUE";
        let params = [];
        if (User_id) {
            query += ` AND User_id = ?`;
            params.push(User_id);
        }

        if (Phone) {
            query += ` AND Phone = ?`;
            params.push(Phone);
        }

        if (Email) {
            query += ` AND Email = ?`;
            params.push(Email);
        }

        if (Username) {
            query += ' AND Username = ?';
            params.push(Username);
        }

        if (params.length === 0) {
            return reject(new Error("参数错误"));
        }

        connection.query(query, params, (err, data) => {
            if (err) return reject(err);
            if (data.length === 0) return reject(resolve({ code: 1, msg: "用户不存在" }));
            if (data[0].Status === 'Deleted') return reject(resolve({ code: 1, msg: "用户已注销" }));
            return resolve({ code: 0, data: data });
        });
    });
}

// 列出曲目列表
const ListOpera = () => {
    return new Promise((resolve, reject) => {
        connection.query("SELECT *, (SELECT COUNT(*) FROM `favorites`  WHERE favorites.Opera_id = opera.Opera_id) AS Favorite_count FROM `opera`", (err, data) => {
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
    let query = "SELECT f.Favorite_id, f.User_id, f.Opera_id, f.Created_time, o.Opera_title, o.Avatar FROM `favorites` f JOIN `opera` o ON o.Opera_id = f.Opera_id WHERE TRUE ";
    let params = [];
    if (User_id) {
        query += `AND f.User_id = ? `;
        params.push(User_id);
    }
    if (Opera_id) {
        query += `AND f.Opera_id = ? `;
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
        let query = "SELECT comments.*, users.Username,Users.Icon FROM `comments` JOIN `users` ON comments.user_id = users.user_id WHERE TRUE ";
        let params = [];
        if (Opera_id) {
            query += `AND Opera_id = ?`;
            params.push(Opera_id);
        }
        if (User_id) {
            query += `AND User_id = ?`;
            params.push(User_id);
        }
        if (Comment_id) {
            query += `AND Comment_id = ?`;
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

// 统计评论数量
const getCommentsCount = (params, value) => {
    return new Promise((resolve, reject) => {
        let query = "SELECT ";
        // 参数校验
        if (isNaN(value)) {
            return reject(new Error("参数错误"));
        }
        if (params === 'Opera_id' || params === 'User_id' || params === 'Comment_id') {
            query += `${params}, Count(*) AS Comments_count FROM \`comments\` WHERE ${params} = ? GROUP BY ${params} `;
        } else {
            return reject(new Error("参数错误"));
        }

        connection.query(query, [value], (err, data) => {
            if (err) return reject(err);
            if (data.length === 0) {
                const valueNum = JSON.parse(value);
                return resolve({ code: 0, data: { [params]: valueNum, Comments_count: 0 } });
            }
            return resolve({ code: 0, data: data[0] });
        });
    });
}

// 发表评论
const postComment = (User_id, Opera_id, Comment_text) => {
    return new Promise((resolve, reject) => {
        connection.query("INSERT INTO `comments` (User_id, Opera_id, Comment_text) VALUES (?, ?, ?)", [User_id, Opera_id, Comment_text], (err, data) => {
            if (err) return reject({ code: 1, msg: err });
            return resolve({ code: 0, msg: "成功" });
        });
    });
}

// 删除评论
const deleteComment = (User_id, Comment_id) => {
    return new Promise((resolve, reject) => {
        connection.query("DELETE FROM `comments` WHERE User_id = ? AND Comment_id = ?", [User_id, Comment_id], (err, data) => {
            if (err) return reject({ code: 1, msg: "失败" });
            return resolve({ code: 0, msg: "成功" });
        });
    });
}

// 新增戏曲
const addOpera = (Opera_title, Artist, Release_date, Duration, Music_path, Video_path, Description, Avatar) => {
    if (!Opera_title || !Video_path) {
        throw { code: 1, msg: "不合法的参数" };
    }
    const cosDomain = process.env.Bucket + '.cos.' + process.env.Region + '.myqcloud.com';
    Video_path = Video_path.replace(cosDomain, process.env.CDN);
    // 动态构建SQL语句和参数数组
    let query = "INSERT INTO `opera` (Opera_title, Video_path";
    let params = [Opera_title, Video_path];
    let placeholders = "?, ?";
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
    if (Music_path) {
        Music_path = Music_path.replace(cosDomain, process.env.CDN);
        query += ", Music_path";
        params.push(Music_path);
        placeholders += ", ?";
    }
    if (Avatar) {
        Avatar = Avatar.replace(cosDomain, process.env.CDN);
        query += ", Avatar";
        params.push(Avatar);
        placeholders += ", ?";
    }

    query += `) VALUES (${placeholders})`;
    return new Promise((resolve, reject) => {
        connection.query(query, params, (err, data) => {
            if (err) return reject({ code: 1, msg: "添加失败" });
            return resolve({ code: 0, msg: "添加成功" });
        });
    });
}

// 删除戏曲
const deleteOpera = (Opera_id) => {
    if (!Opera_id) {
        throw { code: 1, msg: "不合法的参数" };
    }
    return new Promise((resolve, reject) => {
        connection.query("DELETE FROM `opera` WHERE Opera_id = ?", [Opera_id], (err, data) => {
            if (err) return reject({ code: 1, msg: "失败" });
            return resolve({ code: 0, msg: "成功" });
        });
    });
}

// 修改戏曲信息
const updateOpera = async (Opera_id, UpdateFields, Value) => {
    if (!Opera_id || !UpdateFields || !Value) {
        throw { code: 1, msg: "不合法的参数" };
    }
    if (UpdateFields !== 'Opera_title' && UpdateFields !== 'Artist' && UpdateFields !== 'Release_date' && UpdateFields !== 'Duration' && UpdateFields !== 'Music_path' && UpdateFields !== 'Video_path' && UpdateFields !== 'Description') {
        throw { code: 1, msg: "不合法的请求" };
    }

    const data = await getOperaInfo(Opera_id, null);
    if (data.code !== 0) {
        throw { code: 1, msg: "曲目不存在" };
    }
    let query = "UPDATE `opera` SET ";
    let params = [Value, Opera_id];
    query += `${UpdateFields} = ? WHERE Opera_id = ?`;

    return new Promise((resolve, reject) => {
        connection.query(query, params, (err, data) => {
            if (err) return reject({ code: 1, msg: "更新失败" });
            return resolve({ code: 0, msg: "更新成功" });
        });
    });
}

// 新增收藏
const addFavorite = (User_id, Opera_id) => {
    if (!User_id || !Opera_id) {
        throw { code: 1, msg: "不合法的参数" };
    }
    return new Promise((resolve, reject) => {
        connection.query("INSERT INTO `favorites` (User_id, Opera_id) VALUES (?, ?)", [User_id, Opera_id], (err, data) => {
            if (err) return reject({ code: 1, msg: "失败" });
            return resolve({ code: 0, msg: "成功" });
        });
    });
}

// 取消收藏
const deleteFavorite = (User_id, Opera_id) => {
    if (!User_id || !Opera_id) {
        throw { code: 1, msg: "不合法的参数" };
    }
    return new Promise((resolve, reject) => {
        connection.query("DELETE FROM `favorites` WHERE User_id = ? AND Opera_id = ?", [User_id, Opera_id], (err, data) => {
            if (err) return reject({ code: 1, msg: "失败" });
            return resolve({ code: 0, msg: "成功" });
        });
    });
}

// 搜索
const search = (Keywords, Table, Field, Column) => {
    let query = `SELECT `;
    // 参数校验
    if (!Keywords || !Table || !Field || !Array.isArray(Column) || Column.length === 0) {
        throw { code: 1, msg: "不合法的参数" };
    }
    // 动态构建SQL语句
    Column.forEach((item, index) => {
        query += `${item}`;
        if (index !== Column.length - 1) {
            query += `, `;
        }
    });
    query += ` FROM ${Table} WHERE LOCATE(?, ${Field}) > 0`;
    return new Promise((resolve, reject) => {
        connection.query(query, [Keywords], (err, data) => {
            if (err) return reject(err);
            if (data.length === 0) return reject(resolve({ code: 1, msg: "搜索结果为空" }));
            return resolve({ code: 0, data: data });
        });
    });
}

const getAiResult = async (Opera_id) => {
    if (!Opera_id) {
        throw { code: 1, msg: "不合法的参数" };
    }
    const data = await getOperaInfo(Opera_id, null);
    if (data.code !== 0) {
        throw { code: 1, msg: "曲目不存在" };
    }
    return new Promise((resolve, reject) => {
        connection.query("SELECT * FROM `aidata` WHERE Opera_id = ?", [Opera_id], (err, data) => {
            if (err) return reject(err);
            if (data.length === 0) return reject(resolve({ code: 1, msg: "数据正在生成中" }));
            return resolve({ code: 0, data: data });
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
    deleteFavorite,
    search,
    getCommentsCount,
    getAiResult,
}