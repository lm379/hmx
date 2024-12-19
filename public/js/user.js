Cookie.set('Admin', false);
let token = Cookie.get('Token');
var app = new Vue({
    el: '#app',
    data: {
        UserInfo: {},
        Menu: [
            { Title: '个人信息', Icon: 'ri-user-line' },
            // { Title: '评论管理', Icon: 'ri-building-line' },    
            { Title: '我的收藏', Icon: 'ri-star-line' },
        ],
        MenuActive: '个人信息',
        FollowList: [],
        NewUserInfo: { EMail: '', Phone: '', NewPassword: '', Sex: '', OldPassword: '' },
        NewUserAvatar: { File: null, Img: null },
    },
    mounted: function () {
        // 初始化
        this.getUserInfo();
    },
    methods: {
        // 获取用户信息
        getUserInfo: function () {
            let token = Cookie.get('Token');

            if (token) {
                HTTP.get('get_user_info', null, 'info').then(json => {
                    if (json.code === 0) {
                        this.UserInfo = json.data[0];
                        if (!json.data[0].Icon) {
                            this.UserInfo.Icon = './img/logo.png'
                        }
                        this.NewUserInfo = {
                            Username: this.UserInfo.Username,
                            Email: this.UserInfo.Email,
                            Phone: this.UserInfo.Phone,
                            Sex: this.UserInfo.Sex,
                        }
                        this.NewUserAvatar = {
                            File: null,
                            Img: null
                        }
                    } else {
                        Cookie.del('Token');
                        Cookie.del('Admin');
                        window.location.href = 'login.html';
                    }
                });
            } else {
                Cookie.del('Token');
                Cookie.del('Admin');
                window.location.href = 'login.html'
            }
        },

        // 修改用户信息
        setUserInfo: function () {
            if (this.NewUserInfo.Email == '' || this.NewUserInfo.Phone == '') {
                alert('请填写邮箱地址和联系电话')
            }

            else if (!this.checkPhone(this.NewUserInfo.Phone)) {
                alert('请输入正确的手机号码')
            }
            else if (!this.checkEMail(this.NewUserInfo.Email)) {
                alert('请输入正确的邮箱地址')
            } else {
                let options = {};
                if (this.UserInfo.Email != this.NewUserInfo.Email) {
                    options.email = this.NewUserInfo.Email;
                }
                if (this.UserInfo.Phone != this.NewUserInfo.Phone) {
                    options.phone = this.NewUserInfo.Phone;
                }
                if (this.UserInfo.Username != this.NewUserInfo.Username) {
                    options.username = this.NewUserInfo.Username;
                }
                if (this.UserInfo.Sex != this.NewUserInfo.Sex) {
                    options.sex = this.NewUserInfo.Sex;
                }

                if (Object.keys(options).length === 0) {
                    alert('请填写要修改的信息');
                    return;
                }
                HTTP.post('update_user', options, 'user').then(json => {
                    if (json.code === 0) {
                        alert('修改信息成功');
                        if (json.token) {
                            Cookie.del('Token');
                            Cookie.set('Token', json.token);
                        }
                        window.location.href = 'user.html?opt=个人信息';
                    } else {
                        alert(json.msg);
                    }
                })
            }
        },

        setPassword: function () {
            // 如果用户设置了新密码，那么添加到options中
            if (this.NewUserInfo.Password != '') {
                if (!this.checkPassword(this.NewUserInfo.OldPassword) && !this.checkPassword(this.NewUserInfo.NewPassword)) {
                    alert('密码格式有误')
                }
                let options = {};
                options.newpassword = this.NewUserInfo.NewPassword;
                options.oldpassword = this.NewUserInfo.OldPassword;

                HTTP.post('modify_password', options, 'user').then(json => {
                    if (json.code === 0) {
                        alert('修改密码成功');
                        Cookie.del('Token');
                        Cookie.del('Admin');
                        window.location.href = 'login.html';
                    } else {
                        alert(json.msg);
                    }
                });
            }
        },

        // 选择头像
        changeImg: function () {
            let App = this;
            let files = document.getElementById('imgFile').files;
            if (files.length > 0) {
                let file = files[0];
                App.NewUserAvatar.File = file;

                //创建文件读取对象
                let reader = new FileReader();
                // 读取文件
                reader.readAsDataURL(file);
                // 显示选择的图片
                reader.onloadend = function () {
                    App.NewUserAvatar.Img = this.result;
                    // 上传头像
                    uploadFile(file, function (err, data) {
                        if (err) {
                            alert(err);
                        } else {
                            let options = { 'icon': data.url };
                            HTTP.post('update_user', options, 'user').then(json => {
                                if (json.code === 0) {
                                    alert('修改头像成功');
                                    App.getUserInfo();
                                } else {
                                    alert(json.msg);
                                }
                            });
                        }
                    });
                }
            }
        },

        // 获取收藏列表
        getFollowList: function () {
            this.FollowList = []
            HTTP.get('get_favorites', null, 'info').then(json => {
                if (json.code === 0) {
                    this.FollowList = json.data;
                }
            })
        },

        // 取消收藏
        delFollow: function (event) {
            let opera_id = event.target.dataset.id;
            opera_id = parseInt(opera_id);
            let params = { "opera_id": opera_id }
            // 弹出提示框
            let check = confirm("确定取消该收藏吗?");
            // 如果用户点了是，那么取消该收藏
            if (check) {
                HTTP.post('delete_favorite', params, 'info').then(json => {
                    // 取消成功后，刷新收藏列表
                    if (json.code === 0) { this.getFollowList() }
                })
            }
        },

        // 数量加
        num_add: function (event) {
            let index = event.target.dataset.index;
            this.CartList[index].Num++
        },
        // 数量减
        num_remove: function (event) {
            let index = event.target.dataset.index;
            if (this.CartList[index].Num > 1) { this.CartList[index].Num-- }
        },

        // 点击菜单时，切换页面内容
        switchMenu: function (event, page) {
            let menu;
            if (page) { menu = page }
            else { menu = event.target.dataset.menu }
            this.MenuActive = menu;
            switch (menu) {
                case '我的收藏':
                    this.getFollowList()
                    break;
            }
        },

        // 使用正则表达式，检查手机号码是否正确
        checkPhone: function (phone) {
            let pattern = /^1[3456789]\d{9}$/;
            if (!pattern.test(phone)) { return false; }
            else { return true; }
        },

        // 使用正则表达式，检查邮箱是否正确
        checkEMail: function (email) {
            let pattern = /^([A-Za-z0-9_\-\.\u4e00-\u9fa5])+\@([A-Za-z0-9_\-\.])+\.([A-Za-z]{2,8})$/;
            if (pattern.test(email)) { return true; }
            else { return false; }
        },

        // 使用正则表达式，检查密码是否符合正确
        checkPassword: function (password) {
            var pattern = /^(?=.*[a-zA-Z])(?=.*\d|.*[\W_]).{8,16}$/;
            if (pattern.test(password)) { return true; }
            else { return false; }
        },

        // 跳转页面
        jump: function (event) {
            let page = event.target.dataset.page + '.html';
            window.location.href = page;
        },

        // 从url获取参数
        getQuery: function (value) {
            let query = window.location.search.substring(1);
            let vars = query.split("&");
            for (const i of vars) {
                let pair = i.split("=");
                if (pair[0] == value) { return pair[1]; }
            }
            return false;
        },

        // 退出登录
        logout: function () {
            // 弹出提示框
            let check = confirm("确定退出登录吗?");
            if (check) {
                Cookie.del('Token')
                window.location.href = 'login.html'
            }
        },

        // 格式化时间戳
        formatDate: function (dateString) {
            const options = { year: 'numeric', month: '2-digit', day: '2-digit' };
            return new Date(dateString).toLocaleDateString('zh-CN', options);
        },
    }
});


