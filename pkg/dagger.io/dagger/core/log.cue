package core

import "dagger.io/dagger"

#Log: {
	$dagger: task: _name: "Log"

	input: dagger.#Logger

	// Log level; if not available uses the logger's default
	level?: dagger.#Level

	// Message to send
	message: string

  // Additional fields to display
	fields?: [string]: _
}

#TraceLevel: dagger.#TraceLevel
#DebugLevel: dagger.#DebugLevel
#InfoLevel: dagger.#InfoLevel
#WarnLevel: dagger.#WarnLevel
#ErrorLevel: dagger.#ErrorLevel
#FatalLevel: dagger.#Fatallevel
#PanicLevel: dagger.#PanicLevel
