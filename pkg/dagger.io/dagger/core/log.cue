package core

#Log: {
	$dagger: task: _name: "Log"

	// Log level; defaults to info
	level: #Level | *#InfoLevel

	// Message to send
	message: string

	fields?: [string]: string
}

#Level: #TraceLevel | #DebugLevel | #InfoLevel | #WarnLevel | #ErrorLevel | #FatalLevel | #PanicLevel

#TraceLevel: "trace"
#DebugLevel: "debug"
#InfoLevel: "info"
#WarnLevel: "warn"
#ErrorLevel: "error"
#FatalLevel: "fatal"
#PanicLevel: "panic"
