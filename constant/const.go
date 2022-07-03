package constant

const (
	// JaegerAgent jaeger agent 连接地址，负责将服务端的 span 通过 UDP 传输到 agent
	JaegerAgent = "127.0.0.1:5775"

	ClientTracer = "jaeger-client"
	ServerTracer = "jaeger-server"
	NoticeTracer = "jaeger-notice"

	ClientMicroServer = "micro-client"
	ServerMicroServer = "micro-server"
	NoticeMicroServer = "micro-notice"
)


