package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv 加载 .env 文件
// .env 文件是可选的,如果不存在不会报错
// 这个函数应该在加载 config.yaml 之前调用
//
// 工作流程:
//  1. 尝试加载项目根目录的 .env 文件
//  2. 如果文件不存在,静默跳过
//  3. 如果文件存在但格式错误,记录错误但不中断
//
// 使用场景:
//   - 本地开发: 创建 .env 文件存放敏感配置
//   - 生产环境: 不使用 .env 文件,直接使用系统环境变量
//
// 注意事项:
//   - .env 文件不应该提交到 Git
//   - .env 文件中的变量会被加载到进程环境变量中
//   - 如果系统环境变量已存在同名变量,.env 文件的值会被忽略
func LoadEnv() {
	// 尝试加载 .env 文件
	// godotenv.Load() 会:
	// 1. 读取 .env 文件
	// 2. 解析 KEY=VALUE 格式
	// 3. 将变量设置到进程环境变量中
	// 4. 不会覆盖已存在的环境变量
	err := godotenv.Load(EnvFilePath)
	if err != nil {
		// .env 文件不存在或读取失败
		// 这是正常情况,不需要报错
		// 生产环境通常不使用 .env 文件

		// 调试: 打印到 stderr 以便诊断
		fmt.Fprintf(os.Stderr, "[DEBUG] .env file not loaded: %v\n", err)
		return
	}

	// .env 文件加载成功
	// 调试: 打印成功信息
	fmt.Fprintf(os.Stderr, "[DEBUG] .env file loaded successfully\n")

	// 调试: 打印一些关键环境变量
	fmt.Fprintf(os.Stderr, "[DEBUG] DB_DRIVER=%s\n", os.Getenv("DB_DRIVER"))
	fmt.Fprintf(os.Stderr, "[DEBUG] REDIS_HOST=%s\n", os.Getenv("REDIS_HOST"))
	fmt.Fprintf(os.Stderr, "[DEBUG] DB_HOST=%s\n", os.Getenv("DB_HOST"))
	fmt.Fprintf(os.Stderr, "[DEBUG] REDIS_ENABLED=%s\n", os.Getenv("REDIS_ENABLED"))
	fmt.Fprintf(os.Stderr, "[DEBUG] REI_APP_TEST=%s\n", os.Getenv("REI_APP_TEST"))
}

// OverrideWithEnv 使用环境变量覆盖配置
// 优先级: 环境变量 > config.yaml
//
// 参数:
//
//	cfg: 从 config.yaml 加载的配置
//
// 工作流程:
//  1. 检查每个支持的环境变量
//  2. 如果环境变量存在,使用其值覆盖配置
//  3. 如果环境变量不存在,保持 config.yaml 的值
//
// 使用示例:
//
//	config := loadFromYaml()
//	OverrideWithEnv(config)
//	// 此时 config 中的值可能已被环境变量覆盖
func OverrideWithEnv(cfg *Config) {
	// 调试: 显示开始覆盖配置
	fmt.Fprintf(os.Stderr, "[DEBUG] OverrideWithEnv: starting environment variable override\n")

	// 数据库配置
	overrideDatabaseConfig(&cfg.Database)

	// Redis 配置
	overrideRedisConfig(&cfg.Redis)

	// 服务器配置
	overrideServerConfig(&cfg.Server)

	// 日志配置
	overrideLoggerConfig(&cfg.Logger)

	// 国际化配置
	overrideI18nConfig(&cfg.I18n)

	// 调试: 显示覆盖后的值
	fmt.Fprintf(os.Stderr, "[DEBUG] After override - DB_DRIVER=%s, DB_HOST=%s, REDIS_ENABLED=%v\n",
		cfg.Database.Driver, cfg.Database.Host, cfg.Redis.Enabled)
}
