# 轻量级 CI/CD 流水线服务说明

## 项目背景

本项目旨在用自研的轻量级服务替代 Jenkins，实现 CI/CD 流水线的自动化触发。我们的实际需求很简单：所有构建、部署等核心逻辑都已经封装在后端 shell 脚本中，Jenkins 仅仅作为 webhook 触发器，显得过于庞大和复杂。因此，我开发了这个 Go 编写的中间件服务，专注于 webhook 事件接收和脚本调度，极大简化了运维流程。

---

## 设计理念

- **极简触发**：只负责监听 GitLab webhook，根据事件类型（push/tag）自动调用对应 shell 脚本。
- **模块化配置**：前后端流水线分离，支持多环境、多类型部署。
- **易于维护**：所有业务逻辑都在 shell 脚本中，服务本身无业务耦合，便于扩展和排查。
- **轻量高效**：无需引入庞大的 Jenkins，仅用一个 Go 服务即可满足需求。

---

## 流水线配置

### 前端流水线

| 名称 | 环境 | 脚本路径 | 触发事件 |
|------|------|----------|----------|
| bbx_Front_FAT_UAT | 本地测试环境 | /home/jenkins/execute/pc/beta/build-deploy.sh | push |
| bbx_Front_Online  | 线上生产环境 | /home/jenkins/execute/pc/prod/build-prod.sh  | push |

### 后端流水线

| 名称 | 类型 | 环境 | 脚本路径 | 触发事件 |
|------|------|------|----------|----------|
| bbx_bbz_OpenBeta_pipline | 本地部署 | 本地测试环境 | /home/jenkins/execute/beta/beatenv_build/build-beta.sh | push |
| bbx_bbz_Prod_pipline | 一键部署 | 线上生产环境 | /home/jenkins/build/bbx-saas/build/prod/prod_one_click/build-all-prod.sh | tag |
| bbx_bbz_OnlineAlone_pipline | 模块部署 | 线上生产环境 | /home/jenkins/execute/prod/prodenv_build/build-prod.sh | push |

---

## 使用说明

- **GitLab Webhook 配置建议：**
    - URL: `http://your-server-ip:18088/webhook`
    - Secret Token: `your-gitlab-secret-token`
- 服务会自动识别 webhook 事件类型（push/tag），并调度对应脚本。
- 日志输出到控制台或日志文件，方便排查问题。
- 所有流水线配置和脚本路径均可按需扩展。

---

## 示例：Webhook 测试命令

```bash
curl -X POST http://localhost:18088/webhook \
  -H "X-Gitlab-Event: Push Hook" \
  -H "X-Gitlab-Token: your-gitlab-secret-token" \
  -H "Content-Type: application/json" \
  -d '{
    "object_kind": "push",
    "event_name": "push",
    "ref": "refs/heads/main",
    "total_commits_count": 1,
    "user_email": "dev@example.com",
    "user_name": "Developer",
    "project": {
      "name": "your-backend-project",
      "web_url": "http://gitlab.example.com/your/backend"
    },
    "commits": [
      {
        "id": "b489d49dfdab56d79001e83664ac5467ed83f5e7",
        "message": "[build] 修复用户登录接口的BUG",
        "author": {
          "name": "Developer",
          "email": "dev@example.com"
        },
        "timestamp": "2025-07-10T18:30:00+08:00",
        "url": "http://gitlab.example.com/your/backend/-/commit/b489d49dfdab56d79001e83664ac5467ed83f5e7"
      }
    ]
  }'
```

---

**总结：**  
本项目专为“只需触发脚本”的 CI/CD 场景设计，极大简化了传统 Jenkins 方案，适合所有 CI/CD 逻辑已沉淀在脚本、只需自动触发的团队使用。
