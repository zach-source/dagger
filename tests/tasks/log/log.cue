package main

import (
	"dagger.io/dagger"
	"dagger.io/dagger/core"
)

dagger.#Plan & {
	actions: {
		test: {
			helloWorld: core.#Log & {
				message: "hello world"
			}

			fields: core.#Log & {
				message: "a message"
				fields: {
					"hello": "world"
				}
			}

			moreFields: core.#Log & {
				message: "a message"
				fields: {
					"hello": "world"
					"fun": "1"
					"foo": "bar"
				}
			}

			warnLevel: core.#Log & {
				message: "a message"
				fields: {
					"hello": "world"
					"fun": "1"
					"foo": "bar"
				}
				level: core.#WarnLevel
			}
		}
	}
}
