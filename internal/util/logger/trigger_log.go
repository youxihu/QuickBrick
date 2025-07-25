package logger

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"strings"
	"time"

	"QuickBrick/internal/domain"
)

func WriteTriggerLogToFile(pushEvent domain.PushEvent, env string, script string, scriptOutput string, ip string) error {
	envType, subEnv, err := parseEnv(env)
	if err != nil {
		Error("不支持的环境类型", ip, 500,
			zap.String("env", env),
			zap.Error(err),
		)
		return err
	}

	baseDir := "logs"
	dateDir := time.Now().Format("20060102")
	hourMin := time.Now().Format("1504")
	logDir := fmt.Sprintf("%s/%s/%s/%s", baseDir, envType, subEnv, dateDir)

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err = os.MkdirAll(logDir, 0755)
		if err != nil {
			Error("无法创建日志目录", ip, 500,
				zap.String("dir", logDir),
				zap.Error(err),
			)
			return fmt.Errorf("无法创建日志目录: %v", err)
		}
	}

	username := sanitizeUsername(pushEvent.UserName)
	filename := fmt.Sprintf("%s/%s_%s.log", logDir, hourMin, username)

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		Error("无法打开日志文件", ip, 500,
			zap.String("filename", filename),
			zap.Error(err),
		)
		return fmt.Errorf("无法打开日志文件: %v", err)
	}
	defer file.Close()

	content := generateLogContent(pushEvent, env, script, scriptOutput)
	_, err = file.WriteString(content)
	if err != nil {
		Error("写入日志失败", ip, 500,
			zap.String("filename", filename),
			zap.Error(err),
		)
		return fmt.Errorf("写入日志失败: %v", err)
	}

	Info("构建日志写入成功", ip, 200,
		zap.String("env", env),
		zap.String("filename", filename),
		zap.String("project", pushEvent.Project.Name),
		zap.String("ref", pushEvent.Ref),
	)

	return nil
}

// parseEnv 解析 env 字段，返回一级和二级目录
func parseEnv(env string) (envType, subEnv string, err error) {
	switch {
	case strings.HasPrefix(env, "fe-beta"):
		envType = "beta"
		subEnv = "fe"
	case strings.HasPrefix(env, "fe-prod"):
		envType = "prod"
		subEnv = "fe"
	case strings.HasPrefix(env, "be-beta"):
		envType = "beta"
		subEnv = "be"
	case strings.HasPrefix(env, "be-prod") && !strings.Contains(env, "be-all-prod"):
		envType = "prod"
		subEnv = "be"
	case env == "be-all-prod":
		envType = "prod"
		subEnv = "be-tag"
	default:
		return "", "", fmt.Errorf("不支持的环境类型: %s", env)
	}
	return envType, subEnv, nil
}

// 替换非法字符
func sanitizeUsername(username string) string {
	username = strings.ReplaceAll(username, "@", "_at_")
	username = strings.ReplaceAll(username, ".", "_dot_")
	return username
}

// 构建日志内容
func generateLogContent(pushEvent domain.PushEvent, env string, script string, scriptOutput string) string {
	var builder strings.Builder

	// === 触发详情 ===
	builder.WriteString("========================================\n")
	builder.WriteString("            🔔 触发详情（Trigger Info）\n")
	builder.WriteString("========================================\n")
	builder.WriteString(fmt.Sprintf("记录生成时间（Log Time）：%s\n", time.Now().Format("2006-01-02 15:04:05")))
	builder.WriteString(fmt.Sprintf("环境标识（Env）：%s\n", env))
	builder.WriteString(fmt.Sprintf("项目名称（Project Name）：%s\n", pushEvent.Project.Name))
	builder.WriteString(fmt.Sprintf("项目地址（Project URL）：%s\n", pushEvent.Project.URL))
	builder.WriteString(fmt.Sprintf("分支或标签（Git Ref/Tag）：%s\n", pushEvent.Ref))
	builder.WriteString(fmt.Sprintf("事件类型（Event Type）：%s\n", pushEvent.ObjectKind))

	// === 提交记录 ===
	builder.WriteString("\n========================================\n")
	builder.WriteString("            📤 提交记录（Git Commit Info）\n")
	builder.WriteString("========================================\n")
	for i, commit := range pushEvent.Commits {
		builder.WriteString(fmt.Sprintf("提交 #%d：\n", i+1))
		builder.WriteString(fmt.Sprintf("  Commit ID：%s\n", commit.ID))
		builder.WriteString(fmt.Sprintf("  提交人（Committer）：%s <%s>\n", commit.Author.Name, commit.Author.Email))
		builder.WriteString(fmt.Sprintf("  提交说明（Message）：%s\n", commit.Message))
		builder.WriteString(fmt.Sprintf("  提交时间（Timestamp）：%s\n", commit.Timestamp))
		builder.WriteString(fmt.Sprintf("  提交详情链接（URL）：%s\n", commit.URL))
	}

	// === 执行脚本 ===
	builder.WriteString("\n========================================\n")
	builder.WriteString("            🛠️ 执行脚本（Script Execution）\n")
	builder.WriteString("========================================\n")
	builder.WriteString(fmt.Sprintf("脚本路径（Script Path）：%s\n", script))

	builder.WriteString("\n--- 脚本输出（Script Output） ---\n")
	if scriptOutput == "" {
		builder.WriteString("(无输出/No Output)\n")
	} else {
		lines := strings.Split(scriptOutput, "\n")
		for _, line := range lines {
			builder.WriteString("  ")
			builder.WriteString(line)
			builder.WriteString("\n")
		}
	}
	builder.WriteString("========================================\n\n")

	return builder.String()
}
