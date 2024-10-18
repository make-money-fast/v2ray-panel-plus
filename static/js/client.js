const green = '#0bbd87';
const gray = '#cccccc'
const red = '#eb5851'

new Vue({
    el: '#app', data() {
        return {
            configList: [], // 配置列表.
            clientRuntime: {},
            configDrawer: false,
            editConfig: {},
            defaultEditConfig: {
                host: '',
                port: '',
                id: '',
                remark: '',
                network: '',
                alias: '',
                uuid: '',
                settings: {
                    tcpSettings: {},
                    kcpSettings: {
                        uplinkCapacity: 5,
                        downlinkCapacity: 100,
                        congestion: true,
                        header: {
                            type: 'none',
                        }
                    },
                    wsSettings: {
                        path: '',
                    },
                },
            },
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
                network: [
                    {required: true, message: '请选择协议', trigger: 'blur'},
                ]
            },
            fullscreenLoading: false,
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

            // 本地设置
            localConfigDrawer: false,
            defaultLocalConfig: {
                socksAddress: '0.0.0.0',
                socksPort: 10080,
                httpAddress: '0.0.0.0',
                httpPort: 10081,
            },
            localConfig: {},
            localRules: {
                socksPort: [
                    {
                        min: 1000, max: 65536, type: 'number'
                    }
                ]
            },
            stateDrawer: false,
            state: [{
                content: '运行状态',
                color: gray
            }, {
                content: 'socks端口',
                color: gray
            }, {
                content: 'http端口',
                color: gray
            }, {
                content: '能否连上代理服务器',
                color: gray
            }, {
                content: '能否通过代理服务器正常上网',
                color: gray
            }],
            configJson: '',
            configJsonDrawer: false,
            vmessImportDrawer: false,
            vmessImport: {
                vmess: '',
            },
            vmessImportRule: {
                vmess: [
                    {required: true, message: '请输入连接', trigger: 'blur'},
                ]
            },
            vmessShare: false,
            vmessShareItem: {
                link: '',
                qrCode: '',
            },
            configJsonImportDrawer: false,
            configImportConfig: {
                config: '',
            },
            configImportConfigRule: {
                config: [
                    {required: true, message: '请输入json', trigger: 'blur'},
                ]
            },
            proxyStatus: {
                mode: 0,
                status: 0,
            }
        }
    },
    mounted() {
        this.getConfigList();
        this.initLocalConfig();
        this.getSystemProxyStatus();
    },
    methods: {
        async getSystemProxyStatus() {
            let rsp = await this.post('/api/systemProxyStatus')
            this.proxyStatus = {
                mode: rsp.mode,
                status: rsp.status,
                pacAddress: '',
            }
        },
        async setProxy(mode) {
            let rsp = await  this.post("/api/setProxy",{mode: mode})
            if(rsp) {
                await this.getSystemProxyStatus();
            }
        },
        async handleProxyClick(command) {
            if (command === "pac") {
                await this.setProxy(2)
            }
            if (command == 'global') {
                await this.setProxy(1)
            }
            if (command == 'close') {
                await this.clearProxy()
            }
        },
        async stopProgress() {
            let rsp = this.post('/api/shutdown', {})
            if (rsp) {
                this.$message.error("程序已退出");
                setTimeout(() => {
                    window.close();
                },2000)
            }
        },
        onConfigImportFormCancel() {
            this.configImportConfig = {
                config: '',
            }
            this.configJsonImportDrawer = false;
        },
        async onConfigImportFormSubmit(formName) {
            this.$refs[formName].validate(async (valid) => {
                if (!valid) {
                    this.$message.error("表单验证失败, 请填写完整");
                    return
                }
                let rsp = await this.post("/api/config-import", this.configImportConfig)
                if (rsp) {
                    this.$message.success("操作成功");
                    this.configJsonImportDrawer = false;
                    this.getConfigList();
                }
            });
        },
        async importConfigJson() {
            this.configImportConfig = {
                config: '',
            }
            this.configJsonImportDrawer = true
        },
        async getConfigJson(item) {
            let rsp = await this.post('/api/configJSON', {uuid: item.uuid})
            if (rsp) {
                this.configJsonDrawer = true;
                rsp = rsp.replaceAll("\n", "<br />")
                this.configJson = rsp;
            }
        },
        async handleShare(item) {
            this.vmessShare = true;
            this.vmessShareItem.link = item.vmess;
            this.vmessShareItem.qrCode = '/qrcode?code=' + item.vmess;
        },
        async handleStopClick() {
            await this.post("/api/stop", {})
            await this.getConfigList()
        },
        async onVmessImportEdit(formName) {
            this.$refs[formName].validate(async (valid) => {
                if (!valid) {
                    this.$message.error("表单验证失败, 请填写完整");
                    return
                }
                let rsp = await this.post("/api/importVmess", this.vmessImport)
                if (rsp) {
                    this.$message.success("操作成功");
                    this.localConfigDrawer = false;
                    this.getConfigList();
                }
            });
        },
        onVmessImportCancel() {
            this.vmessImportDrawer = false;
            this.vmessImport.vmess = '';
        },
        handleVmessImportClick() {
            this.vmessImportDrawer = true;
        },
        async onLocalConfigEdit(formName) {
            this.$refs[formName].validate(async (valid) => {
                if (!valid) {
                    this.$message.error("表单验证失败, 请填写完整");
                    return
                }
                let rsp = await this.post("/api/updateLocalConfig", this.localConfig)
                if (rsp) {
                    this.$message.success("操作成功");
                    this.localConfigDrawer = false;
                }
            });
        },
        onLocalConfigCancel() {
            this.localConfigDrawer = false;
        },
        async handleLocalConfigClick() {
            this.localConfigDrawer = true;
        },
        async initLocalConfig() {
            let rsp = await this.post('/api/getLocalConfig')
            this.localConfig = rsp;
        },
        async handleReload(item) {
            let rsp = await this.post('/api/reload', {uuid: item.uuid})
            if (rsp) {
                this.getConfigList()
            }
        },
        async handleDelete(item) {
            let rsp = await this.post("/api/del", {uuid: item.uuid});
            if (rsp) {
                this.$message.success('操作成功')
                await this.getConfigList()
            }
        },
        onConfigCancel() {
            this.editConfig = {}
            this.configDrawer = false;
        },
        async onConfigEdit(formName) {
            this.$refs[formName].validate(async (valid) => {
                if (!valid) {
                    this.$message.error("表单验证失败, 请填写完整");
                    return
                }
                let rsp = await this.post("/api/edit", this.editConfig)
                if (rsp) {
                    this.getConfigList()
                    this.$message.success("操作成功");
                    this.editConfig = {};
                    this.configDrawer = false;
                }
            });
        },
        handleEditClick(item) {
            this.configDrawer = true
            this.editConfig = Object.assign({}, this.defaultEditConfig);
            this.editConfig.host = item.host;
            this.editConfig.port = item.port;
            this.editConfig.id = item.id;
            this.editConfig.remark = item.remark;
            this.editConfig.network = item.protocol;
            this.editConfig.alias = item.alias;
            this.editConfig.uuid = item.uuid;
            if (item.network === 'kcp') {
                this.editConfig.network.settings.kcpSettings = item.config.outbounds[0].streamSettings.kcpSettings;
            }
        },
        async handleMenuClick(command) {
            if (command === "init") {
                await this.init()
                await this.getConfigList();
                return
            }
            if (command == 'vmess') {
                await this.handleVmessImportClick();
            }
            if (command == 'import') {
                this.importConfigJson();
            }
        },
        async init() {
            let rsp = await this.post('/api/init')
            console.log(rsp)
        },
        async getConfigList() {
            let rsp = await this.post("/api/list", {})
            this.configList = rsp ?? {};
        },
        async autoTestTimeLine() {
            this.fullscreenLoading = true;
            let state = await this.post("/api/autoTest", {});
            this.fullscreenLoading = false;
            if (state.isRunning) {
                this.state[0].color = green
            } else {
                this.state[0].color = red
            }
            if (state.socks) {
                this.state[1].color = green
            } else {
                this.state[1].color = red
            }
            if (state.http) {
                this.state[2].color = green
            } else {
                this.state[2].color = red
            }
            if (state.connectToServer) {
                this.state[3].color = green
            } else {
                this.state[3].color = red
            }
            if (state.porxyOK) {
                this.state[4].color = green
            } else {
                this.state[4].color = red
            }

            this.stateDrawer = true
        },
        async post(path, data) {
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