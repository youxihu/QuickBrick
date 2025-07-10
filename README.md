# CI/CD 流水线配置说明

本项目采用模块化设计，分别处理前端和后端的 CI/CD 任务。每条流水线对应不同的环境和脚本，并监听特定 GitLab 事件（如 `push` 或 `tag`）。

---

## 前端流水线

| 名称 | 环境 | 脚本路径 | 监听事件 |
|------|------|----------|-----------|
| `bbx_Front_FAT_UAT` | 本地测试环境 | `/home/jenkins/execute/pc/beta/build-deploy.sh` | `push` |
| `bbx_Front_Online` | 线上生产环境 | `/home/jenkins/execute/pc/prod/build-prod.sh` | `push` |

---

## 后端流水线

| 名称 | 类型   | 环境 | 脚本路径                                                                    | 监听事件 |
|------|------|------|-------------------------------------------------------------------------|-----------|
| `bbx_bbz_OpenBeta_pipline` | 本地部署 | 本地测试环境 | `/home/jenkins/execute/beta/beatenv_build/build-beta.sh`                | `push` |
| `bbx_bbz_Prod_pipline` | 一键部署 | 线上生产环境 | `/home/jenkins/build/bbx-saas/build/prod/prod_one_click/build-all-prod.sh` | `tag` |
| `bbx_bbz_OnlineAlone_pipline` | 模块部署 | 线上生产环境 | `/home/jenkins/execute/prod/prodenv_build/build-prod.sh`                | `push` |

---

## 说明

- **GitLab Webhook 配置建议：**
    - URL: `http://your-server-ip:18088/webhook`
    - Secret Token: `your-gitlab-secret-token`
- 所有流水线由 Go 编写的轻量级中间件服务统一接收并调度执行。
- 支持自动判断 `push` 或 `tag` 事件类型，并触发对应脚本。
- 日志会输出到控制台或日志文件中，便于排查问题。

---

## 🛠️ 示例 GitLab Webhook 请求测试命令

### Push 请求：

```bash
curl -X POST http://localhost:18088/webhook \
  -H "X-Gitlab-Event: Push Hook" \
  -H "X-Gitlab-Token: your-gitlab-secret-token" \
  -H "Content-Type: application/json" \
  -d '{"object_kind":"push","ref":"refs/heads/master"}'