package logger

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"strings"
	"time"

	"QuickBrick/internal/domain"
)

// WriteTriggerLogToFile 将完整的构建信息写入日志文件，并使用 zap 记录操作日志
func WriteTriggerLogToFile(pushEvent domain.PushEvent, env string, script string, scriptOutput string) error {
	// 解析环境层级
	envType, subEnv, err := parseEnv(env)
	if err != nil {
		Logger.Error("不支持的环境类型", zap.String("env", env), zap.Error(err))
		return err
	}

	// 构造目录结构
	baseDir := "logs"
	dateDir := time.Now().Format("20060102")
	hourMin := time.Now().Format("1504")

	// 构建完整路径
	logDir := fmt.Sprintf("%s/%s/%s/%s", baseDir, envType, subEnv, dateDir)

	// 创建目录结构
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err = os.MkdirAll(logDir, 0755)
		if err != nil {
			Logger.Error("无法创建日志目录",
				zap.String("dir", logDir),
				zap.Error(err),
			)
			return fmt.Errorf("无法创建日志目录: %v", err)
		}
	}

	// 构造文件名
	username := sanitizeUsername(pushEvent.UserName)
	filename := fmt.Sprintf("%s/%s_%s.log", logDir, hourMin, username)

	// 写入日志内容
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		Logger.Error("无法打开日志文件",
			zap.String("filename", filename),
			zap.Error(err),
		)
		return fmt.Errorf("无法打开日志文件: %v", err)
	}
	defer file.Close()

	content := generateLogContent(pushEvent, env, script, scriptOutput)
	_, err = file.WriteString(content)
	if err != nil {
		Logger.Error("写入日志失败",
			zap.String("filename", filename),
			zap.Error(err),
		)
		return fmt.Errorf("写入日志失败: %v", err)
	}

	Logger.Info("成功写入构建日志",
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
	builder.WriteString("            🔔 触发详情\n")
	builder.WriteString("========================================\n")
	builder.WriteString(fmt.Sprintf("时间:     %s\n", time.Now().Format("2006-01-02 15:04:05")))
	builder.WriteString(fmt.Sprintf("环境:     %s\n", env))
	builder.WriteString(fmt.Sprintf("项目:     %s (%s)\n", pushEvent.Project.Name, pushEvent.Project.URL))
	builder.WriteString(fmt.Sprintf("分支/Tag: %s\n", pushEvent.Ref))
	builder.WriteString(fmt.Sprintf("事件类型: %s\n", pushEvent.ObjectKind))

	// === 提交记录 ===
	builder.WriteString("\n========================================\n")
	builder.WriteString("            📤 提交记录\n")
	builder.WriteString("========================================\n")
	for i, commit := range pushEvent.Commits {
		builder.WriteString(fmt.Sprintf("提交 #%d:\n", i+1))
		builder.WriteString(fmt.Sprintf("  ID:         %s\n", commit.ID))
		builder.WriteString(fmt.Sprintf("  提交者:     %s <%s>\n", commit.Author.Name, commit.Author.Email))
		builder.WriteString(fmt.Sprintf("  提交信息:   %s\n", commit.Message))
		builder.WriteString(fmt.Sprintf("  查看链接:   %s\n", commit.URL))
	}

	// === 执行脚本 ===
	builder.WriteString("\n========================================\n")
	builder.WriteString("            🛠️ 执行脚本\n")
	builder.WriteString("========================================\n")
	builder.WriteString(fmt.Sprintf("脚本路径: %s\n", script))

	builder.WriteString("\n--- 脚本输出 ---\n")
	if scriptOutput == "" {
		builder.WriteString("(无输出)\n")
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
