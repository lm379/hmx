Cookie.set('Admin', true)
var app = new Vue({
    el: '#app',
    data: {
        // Menu: [{ Type: '读取中' }],
        // Menu_New: '',
        ItemList: [],
        NewProduct: { Opera_title: '', Description: '', Release_date: '', Image: null, Video_path: '', Artist: '', },
        NewProduct_Select_Image: null,
        NewProduct_Select_Video: null,
    },
    mounted: function () {
        // 初始化
        let App = this;
        let token = Cookie.get('Token');
        if (!token) {
            window.location.href = 'login_admin.html'
            return;
        }
        loginConfirm()

        function loginConfirm() {
            HTTP.post('login_admin_confirm', null, 'user').then(json => {
                if (json.code === 0) {
                    // 保存Token，以方便日后自动登录
                    Cookie.set('Token', json.token);
                    if (json.role === 'Administrator') {
                        App.getOperaList()
                    } else {
                        alert('您没有权限访问该页面')
                        Cookie.del('Token')
                        window.location.href = 'index.html'
                    }
                } else {
                    alert('登录已过期，请重新登录')
                    Cookie.del('Token')
                    window.location.href = 'login_admin.html'
                }
            })
        }
    },
    methods: {
        // 获取戏曲列表
        getOperaList: function () {
            HTTP.get('list_opera', null, 'info').then(json => {
                if (json.code == 0) {
                    this.ItemList = json.data
                }
            })
        },

        changeImg: function () {
            let App = this;
            let files = document.getElementById('imgFile').files;
            if (files.length > 0) {
                let file = files[0];
                App.NewProduct.Image = file;

                //创建文件读取对象
                let reader = new FileReader();
                // 读取文件
                reader.readAsDataURL(file);
                // 显示选择的图片
                reader.onloadend = function () {
                    App.NewProduct_Select_Image = this.result;
                }
            }
        },

        changeVideo: function () {
            let App = this;
            let files = document.getElementById('videoFile').files;
            if (files.length > 0) {
                let file = files[0];
                App.NewProduct.Video_path = file;

                //创建文件读取对象
                let reader = new FileReader();
                // 读取文件
                reader.readAsDataURL(file);
                // 显示选择的图片
                reader.onloadend = function () {
                    App.NewProduct_Select_Video = this.result;
                }
            }
        },

        addProduct: function () {
            let App = this,
                info = App.NewProduct,
                check = true;
            // 循环新戏曲信息json，如果有空值存在，说明还有信息没有填写
            for (const key in info) {
                if (info[key] == '' || info[key] == null) { check = false }
            }
            // if (!check) {
            // alert('请输入完整戏曲信息')
            // } else {
            // 上传图片
            let options = {
                "opera_title": info.Opera_title,
                "description": info.Description,
                "release_date": info.Release_date,
                "artist": info.Artist,
            }


            uploadFile(info.Image, function (err, data) {
                if (err) {
                    alert(err);
                } else {
                    info.Image = data.url;
                    options['avatar'] = info.Image

                    // 上传视频
                    uploadFile(info.Video_path, function (err, data) {
                        if (err) {
                            alert(err);
                        } else {
                            info.Video_path = data.url;
                            options['video_path'] = info.Video_path
                            HTTP.post('add_opera', options, 'info').then(json => {
                                // 添加成功后，刷新戏曲列表
                                if (json.code == 0) {
                                    App.getOperaList()
                                    App.NewProduct = { Opera_title: '', Description: '', Release_date: '', Image: null, Video_path: '', Artist: '', }
                                    App.NewProduct_Select_Image = null
                                    App.NewProduct_Select_Video = null
                                }
                            })
                        }
                    }, function (progress) {
                        // 上传开始时显示进度条
                        const progressContainer = document.getElementById('progress-container');
                        progressContainer.style.display = 'block';

                        // 更新进度条
                        const progressBar = document.getElementById('progress-bar');
                        progressBar.style.width = progress.percent + '%';
                        progressBar.textContent = progress.percent + '%';
                    });
                }
            });


            // }

        },

        // 删除戏曲
        delProduct: function (event) {
            let Opera_id = event.target.dataset.id;
            Opera_id = parseInt(Opera_id)
            // 弹出提示框
            let check = confirm("确定删除该戏曲吗?");
            // 如果用户点了是，那么删除该戏曲
            let options = { "opera_id": Opera_id }
            if (check) {
                HTTP.post('delete_opera', options, 'info').then(json => {
                    // 删除成功后，刷新戏曲列表
                    if (json.code == 0) { this.getOperaList() }
                })
            }
        },

        jump: function (event) {
            let page = event.target.dataset.page + '.html';
            window.location.href = page;
        },

        // 退出登录
        logout: function () {
            // 弹出提示框
            let check = confirm("确定退出登录吗?");
            if (check) {
                Cookie.del('Token')
                window.location.href = 'login_admin.html'
            }
        },
    }
});


