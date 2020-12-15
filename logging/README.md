# Logging

Integrate context with your logs.

## Get

`go get -u github.com/anz-bank/pkg/logging`

## Overview

### The problem

I want to add contextual data to my logs, but all other packages out there require many
function calls just to perform one log.

### The solution

Use a logger that understands context.

### What this package provides

1. Easy api to add log context to a context and include it in your logs down the line.
2. A standardized approach to how this is achieved, so it can be used across packages
and projects.
3. Integration with many third party libraries that take loggers
