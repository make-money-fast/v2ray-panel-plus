new Vue({
    el: '#app', data() {
        return {
            host: '',
            configList: [], // 配置列表.
            configDrawer: false,
            editConfig: {},
            vmessShare: false,
            vmessShareItem: {
                link: '',
                qrCode: '',
            },
            runtimeConfig: false,
            runtimeString: '',
            rules: {
                alias: [
                    {required: true, message: '请输入别名', trigger: 'blur'},
                ],
                host: [
                    {required: true, message: '请输入地址', trigger: 'blur'},
                ],
                port: [
                    {required: true, message: '请输入端口', trigger: 'blur'},
                ],
                id: [
                    {required: true, message: '请输入id', trigger: 'blur'},
                ],
                protocol: [
                    {required: true, message: '请选择协议', trigger: 'blur'},
                ]
            },
            protocolList: [
                'kcp',
                'ws',
                'tcp',
            ],
            headers: [
                'none',
                'srtp',
                'utp',
                'wechat-video',
                'dtls',
                'wireguard'
            ],
        }
    },
    mounted() {
        const params = new URLSearchParams(window.location.search);
        this.host = params.get('host'); // 'John'
        this.getList();
    },
    methods: {
        async handleMenuClick() {
        },
        async getList() {
            let rsp = await this.post('/api/server/list')
            this.configList = rsp;
        },
        async handleEditClick(item) {
            this.editConfig = Object.assign({}, item)
            console.log(this.editConfig)
            this.configDrawer = true;
            this.editConfig.add = false
            console.log(item)
        },
        async handleAddClick() {
            this.editConfig = {}
            this.editConfig.add = true
            this.editConfig.id = await this.getUUID();
            this.configDrawer = true;
        },
        async handleReloadClick() {
            let rsp = await this.post('/api/server/reload');
            if (rsp) {
                this.$message.success("操作成功");
            }
        },
        async getUUID() {
            let uuid = await this.post("/api/server/uuid")
            return uuid;
        },
        async handleCopy(item) {
            this.editConfig = Object.assign({},item);
            this.editConfig.uuid = await this.getUUID();
            this.editConfig.alias = this.editConfig.alias + "--复制"
            rsp = await this.post("/api/server/add", this.editConfig)
            if (rsp) {
                this.getList()
                this.$message.success("操作成功");
            }
        },
        async handleRuntimeConfig(item) {
            let rsp = await this.post("/api/server/runtime")
            this.runtimeConfig = true;
            this.runtimeString = rsp;
            this.$nextTick(() => {
                console.log(this.runtimeString)
                var editor = ace.edit("editor");
                editor.setTheme("ace/theme/chrome");
                editor.session.setMode("ace/mode/json");
            })
        },
        async onRuntimeConfigCancel() {
            this.runtimeConfig = false;
        },
        async handleShare(item) {
            this.vmessShare = true;
            this.vmessShareItem.link = item.vmess;
            this.vmessShareItem.qrCode = '/qrcode?code=' + item.vmess;
        },
        async onConfigEdit(formName) {
            this.$refs[formName].validate(async (valid) => {
                if (!valid) {
                    this.$message.error("表单验证失败, 请填写完整");
                    return
                }
                let rsp;
                if (this.editConfig.add) {
                    rsp = await this.post("/api/server/add", this.editConfig)
                } else {
                    rsp = await this.post("/api/server/edit", this.editConfig)
                }
                if (rsp) {
                    this.getList()
                    this.$message.success("操作成功");
                    this.editConfig = {};
                    this.configDrawer = false;
                }
            });
        },
        onConfigCancel() {
            this.configDrawer = false;
            this.editConfig = {};
        },
        async handleDelete(item) {
            let rsp = await this.post("/api/server/del", {uuid: item.uuid})
            if (rsp) {
                this.getList();
            }
        },
        async post(path, data) {
            if (!data) {
                data = {}
            }
            data.server_url = this.host;
            let rsp = await axios.post(path, data, {
                headers: {
                    'Content-Type': 'application/json',
                }
            })
            if (rsp.data.code === 0) {
                return rsp.data.data;
            }
            this.$message.error("操作失败:" + rsp.data.msg);
            return null;
        },
    }
})
